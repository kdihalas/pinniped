// Copyright 2024-2025 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package integration

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"go.pinniped.dev/internal/certauthority"
	"go.pinniped.dev/internal/here"
	"go.pinniped.dev/test/testlib"
)

// TestTLSSpecValidationConcierge_Parallel tests kubebuilder and status condition validation
// on the TLSSpec in Pinniped concierge WebhookAuthenticator and JWTAuthenticator CRDs.
func TestTLSSpecValidationConcierge_Parallel(t *testing.T) {
	env := testlib.IntegrationEnv(t)

	ca, err := certauthority.New("pinniped-test", 24*time.Hour)
	require.NoError(t, err)
	indentedCAPEM := indentForHeredoc(string(ca.Bundle()))

	webhookAuthenticatorYamlTemplate := here.Doc(`
		apiVersion: authentication.concierge.%s/v1alpha1
		kind: WebhookAuthenticator
		metadata:
			name: %s
		spec:
			endpoint: %s
			%s
	`)

	jwtAuthenticatorYamlTemplate := here.Doc(`
		apiVersion: authentication.concierge.%s/v1alpha1
		kind: JWTAuthenticator
		metadata:
			name: %s
		spec:
			issuer: %s
			audience: some-audience
			%s
	`)

	testCases := []struct {
		name string

		tlsYAML func(secretOrConfigmapName string) string

		secretOrConfigmapKind     string
		secretType                string
		secretOrConfigmapDataYAML string

		wantErrorSnippets            []string
		wantTLSValidConditionMessage func(namespace string, secretOrConfigmapName string) string
	}{
		{
			name: "should disallow certificate authority data source with missing name",
			tlsYAML: func(secretOrConfigmapName string) string {
				return here.Doc(`
					tls:
						certificateAuthorityDataSource:
							kind: Secret
							key: bar
				`)
			},
			wantErrorSnippets: []string{`The %s "%s" is invalid: spec.tls.certificateAuthorityDataSource.name: Required value`},
		},
		{
			name: "should disallow certificate authority data source with empty value for name",
			tlsYAML: func(secretOrConfigmapName string) string {
				return here.Doc(`
					tls:
						certificateAuthorityDataSource:
							kind: Secret
							name: ""
							key: bar
				`)
			},
			wantErrorSnippets: []string{`The %s "%s" is invalid: spec.tls.certificateAuthorityDataSource.name: Invalid value: "": spec.tls.certificateAuthorityDataSource.name in body should be at least 1 chars long`},
		},
		{
			name: "should disallow certificate authority data source with missing key",
			tlsYAML: func(secretOrConfigmapName string) string {
				return here.Doc(`
					tls:
						certificateAuthorityDataSource:
							kind: Secret
							name: foo
				`)
			},
			wantErrorSnippets: []string{`The %s "%s" is invalid: spec.tls.certificateAuthorityDataSource.key: Required value`},
		},
		{
			name: "should disallow certificate authority data source with empty value for key",
			tlsYAML: func(secretOrConfigmapName string) string {
				return here.Doc(`
					tls:
						certificateAuthorityDataSource:
							kind: Secret
							name: foo
							key: ""
				`)
			},
			wantErrorSnippets: []string{`The %s "%s" is invalid: spec.tls.certificateAuthorityDataSource.key: Invalid value: "": spec.tls.certificateAuthorityDataSource.key in body should be at least 1 chars long`},
		},
		{
			name: "should disallow certificate authority data source with missing kind",
			tlsYAML: func(secretOrConfigmapName string) string {
				return here.Doc(`
					tls:
						certificateAuthorityDataSource:
							name: foo
							key: bar
				`)
			},
			wantErrorSnippets: []string{`The %s "%s" is invalid: spec.tls.certificateAuthorityDataSource.kind: Required value`},
		},
		{
			name: "should disallow certificate authority data source with empty value for kind",
			tlsYAML: func(secretOrConfigmapName string) string {
				return here.Doc(`
					tls:
						certificateAuthorityDataSource:
							kind: ""
							name: foo
							key: bar
				`)
			},
			wantErrorSnippets: []string{`The %s "%s" is invalid: spec.tls.certificateAuthorityDataSource.kind: Unsupported value: "": supported values: "Secret", "ConfigMap"`},
		},
		{
			name: "should disallow certificate authority data source with invalid kind",
			tlsYAML: func(secretOrConfigmapName string) string {
				return here.Doc(`
					tls:
						certificateAuthorityDataSource:
							kind: sorcery
							name: foo
							key: bar
				`)
			},
			wantErrorSnippets: []string{`The %s "%s" is invalid: spec.tls.certificateAuthorityDataSource.kind: Unsupported value: "sorcery": supported values: "Secret", "ConfigMap"`},
		},
		{
			name: "should get error condition when using both fields of the tls spec",
			tlsYAML: func(secretOrConfigmapName string) string {
				return here.Doc(`
					tls:
						certificateAuthorityData: "some CA data"
						certificateAuthorityDataSource:
							kind: ConfigMap
							name: foo
							key: bar
				`)
			},
			wantErrorSnippets: nil,
			wantTLSValidConditionMessage: func(namespace string, secretOrConfigmapName string) string {
				return "spec.tls is invalid: both tls.certificateAuthorityDataSource and tls.certificateAuthorityData provided"
			},
		},
		{
			name: "should get error condition when certificateAuthorityData is not base64 data",
			tlsYAML: func(secretOrConfigmapName string) string {
				return here.Doc(`
					tls:
						certificateAuthorityData: "this is not base64 encoded"
				`)
			},
			wantErrorSnippets: nil,
			wantTLSValidConditionMessage: func(namespace string, secretOrConfigmapName string) string {
				return `spec.tls.certificateAuthorityData is invalid: illegal base64 data at input byte 4`
			},
		},
		{
			name: "should get error condition when certificateAuthorityData does not contain PEM data",
			tlsYAML: func(secretOrConfigmapName string) string {
				return here.Docf(`
					tls:
						certificateAuthorityData: "%s"
				`, base64.StdEncoding.EncodeToString([]byte("this is not PEM data")))
			},
			wantErrorSnippets: nil,
			wantTLSValidConditionMessage: func(namespace string, secretOrConfigmapName string) string {
				return `spec.tls.certificateAuthorityData is invalid: no base64-encoded PEM certificates found in 28 bytes of data (PEM certificates must begin with "-----BEGIN CERTIFICATE-----")`
			},
		},
		{
			name: "should get error condition when using a ConfigMap source and the ConfigMap does not exist",
			tlsYAML: func(secretOrConfigmapName string) string {
				return here.Doc(`
					tls:
						certificateAuthorityDataSource:
							kind: ConfigMap
							name: this-cm-does-not-exist
							key: bar
				`)
			},
			wantErrorSnippets: nil,
			wantTLSValidConditionMessage: func(namespace string, secretOrConfigmapName string) string {
				return fmt.Sprintf(
					`spec.tls.certificateAuthorityDataSource is invalid: failed to get configmap "%s/this-cm-does-not-exist": configmap "this-cm-does-not-exist" not found`,
					namespace)
			},
		},
		{
			name: "should get error condition when using a Secret source and the Secret does not exist",
			tlsYAML: func(secretOrConfigmapName string) string {
				return here.Doc(`
					tls:
						certificateAuthorityDataSource:
							kind: Secret
							name: this-secret-does-not-exist
							key: bar
				`)
			},
			wantErrorSnippets: nil,
			wantTLSValidConditionMessage: func(namespace string, secretOrConfigmapName string) string {
				return fmt.Sprintf(
					`spec.tls.certificateAuthorityDataSource is invalid: failed to get secret "%s/this-secret-does-not-exist": secret "this-secret-does-not-exist" not found`,
					namespace)
			},
		},
		{
			name:                  "should get error condition when using a Secret source and the Secret is the wrong type",
			secretOrConfigmapKind: "Secret",
			secretType:            "wrong-type",
			secretOrConfigmapDataYAML: here.Doc(`
				bar: "does not matter for this test"
			`),
			tlsYAML: func(secretOrConfigmapName string) string {
				return here.Docf(`
					tls:
						certificateAuthorityDataSource:
							kind: Secret
							name: %s
							key: bar
				`, secretOrConfigmapName)
			},
			wantErrorSnippets: nil,
			wantTLSValidConditionMessage: func(namespace string, secretOrConfigmapName string) string {
				return fmt.Sprintf(
					`spec.tls.certificateAuthorityDataSource is invalid: secret "%s/%s" of type "wrong-type" cannot be used as a certificate authority data source`,
					namespace, secretOrConfigmapName)
			},
		},
		{
			name:                  "should get error condition when using a Secret source and the key does not exist",
			secretOrConfigmapKind: "Secret",
			secretType:            string(corev1.SecretTypeOpaque),
			secretOrConfigmapDataYAML: here.Doc(`
				foo: "foo is the wrong key"
			`),
			tlsYAML: func(secretOrConfigmapName string) string {
				return here.Docf(`
					tls:
						certificateAuthorityDataSource:
							kind: Secret
							name: %s
							key: bar
				`, secretOrConfigmapName)
			},
			wantErrorSnippets: nil,
			wantTLSValidConditionMessage: func(namespace string, secretOrConfigmapName string) string {
				return fmt.Sprintf(
					`spec.tls.certificateAuthorityDataSource is invalid: key "bar" not found in secret "%s/%s"`,
					namespace, secretOrConfigmapName)
			},
		},
		{
			name:                  "should get error condition when using a ConfigMap source and the key does not exist",
			secretOrConfigmapKind: "ConfigMap",
			secretOrConfigmapDataYAML: here.Doc(`
				foo: "foo is the wrong key"
			`),
			tlsYAML: func(secretOrConfigmapName string) string {
				return here.Docf(`
					tls:
						certificateAuthorityDataSource:
							kind: ConfigMap
							name: %s
							key: bar
				`, secretOrConfigmapName)
			},
			wantErrorSnippets: nil,
			wantTLSValidConditionMessage: func(namespace string, secretOrConfigmapName string) string {
				return fmt.Sprintf(
					`spec.tls.certificateAuthorityDataSource is invalid: key "bar" not found in configmap "%s/%s"`,
					namespace, secretOrConfigmapName)
			},
		},
		{
			name:                  "should get error condition when using a Secret source and the key has an empty value",
			secretOrConfigmapKind: "Secret",
			secretType:            string(corev1.SecretTypeOpaque),
			secretOrConfigmapDataYAML: here.Doc(`
				bar: ""
			`),
			tlsYAML: func(secretOrConfigmapName string) string {
				return here.Docf(`
					tls:
						certificateAuthorityDataSource:
							kind: Secret
							name: %s
							key: bar
				`, secretOrConfigmapName)
			},
			wantErrorSnippets: nil,
			wantTLSValidConditionMessage: func(namespace string, secretOrConfigmapName string) string {
				return fmt.Sprintf(
					`spec.tls.certificateAuthorityDataSource is invalid: key "bar" has empty value in secret "%s/%s"`,
					namespace, secretOrConfigmapName)
			},
		},
		{
			name:                  "should get error condition when using a ConfigMap source and the key has an empty value",
			secretOrConfigmapKind: "ConfigMap",
			secretOrConfigmapDataYAML: here.Doc(`
				bar: ""
			`),
			tlsYAML: func(secretOrConfigmapName string) string {
				return here.Docf(`
					tls:
						certificateAuthorityDataSource:
							kind: ConfigMap
							name: %s
							key: bar
				`, secretOrConfigmapName)
			},
			wantErrorSnippets: nil,
			wantTLSValidConditionMessage: func(namespace string, secretOrConfigmapName string) string {
				return fmt.Sprintf(
					`spec.tls.certificateAuthorityDataSource is invalid: key "bar" has empty value in configmap "%s/%s"`,
					namespace, secretOrConfigmapName)
			},
		},
		{
			name:                  "should get error condition when using a Secret source and the Secret contains data which is not in PEM format",
			secretOrConfigmapKind: "Secret",
			secretType:            string(corev1.SecretTypeOpaque),
			secretOrConfigmapDataYAML: here.Doc(`
				bar: "this is not a PEM cert"
			`),
			tlsYAML: func(secretOrConfigmapName string) string {
				return here.Docf(`
					tls:
						certificateAuthorityDataSource:
							kind: Secret
							name: %s
							key: bar
				`, secretOrConfigmapName)
			},
			wantErrorSnippets: nil,
			wantTLSValidConditionMessage: func(namespace string, secretOrConfigmapName string) string {
				return fmt.Sprintf(
					`spec.tls.certificateAuthorityDataSource is invalid: key "bar" with 22 bytes of data in secret "%s/%s" is not a PEM-encoded certificate (PEM certificates must begin with "-----BEGIN CERTIFICATE-----")`,
					namespace, secretOrConfigmapName)
			},
		},
		{
			name:                  "should get error condition when using a ConfigMap source and the ConfigMap contains data which is not in PEM format",
			secretOrConfigmapKind: "ConfigMap",
			secretOrConfigmapDataYAML: here.Doc(`
				bar: "this is not a PEM cert"
			`),
			tlsYAML: func(secretOrConfigmapName string) string {
				return here.Docf(`
					tls:
						certificateAuthorityDataSource:
							kind: ConfigMap
							name: %s
							key: bar
				`, secretOrConfigmapName)
			},
			wantErrorSnippets: nil,
			wantTLSValidConditionMessage: func(namespace string, secretOrConfigmapName string) string {
				return fmt.Sprintf(
					`spec.tls.certificateAuthorityDataSource is invalid: key "bar" with 22 bytes of data in configmap "%s/%s" is not a PEM-encoded certificate (PEM certificates must begin with "-----BEGIN CERTIFICATE-----")`,
					namespace, secretOrConfigmapName)
			},
		},
		{
			name:                  "should create a custom resource passing all validations using a Secret source of type Opaque",
			secretOrConfigmapKind: "Secret",
			secretType:            string(corev1.SecretTypeOpaque),
			secretOrConfigmapDataYAML: here.Docf(`
				bar: |
					%s
			`, indentedCAPEM),
			tlsYAML: func(secretOrConfigmapName string) string {
				return here.Docf(`
					tls:
						certificateAuthorityDataSource:
							kind: Secret
							name: %s
							key: bar
				`, secretOrConfigmapName)
			},
			wantErrorSnippets: nil,
			wantTLSValidConditionMessage: func(namespace string, secretOrConfigmapName string) string {
				return `spec.tls is valid: using configured CA bundle`
			},
		},
		{
			name:                  "should create a custom resource passing all validations using a Secret source of type tls",
			secretOrConfigmapKind: "Secret",
			secretType:            string(corev1.SecretTypeTLS),
			secretOrConfigmapDataYAML: here.Docf(`
				tls.crt: foo
				tls.key: foo
				bar: |
					%s
			`, indentedCAPEM),
			tlsYAML: func(secretOrConfigmapName string) string {
				return here.Docf(`
					tls:
						certificateAuthorityDataSource:
							kind: Secret
							name: %s
							key: bar
				`, secretOrConfigmapName)
			},
			wantErrorSnippets: nil,
			wantTLSValidConditionMessage: func(namespace string, secretOrConfigmapName string) string {
				return `spec.tls is valid: using configured CA bundle`
			},
		},
		{
			name:                  "should create a custom resource passing all validations using a ConfigMap source",
			secretOrConfigmapKind: "ConfigMap",
			secretOrConfigmapDataYAML: here.Docf(`
				bar: |
					%s
			`, indentedCAPEM),
			tlsYAML: func(secretOrConfigmapName string) string {
				return here.Docf(`
					tls:
						certificateAuthorityDataSource:
							kind: ConfigMap
							name: %s
							key: bar
				`, secretOrConfigmapName)
			},
			wantErrorSnippets: nil,
			wantTLSValidConditionMessage: func(namespace string, secretOrConfigmapName string) string {
				return `spec.tls is valid: using configured CA bundle`
			},
		},
		{
			name:              "should create a custom resource without any tls spec",
			tlsYAML:           func(secretOrConfigmapName string) string { return "" },
			wantErrorSnippets: nil,
			wantTLSValidConditionMessage: func(namespace string, secretOrConfigmapName string) string {
				return "spec.tls is valid: no TLS configuration provided: using default root CA bundle from container image"
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			t.Run("apply webhook authenticator", func(t *testing.T) {
				resourceName := "test-webhook-authenticator-" + testlib.RandHex(t, 7)

				secretOrConfigmapResourceName := createSecretOrConfigMapFromData(t,
					resourceName,
					env.ConciergeNamespace,
					tc.secretOrConfigmapKind,
					tc.secretType,
					tc.secretOrConfigmapDataYAML,
				)

				yamlBytes := []byte(fmt.Sprintf(webhookAuthenticatorYamlTemplate,
					env.APIGroupSuffix, resourceName, env.TestWebhook.Endpoint,
					indentForHeredoc(tc.tlsYAML(secretOrConfigmapResourceName))))

				stdOut, stdErr, err := performKubectlApply(t, resourceName, yamlBytes)
				requireKubectlApplyResult(t, stdOut, stdErr, err,
					fmt.Sprintf(`webhookauthenticator.authentication.concierge.%s`, env.APIGroupSuffix),
					tc.wantErrorSnippets,
					"WebhookAuthenticator",
					resourceName,
				)

				if tc.wantErrorSnippets == nil {
					requireTLSValidConditionMessageOnResource(t,
						resourceName,
						env.ConciergeNamespace,
						"WebhookAuthenticator",
						tc.wantTLSValidConditionMessage(env.ConciergeNamespace, secretOrConfigmapResourceName),
					)
				}
			})

			t.Run("apply jwt authenticator", func(t *testing.T) {
				supervisorIssuer := env.InferSupervisorIssuerURL(t)

				resourceName := "test-jwt-authenticator-" + testlib.RandHex(t, 7)

				secretOrConfigmapResourceName := createSecretOrConfigMapFromData(t,
					resourceName,
					env.ConciergeNamespace,
					tc.secretOrConfigmapKind,
					tc.secretType,
					tc.secretOrConfigmapDataYAML,
				)

				yamlBytes := []byte(fmt.Sprintf(jwtAuthenticatorYamlTemplate,
					env.APIGroupSuffix, resourceName, supervisorIssuer.Issuer(),
					indentForHeredoc(tc.tlsYAML(secretOrConfigmapResourceName))))

				stdOut, stdErr, err := performKubectlApply(t, resourceName, yamlBytes)
				requireKubectlApplyResult(t, stdOut, stdErr, err,
					fmt.Sprintf(`jwtauthenticator.authentication.concierge.%s`, env.APIGroupSuffix),
					tc.wantErrorSnippets,
					"JWTAuthenticator",
					resourceName,
				)

				if tc.wantErrorSnippets == nil {
					requireTLSValidConditionMessageOnResource(t,
						resourceName,
						env.ConciergeNamespace,
						"JWTAuthenticator",
						tc.wantTLSValidConditionMessage(env.ConciergeNamespace, secretOrConfigmapResourceName),
					)
				}
			})
		})
	}
}

