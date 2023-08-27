//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
	v1 "open-cluster-management.io/api/operator/v1"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KlusterletConfig) DeepCopyInto(out *KlusterletConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KlusterletConfig.
func (in *KlusterletConfig) DeepCopy() *KlusterletConfig {
	if in == nil {
		return nil
	}
	out := new(KlusterletConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KlusterletConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KlusterletConfigList) DeepCopyInto(out *KlusterletConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]KlusterletConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KlusterletConfigList.
func (in *KlusterletConfigList) DeepCopy() *KlusterletConfigList {
	if in == nil {
		return nil
	}
	out := new(KlusterletConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KlusterletConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KlusterletConfigSpec) DeepCopyInto(out *KlusterletConfigSpec) {
	*out = *in
	if in.Registries != nil {
		in, out := &in.Registries, &out.Registries
		*out = make([]Registries, len(*in))
		copy(*out, *in)
	}
	out.PullSecret = in.PullSecret
	if in.NodePlacement != nil {
		in, out := &in.NodePlacement, &out.NodePlacement
		*out = new(v1.NodePlacement)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KlusterletConfigSpec.
func (in *KlusterletConfigSpec) DeepCopy() *KlusterletConfigSpec {
	if in == nil {
		return nil
	}
	out := new(KlusterletConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KlusterletConfigStatus) DeepCopyInto(out *KlusterletConfigStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KlusterletConfigStatus.
func (in *KlusterletConfigStatus) DeepCopy() *KlusterletConfigStatus {
	if in == nil {
		return nil
	}
	out := new(KlusterletConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Registries) DeepCopyInto(out *Registries) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Registries.
func (in *Registries) DeepCopy() *Registries {
	if in == nil {
		return nil
	}
	out := new(Registries)
	in.DeepCopyInto(out)
	return out
}