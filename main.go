package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"octavia-tls-controller/cmd"
)

func init() {

}

func main() {
	namespace := flag.String("namespace", "default", "Namespace of the controller")
	secretName := flag.String("secret-name", "", "Name of the kubernetes secret secret")
	certificateName := flag.String("certificate-name", "", "Name of the octavia certificate")
	serviceName := flag.String("service-name", "", "Name of the load balancer service")

	flag.Parse()

	err := validateFlags(*namespace, *secretName, *certificateName, *serviceName)
	if err != nil {
		log.Fatalf("Error validating flags: %v", err)
	}

	//var certificates []*containers.SecretRef

	//octaviaContainerName := "ru-1-op-tls-container"

	log.Info("Starting the octavia TLS controller...")

	err = cmd.New(*namespace, *secretName, *certificateName, *serviceName)
	if err != nil {
		log.Fatalf("Error starting the controller: %v", err)
	}

	//for _, secret := range controller.GetSecrets() {

	//log.Infof("Secret %s found", secret.Name)
	//cert, key, err := controller.CreateCertificate(secret.Name, secret.TLS.Algorithm, secret.TLS.BitLength, secret.TLS.Mode, secret.TLS.Crt, secret.TLS.Key)
	//if err != nil {
	//	log.Errorf("Error creating certificate: %v", err)
	//} else {
	//	log.Infof("Certificate %s created: %v", secret.Name, cert)
	//	log.Infof("Key %s created: %v", secret.Name, key)
	//}

	//secret := &secrets.Secret{
	//	Name: secret.Name,
	//	Algorithm:
	//}
	//
	//containers.SecretRef{
	//	Name:      fmt.Sprintf("%s-%s-certificate", secret.Namespace, secret.Name),
	//	SecretRef: certificates.SecretRef{},
	//}
	//
	//certificates = append(certificates, &KubernetesCert{
	//	Cert: secret.TLS.Crt,
	//	Key:  secret.TLS.Key,
	//})
}

//os.Exit(1)
//
//if controller.ContainerByNameIsExist(octaviaContainerName) {
//	log.Infof("Container %s exists", octaviaContainerName)
//} else {
//	log.Infof("Container %s does not exist", octaviaContainerName)
//	container, err := controller.CreateContainerByName(octaviaContainerName)
//	if err != nil {
//		log.Fatalf("Error creating container: %v", err)
//	}
//	log.Infof("Container %s created: %v", octaviaContainerName, container)
//}

//spew.Dump(controller.GetSecrets())
//}

func validateFlags(namespace, secretName, certificateName, serviceName string) error {
	if namespace == "" {
		return fmt.Errorf("namespace flag is required")
	}
	if secretName == "" {
		return fmt.Errorf("secret-name flag is required")
	}
	if certificateName == "" {
		return fmt.Errorf("certificate-name flag is required")
	}
	if serviceName == "" {
		return fmt.Errorf("service-name flag is required")
	}
	return nil
}