func indentForHeredoc(s string) string {
	// Further indent every line except for the first line by four spaces.
	// Use four spaces because that's what here.Doc uses.
	// Do not indent the first line because the template already indents it.
	return strings.ReplaceAll(s, "\n", "\n    ")
}

func requireTLSValidConditionMessageOnResource(t *testing.T, resourceName string, namespace string, resourceType string, wantMessage string) {
	t.Helper()

	require.NotEmpty(t, resourceName, "bad test setup: empty resourceName")
	require.NotEmpty(t, resourceType, "bad test setup: empty resourceType")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	t.Cleanup(cancel)

	conciergeAuthClient := testlib.NewConciergeClientset(t).AuthenticationV1alpha1()
	supervisorIDPClient := testlib.NewSupervisorClientset(t).IDPV1alpha1()

	switch resourceType {
	case "JWTAuthenticator":
		testlib.RequireEventuallyf(t, func(requireEventually *require.Assertions) {
			got, err := conciergeAuthClient.JWTAuthenticators().Get(ctx, resourceName, metav1.GetOptions{})
			requireEventually.NoError(err)
			requireConditionHasMessage(requireEventually, got.Status.Conditions, "TLSConfigurationValid", wantMessage)
		}, 10*time.Second, 1*time.Second, "expected resource %s to have condition message %q", resourceName, wantMessage)
	case "WebhookAuthenticator":
		testlib.RequireEventuallyf(t, func(requireEventually *require.Assertions) {
			got, err := conciergeAuthClient.WebhookAuthenticators().Get(ctx, resourceName, metav1.GetOptions{})
			requireEventually.NoError(err)
			requireConditionHasMessage(requireEventually, got.Status.Conditions, "TLSConfigurationValid", wantMessage)
		}, 10*time.Second, 1*time.Second, "expected resource %s to have condition message %q", resourceName, wantMessage)
	case "OIDCIdentityProvider":
		require.NotEmpty(t, namespace, "bad test setup: empty namespace")
		testlib.RequireEventuallyf(t, func(requireEventually *require.Assertions) {
			got, err := supervisorIDPClient.OIDCIdentityProviders(namespace).Get(ctx, resourceName, metav1.GetOptions{})
			requireEventually.NoError(err)
			requireConditionHasMessage(requireEventually, got.Status.Conditions, "TLSConfigurationValid", wantMessage)
		}, 10*time.Second, 1*time.Second, "expected resource %s to have condition message %q", resourceName, wantMessage)
	case "LDAPIdentityProvider":
		require.NotEmpty(t, namespace, "bad test setup: empty namespace")
		testlib.RequireEventuallyf(t, func(requireEventually *require.Assertions) {
			got, err := supervisorIDPClient.LDAPIdentityProviders(namespace).Get(ctx, resourceName, metav1.GetOptions{})
			requireEventually.NoError(err)
			requireConditionHasMessage(requireEventually, got.Status.Conditions, "TLSConfigurationValid", wantMessage)
		}, 10*time.Second, 1*time.Second, "expected resource %s to have condition message %q", resourceName, wantMessage)
	case "ActiveDirectoryIdentityProvider":
		require.NotEmpty(t, namespace, "bad test setup: empty namespace")
		testlib.RequireEventuallyf(t, func(requireEventually *require.Assertions) {
			got, err := supervisorIDPClient.ActiveDirectoryIdentityProviders(namespace).Get(ctx, resourceName, metav1.GetOptions{})
			requireEventually.NoError(err)
			requireConditionHasMessage(requireEventually, got.Status.Conditions, "TLSConfigurationValid", wantMessage)
		}, 10*time.Second, 1*time.Second, "expected resource %s to have condition message %q", resourceName, wantMessage)
	case "GitHubIdentityProvider":
		require.NotEmpty(t, namespace, "bad test setup: empty namespace")
		testlib.RequireEventuallyf(t, func(requireEventually *require.Assertions) {
			got, err := supervisorIDPClient.GitHubIdentityProviders(namespace).Get(ctx, resourceName, metav1.GetOptions{})
			requireEventually.NoError(err)
			requireConditionHasMessage(requireEventually, got.Status.Conditions, "TLSConfigurationValid", wantMessage)
		}, 10*time.Second, 1*time.Second, "expected resource %s to have condition message %q", resourceName, wantMessage)
	default:
		require.Failf(t, "unexpected resource type", "type %q", resourceType)
	}
}

