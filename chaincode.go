package main

import (
	"errors"
	"fmt"
	"encoding/json"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChaincode struct {
}

var thingsIndexStr = "_thingsIndex"			//name for the key/value that will store all things
var accountsIndexStr = "_accountsIndex"		//name for the key/value that will store all accounts

type Thing struct {
	Id 			string	`json:"id"`
	Name		string  `json:"name"`
	Price		int     `json:"price"`
	InMarket	bool	`json:"inMarket"`
	Owner		string	`json:"owner"`
}

type Account struct {
	Username	string  `json:"username"`
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
	} else if function == "createThing" {
		return t.createThing(stub, args)
	} else if function == "sellThing" {
		return t.sellThing(stub, args)
	} else if function == "buyThing" {
		return t.buyThing(stub, args)
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
	err = append_id(stub, accountsIndexStr, account.Username)
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

func (t *SimpleChaincode) createThing(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	// Args
	// 0 	username	json
	// 1 	thing		json

	if len(args) != 2 {
		fmt.Println("Incorrect number of arguments. Expecting 1")
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	var username = args[0]
	var thingJson = args[1]

	var thing Thing
	json.Unmarshal([]byte(thingJson), &thing)

	// Set the owner for the thing
	thing.Owner = username
	thing.Price = 0
	thing.InMarket = false

	// Checking thing existence.
	res, err := stub.GetState(thing.Id)
	if res != nil {
		fmt.Println("Thing already exists.")
		return nil, errors.New("Thing already exists.")
	}
	if err != nil {
		fmt.Println("Error fetching thing.")
		return nil, errors.New("Error fetching thing.")
	}

	// Append the thing to the array of indexes.
	err = append_id(stub, thingsIndexStr, thing.Id)
	if err != nil {
		fmt.Println("Error appending new thing.")
		return nil, errors.New("Error appending new thing.")
	}

	thingBytes, err := json.Marshal(thing)
	err = stub.PutState(thing.Id, thingBytes)
	if err != nil {
		fmt.Println("Error putting thing on ledger.")
		return nil, errors.New("Error putting thing on ledger.")
	}

	return nil, nil
}

func (t *SimpleChaincode) sellThing(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	// Args
	// 0 	username	json
	// 1 	thing		json
	// 2	price		json

	if len(args) != 3 {
		fmt.Println("Incorrect number of arguments. Expecting 3")
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	// var username = args[0]
	var thingId = args[1]
	var price = args[2]

	// Get the thing
	thingAsBytes, err := stub.GetState(thingId)
	if err != nil {
		fmt.Println("Error fetching thingId" + thingId)
		return nil, errors.New("Error fetching thingId " + thingId)
	}

	var thing Thing
	err = json.Unmarshal(thingAsBytes, &thing)
	if err != nil {
		fmt.Println("Error unmarshalling bytes into thing")
		return nil, errors.New("Error unmarshalling bytes into thing")
	}

	p, err := strconv.Atoi(price);
	if err != nil {
		fmt.Println("Invalid price.")
		return nil, errors.New("Invalid price")
	}

	thing.Price = p
	thing.InMarket = true

	thingBytes, err := json.Marshal(thing)
	err = stub.PutState(thing.Id, thingBytes)
	if err != nil {
		fmt.Println("Error putting thing on ledger.")
		return nil, errors.New("Error putting thing on ledger.")
	}

	return nil, nil
}

func (t *SimpleChaincode) buyThing(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	// Args
	// 0 	username	json
	// 1 	thing		json

	if len(args) != 2 {
		fmt.Println("Incorrect number of arguments. Expecting 2")
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	var username = args[0]
	var thingId = args[1]

	// Get the thing
	thingAsBytes, err := stub.GetState(thingId)
	if err != nil {
		fmt.Println("Error fetching thingId" + thingId)
		return nil, errors.New("Error fetching thingId " + thingId)
	}

	var thing Thing
	err = json.Unmarshal(thingAsBytes, &thing)
	if err != nil {
		fmt.Println("Error unmarshalling bytes into thing")
		return nil, errors.New("Error unmarshalling bytes into thing")
	}

	toAsBytes, err := stub.GetState(username)
	if err != nil {
		fmt.Println("Error fetching account" + username)
		return nil, errors.New("Error fetching account " + username)
	}

	// Get current owner account
	var to Account
	err = json.Unmarshal(toAsBytes, &to)
	if err != nil {
		fmt.Println("Error unmarshalling bytes into account")
		return nil, errors.New("Error unmarshalling bytes into account")
	}

	fromAsBytes, err := stub.GetState(thing.Owner)
	if err != nil {
		fmt.Println("Error fetching account" + thing.Owner)
		return nil, errors.New("Error fetching account " + thing.Owner)
	}

	// Get future owner account
	var from Account
	err = json.Unmarshal(fromAsBytes, &from)
	if err != nil {
		fmt.Println("Error unmarshalling bytes into account")
		return nil, errors.New("Error unmarshalling bytes into account")
	}

	if thing.Price >= to.Money {
		fmt.Println("Not enough money")
		return nil, errors.New("Not enough money")
	}

	// Perform trade
	to.Money -= thing.Price
	from.Money += thing.Price

	thing.Owner = username
	thing.Price = 0
	thing.InMarket = false

	// Save thing
	thingBytes, err := json.Marshal(thing)
	err = stub.PutState(thing.Id, thingBytes)
	if err != nil {
		fmt.Println("Error putting thing on ledger.")
		return nil, errors.New("Error putting thing on ledger.")
	}

	toBytes, err := json.Marshal(to)
	err = stub.PutState(to.Username, toBytes)
	if err != nil {
		fmt.Println("Error putting account on ledger.")
		return nil, errors.New("Error putting account on ledger.")
	}

	fromBytes, err := json.Marshal(from)
	err = stub.PutState(from.Username, fromBytes)
	if err != nil {
		fmt.Println("Error putting account on ledger.")
		return nil, errors.New("Error putting account on ledger.")
	}

	value, err := json.Marshal(to)
	return value, nil
}

// =============================================================================
// Query
// =============================================================================
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running: " + function)

	// Handle different functions
	if function == "getAccount" {
		return t.getAccount(stub, args)
	} else if function == "getThings" {
		return t.getThings(stub, args)
	}


	fmt.Println("Received unknown query function: " + function)
	return nil, errors.New("Received unknown query function" + function)
}

func (t *SimpleChaincode) getAccount(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	// Args
	//	0    username

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

func (t *SimpleChaincode) getThings(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	//Args
	//  market
	//  username


	if len(args) != 1 && len(args) != 2 {
		fmt.Println("Incorrect number of arguments. Expecting 1 or 2")
		return nil, errors.New("Incorrect number of arguments. Expecting 1 or 2")
	}

	market, err := strconv.ParseBool(args[0])
	var me = len(args) == 2
	var username string
	if me {
		username = args[1]
	}

	thingsBytes, err := stub.GetState(thingsIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get " + thingsIndexStr)
	}

	// Unmarshal the index
	var thingsIndex []string
	json.Unmarshal(thingsBytes, &thingsIndex)

	var things []Thing
	for _, thingId := range thingsIndex {
		bytes, err := stub.GetState(thingId)
		if err != nil {
			fmt.Println("Unable to get thing with ID: " + thingId)
			return nil, errors.New("Unable to get thing with ID: " + thingId)
		}

		var thing Thing
		json.Unmarshal(bytes, &thing)
		if (me && thing.Owner == username || !me) && (market && thing.InMarket || !market) {
			things = append(things, thing)
		}
	}

	thingsJson, err := json.Marshal(things)
	if err != nil {
		fmt.Println("Could not convert clients to JSON ")
		return nil, errors.New("Could not convert clients to JSON ")
	}

	return thingsJson, nil
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

func remove_id(stub *shim.ChaincodeStub, indexStr string, id string) (error) {

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

	//remove marble from index
	for i,val := range tmpIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + id)
		if val == id{															//find the correct marble
			fmt.Println("found item")
			tmpIndex = append(tmpIndex[:i], tmpIndex[i+1:]...)			//remove it
			for x:= range tmpIndex{											//debug prints...
				fmt.Println(string(x) + " - " + tmpIndex[x])
			}
			break
		}
	}

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
