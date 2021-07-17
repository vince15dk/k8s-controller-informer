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

type InstanceController struct {
	instance *v1alpha1.Instance
}

func SettingAuthHeader(h *http.Header, token model.CreateAccessResponse) *http.Header {
	h.Set("Content-Type", "application/json")
	h.Set("X-Auth-Token", token.Access.Token.ID)
	return h
}

func AddBodyInstance(a *model.Instance, b *v1alpha1.Instance){
	a.Server.Name = b.Spec.InstName
	a.Server.ImageRef = model.Images[b.Spec.ImageRef]
	a.Server.FlavorRef = model.Flavors[b.Spec.FlavorRef]
	a.Server.Networks = []model.Subnet{{b.Spec.SubnetId}}
	a.Server.KeyName = b.Spec.KeyName
	a.Server.MinCount = b.Spec.Count
	a.Server.BlockDeviceMappingV2 = []model.BlockDevice{{UUID: model.Images[b.Spec.ImageRef], BootIndex: 0,
		VolumeSize: b.Spec.BlockSize, DeviceName: "vda", SourceType: "image", DestinationType: "volume", DeleteOnTermination: 1}}
}

func (c *InstanceController)CreateInstance() {
	// Get Token
	token := GetToken(c.instance)
	// Setting Auth Header
	newHeader := SettingAuthHeader(&http.Header{}, token)

	// Creating Instance
	url := baseUrl + c.instance.Spec.TenantId + "/servers"
	// Mutating Instance object
	AddBodyInstance(inst, c.instance)
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

func (c *InstanceController)DeleteInstance() {
	// Get Token
	token := GetToken(c.instance)
	// Setting Auth Header
	newHeader := SettingAuthHeader(&http.Header{}, token)
	urlGetInstance := baseUrl + c.instance.Spec.TenantId + "/servers/detail"
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
		for i := num; i <= c.instance.Spec.Count; i++ {
			if v.Name == fmt.Sprintf("%s-%d", c.instance.Spec.InstName, i) {
				serverIds = append(serverIds, v.ID)
			}
		}
	}
	for _, v := range serverIds{
		urlDeleteInstance := baseUrl + c.instance.Spec.TenantId + "/servers/" + v
		resp, err := DeleteHandelFunc(urlDeleteInstance, *newHeader)
		if err != nil {
			fmt.Println(err)
		}
		resp.Body.Close()
	}
}

func ListInstance(instance *v1alpha1.Instance){
	// Get Token
	token := GetToken(instance)
	// Setting Auth Header
	newHeader := SettingAuthHeader(&http.Header{}, token)
	urlGetInstance := baseUrl + instance.Spec.TenantId + "/servers/detail"
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
	err = json.Unmarshal(newBytes, servers)

	diff := len(servers.Servers) - instance.Spec.Count
	url := baseUrl + instance.Spec.TenantId + "/servers"

	if diff < 0 {
		AddBodyInstance(inst, instance)
		diff *= -1
		if diff == 1{
			inst.Server.Name = instance.Spec.InstName+ fmt.Sprintf("-%d", diff + 1)
		}else{
			inst.Server.Name = instance.Spec.InstName
		}
		inst.Server.MinCount = diff
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

}
