// Copyright 2020-2021 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package concierge contains functionality to load/store Config's from/to
// some source.
package concierge

import (
	"fmt"
	"io/ioutil"
	"strings"

	"k8s.io/utils/pointer"
	"sigs.k8s.io/yaml"

	"go.pinniped.dev/internal/constable"
	"go.pinniped.dev/internal/groupsuffix"
	"go.pinniped.dev/internal/plog"
)

const (
	aboutAYear   = 60 * 60 * 24 * 365
	about9Months = 60 * 60 * 24 * 30 * 9

	// Use 10250 because it happens to be the same port on which the Kubelet listens, so some cluster types
	// are more permissive with servers that run on this port. For example, GKE private clusters do not
	// allow traffic from the control plane to most ports, but do allow traffic to port 10250. This allows
	// the Concierge to work without additional configuration on these types of clusters.
	aggregatedAPIServerPortDefault = 10250

	// Use port 8444 because that is the port that was selected for the first released version of the
	// impersonation proxy, and has been the value since. It was originally selected because the
	// aggregated API server used to run on 8443 (has since changed), so 8444 was the next available port.
	impersonationProxyPortDefault = 8444
)

// FromPath loads an Config from a provided local file path, inserts any
// defaults (from the Config documentation), and verifies that the config is
// valid (per the Config documentation).
//
// Note! The Config file should contain base64-encoded WebhookCABundle data.
// This function will decode that base64-encoded data to PEM bytes to be stored
// in the Config.
func FromPath(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("decode yaml: %w", err)
	}

	maybeSetAPIDefaults(&config.APIConfig)
	maybeSetAggregatedAPIServerPortDefaults(&config.AggregatedAPIServerPort)
	maybeSetImpersonationProxyServerPortDefaults(&config.ImpersonationProxyServerPort)
	maybeSetAPIGroupSuffixDefault(&config.APIGroupSuffix)
	maybeSetKubeCertAgentDefaults(&config.KubeCertAgentConfig)

	if err := validateAPI(&config.APIConfig); err != nil {
		return nil, fmt.Errorf("validate api: %w", err)
	}

	if err := validateAPIGroupSuffix(*config.APIGroupSuffix); err != nil {
		return nil, fmt.Errorf("validate apiGroupSuffix: %w", err)
	}

	if err := validateServerPort(config.AggregatedAPIServerPort); err != nil {
		return nil, fmt.Errorf("validate aggregatedAPIServerPort: %w", err)
	}

	if err := validateServerPort(config.ImpersonationProxyServerPort); err != nil {
		return nil, fmt.Errorf("validate impersonationProxyServerPort: %w", err)
	}

	if err := validateNames(&config.NamesConfig); err != nil {
		return nil, fmt.Errorf("validate names: %w", err)
	}

	if err := plog.ValidateAndSetLogLevelGlobally(config.LogLevel); err != nil {
		return nil, fmt.Errorf("validate log level: %w", err)
	}

	if config.Labels == nil {
		config.Labels = make(map[string]string)
	}

	return &config, nil
}

func maybeSetAPIDefaults(apiConfig *APIConfigSpec) {
	if apiConfig.ServingCertificateConfig.DurationSeconds == nil {
		apiConfig.ServingCertificateConfig.DurationSeconds = pointer.Int64Ptr(aboutAYear)
	}

	if apiConfig.ServingCertificateConfig.RenewBeforeSeconds == nil {
		apiConfig.ServingCertificateConfig.RenewBeforeSeconds = pointer.Int64Ptr(about9Months)
	}
}

func maybeSetAPIGroupSuffixDefault(apiGroupSuffix **string) {
	if *apiGroupSuffix == nil {
		*apiGroupSuffix = pointer.StringPtr(groupsuffix.PinnipedDefaultSuffix)
	}
}

func maybeSetAggregatedAPIServerPortDefaults(port **int64) {
	if *port == nil {
		*port = pointer.Int64Ptr(aggregatedAPIServerPortDefault)
	}
}

func maybeSetImpersonationProxyServerPortDefaults(port **int64) {
	if *port == nil {
		*port = pointer.Int64Ptr(impersonationProxyPortDefault)
	}
}

func maybeSetKubeCertAgentDefaults(cfg *KubeCertAgentSpec) {
	if cfg.NamePrefix == nil {
		cfg.NamePrefix = pointer.StringPtr("pinniped-kube-cert-agent-")
	}

	if cfg.Image == nil {
		cfg.Image = pointer.StringPtr("debian:latest")
	}
}

func validateNames(names *NamesConfigSpec) error {
	missingNames := []string{}
	if names == nil {
		names = &NamesConfigSpec{}
	}
	if names.ServingCertificateSecret == "" {
		missingNames = append(missingNames, "servingCertificateSecret")
	}
	if names.CredentialIssuer == "" {
		missingNames = append(missingNames, "credentialIssuer")
	}
	if names.APIService == "" {
		missingNames = append(missingNames, "apiService")
	}
	if names.ImpersonationLoadBalancerService == "" {
		missingNames = append(missingNames, "impersonationLoadBalancerService")
	}
	if names.ImpersonationClusterIPService == "" {
		missingNames = append(missingNames, "impersonationClusterIPService")
	}
	if names.ImpersonationTLSCertificateSecret == "" {
		missingNames = append(missingNames, "impersonationTLSCertificateSecret")
	}
	if names.ImpersonationCACertificateSecret == "" {
		missingNames = append(missingNames, "impersonationCACertificateSecret")
	}
	if names.ImpersonationSignerSecret == "" {
		missingNames = append(missingNames, "impersonationSignerSecret")
	}
	if names.AgentServiceAccount == "" {
		missingNames = append(missingNames, "agentServiceAccount")
	}
	if len(missingNames) > 0 {
		return constable.Error("missing required names: " + strings.Join(missingNames, ", "))
	}
	return nil
}

func validateAPI(apiConfig *APIConfigSpec) error {
	if *apiConfig.ServingCertificateConfig.DurationSeconds < *apiConfig.ServingCertificateConfig.RenewBeforeSeconds {
		return constable.Error("durationSeconds cannot be smaller than renewBeforeSeconds")
	}

	if *apiConfig.ServingCertificateConfig.RenewBeforeSeconds <= 0 {
		return constable.Error("renewBefore must be positive")
	}

	return nil
}

func validateAPIGroupSuffix(apiGroupSuffix string) error {
	return groupsuffix.Validate(apiGroupSuffix)
}

func validateServerPort(port *int64) error {
	// It cannot be below 1024 because the container is not running as root.
	if *port < 1024 || *port > 65535 {
		return constable.Error("must be within range 1024 to 65535")
	}
	return nil
}
