package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type InstanceSpec struct {
	TenantId string `json:"tenantId"`
	UserName string `json:"userName"`
	PassWord string `json:"password"`
	InstName string `json:"instName"`
	ImageRef string `json:"imageRef"`
	FlavorRef string `json:"flavorRef"`
	SubnetId string `json:"subnetId"`
	KeyName string `json:"keyName"`
	BlockSize int `json:"blockSize"`
	Count int `json:"count"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Instance struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec InstanceSpec `json:"spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type InstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Instance `json:"items"`
}