package main

import (
	"fmt"
	"github.com/vince15dk/k8s-controller-informer/api/types/v1alpha1"
	"github.com/vince15dk/k8s-controller-informer/model"
	"io/ioutil"
	"net/http"
)

var (
	urlCreateRepo = "https://api-identity.infrastructure.cloud.toast.com/v2.0/tokens"
	inst = &model.Instance{}
)

func SettingAuthHeader(h *http.Header, token model.CreateAccessResponse)*http.Header{
	h.Set("Content-Type", "application/json")
	h.Set("X-Auth-Token", token.Access.Token.ID)
	return h
}

func CreateInstance(instance *v1alpha1.Instance){
	token := GetToken(instance)
	// Setting Auth Header
	newHeader := SettingAuthHeader(&http.Header{}, token)

	// Creating Instance
	url := "https://kr1-api-instance.infrastructure.cloud.toast.com/v2/" + instance.Spec.TenantId + "/servers"
		// Mutating Instance object
	inst.Server.Name = instance.Spec.InstName
	inst.Server.ImageRef = model.Images[instance.Spec.ImageRef]
	inst.Server.FlavorRef = model.Flavors[instance.Spec.FlavorRef]
	inst.Server.Networks = []model.Subnet{{instance.Spec.SubnetId}}
	inst.Server.KeyName = instance.Spec.KeyName
	inst.Server.MinCount = instance.Spec.Count
	inst.Server.BlockDeviceMappingV2 = []model.BlockDevice{{UUID: model.Images[instance.Spec.ImageRef], BootIndex: 0,
		VolumeSize: instance.Spec.BlockSize, DeviceName: "vda", SourceType: "image", DestinationType: "volume", DeleteOnTermination: 1}}

	resp, err := PostHandleFunc(url, inst, *newHeader)
	if err != nil {
		fmt.Println(err)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	fmt.Println(string(bytes))
}

func DeleteInstance(instance *v1alpha1.Instance){
	fmt.Println("deleting instance", instance.Spec.ImageRef)
}