// Copyright 2020-2025 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package callback provides a handler for the OIDC callback endpoint.
package callback

import (
	"mime"
	"net/http"
	"net/url"
	"strings"

	"github.com/ory/fosite"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/audit"

	"go.pinniped.dev/internal/auditevent"
	"go.pinniped.dev/internal/federationdomain/downstreamsession"
	"go.pinniped.dev/internal/federationdomain/federationdomainproviders"
	"go.pinniped.dev/internal/federationdomain/formposthtml"
	"go.pinniped.dev/internal/federationdomain/oidc"
	"go.pinniped.dev/internal/httputil/httperr"
	"go.pinniped.dev/internal/httputil/securityheader"
	"go.pinniped.dev/internal/plog"
)

func paramsSafeToLog() sets.Set[string] {
	return sets.New[string](
		// Due to https://datatracker.ietf.org/doc/html/rfc6749#section-4.1.2.1,
		// authorize errors can have these parameters, which should not contain PII or secrets and are safe to log.
		"error", "error_description", "error_uri",
		// Note that this endpoint also receives 'code' and 'state' params, which are not safe to log.
	)
}

func NewHandler(
	upstreamIDPs federationdomainproviders.FederationDomainIdentityProvidersFinderI,
	oauthHelper fosite.OAuth2Provider,
	stateDecoder, cookieDecoder oidc.Decoder,
	redirectURI string,
	auditLogger plog.AuditLogger,
) http.Handler {
	handler := httperr.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		if err := auditLogger.AuditRequestParams(r, paramsSafeToLog()); err != nil {
			plog.DebugErr("error parsing callback request params", err)
			return httperr.New(http.StatusBadRequest, "error parsing request params")
		}

		decodedState, err := validateRequest(r, stateDecoder, cookieDecoder, auditLogger)
		if err != nil {
			return err
		}

		idp, err := upstreamIDPs.FindUpstreamIDPByDisplayName(decodedState.UpstreamName)
		if err != nil || idp == nil {
			plog.Warning("upstream provider not found")
			return httperr.New(http.StatusUnprocessableEntity, "upstream provider not found")
		}

		auditLogger.Audit(auditevent.UsingUpstreamIDP, &plog.AuditParams{
			ReqCtx: r.Context(),
			KeysAndValues: []any{
				"displayName", idp.GetDisplayName(),
				"resourceName", idp.GetProvider().GetResourceName(),
				"resourceUID", idp.GetProvider().GetResourceUID(),
				"type", idp.GetSessionProviderType(),
			},
		})

		downstreamAuthParams, err := url.ParseQuery(decodedState.AuthParams)
		if err != nil {
			plog.Error("error reading state downstream auth params", err)
			return httperr.New(http.StatusBadRequest, "error reading state downstream auth params")
		}

		// Recreate enough of the original authorize request, so we can pass it to NewAuthorizeRequest().
		reconstitutedAuthRequest := &http.Request{Form: downstreamAuthParams}
		authorizeRequester, err := oauthHelper.NewAuthorizeRequest(r.Context(), reconstitutedAuthRequest)
		if err != nil {
			plog.Error("error using state downstream auth params", err,
				"identityProviderDisplayName", idp.GetDisplayName(),
				"identityProviderResourceName", idp.GetProvider().GetResourceName(),
				"supervisorCallbackURL", redirectURI,
				"fositeErr", oidc.FositeErrorForLog(err))
			return httperr.New(http.StatusBadRequest, "error using state downstream auth params")
		}

		// Automatically grant certain scopes, but only if they were requested.
		// This is instead of asking the user to approve these scopes. Note that `NewAuthorizeRequest` would have returned
		// an error if the client requested a scope that they are not allowed to request, so we don't need to worry about that here.
		downstreamsession.AutoApproveScopes(authorizeRequester)

		identity, loginExtras, err := idp.LoginFromCallback(r.Context(), authcode(r), decodedState.PKCECode, decodedState.Nonce, redirectURI)
		if err != nil {
			plog.WarningErr("unable to complete login from callback", err,
				"identityProviderDisplayName", idp.GetDisplayName(),
				"identityProviderResourceName", idp.GetProvider().GetResourceName(),
				"supervisorCallbackURL", redirectURI)
			return err
		}

		session, err := downstreamsession.NewPinnipedSession(r.Context(), auditLogger, &downstreamsession.SessionConfig{
			UpstreamIdentity:    identity,
			UpstreamLoginExtras: loginExtras,
			ClientID:            authorizeRequester.GetClient().GetID(),
			GrantedScopes:       authorizeRequester.GetGrantedScopes(),
			IdentityProvider:    idp,
			SessionIDGetter:     authorizeRequester,
		})
		if err != nil {
			plog.WarningErr("unable to create a Pinniped session", err,
				"identityProviderDisplayName", idp.GetDisplayName(),
				"identityProviderResourceName", idp.GetProvider().GetResourceName(),
				"supervisorCallbackURL", redirectURI)
			return httperr.Wrap(http.StatusUnprocessableEntity, err.Error(), err)
		}

		authorizeResponder, err := oauthHelper.NewAuthorizeResponse(r.Context(), authorizeRequester, session)
		if err != nil {
			plog.WarningErr("error while generating and saving authcode", err,
				"identityProviderDisplayName", idp.GetDisplayName(),
				"identityProviderResourceName", idp.GetProvider().GetResourceName(),
				"supervisorCallbackURL", redirectURI,
				"fositeErr", oidc.FositeErrorForLog(err))
			return httperr.Wrap(http.StatusInternalServerError, "error while generating and saving authcode", err)
		}

		oauthHelper.WriteAuthorizeResponse(r.Context(), w, authorizeRequester, authorizeResponder)

		return nil
	})
	return securityheader.WrapWithCustomCSP(handler, formposthtml.ContentSecurityPolicy())
}

