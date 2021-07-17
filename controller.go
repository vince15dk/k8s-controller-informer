package main

import (
	"encoding/json"
	"fmt"
	"github.com/vince15dk/k8s-controller-informer/api/types/v1alpha1"
	"github.com/vince15dk/k8s-controller-informer/model"
	"io/ioutil"
	"net/http"
)

var (
	urlCreateRepo = "https://api-identity.infrastructure.cloud.toast.com/v2.0/tokens"
	inst          = &model.Instance{}
	baseUrl = "https://kr1-api-instance.infrastructure.cloud.toast.com/v2/"
)

func SettingAuthHeader(h *http.Header, token model.CreateAccessResponse) *http.Header {
	h.Set("Content-Type", "application/json")
	h.Set("X-Auth-Token", token.Access.Token.ID)
	return h
}

func CreateInstance(instance *v1alpha1.Instance) {
	// Get Token
	token := GetToken(instance)
	// Setting Auth Header
	newHeader := SettingAuthHeader(&http.Header{}, token)

	// Creating Instance
	url := baseUrl + instance.Spec.TenantId + "/servers"
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
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

}

func DeleteInstance(instance *v1alpha1.Instance) {
	// Get Token
	token := GetToken(instance)
	// Setting Auth Header
	newHeader := SettingAuthHeader(&http.Header{}, token)
	urlGetInstance := "https://kr1-api-instance.infrastructure.cloud.toast.com/v2/" + instance.Spec.TenantId + "/servers/detail"
	newResponse, err := ListHandleFunc(urlGetInstance, *newHeader)
	if err != nil {
		fmt.Println(err)
	}
	//servers := &model.ServerInfo{}
	newBytes, err := ioutil.ReadAll(newResponse.Body)
	if err != nil {
		fmt.Println(err)
	}
	defer newResponse.Body.Close()
	servers := &model.ServerInfo{}
	var serverIds []string
	var num = 1
	err = json.Unmarshal(newBytes, servers)
	for _, v := range servers.Servers {
		for i := num; i <= instance.Spec.Count; i++ {
			if v.Name == fmt.Sprintf("%s-%d", instance.Spec.InstName, i) {
				serverIds = append(serverIds, v.ID)
			}
		}
	}
	for _, v := range serverIds{
		urlDeleteInstance := baseUrl + instance.Spec.TenantId + "/servers/" + v
		resp, err := DeleteHandelFunc(urlDeleteInstance, *newHeader)
		if err != nil {
			fmt.Println(err)
		}
		resp.Body.Close()
	}
}