func requireConditionHasMessage(assertions *require.Assertions, actualConditions []metav1.Condition, conditionType string, wantMessage string) {
	assertions.NotEmpty(actualConditions, "wanted to have conditions but was empty")
	for _, c := range actualConditions {
		if c.Type == conditionType {
			assertions.Equal(wantMessage, c.Message)
			return
		}
	}
	assertions.Failf("did not find condition with expected type",
		"type %q, actual conditions: %#v", conditionType, actualConditions)
}

func createSecretOrConfigMapFromData(
	t *testing.T,
	resourceNameSuffix string,
	namespace string,
	kind string,
	secretType string,
	dataYAML string,
) string {
	t.Helper()

	if kind == "" {
		// Nothing to create.
		return ""
	}

	require.NotEmpty(t, resourceNameSuffix, "bad test setup: empty resourceNameSuffix")
	require.NotEmpty(t, namespace, "bad test setup: empty namespace")

	var resourceYAML string
	lowerKind := strings.ToLower(kind)
	resourceName := lowerKind + "-" + resourceNameSuffix

	// Further indent every line except for the first line by four spaces.
	// Use four spaces because that's what here.Doc uses.
	// Do not indent the first line because the template already indents it.
	indentedDataYAML := strings.ReplaceAll(dataYAML, "\n", "\n    ")

	switch lowerKind {
	case "secret":
		require.NotEmpty(t, secretType, "bad test setup: empty secret type")
		resourceYAML = here.Docf(`
			apiVersion: v1
			kind: Secret
			metadata:
				name: %s
				namespace: %s
			type: %s
			stringData:
				%s
		`, resourceName, namespace, secretType, indentedDataYAML)
	case "configmap":
		resourceYAML = here.Docf(`
			apiVersion: v1
			kind: ConfigMap
			metadata:
				name: %s
				namespace: %s
			data:
				%s
		`, resourceName, namespace, indentedDataYAML)
	default:
		require.Failf(t, "unexpected kind in test setup", "kind was %q", kind)
	}

	stdOut, stdErr, err := performKubectlApply(t, resourceName, []byte(resourceYAML))
	require.NoErrorf(t, err,
		"expected kubectl apply to succeed but got: %s\nstdout: %s\nstderr: %s\nyaml:\n%s",
		err, stdOut, stdErr, resourceYAML)

	return resourceName
}

