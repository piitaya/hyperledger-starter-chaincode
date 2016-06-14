package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	"time"
	"strings"
	"github.com/openblockchain/obc-peer/openchain/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

var itemIndexStr = "_itemindex"                             //name for the key/value that will store a list of all known marbles
var userIndexStr = "_userindex"	                            //name for the key/value that will store all open trades

type Item struct{
	Name       string  `json:"name"`	                    //the fieldtags are needed to keep case from bouncing around
	Price      int     `json:"price"`
}

type User struct{
	Name       string  `json:"name"`
	Size       int     `json:"size"`
}
