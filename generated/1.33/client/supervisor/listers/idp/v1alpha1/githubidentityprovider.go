// Copyright 2020-2025 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	idpv1alpha1 "go.pinniped.dev/generated/1.33/apis/supervisor/idp/v1alpha1"
	labels "k8s.io/apimachinery/pkg/labels"
	listers "k8s.io/client-go/listers"
	cache "k8s.io/client-go/tools/cache"
)

// GitHubIdentityProviderLister helps list GitHubIdentityProviders.
// All objects returned here must be treated as read-only.
type GitHubIdentityProviderLister interface {
	// List lists all GitHubIdentityProviders in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*idpv1alpha1.GitHubIdentityProvider, err error)
	// GitHubIdentityProviders returns an object that can list and get GitHubIdentityProviders.
	GitHubIdentityProviders(namespace string) GitHubIdentityProviderNamespaceLister
	GitHubIdentityProviderListerExpansion
}

// gitHubIdentityProviderLister implements the GitHubIdentityProviderLister interface.
type gitHubIdentityProviderLister struct {
	listers.ResourceIndexer[*idpv1alpha1.GitHubIdentityProvider]
}

// NewGitHubIdentityProviderLister returns a new GitHubIdentityProviderLister.
func NewGitHubIdentityProviderLister(indexer cache.Indexer) GitHubIdentityProviderLister {
	return &gitHubIdentityProviderLister{listers.New[*idpv1alpha1.GitHubIdentityProvider](indexer, idpv1alpha1.Resource("githubidentityprovider"))}
}

// GitHubIdentityProviders returns an object that can list and get GitHubIdentityProviders.
func (s *gitHubIdentityProviderLister) GitHubIdentityProviders(namespace string) GitHubIdentityProviderNamespaceLister {
	return gitHubIdentityProviderNamespaceLister{listers.NewNamespaced[*idpv1alpha1.GitHubIdentityProvider](s.ResourceIndexer, namespace)}
}

// GitHubIdentityProviderNamespaceLister helps list and get GitHubIdentityProviders.
// All objects returned here must be treated as read-only.
type GitHubIdentityProviderNamespaceLister interface {
	// List lists all GitHubIdentityProviders in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*idpv1alpha1.GitHubIdentityProvider, err error)
	// Get retrieves the GitHubIdentityProvider from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*idpv1alpha1.GitHubIdentityProvider, error)
	GitHubIdentityProviderNamespaceListerExpansion
}

// gitHubIdentityProviderNamespaceLister implements the GitHubIdentityProviderNamespaceLister
// interface.
type gitHubIdentityProviderNamespaceLister struct {
	listers.ResourceIndexer[*idpv1alpha1.GitHubIdentityProvider]
}
