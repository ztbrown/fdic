package main

import (
  "flag"
  "fmt"
  "log"
  "path/filepath"

  "k8s.io/apimachinery/pkg/labels"
  "k8s.io/apimachinery/pkg/util/runtime"
  "k8s.io/client-go/informers"
  "k8s.io/client-go/kubernetes"
  "k8s.io/client-go/tools/cache"
  "k8s.io/client-go/tools/clientcmd"
  "k8s.io/client-go/util/homedir"
)

func main() {

  var kubeconfig *string

  if home := homedir.HomeDir(); home != "" {
    kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
  } else {
    kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
  }

  flag.Parse()

  config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
  if err != nil {
    panic(err)
  }

  clientset, err := kubernetes.NewForConfig(config)
  if err != nil {
    log.Panic(err.Error())
  }

  // stop signal for the informer
  stopper := make(chan struct{})
  defer close(stopper)

  factory := informers.NewSharedInformerFactory(clientset, 0)
  ingressInformer := factory.Networking().V1().Ingresses()
  informer := ingressInformer.Informer()

  defer runtime.HandleCrash()

  // start informer ->
  go factory.Start(stopper)

  // start to sync and call list
  if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
    runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
    return
  }

  informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
    AddFunc: onAdd,
    UpdateFunc: func(interface{}, interface{}) { fmt.Println("update not implemented") },
    DeleteFunc: func(interface{}) { fmt.Println("delete not implemented") },
  })

  // find pods in one ns, or find pods from --all-namespaces
  lister := ingressInformer.Lister().Ingresses("dev")

  ingress, err := lister.List(labels.Everything())

  if err != nil {
    fmt.Println(err)
  }

  fmt.Println("ingress:", ingress)

  <-stopper
}


func onAdd(obj interface{}) {
  fmt.Println("add ingress not implemented")
}
