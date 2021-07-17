package main

import (
	"fmt"
	"github.com/vince15dk/k8s-controller-informer/api/types/v1alpha1"
)

var (
	urlCreateRepo = "https://api-identity.infrastructure.cloud.toast.com/v2.0/tokens"
)

func CreateInstance(instance *v1alpha1.Instance){
	fmt.Println("creating instance", instance.Spec.ImageRef)
}

func DeleteInstance(instance *v1alpha1.Instance){
	fmt.Println("deleting instance", instance.Spec.ImageRef)
}