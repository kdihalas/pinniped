//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Copyright 2020-2025 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by conversion-gen. DO NOT EDIT.

package v1alpha1

import (
	unsafe "unsafe"

	clientsecret "go.pinniped.dev/generated/1.31/apis/supervisor/clientsecret"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*OIDCClientSecretRequest)(nil), (*clientsecret.OIDCClientSecretRequest)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_OIDCClientSecretRequest_To_clientsecret_OIDCClientSecretRequest(a.(*OIDCClientSecretRequest), b.(*clientsecret.OIDCClientSecretRequest), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*clientsecret.OIDCClientSecretRequest)(nil), (*OIDCClientSecretRequest)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_clientsecret_OIDCClientSecretRequest_To_v1alpha1_OIDCClientSecretRequest(a.(*clientsecret.OIDCClientSecretRequest), b.(*OIDCClientSecretRequest), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*OIDCClientSecretRequestList)(nil), (*clientsecret.OIDCClientSecretRequestList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_OIDCClientSecretRequestList_To_clientsecret_OIDCClientSecretRequestList(a.(*OIDCClientSecretRequestList), b.(*clientsecret.OIDCClientSecretRequestList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*clientsecret.OIDCClientSecretRequestList)(nil), (*OIDCClientSecretRequestList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_clientsecret_OIDCClientSecretRequestList_To_v1alpha1_OIDCClientSecretRequestList(a.(*clientsecret.OIDCClientSecretRequestList), b.(*OIDCClientSecretRequestList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*OIDCClientSecretRequestSpec)(nil), (*clientsecret.OIDCClientSecretRequestSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_OIDCClientSecretRequestSpec_To_clientsecret_OIDCClientSecretRequestSpec(a.(*OIDCClientSecretRequestSpec), b.(*clientsecret.OIDCClientSecretRequestSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*clientsecret.OIDCClientSecretRequestSpec)(nil), (*OIDCClientSecretRequestSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_clientsecret_OIDCClientSecretRequestSpec_To_v1alpha1_OIDCClientSecretRequestSpec(a.(*clientsecret.OIDCClientSecretRequestSpec), b.(*OIDCClientSecretRequestSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*OIDCClientSecretRequestStatus)(nil), (*clientsecret.OIDCClientSecretRequestStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_OIDCClientSecretRequestStatus_To_clientsecret_OIDCClientSecretRequestStatus(a.(*OIDCClientSecretRequestStatus), b.(*clientsecret.OIDCClientSecretRequestStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*clientsecret.OIDCClientSecretRequestStatus)(nil), (*OIDCClientSecretRequestStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_clientsecret_OIDCClientSecretRequestStatus_To_v1alpha1_OIDCClientSecretRequestStatus(a.(*clientsecret.OIDCClientSecretRequestStatus), b.(*OIDCClientSecretRequestStatus), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha1_OIDCClientSecretRequest_To_clientsecret_OIDCClientSecretRequest(in *OIDCClientSecretRequest, out *clientsecret.OIDCClientSecretRequest, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_OIDCClientSecretRequestSpec_To_clientsecret_OIDCClientSecretRequestSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1alpha1_OIDCClientSecretRequestStatus_To_clientsecret_OIDCClientSecretRequestStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha1_OIDCClientSecretRequest_To_clientsecret_OIDCClientSecretRequest is an autogenerated conversion function.
func Convert_v1alpha1_OIDCClientSecretRequest_To_clientsecret_OIDCClientSecretRequest(in *OIDCClientSecretRequest, out *clientsecret.OIDCClientSecretRequest, s conversion.Scope) error {
	return autoConvert_v1alpha1_OIDCClientSecretRequest_To_clientsecret_OIDCClientSecretRequest(in, out, s)
}

func autoConvert_clientsecret_OIDCClientSecretRequest_To_v1alpha1_OIDCClientSecretRequest(in *clientsecret.OIDCClientSecretRequest, out *OIDCClientSecretRequest, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_clientsecret_OIDCClientSecretRequestSpec_To_v1alpha1_OIDCClientSecretRequestSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_clientsecret_OIDCClientSecretRequestStatus_To_v1alpha1_OIDCClientSecretRequestStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_clientsecret_OIDCClientSecretRequest_To_v1alpha1_OIDCClientSecretRequest is an autogenerated conversion function.
func Convert_clientsecret_OIDCClientSecretRequest_To_v1alpha1_OIDCClientSecretRequest(in *clientsecret.OIDCClientSecretRequest, out *OIDCClientSecretRequest, s conversion.Scope) error {
	return autoConvert_clientsecret_OIDCClientSecretRequest_To_v1alpha1_OIDCClientSecretRequest(in, out, s)
}

func autoConvert_v1alpha1_OIDCClientSecretRequestList_To_clientsecret_OIDCClientSecretRequestList(in *OIDCClientSecretRequestList, out *clientsecret.OIDCClientSecretRequestList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]clientsecret.OIDCClientSecretRequest)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_v1alpha1_OIDCClientSecretRequestList_To_clientsecret_OIDCClientSecretRequestList is an autogenerated conversion function.
func Convert_v1alpha1_OIDCClientSecretRequestList_To_clientsecret_OIDCClientSecretRequestList(in *OIDCClientSecretRequestList, out *clientsecret.OIDCClientSecretRequestList, s conversion.Scope) error {
	return autoConvert_v1alpha1_OIDCClientSecretRequestList_To_clientsecret_OIDCClientSecretRequestList(in, out, s)
}

func autoConvert_clientsecret_OIDCClientSecretRequestList_To_v1alpha1_OIDCClientSecretRequestList(in *clientsecret.OIDCClientSecretRequestList, out *OIDCClientSecretRequestList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	out.Items = *(*[]OIDCClientSecretRequest)(unsafe.Pointer(&in.Items))
	return nil
}

// Convert_clientsecret_OIDCClientSecretRequestList_To_v1alpha1_OIDCClientSecretRequestList is an autogenerated conversion function.
func Convert_clientsecret_OIDCClientSecretRequestList_To_v1alpha1_OIDCClientSecretRequestList(in *clientsecret.OIDCClientSecretRequestList, out *OIDCClientSecretRequestList, s conversion.Scope) error {
	return autoConvert_clientsecret_OIDCClientSecretRequestList_To_v1alpha1_OIDCClientSecretRequestList(in, out, s)
}

func autoConvert_v1alpha1_OIDCClientSecretRequestSpec_To_clientsecret_OIDCClientSecretRequestSpec(in *OIDCClientSecretRequestSpec, out *clientsecret.OIDCClientSecretRequestSpec, s conversion.Scope) error {
	out.GenerateNewSecret = in.GenerateNewSecret
	out.RevokeOldSecrets = in.RevokeOldSecrets
	return nil
}

// Convert_v1alpha1_OIDCClientSecretRequestSpec_To_clientsecret_OIDCClientSecretRequestSpec is an autogenerated conversion function.
func Convert_v1alpha1_OIDCClientSecretRequestSpec_To_clientsecret_OIDCClientSecretRequestSpec(in *OIDCClientSecretRequestSpec, out *clientsecret.OIDCClientSecretRequestSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_OIDCClientSecretRequestSpec_To_clientsecret_OIDCClientSecretRequestSpec(in, out, s)
}

func autoConvert_clientsecret_OIDCClientSecretRequestSpec_To_v1alpha1_OIDCClientSecretRequestSpec(in *clientsecret.OIDCClientSecretRequestSpec, out *OIDCClientSecretRequestSpec, s conversion.Scope) error {
	out.GenerateNewSecret = in.GenerateNewSecret
	out.RevokeOldSecrets = in.RevokeOldSecrets
	return nil
}

// Convert_clientsecret_OIDCClientSecretRequestSpec_To_v1alpha1_OIDCClientSecretRequestSpec is an autogenerated conversion function.
func Convert_clientsecret_OIDCClientSecretRequestSpec_To_v1alpha1_OIDCClientSecretRequestSpec(in *clientsecret.OIDCClientSecretRequestSpec, out *OIDCClientSecretRequestSpec, s conversion.Scope) error {
	return autoConvert_clientsecret_OIDCClientSecretRequestSpec_To_v1alpha1_OIDCClientSecretRequestSpec(in, out, s)
}

func autoConvert_v1alpha1_OIDCClientSecretRequestStatus_To_clientsecret_OIDCClientSecretRequestStatus(in *OIDCClientSecretRequestStatus, out *clientsecret.OIDCClientSecretRequestStatus, s conversion.Scope) error {
	out.GeneratedSecret = in.GeneratedSecret
	out.TotalClientSecrets = in.TotalClientSecrets
	return nil
}

// Convert_v1alpha1_OIDCClientSecretRequestStatus_To_clientsecret_OIDCClientSecretRequestStatus is an autogenerated conversion function.
func Convert_v1alpha1_OIDCClientSecretRequestStatus_To_clientsecret_OIDCClientSecretRequestStatus(in *OIDCClientSecretRequestStatus, out *clientsecret.OIDCClientSecretRequestStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_OIDCClientSecretRequestStatus_To_clientsecret_OIDCClientSecretRequestStatus(in, out, s)
}

func autoConvert_clientsecret_OIDCClientSecretRequestStatus_To_v1alpha1_OIDCClientSecretRequestStatus(in *clientsecret.OIDCClientSecretRequestStatus, out *OIDCClientSecretRequestStatus, s conversion.Scope) error {
	out.GeneratedSecret = in.GeneratedSecret
	out.TotalClientSecrets = in.TotalClientSecrets
	return nil
}

// Convert_clientsecret_OIDCClientSecretRequestStatus_To_v1alpha1_OIDCClientSecretRequestStatus is an autogenerated conversion function.
func Convert_clientsecret_OIDCClientSecretRequestStatus_To_v1alpha1_OIDCClientSecretRequestStatus(in *clientsecret.OIDCClientSecretRequestStatus, out *OIDCClientSecretRequestStatus, s conversion.Scope) error {
	return autoConvert_clientsecret_OIDCClientSecretRequestStatus_To_v1alpha1_OIDCClientSecretRequestStatus(in, out, s)
}
