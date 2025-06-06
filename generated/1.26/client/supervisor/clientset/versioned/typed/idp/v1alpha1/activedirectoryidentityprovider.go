// Copyright 2020-2025 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "go.pinniped.dev/generated/1.26/apis/supervisor/idp/v1alpha1"
	scheme "go.pinniped.dev/generated/1.26/client/supervisor/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ActiveDirectoryIdentityProvidersGetter has a method to return a ActiveDirectoryIdentityProviderInterface.
// A group's client should implement this interface.
type ActiveDirectoryIdentityProvidersGetter interface {
	ActiveDirectoryIdentityProviders(namespace string) ActiveDirectoryIdentityProviderInterface
}

// ActiveDirectoryIdentityProviderInterface has methods to work with ActiveDirectoryIdentityProvider resources.
type ActiveDirectoryIdentityProviderInterface interface {
	Create(ctx context.Context, activeDirectoryIdentityProvider *v1alpha1.ActiveDirectoryIdentityProvider, opts v1.CreateOptions) (*v1alpha1.ActiveDirectoryIdentityProvider, error)
	Update(ctx context.Context, activeDirectoryIdentityProvider *v1alpha1.ActiveDirectoryIdentityProvider, opts v1.UpdateOptions) (*v1alpha1.ActiveDirectoryIdentityProvider, error)
	UpdateStatus(ctx context.Context, activeDirectoryIdentityProvider *v1alpha1.ActiveDirectoryIdentityProvider, opts v1.UpdateOptions) (*v1alpha1.ActiveDirectoryIdentityProvider, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.ActiveDirectoryIdentityProvider, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.ActiveDirectoryIdentityProviderList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ActiveDirectoryIdentityProvider, err error)
	ActiveDirectoryIdentityProviderExpansion
}

// activeDirectoryIdentityProviders implements ActiveDirectoryIdentityProviderInterface
type activeDirectoryIdentityProviders struct {
	client rest.Interface
	ns     string
}

// newActiveDirectoryIdentityProviders returns a ActiveDirectoryIdentityProviders
func newActiveDirectoryIdentityProviders(c *IDPV1alpha1Client, namespace string) *activeDirectoryIdentityProviders {
	return &activeDirectoryIdentityProviders{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the activeDirectoryIdentityProvider, and returns the corresponding activeDirectoryIdentityProvider object, and an error if there is any.
func (c *activeDirectoryIdentityProviders) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ActiveDirectoryIdentityProvider, err error) {
	result = &v1alpha1.ActiveDirectoryIdentityProvider{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("activedirectoryidentityproviders").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ActiveDirectoryIdentityProviders that match those selectors.
func (c *activeDirectoryIdentityProviders) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ActiveDirectoryIdentityProviderList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.ActiveDirectoryIdentityProviderList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("activedirectoryidentityproviders").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested activeDirectoryIdentityProviders.
func (c *activeDirectoryIdentityProviders) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("activedirectoryidentityproviders").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a activeDirectoryIdentityProvider and creates it.  Returns the server's representation of the activeDirectoryIdentityProvider, and an error, if there is any.
func (c *activeDirectoryIdentityProviders) Create(ctx context.Context, activeDirectoryIdentityProvider *v1alpha1.ActiveDirectoryIdentityProvider, opts v1.CreateOptions) (result *v1alpha1.ActiveDirectoryIdentityProvider, err error) {
	result = &v1alpha1.ActiveDirectoryIdentityProvider{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("activedirectoryidentityproviders").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(activeDirectoryIdentityProvider).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a activeDirectoryIdentityProvider and updates it. Returns the server's representation of the activeDirectoryIdentityProvider, and an error, if there is any.
func (c *activeDirectoryIdentityProviders) Update(ctx context.Context, activeDirectoryIdentityProvider *v1alpha1.ActiveDirectoryIdentityProvider, opts v1.UpdateOptions) (result *v1alpha1.ActiveDirectoryIdentityProvider, err error) {
	result = &v1alpha1.ActiveDirectoryIdentityProvider{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("activedirectoryidentityproviders").
		Name(activeDirectoryIdentityProvider.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(activeDirectoryIdentityProvider).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *activeDirectoryIdentityProviders) UpdateStatus(ctx context.Context, activeDirectoryIdentityProvider *v1alpha1.ActiveDirectoryIdentityProvider, opts v1.UpdateOptions) (result *v1alpha1.ActiveDirectoryIdentityProvider, err error) {
	result = &v1alpha1.ActiveDirectoryIdentityProvider{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("activedirectoryidentityproviders").
		Name(activeDirectoryIdentityProvider.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(activeDirectoryIdentityProvider).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the activeDirectoryIdentityProvider and deletes it. Returns an error if one occurs.
func (c *activeDirectoryIdentityProviders) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("activedirectoryidentityproviders").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *activeDirectoryIdentityProviders) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("activedirectoryidentityproviders").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched activeDirectoryIdentityProvider.
func (c *activeDirectoryIdentityProviders) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ActiveDirectoryIdentityProvider, err error) {
	result = &v1alpha1.ActiveDirectoryIdentityProvider{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("activedirectoryidentityproviders").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
