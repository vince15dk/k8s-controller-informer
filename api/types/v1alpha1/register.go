package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "nhn.cloud.io"
const GroupVersion = "v1alpha1"

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: GroupVersion}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownType)
	AddToScheme = SchemeBuilder.AddToScheme
)

func addKnownType(scheme *runtime.Scheme)error{
	scheme.AddKnownTypes(SchemeGroupVersion, &Instance{}, &InstanceList{})

	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}