package web

import (
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"net/http"
	"rest-api-go/internal/controller"
)

func (setup *OrgSetup) Invoke(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Invoke request")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %s", err)
		return
	}
	chainCodeName := r.FormValue("chaincodeid")
	channelID := r.FormValue("channelid")
	function := r.FormValue("function")
	args := r.Form["args"]
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)

	chainCodeProxy := controller.NewChaincodeProxy(&setup.UserRepo)

	forward := func() (string, error) {
		txn_proposal, err := contract.NewProposal(function, client.WithArguments(args...))
		if err != nil {
			return "", fmt.Errorf("Error creating txn proposal: %s", err)
		}
		txn_endorsed, err := txn_proposal.Endorse()
		if err != nil {
			return "", fmt.Errorf("Error endorsing txn: %s", err)
		}
		txn_committed, err := txn_endorsed.Submit()
		if err != nil {
			return "", fmt.Errorf("Error submitting transaction: %s", err)
		}
		return fmt.Sprintf("Transaction ID : %s Response: %s", txn_committed.TransactionID(), txn_endorsed.Result()), nil
	}

	response, err := chainCodeProxy.ValidateAndForward(function, args, forward)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	fmt.Fprintf(w, "%s", response)
}
