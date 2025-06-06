// Copyright 2020-2025 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "go.pinniped.dev/generated/1.30/apis/concierge/authentication/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeWebhookAuthenticators implements WebhookAuthenticatorInterface
type FakeWebhookAuthenticators struct {
	Fake *FakeAuthenticationV1alpha1
}

var webhookauthenticatorsResource = v1alpha1.SchemeGroupVersion.WithResource("webhookauthenticators")

var webhookauthenticatorsKind = v1alpha1.SchemeGroupVersion.WithKind("WebhookAuthenticator")

// Get takes name of the webhookAuthenticator, and returns the corresponding webhookAuthenticator object, and an error if there is any.
func (c *FakeWebhookAuthenticators) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.WebhookAuthenticator, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(webhookauthenticatorsResource, name), &v1alpha1.WebhookAuthenticator{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.WebhookAuthenticator), err
}

// List takes label and field selectors, and returns the list of WebhookAuthenticators that match those selectors.
func (c *FakeWebhookAuthenticators) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.WebhookAuthenticatorList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(webhookauthenticatorsResource, webhookauthenticatorsKind, opts), &v1alpha1.WebhookAuthenticatorList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.WebhookAuthenticatorList{ListMeta: obj.(*v1alpha1.WebhookAuthenticatorList).ListMeta}
	for _, item := range obj.(*v1alpha1.WebhookAuthenticatorList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested webhookAuthenticators.
func (c *FakeWebhookAuthenticators) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(webhookauthenticatorsResource, opts))
}

// Create takes the representation of a webhookAuthenticator and creates it.  Returns the server's representation of the webhookAuthenticator, and an error, if there is any.
func (c *FakeWebhookAuthenticators) Create(ctx context.Context, webhookAuthenticator *v1alpha1.WebhookAuthenticator, opts v1.CreateOptions) (result *v1alpha1.WebhookAuthenticator, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(webhookauthenticatorsResource, webhookAuthenticator), &v1alpha1.WebhookAuthenticator{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.WebhookAuthenticator), err
}

// Update takes the representation of a webhookAuthenticator and updates it. Returns the server's representation of the webhookAuthenticator, and an error, if there is any.
func (c *FakeWebhookAuthenticators) Update(ctx context.Context, webhookAuthenticator *v1alpha1.WebhookAuthenticator, opts v1.UpdateOptions) (result *v1alpha1.WebhookAuthenticator, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(webhookauthenticatorsResource, webhookAuthenticator), &v1alpha1.WebhookAuthenticator{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.WebhookAuthenticator), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeWebhookAuthenticators) UpdateStatus(ctx context.Context, webhookAuthenticator *v1alpha1.WebhookAuthenticator, opts v1.UpdateOptions) (*v1alpha1.WebhookAuthenticator, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(webhookauthenticatorsResource, "status", webhookAuthenticator), &v1alpha1.WebhookAuthenticator{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.WebhookAuthenticator), err
}

// Delete takes name of the webhookAuthenticator and deletes it. Returns an error if one occurs.
func (c *FakeWebhookAuthenticators) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(webhookauthenticatorsResource, name, opts), &v1alpha1.WebhookAuthenticator{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeWebhookAuthenticators) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(webhookauthenticatorsResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.WebhookAuthenticatorList{})
	return err
}

// Patch applies the patch and returns the patched webhookAuthenticator.
func (c *FakeWebhookAuthenticators) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.WebhookAuthenticator, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(webhookauthenticatorsResource, name, pt, data, subresources...), &v1alpha1.WebhookAuthenticator{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.WebhookAuthenticator), err
}
