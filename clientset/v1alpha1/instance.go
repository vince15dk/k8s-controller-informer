package v1alpha1

import (
	"context"
	"github.com/vince15dk/k8s-controller-informer/api/types/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"time"
)

type InstanceInterface interface {
	List(ctx context.Context, opts metav1.ListOptions) (*v1alpha1.InstanceList, error)
	Get(ctx context.Context, name string, options metav1.GetOptions) (*v1alpha1.Instance, error)
	Create(ctx context.Context, instance *v1alpha1.Instance, opts metav1.CreateOptions) (*v1alpha1.Instance, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Update(ctx context.Context, instance *v1alpha1.Instance, opts metav1.UpdateOptions) (*v1alpha1.Instance, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
}


type instanceClient struct {
	restClient rest.Interface
	ns         string
}

func (c *instanceClient) List(ctx context.Context, opts metav1.ListOptions) (*v1alpha1.InstanceList, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil{
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result := v1alpha1.InstanceList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("instances").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *instanceClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1alpha1.Instance, error) {
	result := v1alpha1.Instance{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("instances").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *instanceClient) Create(ctx context.Context, instance *v1alpha1.Instance, opts metav1.CreateOptions) (*v1alpha1.Instance, error) {
	result := v1alpha1.Instance{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("instances").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(instance).
		Do(ctx).
		Into(&result)
	return &result, err
}

func (c *instanceClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil{
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("instances").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

func (c *instanceClient) Update(ctx context.Context, instance *v1alpha1.Instance, opts metav1.UpdateOptions) (*v1alpha1.Instance, error) {
	result := v1alpha1.Instance{}
	err := c.restClient.Put().
		Namespace(c.ns).
		Resource("instances").
		Name(instance.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(instance).
		Do(ctx).
		Into(&result)
	return &result, err
}

func (c *instanceClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.restClient.Delete().
		Namespace(c.ns).
		Resource("instances").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}