// Copyright 2020-2025 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	clientset "go.pinniped.dev/generated/1.26/client/supervisor/clientset/versioned"
	clientsecretv1alpha1 "go.pinniped.dev/generated/1.26/client/supervisor/clientset/versioned/typed/clientsecret/v1alpha1"
	fakeclientsecretv1alpha1 "go.pinniped.dev/generated/1.26/client/supervisor/clientset/versioned/typed/clientsecret/v1alpha1/fake"
	configv1alpha1 "go.pinniped.dev/generated/1.26/client/supervisor/clientset/versioned/typed/config/v1alpha1"
	fakeconfigv1alpha1 "go.pinniped.dev/generated/1.26/client/supervisor/clientset/versioned/typed/config/v1alpha1/fake"
	idpv1alpha1 "go.pinniped.dev/generated/1.26/client/supervisor/clientset/versioned/typed/idp/v1alpha1"
	fakeidpv1alpha1 "go.pinniped.dev/generated/1.26/client/supervisor/clientset/versioned/typed/idp/v1alpha1/fake"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/discovery"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/testing"
)

// NewSimpleClientset returns a clientset that will respond with the provided objects.
// It's backed by a very simple object tracker that processes creates, updates and deletions as-is,
// without applying any validations and/or defaults. It shouldn't be considered a replacement
// for a real clientset and is mostly useful in simple unit tests.
func NewSimpleClientset(objects ...runtime.Object) *Clientset {
	o := testing.NewObjectTracker(scheme, codecs.UniversalDecoder())
	for _, obj := range objects {
		if err := o.Add(obj); err != nil {
			panic(err)
		}
	}

	cs := &Clientset{tracker: o}
	cs.discovery = &fakediscovery.FakeDiscovery{Fake: &cs.Fake}
	cs.AddReactor("*", "*", testing.ObjectReaction(o))
	cs.AddWatchReactor("*", func(action testing.Action) (handled bool, ret watch.Interface, err error) {
		gvr := action.GetResource()
		ns := action.GetNamespace()
		watch, err := o.Watch(gvr, ns)
		if err != nil {
			return false, nil, err
		}
		return true, watch, nil
	})

	return cs
}

// Clientset implements clientset.Interface. Meant to be embedded into a
// struct to get a default implementation. This makes faking out just the method
// you want to test easier.
type Clientset struct {
	testing.Fake
	discovery *fakediscovery.FakeDiscovery
	tracker   testing.ObjectTracker
}

func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	return c.discovery
}

func (c *Clientset) Tracker() testing.ObjectTracker {
	return c.tracker
}

var (
	_ clientset.Interface = &Clientset{}
	_ testing.FakeClient  = &Clientset{}
)

// ClientsecretV1alpha1 retrieves the ClientsecretV1alpha1Client
func (c *Clientset) ClientsecretV1alpha1() clientsecretv1alpha1.ClientsecretV1alpha1Interface {
	return &fakeclientsecretv1alpha1.FakeClientsecretV1alpha1{Fake: &c.Fake}
}

// ConfigV1alpha1 retrieves the ConfigV1alpha1Client
func (c *Clientset) ConfigV1alpha1() configv1alpha1.ConfigV1alpha1Interface {
	return &fakeconfigv1alpha1.FakeConfigV1alpha1{Fake: &c.Fake}
}

// IDPV1alpha1 retrieves the IDPV1alpha1Client
func (c *Clientset) IDPV1alpha1() idpv1alpha1.IDPV1alpha1Interface {
	return &fakeidpv1alpha1.FakeIDPV1alpha1{Fake: &c.Fake}
}
