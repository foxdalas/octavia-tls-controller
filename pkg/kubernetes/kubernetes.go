package kubernetes

import (
	"context"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"octavia-tls-controller/pkg/helpers"
	"path/filepath"
)

type Kubernetes struct {
	client *kubernetes.Clientset
}

type Secret struct {
	Name      string
	Namespace string
	Data      struct {
		Crt []byte
		Key []byte
	}
}

//func NewOLD() (*Kubernetes, error) {
//	k8s, err := kubernetesInit()
//
//
//	if err != nil {
//		return &Kubernetes{}, err
//	}
//
//	secretmanager, err := secretmanager.New(os.Getenv("SELECTEL_TOKEN")
//	if err != nil {
//		return &Controller{}, err
//	}
//
//	return &Controller{
//		client:  k8s,
//		secretmanager: secretmanager,
//	}, nil
//}
//
//func (controller *Controller) CreateContainerByName(name string) (*containers.Container, error) {
//	createOpts := containers.CreateOpts{
//		Name:       name,
//		Type:       containers.CertificateContainer,
//		SecretRefs: []containers.SecretRef{},
//	}
//	return containers.Create(context.TODO(), controller.octavia, createOpts).Extract()
//}
//
//func (controller *Controller) ContainerByNameIsExist(name string) bool {
//	listOpts := containers.ListOpts{
//		Name: name,
//	}
//	_, err := containers.List(controller.octavia, listOpts).AllPages(context.TODO())
//	if err != nil {
//		return false
//	}
//	return true
//}

func (kubernetes *Kubernetes) GetSecret(namespace string, name string) (*Secret, error) {
	var secret *Secret

	s, err := kubernetes.client.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return secret, err
	}
	log.Infof("Reading secret: %s\n", s.Name)
	if !kubernetes.isValid(s) {
		return &Secret{}, fmt.Errorf("secret %s is not valid", secret.Name)
	}
	return &Secret{
		Name:      s.Name,
		Namespace: s.Namespace,
		Data: struct {
			Crt []byte
			Key []byte
		}{
			Crt: s.Data["tls.crt"],
			Key: s.Data["tls.key"],
		},
	}, nil
}

func (kubernetes *Kubernetes) isValid(secret *v1.Secret) bool {
	if secret.Data["tls.crt"] != nil && secret.Data["tls.key"] != nil {
		tlsCrt := secret.Data["tls.crt"]
		tlsKey := secret.Data["tls.key"]

		err := helpers.ValidateCertificate(tlsCrt, tlsKey)
		if err != nil {
			return false
		}
	}
	return true
}

func New() (*Kubernetes, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return &Kubernetes{}, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return &Kubernetes{}, err
	}
	// create the client
	return &Kubernetes{clientset}, nil
}

func (kubernetes *Kubernetes) GetService(namespace string, name string) (*v1.Service, error) {
	return kubernetes.client.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (kubernetes *Kubernetes) UpdateServiceAnnotation(service *v1.Service, annotation map[string]string) (*v1.Service, error) {
	service.Annotations = annotation
	return kubernetes.client.CoreV1().Services(service.Namespace).Update(context.TODO(), service, metav1.UpdateOptions{})
}
