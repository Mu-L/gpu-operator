/**
# Copyright (c) NVIDIA CORPORATION.  All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
**/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/NVIDIA/gpu-operator/api/nvidia/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeNVIDIADrivers implements NVIDIADriverInterface
type FakeNVIDIADrivers struct {
	Fake *FakeNvidiaV1alpha1
}

var nvidiadriversResource = v1alpha1.SchemeGroupVersion.WithResource("nvidiadrivers")

var nvidiadriversKind = v1alpha1.SchemeGroupVersion.WithKind("NVIDIADriver")

// Get takes name of the nVIDIADriver, and returns the corresponding nVIDIADriver object, and an error if there is any.
func (c *FakeNVIDIADrivers) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.NVIDIADriver, err error) {
	emptyResult := &v1alpha1.NVIDIADriver{}
	obj, err := c.Fake.
		Invokes(testing.NewRootGetActionWithOptions(nvidiadriversResource, name, options), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.NVIDIADriver), err
}

// List takes label and field selectors, and returns the list of NVIDIADrivers that match those selectors.
func (c *FakeNVIDIADrivers) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.NVIDIADriverList, err error) {
	emptyResult := &v1alpha1.NVIDIADriverList{}
	obj, err := c.Fake.
		Invokes(testing.NewRootListActionWithOptions(nvidiadriversResource, nvidiadriversKind, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.NVIDIADriverList{ListMeta: obj.(*v1alpha1.NVIDIADriverList).ListMeta}
	for _, item := range obj.(*v1alpha1.NVIDIADriverList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested nVIDIADrivers.
func (c *FakeNVIDIADrivers) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchActionWithOptions(nvidiadriversResource, opts))
}

// Create takes the representation of a nVIDIADriver and creates it.  Returns the server's representation of the nVIDIADriver, and an error, if there is any.
func (c *FakeNVIDIADrivers) Create(ctx context.Context, nVIDIADriver *v1alpha1.NVIDIADriver, opts v1.CreateOptions) (result *v1alpha1.NVIDIADriver, err error) {
	emptyResult := &v1alpha1.NVIDIADriver{}
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateActionWithOptions(nvidiadriversResource, nVIDIADriver, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.NVIDIADriver), err
}

// Update takes the representation of a nVIDIADriver and updates it. Returns the server's representation of the nVIDIADriver, and an error, if there is any.
func (c *FakeNVIDIADrivers) Update(ctx context.Context, nVIDIADriver *v1alpha1.NVIDIADriver, opts v1.UpdateOptions) (result *v1alpha1.NVIDIADriver, err error) {
	emptyResult := &v1alpha1.NVIDIADriver{}
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateActionWithOptions(nvidiadriversResource, nVIDIADriver, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.NVIDIADriver), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeNVIDIADrivers) UpdateStatus(ctx context.Context, nVIDIADriver *v1alpha1.NVIDIADriver, opts v1.UpdateOptions) (result *v1alpha1.NVIDIADriver, err error) {
	emptyResult := &v1alpha1.NVIDIADriver{}
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceActionWithOptions(nvidiadriversResource, "status", nVIDIADriver, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.NVIDIADriver), err
}

// Delete takes name of the nVIDIADriver and deletes it. Returns an error if one occurs.
func (c *FakeNVIDIADrivers) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(nvidiadriversResource, name, opts), &v1alpha1.NVIDIADriver{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeNVIDIADrivers) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionActionWithOptions(nvidiadriversResource, opts, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.NVIDIADriverList{})
	return err
}

// Patch applies the patch and returns the patched nVIDIADriver.
func (c *FakeNVIDIADrivers) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.NVIDIADriver, err error) {
	emptyResult := &v1alpha1.NVIDIADriver{}
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceActionWithOptions(nvidiadriversResource, name, pt, data, opts, subresources...), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.NVIDIADriver), err
}
