// Copyright 2020-2025 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"

	v1alpha1 "go.pinniped.dev/generated/1.31/apis/concierge/authentication/v1alpha1"
	scheme "go.pinniped.dev/generated/1.31/client/concierge/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
)

// WebhookAuthenticatorsGetter has a method to return a WebhookAuthenticatorInterface.
// A group's client should implement this interface.
type WebhookAuthenticatorsGetter interface {
	WebhookAuthenticators() WebhookAuthenticatorInterface
}

// WebhookAuthenticatorInterface has methods to work with WebhookAuthenticator resources.
type WebhookAuthenticatorInterface interface {
	Create(ctx context.Context, webhookAuthenticator *v1alpha1.WebhookAuthenticator, opts v1.CreateOptions) (*v1alpha1.WebhookAuthenticator, error)
	Update(ctx context.Context, webhookAuthenticator *v1alpha1.WebhookAuthenticator, opts v1.UpdateOptions) (*v1alpha1.WebhookAuthenticator, error)
	// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
	UpdateStatus(ctx context.Context, webhookAuthenticator *v1alpha1.WebhookAuthenticator, opts v1.UpdateOptions) (*v1alpha1.WebhookAuthenticator, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.WebhookAuthenticator, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.WebhookAuthenticatorList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.WebhookAuthenticator, err error)
	WebhookAuthenticatorExpansion
}

// webhookAuthenticators implements WebhookAuthenticatorInterface
type webhookAuthenticators struct {
	*gentype.ClientWithList[*v1alpha1.WebhookAuthenticator, *v1alpha1.WebhookAuthenticatorList]
}

// newWebhookAuthenticators returns a WebhookAuthenticators
func newWebhookAuthenticators(c *AuthenticationV1alpha1Client) *webhookAuthenticators {
	return &webhookAuthenticators{
		gentype.NewClientWithList[*v1alpha1.WebhookAuthenticator, *v1alpha1.WebhookAuthenticatorList](
			"webhookauthenticators",
			c.RESTClient(),
			scheme.ParameterCodec,
			"",
			func() *v1alpha1.WebhookAuthenticator { return &v1alpha1.WebhookAuthenticator{} },
			func() *v1alpha1.WebhookAuthenticatorList { return &v1alpha1.WebhookAuthenticatorList{} }),
	}
}
