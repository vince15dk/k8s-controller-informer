package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/vince15dk/k8s-controller-informer/api/types/v1alpha1"
	clientV1 "github.com/vince15dk/k8s-controller-informer/clientset/v1alpha1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
)

var state = "create"

func main() {
	kubeconfig := flag.String("kubeconfig", "/Users/nhn/.kube/config", "location to your kubeconfig file")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		// handle error
		fmt.Printf("error %s, building config from flags\n", err.Error())
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("error %s, getting inclusterconfig", err.Error())
		}
	}

	v1alpha1.AddToScheme(scheme.Scheme)
	_, err = kubernetes.NewForConfig(config)
	instClientSet, err := clientV1.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	queue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "instances")
	store := WatchResources(ctx, instClientSet, queue)
	for {
		instancesFromStore := store.List()
		fmt.Printf("instances in store: %d\n", len(instancesFromStore))
		item, shutdown := queue.Get()
		if shutdown {
			return
		}
		instance := item.(*v1alpha1.Instance)
		c := &InstanceController{instance: instance}
		defer queue.Forget(item)
		key, err := cache.MetaNamespaceKeyFunc(item)
		//ns, name, err := cache.SplitMetaNamespaceKey(key)
		if err != nil {
			fmt.Printf("getting key from cache %s\n", err.Error())
			return
		}
		a,_,_ := store.GetByKey(key)
		if state == "check" && a != nil{
			c.ListInstance()
			queue.Done(item)
		}else if state == "create" && a != nil{
			c.CreateInstance()
			queue.Done(item)
		}else if state == "delete" && a == nil{
			c.DeleteInstance()
			queue.Done(item)
		}else{
			panic("you should never be here")
		}
	}
}
