package client_identity

import (
	"crypto/x509"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"log"
	"sync"
)

var lock = &sync.Mutex{}

type X509Identity struct {
	mspID       string
	certificate []byte
}

var id *identity.X509Identity

func GetInstance(mspID string, certificate *x509.Certificate) (*identity.X509Identity, error) {
	if id == nil {
		lock.Lock()
		defer lock.Unlock()
		if id == nil {
			var err error
			log.Println("Creating a client identity.")
			id, err = identity.NewX509Identity(mspID, certificate)
			if err != nil {
				return nil, err
			}
		} else {
			log.Println("Single instance already created.")
		}
	} else {
		log.Println("Single instance already created.")
	}
	return id, nil
}
