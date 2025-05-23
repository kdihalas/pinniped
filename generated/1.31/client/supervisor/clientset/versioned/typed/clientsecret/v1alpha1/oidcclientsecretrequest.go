// Copyright 2020-2025 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"

	v1alpha1 "go.pinniped.dev/generated/1.31/apis/supervisor/clientsecret/v1alpha1"
	scheme "go.pinniped.dev/generated/1.31/client/supervisor/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gentype "k8s.io/client-go/gentype"
)

// OIDCClientSecretRequestsGetter has a method to return a OIDCClientSecretRequestInterface.
// A group's client should implement this interface.
type OIDCClientSecretRequestsGetter interface {
	OIDCClientSecretRequests(namespace string) OIDCClientSecretRequestInterface
}

// OIDCClientSecretRequestInterface has methods to work with OIDCClientSecretRequest resources.
type OIDCClientSecretRequestInterface interface {
	Create(ctx context.Context, oIDCClientSecretRequest *v1alpha1.OIDCClientSecretRequest, opts v1.CreateOptions) (*v1alpha1.OIDCClientSecretRequest, error)
	OIDCClientSecretRequestExpansion
}

// oIDCClientSecretRequests implements OIDCClientSecretRequestInterface
type oIDCClientSecretRequests struct {
	*gentype.Client[*v1alpha1.OIDCClientSecretRequest]
}

// newOIDCClientSecretRequests returns a OIDCClientSecretRequests
func newOIDCClientSecretRequests(c *ClientsecretV1alpha1Client, namespace string) *oIDCClientSecretRequests {
	return &oIDCClientSecretRequests{
		gentype.NewClient[*v1alpha1.OIDCClientSecretRequest](
			"oidcclientsecretrequests",
			c.RESTClient(),
			scheme.ParameterCodec,
			namespace,
			func() *v1alpha1.OIDCClientSecretRequest { return &v1alpha1.OIDCClientSecretRequest{} }),
	}
}
