package main

import (
	"fmt"
	"log"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

func main() {

	//Initialise client
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error()) //panic: unable to load in-cluster configuration, KUBERNETES_SERVICE_HOST and KUBERNETES_SERVICE_PORT must be defined
	}

	client, err := dynamic.NewForConfig(config) //create the client
	if err != nil {
		log.Fatalf("Failed to create DYnamic client")
		panic(err.Error())
	}

	//Create Channel
	channel := make(chan struct{})
	defer close(channel)

	//Create a Informer Factory
	gvr := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	factory := dynamicinformer.NewDynamicSharedInformerFactory(client, time.Second*30)

	//Use it to create the informer

	podinformer := factory.ForResource(gvr).Informer()

	//Event Handelers:

	podinformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Println("Add event called")
		},
		UpdateFunc: func(old, new interface{}) {
			fmt.Println("Update event called")
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("delete event")
		},
	})

	//Start the informer
	go factory.Start(channel)

	<-channel

}
