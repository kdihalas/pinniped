// Copyright 2022-2025 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// This file overrides profiles.go when Pinniped is built in FIPS-only mode.
//go:build fips_strict

package ptls

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"path/filepath"
	"runtime"

	"k8s.io/apiserver/pkg/server/options"

	// Cause fipsonly tls mode with this side effect import.
	_ "go.pinniped.dev/internal/crypto/fips"
	"go.pinniped.dev/internal/plog"
)

// The union of these three variables is all the FIPS-approved TLS 1.2 ciphers.
// If this list does not match the boring crypto compiler's list then the TestFIPSCipherSuites integration
// test should fail, which indicates that this list needs to be updated.
var (
	// secureCipherSuiteIDs is the list of TLS ciphers to use for both clients and servers when using TLS 1.2.
	//
	// FIPS allows the use of these ciphers which golang considers secure.
	secureCipherSuiteIDs = []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	}

	// insecureCipherSuiteIDs is a list of additional ciphers that should be allowed for both clients
	// and servers when using TLS 1.2.
	//
	// Previous versions of FIPS allowed the use of some specific ciphers that golang considers insecure.
	// Go 1.24 does not anymore, so now this list is empty.
	insecureCipherSuiteIDs []uint16

	// additionalSecureCipherSuiteIDsOnlyForLDAPClients are additional ciphers to use only for LDAP clients
	// when using TLS 1.2. These can be used when the Pinniped Supervisor is making calls to an LDAP server
	// configured by an LDAPIdentityProvider or ActiveDirectoryIdentityProvider.
	//
	// When compiled in FIPS mode, there are no extras for LDAP clients.
	additionalSecureCipherSuiteIDsOnlyForLDAPClients []uint16
)

// init: see comment in profiles.go.
func init() {
	switch filepath.Base(os.Args[0]) {
	case "pinniped-server", "pinniped-supervisor", "pinniped-concierge", "pinniped-concierge-kube-cert-agent":
	default:
		return // do not print FIPS logs if we cannot confirm that we are running a server binary
	}

	// this init runs before we have parsed our config to determine our log level
	// thus we must use a log statement that will always print instead of conditionally print
	plog.Always("this server was compiled to use boring crypto in FIPS-only mode",
		"go version", runtime.Version())
}

// Default: see comment in profiles.go.
// This chooses different cipher suites and/or TLS versions compared to non-FIPS mode.
// In FIPS mode, this will use the union of the secureCipherSuiteIDs, additionalSecureCipherSuiteIDsOnlyForLDAPClients,
// and insecureCipherSuiteIDs values defined above.
func Default(rootCAs *x509.CertPool) *tls.Config {
	config := buildTLSConfig(rootCAs, allHardcodedAllowedCipherSuites(), getUserConfiguredAllowedCipherSuitesForTLSOneDotTwo())
	// Note: starting in Go 1.24, boringcrypto supports TLS 1.3, so we allow it here.
	config.MaxVersion = tls.VersionTLS13
	return config
}

// DefaultLDAP: see comment in profiles.go.
// This chooses different cipher suites and/or TLS versions compared to non-FIPS mode.
// In FIPS mode, this is not any different from the Default profile.
func DefaultLDAP(rootCAs *x509.CertPool) *tls.Config {
	return Default(rootCAs)
}

// Secure: see comment in profiles.go.
// This chooses different cipher suites and/or TLS versions compared to non-FIPS mode.
// Note: starting in Go 1.24, boringcrypto supports TLS 1.3, so we allow it here.
// However, until it is safe to assume that a FIPS-compiled k8s server supports TLS 1.3, continue to
// make the Secure profile the same as the Default profile in FIPS mode, to allow both TLS 1.2 and 1.3.
func Secure(rootCAs *x509.CertPool) *tls.Config {
	return Default(rootCAs)
}

// SecureServing: see comment in profiles.go.
// This chooses different cipher suites and/or TLS versions compared to non-FIPS mode.
// Note: starting in Go 1.24, boringcrypto supports TLS 1.3, so we allow it here.
// However, until it is safe to assume that a FIPS-compiled k8s server supports TLS 1.3, continue to
// make SecureServing use the same as the defaultServing profile in FIPS mode, to allow both TLS 1.2 and 1.3.
func SecureServing(opts *options.SecureServingOptionsWithLoopback) {
	defaultServing(opts)
}
