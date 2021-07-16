package v1alpha1

import (
	"github.com/vince15dk/k8s-controller-informer/api/types/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type InstanceV1Interface interface {
	Instances(namespace string) InstanceInterface
}

type InstanceV1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*InstanceV1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1alpha1.GroupName, Version: v1alpha1.GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &InstanceV1Client{restClient: client}, nil
}

func (c *InstanceV1Client) Instances(namespace string) InstanceInterface {
	return &instanceClient{
		restClient: c.restClient,
		ns: namespace,
	}
}