/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	context "context"

	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	applyconfigurationsstoragev1 "k8s.io/client-go/applyconfigurations/storage/v1"
	gentype "k8s.io/client-go/gentype"
	scheme "k8s.io/client-go/kubernetes/scheme"
)

// StorageClassesGetter has a method to return a StorageClassInterface.
// A group's client should implement this interface.
type StorageClassesGetter interface {
	StorageClasses() StorageClassInterface
}

// StorageClassInterface has methods to work with StorageClass resources.
type StorageClassInterface interface {
	Create(ctx context.Context, storageClass *storagev1.StorageClass, opts metav1.CreateOptions) (*storagev1.StorageClass, error)
	Update(ctx context.Context, storageClass *storagev1.StorageClass, opts metav1.UpdateOptions) (*storagev1.StorageClass, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*storagev1.StorageClass, error)
	List(ctx context.Context, opts metav1.ListOptions) (*storagev1.StorageClassList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *storagev1.StorageClass, err error)
	Apply(ctx context.Context, storageClass *applyconfigurationsstoragev1.StorageClassApplyConfiguration, opts metav1.ApplyOptions) (result *storagev1.StorageClass, err error)
	StorageClassExpansion
}

// storageClasses implements StorageClassInterface
type storageClasses struct {
	*gentype.ClientWithListAndApply[*storagev1.StorageClass, *storagev1.StorageClassList, *applyconfigurationsstoragev1.StorageClassApplyConfiguration]
}

// newStorageClasses returns a StorageClasses
func newStorageClasses(c *StorageV1Client) *storageClasses {
	return &storageClasses{
		gentype.NewClientWithListAndApply[*storagev1.StorageClass, *storagev1.StorageClassList, *applyconfigurationsstoragev1.StorageClassApplyConfiguration](
			"storageclasses",
			c.RESTClient(),
			scheme.ParameterCodec,
			"",
			func() *storagev1.StorageClass { return &storagev1.StorageClass{} },
			func() *storagev1.StorageClassList { return &storagev1.StorageClassList{} },
			gentype.PrefersProtobuf[*storagev1.StorageClass](),
		),
	}
}
