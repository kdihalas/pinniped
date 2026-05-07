// Copyright 2024-2026 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package admissionpluginconfig

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/server/options"
	"k8s.io/client-go/discovery"
	kubefake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestValidateAdmissionPluginNames(t *testing.T) {
	tests := []struct {
		name        string
		pluginNames []string
		wantErr     string
	}{
		{
			name:        "empty",
			pluginNames: []string{},
		},
		{
			name:        "nil",
			pluginNames: nil,
		},
		{
			name:        "only invalid values",
			pluginNames: []string{"foo", "bar"},
			wantErr:     "admission plugin names not recognized: [foo bar] (each must be one of [NamespaceLifecycle MutatingAdmissionPolicy MutatingAdmissionWebhook ValidatingAdmissionPolicy ValidatingAdmissionWebhook])",
		},
		{
			name:        "duplicate invalid names are reported once per occurrence",
			pluginNames: []string{"foobar", "NamespaceLifecycle", "foobar"},
			wantErr:     "admission plugin names not recognized: [foobar foobar] (each must be one of [NamespaceLifecycle MutatingAdmissionPolicy MutatingAdmissionWebhook ValidatingAdmissionPolicy ValidatingAdmissionWebhook])",
		},
		{
			name:        "comparison is case-sensitive",
			pluginNames: []string{"namespacelifecycle"},
			wantErr:     "admission plugin names not recognized: [namespacelifecycle] (each must be one of [NamespaceLifecycle MutatingAdmissionPolicy MutatingAdmissionWebhook ValidatingAdmissionPolicy ValidatingAdmissionWebhook])",
		},
		{
			name:        "empty string entry is not recognized",
			pluginNames: []string{""},
			wantErr:     "admission plugin names not recognized: [] (each must be one of [NamespaceLifecycle MutatingAdmissionPolicy MutatingAdmissionWebhook ValidatingAdmissionPolicy ValidatingAdmissionWebhook])",
		},
		{
			name: "all current valid values (this list may change in future versions of Kubernetes packages)",
			pluginNames: []string{
				"NamespaceLifecycle",
				"MutatingAdmissionWebhook",
				"ValidatingAdmissionPolicy",
				"ValidatingAdmissionWebhook",
				"MutatingAdmissionPolicy",
			},
		},
		{
			name: "one invalid value",
			pluginNames: []string{
				"NamespaceLifecycle",
				"MutatingAdmissionWebhook",
				"ValidatingAdmissionPolicy",
				"foobar",
				"ValidatingAdmissionWebhook",
				"MutatingAdmissionPolicy",
			},
			wantErr: "admission plugin names not recognized: [foobar] (each must be one of [NamespaceLifecycle MutatingAdmissionPolicy MutatingAdmissionWebhook ValidatingAdmissionPolicy ValidatingAdmissionWebhook])",
		},
		{
			name: "multiple invalid values",
			pluginNames: []string{
				"NamespaceLifecycle",
				"MutatingAdmissionWebhook",
				"foobat",
				"ValidatingAdmissionPolicy",
				"foobar",
				"ValidatingAdmissionWebhook",
				"foobaz",
				"MutatingAdmissionPolicy",
			},
			wantErr: "admission plugin names not recognized: [foobat foobar foobaz] (each must be one of [NamespaceLifecycle MutatingAdmissionPolicy MutatingAdmissionWebhook ValidatingAdmissionPolicy ValidatingAdmissionWebhook])",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateAdmissionPluginNames(tt.pluginNames)
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestConfigureAdmissionPlugins(t *testing.T) {
	coreResources := &metav1.APIResourceList{
		GroupVersion: corev1.SchemeGroupVersion.String(),
		APIResources: []metav1.APIResource{
			{Name: "pods", Namespaced: true, Kind: "Pod"},
		},
	}

	appsResources := &metav1.APIResourceList{
		GroupVersion: appsv1.SchemeGroupVersion.String(),
		APIResources: []metav1.APIResource{
			{Name: "deployments", Namespaced: true, Kind: "Deployment"},
			{Name: "deployments/scale", Namespaced: true, Kind: "Scale", Group: "apps", Version: "v1"},
		},
	}

	newStyleAdmissionResources := &metav1.APIResourceList{
		GroupVersion: admissionregistrationv1.SchemeGroupVersion.String(),
		APIResources: []metav1.APIResource{
			{Name: "validatingwebhookconfigurations", Kind: "ValidatingWebhookConfiguration"},
			{Name: "validatingadmissionpolicies", Kind: "ValidatingAdmissionPolicy"},
			{Name: "mutatingadmissionpolicies", Kind: "MutatingAdmissionPolicy"},
		},
	}

	oldStyleAdmissionResources := &metav1.APIResourceList{
		GroupVersion: admissionregistrationv1.SchemeGroupVersion.String(),
		APIResources: []metav1.APIResource{
			{Name: "validatingwebhookconfigurations", Kind: "ValidatingWebhookConfiguration"},
		},
	}

	tests := []struct {
		name                  string
		disabledPlugins       []string
		availableAPIResources []*metav1.APIResourceList
		discoveryErr          error
		wantErr               string
		wantDisabledPlugins   []string
	}{
		{
			name: "when there are modern resource types and nil disabled list, then we do not change the plugin configuration",
			availableAPIResources: []*metav1.APIResourceList{
				coreResources,
				newStyleAdmissionResources,
				appsResources,
			},
			disabledPlugins:     nil,
			wantDisabledPlugins: nil,
		},
		{
			name: "when there are modern resource types resource and empty disabled list, then we do not change the plugin configuration",
			availableAPIResources: []*metav1.APIResourceList{
				coreResources,
				newStyleAdmissionResources,
				appsResources,
			},
			disabledPlugins:     []string{},
			wantDisabledPlugins: nil,
		},
		{
			name: "when are modern resource types are missing, as there would not be in an old Kubernetes cluster, then we disable those admission plugins",
			availableAPIResources: []*metav1.APIResourceList{
				coreResources,
				oldStyleAdmissionResources,
				appsResources,
			},
			disabledPlugins:     nil,
			wantDisabledPlugins: []string{"ValidatingAdmissionPolicy", "MutatingAdmissionPolicy"},
		},
		{
			name: "when only ValidatingAdmissionPolicy is missing, then we only automatically disable that admission plugin",
			availableAPIResources: []*metav1.APIResourceList{
				coreResources,
				{
					GroupVersion: admissionregistrationv1.SchemeGroupVersion.String(),
					APIResources: []metav1.APIResource{
						{Name: "validatingwebhookconfigurations", Kind: "ValidatingWebhookConfiguration"},
						{Name: "mutatingadmissionpolicies", Kind: "MutatingAdmissionPolicy"},
					},
				},
				appsResources,
			},
			disabledPlugins:     nil,
			wantDisabledPlugins: []string{"ValidatingAdmissionPolicy"},
		},
		{
			name: "when only MutatingAdmissionPolicy is missing, then we only automatically disable that admission plugin",
			availableAPIResources: []*metav1.APIResourceList{
				coreResources,
				{
					GroupVersion: admissionregistrationv1.SchemeGroupVersion.String(),
					APIResources: []metav1.APIResource{
						{Name: "validatingwebhookconfigurations", Kind: "ValidatingWebhookConfiguration"},
						{Name: "validatingadmissionpolicies", Kind: "ValidatingAdmissionPolicy"},
					},
				},
				appsResources,
			},
			disabledPlugins:     nil,
			wantDisabledPlugins: []string{"MutatingAdmissionPolicy"},
		},
		{
			name: "when only ValidatingAdmissionPolicy is explicitly disabled, then we may still need to automatically disable other admission plugins",
			availableAPIResources: []*metav1.APIResourceList{
				coreResources,
				{
					GroupVersion: admissionregistrationv1.SchemeGroupVersion.String(),
					APIResources: []metav1.APIResource{
						{Name: "validatingwebhookconfigurations", Kind: "ValidatingWebhookConfiguration"},
						{Name: "validatingadmissionpolicies", Kind: "ValidatingAdmissionPolicy"},
					},
				},
				appsResources,
			},
			disabledPlugins:     []string{"ValidatingAdmissionPolicy"},
			wantDisabledPlugins: []string{"ValidatingAdmissionPolicy", "MutatingAdmissionPolicy"},
		},
		{
			name: "when only MutatingAdmissionPolicy is explicitly disabled, then we may still need to automatically disable other admission plugins",
			availableAPIResources: []*metav1.APIResourceList{
				coreResources,
				{
					GroupVersion: admissionregistrationv1.SchemeGroupVersion.String(),
					APIResources: []metav1.APIResource{
						{Name: "validatingwebhookconfigurations", Kind: "ValidatingWebhookConfiguration"},
						{Name: "mutatingadmissionpolicies", Kind: "MutatingAdmissionPolicy"},
					},
				},
				appsResources,
			},
			disabledPlugins:     []string{"MutatingAdmissionPolicy"},
			wantDisabledPlugins: []string{"MutatingAdmissionPolicy", "ValidatingAdmissionPolicy"},
		},
		{
			name: "when there are only older version of are modern resource types, as there would be in an old Kubernetes cluster with the feature flag enabled, then we disable those plugins (because the admission code wants to watch v1)",
			availableAPIResources: []*metav1.APIResourceList{
				coreResources,
				{
					GroupVersion: admissionregistrationv1.SchemeGroupVersion.Group + "/v1beta1",
					APIResources: []metav1.APIResource{
						{Name: "validatingwebhookconfigurations", Kind: "ValidatingWebhookConfiguration"},
						{Name: "validatingadmissionpolicies", Kind: "ValidatingAdmissionPolicy"},
						{Name: "mutatingadmissionpolicies", Kind: "MutatingAdmissionPolicy"},
					},
				},
				appsResources,
			},
			disabledPlugins:     []string{},
			wantDisabledPlugins: []string{"ValidatingAdmissionPolicy", "MutatingAdmissionPolicy"},
		},
		{
			name:                  "when there are no are modern resource types, and all the modern resource type plugins were explicitly disabled, then do not perform discovery, and just disable them",
			availableAPIResources: []*metav1.APIResourceList{},
			discoveryErr:          errors.New("total error from API discovery client"), // shouldn't matter because discovery should have been skipped
			disabledPlugins:       []string{"MutatingAdmissionWebhook", "ValidatingAdmissionPolicy", "MutatingAdmissionPolicy"},
			wantDisabledPlugins:   []string{"MutatingAdmissionWebhook", "ValidatingAdmissionPolicy", "MutatingAdmissionPolicy"},
		},
		{
			name: "when there are no modern resource types, and the modern resource type plugins were not explicitly disabled, still disable them",
			availableAPIResources: []*metav1.APIResourceList{
				coreResources,
				oldStyleAdmissionResources,
				appsResources,
			},
			disabledPlugins:     []string{"MutatingAdmissionWebhook", "NamespaceLifecycle"},
			wantDisabledPlugins: []string{"MutatingAdmissionWebhook", "NamespaceLifecycle", "ValidatingAdmissionPolicy", "MutatingAdmissionPolicy"},
		},
		{
			name:                "when there is a total error returned by discovery",
			discoveryErr:        errors.New("total error from API discovery client"),
			wantErr:             "failed to perform k8s API discovery for purpose of checking availability of admissionregistration.k8s.io resource types: total error from API discovery client",
			wantDisabledPlugins: nil,
		},
		{
			name: "when there is a partial error returned by discovery which does include the group of interest, then we cannot ignore the error, because we could not discover anything about that group",
			availableAPIResources: []*metav1.APIResourceList{
				coreResources,
				oldStyleAdmissionResources,
				appsResources,
			},
			discoveryErr: &discovery.ErrGroupDiscoveryFailed{Groups: map[schema.GroupVersion]error{
				schema.GroupVersion{Group: "someGroup", Version: "v1"}:                    errors.New("fake error for someGroup"),
				schema.GroupVersion{Group: "admissionregistration.k8s.io", Version: "v1"}: errors.New("fake error for admissionregistration"),
			}},
			wantErr:             "failed to perform k8s API discovery for purpose of checking availability of admissionregistration.k8s.io resource types: unable to retrieve the complete list of server APIs: admissionregistration.k8s.io/v1: fake error for admissionregistration, someGroup/v1: fake error for someGroup",
			wantDisabledPlugins: nil,
		},
		{
			name: "when there is a partial error returned by discovery on an new-style cluster which does not include the group of interest, then we can ignore the error and use the default plugins",
			availableAPIResources: []*metav1.APIResourceList{
				coreResources,
				newStyleAdmissionResources,
				appsResources,
			},
			discoveryErr: &discovery.ErrGroupDiscoveryFailed{Groups: map[schema.GroupVersion]error{
				schema.GroupVersion{Group: "someGroup", Version: "v1"}:      errors.New("fake error for someGroup"),
				schema.GroupVersion{Group: "someOtherGroup", Version: "v1"}: errors.New("fake error for someOtherGroup"),
			}},
			wantDisabledPlugins: nil,
		},
		{
			name: "when there is a partial error returned by discovery on an old-style cluster which does not include the group of interest, then we can ignore the error and customize the plugins",
			availableAPIResources: []*metav1.APIResourceList{
				coreResources,
				oldStyleAdmissionResources,
				appsResources,
			},
			discoveryErr: &discovery.ErrGroupDiscoveryFailed{Groups: map[schema.GroupVersion]error{
				schema.GroupVersion{Group: "someGroup", Version: "v1"}:      errors.New("fake error for someGroup"),
				schema.GroupVersion{Group: "someOtherGroup", Version: "v1"}: errors.New("fake error for someOtherGroup"),
			}},
			wantDisabledPlugins: []string{"ValidatingAdmissionPolicy", "MutatingAdmissionPolicy"},
		},
		{
			name: "when both modern resource type plugins are explicitly disabled on a modern cluster, discovery is skipped",
			availableAPIResources: []*metav1.APIResourceList{
				coreResources,
				newStyleAdmissionResources,
				appsResources,
			},
			discoveryErr:        errors.New("discovery should not have been called"), // would surface if the short-circuit were removed
			disabledPlugins:     []string{"ValidatingAdmissionPolicy", "MutatingAdmissionPolicy"},
			wantDisabledPlugins: []string{"ValidatingAdmissionPolicy", "MutatingAdmissionPolicy"},
		},
		{
			name: "when only ValidatingAdmissionPolicy is explicitly disabled on a fully modern cluster, then MutatingAdmissionPolicy is not auto-disabled",
			availableAPIResources: []*metav1.APIResourceList{
				coreResources,
				newStyleAdmissionResources,
				appsResources,
			},
			disabledPlugins:     []string{"ValidatingAdmissionPolicy"},
			wantDisabledPlugins: []string{"ValidatingAdmissionPolicy"},
		},
		{
			name: "when only MutatingAdmissionPolicy is explicitly disabled on a fully modern cluster, then ValidatingAdmissionPolicy is not auto-disabled",
			availableAPIResources: []*metav1.APIResourceList{
				coreResources,
				newStyleAdmissionResources,
				appsResources,
			},
			disabledPlugins:     []string{"MutatingAdmissionPolicy"},
			wantDisabledPlugins: []string{"MutatingAdmissionPolicy"},
		},
		{
			name: "when a non-policy plugin is explicitly disabled on a modern cluster, only that plugin is disabled",
			availableAPIResources: []*metav1.APIResourceList{
				coreResources,
				newStyleAdmissionResources,
				appsResources,
			},
			disabledPlugins:     []string{"NamespaceLifecycle"},
			wantDisabledPlugins: []string{"NamespaceLifecycle"},
		},
		{
			name:                  "when discovery returns nil resources along with a partial error, then the error is returned",
			availableAPIResources: nil,
			discoveryErr: &discovery.ErrGroupDiscoveryFailed{Groups: map[schema.GroupVersion]error{
				schema.GroupVersion{Group: "someGroup", Version: "v1"}: errors.New("fake error for someGroup"),
			}},
			wantErr:             "failed to perform k8s API discovery for purpose of checking availability of admissionregistration.k8s.io resource types: unable to retrieve the complete list of server APIs: someGroup/v1: fake error for someGroup",
			wantDisabledPlugins: nil,
		},
		{
			name: "when there is a partial error returned by discovery which includes the admissionregistration group at a non-v1 version, then we cannot ignore the error",
			availableAPIResources: []*metav1.APIResourceList{
				coreResources,
				oldStyleAdmissionResources,
				appsResources,
			},
			discoveryErr: &discovery.ErrGroupDiscoveryFailed{Groups: map[schema.GroupVersion]error{
				schema.GroupVersion{Group: "admissionregistration.k8s.io", Version: "v1beta1"}: errors.New("fake error for admissionregistration v1beta1"),
			}},
			wantErr:             "failed to perform k8s API discovery for purpose of checking availability of admissionregistration.k8s.io resource types: unable to retrieve the complete list of server APIs: admissionregistration.k8s.io/v1beta1: fake error for admissionregistration v1beta1",
			wantDisabledPlugins: nil,
		},
		{
			name: "when the admissionregistration group is entirely absent from discovery, then we disable both modern admission plugins",
			availableAPIResources: []*metav1.APIResourceList{
				coreResources,
				appsResources,
			},
			disabledPlugins:     nil,
			wantDisabledPlugins: []string{"ValidatingAdmissionPolicy", "MutatingAdmissionPolicy"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			kubeClient := kubefake.NewClientset()
			kubeClient.Resources = tt.availableAPIResources

			// Unfortunately, NewClientset() does not support using reactors to
			// cause discovery to return errors. Instead, we will make our own fake implementation of the
			// discovery client's interface and only mock the parts that we need for this test.
			discoveryClient := newFakeDiscoveryClient(kubeClient)

			if tt.discoveryErr != nil {
				kubeClient.PrependReactor(
					"get",
					"resource",
					func(a k8stesting.Action) (bool, runtime.Object, error) {
						return true, nil, tt.discoveryErr
					},
				)
			}

			opts := &options.RecommendedOptions{
				Admission: options.NewAdmissionOptions(),
			}
			// Sanity checks on opts before we use it.
			require.Empty(t, opts.Admission.DisablePlugins)

			// Call the function under test.
			err := configureAdmissionPlugins(discoveryClient, opts, tt.disabledPlugins)

			if tt.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.wantErr)
			}

			// Check the expected side effects of the function under test, if any.
			require.Equal(t, tt.wantDisabledPlugins, opts.Admission.DisablePlugins)
		})
	}
}

type fakeDiscoveryClient struct {
	fakeClientSet *kubefake.Clientset
}

var _ discovery.ServerResourcesInterface = &fakeDiscoveryClient{}

func newFakeDiscoveryClient(fakeClientSet *kubefake.Clientset) *fakeDiscoveryClient {
	return &fakeDiscoveryClient{
		fakeClientSet: fakeClientSet,
	}
}

// This is the only function from the discovery.DiscoveryInterface that we care to fake for this test.
// The rest of the functions are here only to satisfy the interface.
func (f *fakeDiscoveryClient) ServerPreferredResources() ([]*metav1.APIResourceList, error) {
	action := k8stesting.ActionImpl{
		Verb:     "get",
		Resource: schema.GroupVersionResource{Resource: "resource"},
	}
	// Wire in actions just enough that we can cause errors for the test when we want them.
	// Ignoring the first return value because we don't need it for this test.
	_, err := f.fakeClientSet.Invokes(action, nil)
	// Still return the "partial" results even where there was an error, similar enough to how the real API works.
	return f.fakeClientSet.Resources, err
}

func (f *fakeDiscoveryClient) ServerResourcesForGroupVersion(_ string) (*metav1.APIResourceList, error) {
	return nil, nil
}

func (f *fakeDiscoveryClient) ServerGroupsAndResources() ([]*metav1.APIGroup, []*metav1.APIResourceList, error) {
	return nil, nil, nil
}

func (f *fakeDiscoveryClient) ServerPreferredNamespacedResources() ([]*metav1.APIResourceList, error) {
	return nil, nil
}
