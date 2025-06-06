// Copyright 2020-2025 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by informer-gen. DO NOT EDIT.

package externalversions

import (
	fmt "fmt"

	v1alpha1 "go.pinniped.dev/generated/1.32/apis/supervisor/config/v1alpha1"
	idpv1alpha1 "go.pinniped.dev/generated/1.32/apis/supervisor/idp/v1alpha1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	cache "k8s.io/client-go/tools/cache"
)

// GenericInformer is type of SharedIndexInformer which will locate and delegate to other
// sharedInformers based on type
type GenericInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() cache.GenericLister
}

type genericInformer struct {
	informer cache.SharedIndexInformer
	resource schema.GroupResource
}

// Informer returns the SharedIndexInformer.
func (f *genericInformer) Informer() cache.SharedIndexInformer {
	return f.informer
}

// Lister returns the GenericLister.
func (f *genericInformer) Lister() cache.GenericLister {
	return cache.NewGenericLister(f.Informer().GetIndexer(), f.resource)
}

// ForResource gives generic access to a shared informer of the matching type
// TODO extend this to unknown resources with a client pool
func (f *sharedInformerFactory) ForResource(resource schema.GroupVersionResource) (GenericInformer, error) {
	switch resource {
	// Group=config.supervisor.pinniped.dev, Version=v1alpha1
	case v1alpha1.SchemeGroupVersion.WithResource("federationdomains"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Config().V1alpha1().FederationDomains().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("oidcclients"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Config().V1alpha1().OIDCClients().Informer()}, nil

		// Group=idp.supervisor.pinniped.dev, Version=v1alpha1
	case idpv1alpha1.SchemeGroupVersion.WithResource("activedirectoryidentityproviders"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.IDP().V1alpha1().ActiveDirectoryIdentityProviders().Informer()}, nil
	case idpv1alpha1.SchemeGroupVersion.WithResource("githubidentityproviders"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.IDP().V1alpha1().GitHubIdentityProviders().Informer()}, nil
	case idpv1alpha1.SchemeGroupVersion.WithResource("ldapidentityproviders"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.IDP().V1alpha1().LDAPIdentityProviders().Informer()}, nil
	case idpv1alpha1.SchemeGroupVersion.WithResource("oidcidentityproviders"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.IDP().V1alpha1().OIDCIdentityProviders().Informer()}, nil

	}

	return nil, fmt.Errorf("no informer found for %v", resource)
}
