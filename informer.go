package main

import (
	"context"
	"fmt"
	"github.com/vince15dk/k8s-controller-informer/api/types/v1alpha1"
	client_v1 "github.com/vince15dk/k8s-controller-informer/clientset/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"time"
)

func WatchResources(ctx context.Context, clientSet client_v1.InstanceV1Interface, q workqueue.RateLimitingInterface) cache.Store {
	instStore, instController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return clientSet.Instances("default").List(ctx, lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return clientSet.Instances("default").Watch(ctx, lo)
			},
		},
		&v1alpha1.Instance{},
		30*time.Second,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				fmt.Printf("instance has been created\n")
				state = "create"
				q.Add(obj)
			},
			UpdateFunc: func(old, new interface{}) {
				fmt.Printf("instance has been updated\n")
				state = "check"
				q.Add(new)
			},
			DeleteFunc: func(obj interface{}) {
				fmt.Printf("instance has been deleted\n")
				state = "delete"
				q.Add(obj)
			},
		},
	)
	go instController.Run(wait.NeverStop)
	return instStore
}

