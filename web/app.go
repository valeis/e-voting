package web

import (
	"fmt"
	"net/http"
	"rest-api-go/internal/repository"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// OrgSetup contains organization's config to interact with the network.
type OrgSetup struct {
	OrgName      string
	MSPID        string
	CryptoPath   string
	CertPath     string
	KeyPath      string
	TLSCertPath  string
	PeerEndpoint string
	GatewayPeer  string
	Gateway      client.Gateway
	UserRepo     repository.UserRepository
}

type HandlerFactory interface {
	CreateHandler(action string) http.HandlerFunc
}

type OrgHandlerFactory struct {
	setups OrgSetup
}

func (o OrgHandlerFactory) CreateHandler(action string) http.HandlerFunc {
	switch action {
	case "query":
		return o.setups.Query
	case "invoke":
		return o.setups.Invoke
	default:
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Invalid action", http.StatusNotFound)
		}
	}
}

func Serve(setups OrgSetup) {
	factory := OrgHandlerFactory{setups: setups}
	http.HandleFunc("/query", factory.CreateHandler("query"))
	http.HandleFunc("/invoke", factory.CreateHandler("invoke"))
	fmt.Println("Listening (http://localhost:3000/)...")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		fmt.Println(err)
	}
}
