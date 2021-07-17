package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/vince15dk/k8s-controller-informer/api/types/v1alpha1"
	"github.com/vince15dk/k8s-controller-informer/model"
	"io/ioutil"
	"log"
	"net/http"
)

func GetToken(instance *v1alpha1.Instance)model.CreateAccessResponse{
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	b := &model.CreateAccessRequest{Auth: model.Tenant{
		TenantId: instance.Spec.TenantId,
		PasswordCredentials: model.UserInfo{
			UserName: instance.Spec.UserName,
			Password: instance.Spec.PassWord,
		},
	}}
	response, _ := PostHandleFunc(urlCreateRepo, b, headers)
	bytes, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	var result model.CreateAccessResponse
	if err := json.Unmarshal(bytes, &result); err != nil {
		log.Println(fmt.Sprintf("error when trying to unmarshal create repo successful response: %s", err.Error()))
	}
	return result
}

func PostHandleFunc(url string, body interface{}, headers http.Header)(*http.Response, error){
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	request.Header = headers
	client := http.Client{}
	return client.Do(request)
}


func ListHandleFunc(url string, headers http.Header)(*http.Response, error){
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header = headers
	client := http.Client{}
	return client.Do(request)
}

func DeleteHandelFunc(url string, headers http.Header)(*http.Response, error){
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header = headers
	client := http.Client{}
	return client.Do(request)
}

