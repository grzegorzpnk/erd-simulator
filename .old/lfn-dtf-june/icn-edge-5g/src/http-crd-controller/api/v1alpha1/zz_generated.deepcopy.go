// +build !ignore_autogenerated

/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AMBRConfig) DeepCopyInto(out *AMBRConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AMBRConfig.
func (in *AMBRConfig) DeepCopy() *AMBRConfig {
	if in == nil {
		return nil
	}
	out := new(AMBRConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DNNConfig) DeepCopyInto(out *DNNConfig) {
	*out = *in
	out.AMBR = in.AMBR
	out.Flow = in.Flow
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DNNConfig.
func (in *DNNConfig) DeepCopy() *DNNConfig {
	if in == nil {
		return nil
	}
	out := new(DNNConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlowRule) DeepCopyInto(out *FlowRule) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlowRule.
func (in *FlowRule) DeepCopy() *FlowRule {
	if in == nil {
		return nil
	}
	out := new(FlowRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Mobilesubscriber) DeepCopyInto(out *Mobilesubscriber) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Mobilesubscriber.
func (in *Mobilesubscriber) DeepCopy() *Mobilesubscriber {
	if in == nil {
		return nil
	}
	out := new(Mobilesubscriber)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Mobilesubscriber) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MobilesubscriberList) DeepCopyInto(out *MobilesubscriberList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Mobilesubscriber, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MobilesubscriberList.
func (in *MobilesubscriberList) DeepCopy() *MobilesubscriberList {
	if in == nil {
		return nil
	}
	out := new(MobilesubscriberList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MobilesubscriberList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MobilesubscriberSpec) DeepCopyInto(out *MobilesubscriberSpec) {
	*out = *in
	if in.SNSSAI != nil {
		in, out := &in.SNSSAI, &out.SNSSAI
		*out = make([]SNSSAIConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MobilesubscriberSpec.
func (in *MobilesubscriberSpec) DeepCopy() *MobilesubscriberSpec {
	if in == nil {
		return nil
	}
	out := new(MobilesubscriberSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MobilesubscriberStatus) DeepCopyInto(out *MobilesubscriberStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MobilesubscriberStatus.
func (in *MobilesubscriberStatus) DeepCopy() *MobilesubscriberStatus {
	if in == nil {
		return nil
	}
	out := new(MobilesubscriberStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SNSSAIConfig) DeepCopyInto(out *SNSSAIConfig) {
	*out = *in
	if in.DNN != nil {
		in, out := &in.DNN, &out.DNN
		*out = make([]DNNConfig, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SNSSAIConfig.
func (in *SNSSAIConfig) DeepCopy() *SNSSAIConfig {
	if in == nil {
		return nil
	}
	out := new(SNSSAIConfig)
	in.DeepCopyInto(out)
	return out
}