func performKubectlApply(t *testing.T, resourceName string, yamlBytes []byte) (string, string, error) {
	t.Helper()

	yamlFilepath := filepath.Join(t.TempDir(), fmt.Sprintf("test-perform-kubectl-apply-%s.yaml", resourceName))

	require.NoError(t, os.WriteFile(yamlFilepath, yamlBytes, 0600))

	// Use --validate=false to disable old client-side validations to avoid getting different error messages in Kube 1.24 and older.
	// Note that this also disables validations of unknown and duplicate fields, but that's not what this test is about.
	//nolint:gosec // this is test code.
	cmd := exec.CommandContext(context.Background(), "kubectl", []string{"apply", "--validate=false", "-f", yamlFilepath}...)

	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr
	err := cmd.Run()

	t.Cleanup(func() {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		require.NoError(t, exec.CommandContext(ctx, "kubectl", "delete", "--ignore-not-found", "-f", yamlFilepath).Run())
	})

	return stdOut.String(), stdErr.String(), err
}

func requireKubectlApplyResult(
	t *testing.T,
	kubectlStdOut string,
	kubectlStdErr string,
	kubectlErr error,
	wantSuccessPrefix string,
	wantErrorSnippets []string,
	wantResourceType string,
	wantResourceName string,
) {
	t.Helper()

	if len(wantErrorSnippets) > 0 {
		require.Error(t, kubectlErr)
		actualErrorString := strings.TrimSuffix(kubectlStdErr, "\n")
		for i, snippet := range wantErrorSnippets {
			if i == 0 {
				snippet = fmt.Sprintf(snippet, wantResourceType, wantResourceName)
			}
			require.Contains(t, actualErrorString, snippet)
		}
	} else {
		require.Empty(t, kubectlStdErr)
		require.Regexp(t, regexp.QuoteMeta(wantSuccessPrefix)+regexp.QuoteMeta(fmt.Sprintf("/%s created\n", wantResourceName)), kubectlStdOut)
		require.NoError(t, kubectlErr)
	}
}