func authcode(r *http.Request) string {
	return r.FormValue("code")
}

func validateRequest(r *http.Request, stateDecoder, cookieDecoder oidc.Decoder, auditLogger plog.AuditLogger) (*oidc.UpstreamStateParamData, error) {
	// An upstream OIDC provider will typically return a redirect with the authcode,
	// which the user's browser will follow and therefore perform a GET to this endpoint.
	// See https://openid.net/specs/openid-connect-core-1_0.html#AuthResponse.
	//
	// If the upstream provider supports using response_mode=form_post,
	// and if the admin configured OIDCIdentityProvider.spec.authorizationConfig.additionalAuthorizeParameters
	// to include response_mode=form_post, then the user's browser will POST the authcode to this endpoint.
	// See https://openid.net/specs/oauth-v2-form-post-response-mode-1_0.html.
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		return nil, httperr.Newf(http.StatusMethodNotAllowed, "%s (try GET or POST)", r.Method)
	}

	if r.Method == http.MethodPost {
		// https://openid.net/specs/oauth-v2-form-post-response-mode-1_0.html says
		// "the result parameters being encoded in the body using the application/x-www-form-urlencoded format",
		// so only accept that format. Since "multipart/form-data" is intended for binary data, there's no need to
		// support it here.
		contentType := r.Header.Get("Content-Type")
		if contentType == "" {
			return nil, httperr.New(http.StatusUnsupportedMediaType, "no Content-Type header (try Content-Type: application/x-www-form-urlencoded)")
		}
		parsedContentType, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			// Note that it should not be possible to reach this line since	AuditRequestParams() is already
			// parsing the form and handling unparseable form types before we get here. This is just defensive coding.
			return nil, httperr.Newf(http.StatusUnsupportedMediaType, "cannot parse Content-Type: %s (try Content-Type: application/x-www-form-urlencoded)", contentType)
		}
		if parsedContentType != "application/x-www-form-urlencoded" {
			return nil, httperr.Newf(http.StatusUnsupportedMediaType, "%s (try application/x-www-form-urlencoded)", contentType)
		}
	}

	encodedState, decodedState, err := oidc.ReadStateParamAndValidateCSRFCookie(r, cookieDecoder, stateDecoder)
	if err != nil {
		plog.InfoErr("state or CSRF error", err)
		return nil, err
	}

	auditLogger.Audit(auditevent.AuthorizeIDFromParameters, &plog.AuditParams{
		ReqCtx:        r.Context(),
		KeysAndValues: []any{"authorizeID", encodedState.AuthorizeID()},
	})

	if authcode(r) == "" {
		plog.Info("code param not found")
		return nil, httperr.New(http.StatusBadRequest, errorMsgForNoCodeParam(r))
	}

	return decodedState, nil
}

func errorMsgForNoCodeParam(r *http.Request) string {
	msg := strings.Builder{}

	msg.WriteString("code param not found\n\n")

	errorParam, hasError := r.Form["error"]
	errorDescParam, hasErrorDesc := r.Form["error_description"]
	errorURIParam, hasErrorURI := r.Form["error_uri"]

	if hasError {
		msg.WriteString("error from external identity provider: ")
		msg.WriteString(errorParam[0])
		msg.WriteByte('\n')
	}
	if hasErrorDesc {
		msg.WriteString("error_description from external identity provider: ")
		msg.WriteString(errorDescParam[0])
		msg.WriteByte('\n')
	}
	if hasErrorURI {
		msg.WriteString("error_uri from external identity provider: ")
		msg.WriteString(errorURIParam[0])
		msg.WriteByte('\n')
	}
	if !hasError && !hasErrorDesc && !hasErrorURI {
		msg.WriteString("Something went wrong with your authentication attempt at your external identity provider.\n")
	}

	msg.WriteByte('\n')
	msg.WriteString("Pinniped AuditID: ")
	msg.WriteString(audit.GetAuditIDTruncated(r.Context()))

	return msg.String()
}
