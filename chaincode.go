package main

import (
	"errors"
	"fmt"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChaincode struct {
}

var itemIndexStr = "_itemindex"          //name for the key/value that will store all items
var accountIndexStr = "_accountindex"	 //name for the key/value that will store all accounts

type Item struct{
	Name       string  `json:"name"`
	Price      int     `json:"price"`
}

type Account struct{
	Username    string  `json:"username"`
	Money       int     `json:"money"`
}

// =============================================================================
// Main
// =============================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// =============================================================================
// Init
// =============================================================================
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	return nil, nil
}

// =============================================================================
// Invoke
// =============================================================================
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running: " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "createAccount" {
		return t.createAccount(stub, args)
	} else if function == "createItem" {
		return t.createItem(stub, args)
	} else if function == "sellItem" {
		return t.sellItem(stub, args)
	} else if function == "buyItem" {
		return t.buyItem(stub, args)
	}

	fmt.Println("Received unknown invoke function: " + function)
	return nil, errors.New("Received unknown invoke function" + function)
}

func (t *SimpleChaincode) createAccount(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	// Args
	// 0 	account		json

	if len(args) != 1 {
		fmt.Println("Incorrect number of arguments. Expecting 1")
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	var accountJson = args[0]
	var account Account
	json.Unmarshal([]byte(accountJson), &account)

	// Checking username existence.
	res, err := stub.GetState(account.Username)
	if res != nil {
		fmt.Println("Username already exists.")
		return nil, errors.New("Username already exists.")
	}
	if err != nil {
		fmt.Println("Error fetching username.")
		return nil, errors.New("Error fetching username.")
	}

	// Append the username to the array of indexes.
	err = append_id(stub, accountIndexStr, account.Username)
	if err != nil {
		fmt.Println("Error appending new username.")
		return nil, errors.New("Error appending new username.")
	}

	accountBytes, err := json.Marshal(account)
	err = stub.PutState(account.Username, accountBytes)
	if err != nil {
		fmt.Println("Error putting account on ledger.")
		return nil, errors.New("Error putting account on ledger.")
	}

	return nil, nil
}

func (t *SimpleChaincode) createItem(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	return nil, nil
}

func (t *SimpleChaincode) sellItem(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	return nil, nil
}

func (t *SimpleChaincode) buyItem(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	return nil, nil
}

// =============================================================================
// Query
// =============================================================================
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running: " + function)

	// Handle different functions
	if function == "getAccount" {
		return t.getAccount(stub, args)
	}
	fmt.Println("Received unknown query function: " + function)
	return nil, errors.New("Received unknown query function" + function)
}

func (t *SimpleChaincode) getAccount(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	// Args
	//	0
	//	username

	if len(args) != 1 {
		fmt.Println("Incorrect number of arguments. Expecting 1")
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	var username = args[0]

	// Get the client
	accountBytes, err := stub.GetState(username)
	if err != nil {
		fmt.Println("Error fetching username.")
		return nil, errors.New("Error fetching username")
	}

	if accountBytes == nil {
		accountBytes, err = json.Marshal(nil)
	}

	return accountBytes, nil
}

// =============================================================================
// Utils
// =============================================================================
func append_id(stub *shim.ChaincodeStub, indexStr string, id string) (error) {

	// Retrieve existing index
	indexAsBytes, err := stub.GetState(indexStr)
	if err != nil {
		return errors.New("Failed to get " + indexStr)
	}
	fmt.Println(indexStr + " retrieved")

	// Unmarshal the index
	var tmpIndex []string
	json.Unmarshal(indexAsBytes, &tmpIndex)
	fmt.Println(indexStr + " unmarshalled")

	// Append the id to the index
	tmpIndex = append(tmpIndex, id)

	// Marshal the index
	jsonAsBytes, err := json.Marshal(tmpIndex)
	if err != nil {
		fmt.Println("Error storing new " + indexStr + " into ledger")
		return errors.New("Error storing new " + indexStr + " into ledger")
	}

	// Store the index into the ledger
	err = stub.PutState(indexStr, jsonAsBytes)
	if err != nil {
		fmt.Println("Error storing new " + indexStr + " into ledger")
		return errors.New("Error storing new " + indexStr + " into ledger")
	}

	return nil
}
