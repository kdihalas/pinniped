// Copyright 2020-2025 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "go.pinniped.dev/generated/1.32/apis/concierge/authentication/v1alpha1"
	authenticationv1alpha1 "go.pinniped.dev/generated/1.32/client/concierge/clientset/versioned/typed/authentication/v1alpha1"
	gentype "k8s.io/client-go/gentype"
)

// fakeWebhookAuthenticators implements WebhookAuthenticatorInterface
type fakeWebhookAuthenticators struct {
	*gentype.FakeClientWithList[*v1alpha1.WebhookAuthenticator, *v1alpha1.WebhookAuthenticatorList]
	Fake *FakeAuthenticationV1alpha1
}

func newFakeWebhookAuthenticators(fake *FakeAuthenticationV1alpha1) authenticationv1alpha1.WebhookAuthenticatorInterface {
	return &fakeWebhookAuthenticators{
		gentype.NewFakeClientWithList[*v1alpha1.WebhookAuthenticator, *v1alpha1.WebhookAuthenticatorList](
			fake.Fake,
			"",
			v1alpha1.SchemeGroupVersion.WithResource("webhookauthenticators"),
			v1alpha1.SchemeGroupVersion.WithKind("WebhookAuthenticator"),
			func() *v1alpha1.WebhookAuthenticator { return &v1alpha1.WebhookAuthenticator{} },
			func() *v1alpha1.WebhookAuthenticatorList { return &v1alpha1.WebhookAuthenticatorList{} },
			func(dst, src *v1alpha1.WebhookAuthenticatorList) { dst.ListMeta = src.ListMeta },
			func(list *v1alpha1.WebhookAuthenticatorList) []*v1alpha1.WebhookAuthenticator {
				return gentype.ToPointerSlice(list.Items)
			},
			func(list *v1alpha1.WebhookAuthenticatorList, items []*v1alpha1.WebhookAuthenticator) {
				list.Items = gentype.FromPointerSlice(items)
			},
		),
		fake,
	}
}
