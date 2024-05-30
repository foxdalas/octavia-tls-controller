package cmd

import (
	log "github.com/sirupsen/logrus"
	"octavia-tls-controller/pkg/kubernetes"
	"octavia-tls-controller/pkg/secretmanager"
	"octavia-tls-controller/pkg/token"
)

type Contoller struct {
	Selectel struct {
		Token string
	}
	Certificates        map[string]Certificate
	SelectelCertificate string
}

type Certificate struct {
	Cert []byte
	Key  []byte
	Type string
}

func New(namespace, secretName, certificateName, serviceName string) error {
	controller := &Contoller{
		Certificates: make(map[string]Certificate),
	}

	k8s, err := kubernetes.New()
	if err != nil {
		return err
	}

	t, err := token.GetAuthToken()
	if err != nil {
		return err
	}
	controller.Selectel.Token = t

	sm, err := secretmanager.New(t)
	if err != nil {
		return err
	}

	secret, err := k8s.GetSecret(namespace, secretName)
	if err != nil {
		return err
	}
	log.Infof("Kubernetes Secret %s found", secret.Name)

	// Create or update certificates
	log.Infof("Creating SecretManager Certificate %s", certificateName)
	crt, err := sm.CreateOrUpdate(certificateName, secret.Data.Crt, secret.Data.Key)
	if err != nil {
		return err
	}
	log.Infof("Certificate %s created with id %s", crt.Name, crt.ID)

	// Get and update Kubernetes service annotation
	log.Infof("Updating Kubernetes Service %s", serviceName)
	service, err := k8s.GetService(namespace, serviceName)
	if err != nil {
		return err
	}
	log.Infof("Kubernetes Service %s found", service.Name)
	annotations := service.Annotations
	annotations["loadbalancer.openstack.org/default-tls-container-ref"] = crt.ID
	updatedService, err := k8s.UpdateServiceAnnotation(service, annotations)
	log.Infof("Kubernetes Service %s updated", updatedService.Name)
	log.Infof("Kubernetes Service %s annotations: %v", updatedService.Name, updatedService.Annotations)

	return nil
}
