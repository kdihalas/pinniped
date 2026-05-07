// Copyright 2021-2026 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package kubecertagent

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/client-go/rest"

	"go.pinniped.dev/internal/crypto/ptls"
	"go.pinniped.dev/internal/kubeclient"
	"go.pinniped.dev/internal/testutil/tlsserver"
)

func TestSecureTLS(t *testing.T) {
	var sawRequest bool
	server, serverCA := tlsserver.TestServerIPv4(t, http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		tlsserver.AssertTLS(t, r, ptls.Secure)
		sawRequest = true
	}), tlsserver.RecordTLSHello)

	config := &rest.Config{
		Host: server.URL,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: serverCA,
		},
	}

	client, err := kubeclient.New(kubeclient.WithConfig(config))
	require.NoError(t, err)

	// Build this exactly like our production code does.
	podCommandExecutor := NewPodCommandExecutor(client.JSONConfig, client.Kubernetes)

	got, err := podCommandExecutor.Exec(context.Background(), "podNamespace", "podName", "containerName", "command", "arg1", "arg2")
	// Expect to get an error because the fake server above does not allow upgrade to spdy.
	// This doesn't matter because all we really care about in this test is the results of AssertTLS.
	require.EqualError(t, err, "unable to upgrade connection: empty server response")
	require.Empty(t, got)

	require.True(t, sawRequest)
}
