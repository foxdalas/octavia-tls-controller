package secretmanager

import (
	"context"
	"fmt"
	sm "github.com/selectel/secretsmanager-go"
	"github.com/selectel/secretsmanager-go/service/certs"
	log "github.com/sirupsen/logrus"
)

type SecretManager struct {
	Client *sm.Client
}

func New(token string) (*SecretManager, error) {
	clinet, err := sm.New(
		sm.WithAuthOpts(&sm.AuthOpts{KeystoneToken: token}),
	)

	return &SecretManager{
		Client: clinet,
	}, err

}

func (secretmanager *SecretManager) List() (certs.GetCertificatesResponse, error) {
	return secretmanager.Client.Certificates.List(context.TODO())
}

func (secretmanager *SecretManager) Get(name string) (certs.Certificate, error) {
	crt, err := secretmanager.List()
	if err != nil {
		return certs.Certificate{}, err
	}
	for _, c := range crt {
		if c.Name == name {
			return secretmanager.Client.Certificates.Get(context.TODO(), c.ID)
		}
	}
	return certs.Certificate{}, fmt.Errorf("certificate not found")
}

func (secretmanager *SecretManager) CreateOrUpdate(name string, cert []byte, key []byte) (certs.Certificate, error) {

	crt, err := secretmanager.Get(name)
	if err != nil {
		if err.Error() == "certificate not found" {
			log.Infof("Certificate %s not found, creating...", name)
			return secretmanager.Create(name, cert, key)
		}
	}
	log.Infof("Certificate %s found with id %s, updating...", name, crt.ID)
	return secretmanager.Update(name, crt.ID, cert, key)
}

func (secretmanager *SecretManager) Create(name string, cert []byte, key []byte) (certs.Certificate, error) {
	pem := certs.Pem{
		Certificates: []string{string(cert)},
		PrivateKey:   string(key),
	}

	crt := certs.CreateCertificateRequest{
		Name: name,
		Pem:  pem,
	}

	return secretmanager.Client.Certificates.Create(context.TODO(), crt)
}

func (secretmanager *SecretManager) Update(name string, id string, cert []byte, key []byte) (certs.Certificate, error) {
	pem := certs.Pem{
		Certificates: []string{string(cert)},
		PrivateKey:   string(key),
	}

	crt := certs.UpdateCertificateVersionRequest{
		Pem: pem,
	}

	err := secretmanager.Client.Certificates.UpdateVersion(context.TODO(), id, crt)
	if err != nil {
		return certs.Certificate{}, err
	}
	return secretmanager.Get(name)
}
