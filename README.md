# informer
This is a GO File
An informer Need to keep informing about an object's state to the CR 
The continuous polling causes overhead. So informer: 
Query the state info and store it in a cache 
It publishes only when there is a change in the state of the object 
Every informer has: 
Reflector- watches and adds the events in Queue 
A FIFO Queue 
Indexer 

Shared informer: 
A single informer it uses a single cache, cache is limited so shared informers help in sharing cache between all the informers. 
By using shared informer we can react to changes- without manually polling for k8s api server 

Create a shared informer 
Steps to create informer: 
[To create informer , you need channel and client ] 
Create a K8s client ( needed for creating informer) 
Create a channel ( for communication to happen) 
Create an informer - add all event handlers 
Start the channel and therefore the informer 

 
import client-go/dynamic/dynamicinformer to create a shared informer 
Initialisation: 
config, err := rest.InClusterConfig()  //  import " k8s.io/client-go/rest" 
// retrieve current config from k8 cluster environment 
if err != nil { 
// Handle error 
} 
clientset, err := kubernetes.NewForConfig(config) // import " k8s.io/client-go/k8s" 
//Create clientset using the necessary config 
if err != nil { 
// Handle error 
} 
Now the created clientset provides access to the API server 
Create a channel( this will later  be used for starting and graceful stopping the informers ) 
stopCh := make(chan struct{}) 
defer close(stopCh)// gurantee the close when function exists 

Define an event handler to handle events 
func handleEvent(obj interface{}) { 
// obj- variable name, type- Go interface of no specific  type 
// Process the event and perform any desired actions 
} 

Create a DynamicsharedInformerFactory -helps to create and manage informers 
 sharedInformerFactory := dynamicinformer.NewDynamicSharedInformerFactory(clientset, time.Second*30) // time:  
how often the informer should check for updates from the API server. 

Create a shared informer 
podInformer := sharedInformerFactory.Core().V1().Pods() 

Setup event handler for the informer: 
podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{ 
AddFunc:    handleEvent, 
UpdateFunc: handleEvent, 
DeleteFunc: handleEvent, 
}) 

Start the informer: 
sharedInformerFactory.Start(stopCh); //Here we initiate the informer's event loop 
 informer will invoke registered event handlers whenever event occurs. 

Informer will continue processing events until stopCh is closed or error occurred. 

Wait for Informer's cache to sync: 
if !cache.WaitForCacheSync(stopCh, podInformer.Informer().HasSynced) { 
// Handle error 
} 


 

 
