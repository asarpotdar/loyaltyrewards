/*
Copyright IBM Corp 2016 All Rights Reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
   http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main
import (
 "errors"
 "fmt"
 "encoding/json"
 "github.com/hyperledger/fabric/core/chaincode/shim"
)
// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}
var customerTxStr = "_customerTxStr"
var customerMerchnatTxStr = "_customerMerchnatTxStr"
type CustomerData struct{
 customer_first_name string `json:"CUSTOMER_FIRST_NAME"`
 customer_last_name string `json:"CUSTOMER_LAST_NAME"`
 customer_id string `json:"CUSTOMER_ID"`
}
type CustomerMerchantData struct{
 customer_id string `json:"CUSTOMER_ID"`
 merchant_id string `json:"MERCHANT_ID"`
 points string `json:"POINTS"`
}
// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
 err := shim.Start(new(SimpleChaincode))
 if err != nil {
  fmt.Printf("Error starting Simple chaincode: %s", err)
 }
}
// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
 var err error
 if len(args) != 1 {
  return nil, errors.New("Incorrect number of arguments. Expecting 1")
 }
 if function == "customer"{
  t.InitCustData(stub, args);
 }
 if function == "custForMerchant" {
  t.InitCustMerchantData(stub, args);
 }
 if err != nil {
  return nil, err
 }
 return nil, nil
}
func (t *SimpleChaincode) InitCustData(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
 var err error
 if len(args) != 1 {
  return nil, errors.New("Incorrect number of arguments. Expecting 1")
 }
 fmt.Println("Init for Customer Data")
 var emptyCustTxs []CustomerData
 jsonAsBytes, _ := json.Marshal(emptyCustTxs)
 err = stub.PutState(customerTxStr, jsonAsBytes)
 if err != nil {
  return nil, err
 }
 return nil, nil
}
func (t *SimpleChaincode) InitCustMerchantData(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
 var err error
 if len(args) != 1 {
  return nil, errors.New("Incorrect number of arguments. Expecting 1")
 }
 fmt.Println("Init for Customer For a Merchant")
 var emptyCustMerchantTxs []CustomerMerchantData
 jsonAsBytes, _ := json.Marshal(emptyCustMerchantTxs)
 err = stub.PutState(customerMerchnatTxStr, jsonAsBytes)
 if err != nil {
  return nil, err
 }
 return nil, nil
}
// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
 fmt.Println("invoke is running " + function)
 // register customer
 if function == customerTxStr {
  return t.CustomerSignUp(stub, args)
 }
 fmt.Println("invoke did not find func: " + function)     //error
 return nil, errors.New("Received unknown function invocation: " + function)
}
func (t *SimpleChaincode) CustomerSignUp(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
  var CustomerDataObj CustomerData
  var CustomerDataList []CustomerData
  var err error
  if len(args) != 3 {
   return nil, errors.New("Incorrect number of arguments. Need 3 arguments")
  }
  // Initialize the chaincode
  CustomerDataObj.customer_first_name = args[0]
  CustomerDataObj.customer_last_name = args[1]
  CustomerDataObj.balance = args[2]
  fmt.Printf("Input from user:%s\n", CustomerDataObj)
  customerTxsAsBytes, err := stub.GetState(customerTxStr)
  if err != nil {
   return nil, errors.New("Failed to get consumer Transactions")
  }
  json.Unmarshal(customerTxsAsBytes, &CustomerDataList)
  CustomerDataList = append(CustomerDataList, CustomerDataObj)
  jsonAsBytes, _ := json.Marshal(CustomerDataList)
  err = stub.PutState(customerTxStr, jsonAsBytes)
  if err != nil {
   return nil, err
  }
 return nil, errors.New("Received unknown function invocation: ")
}
func (t *SimpleChaincode) CustomerMerchantTransaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
  var CustomerTxnDataObj CustomerMerchantData
  var CustomerTxnList []CustomerMerchantData
  var err error
  if len(args) != 3 {
   return nil, errors.New("Incorrect number of arguments. Need 3 arguments")
  }
  // Initialize the chaincode
  CustomerDataObj.customer_id = args[0]
  CustomerDataObj.merchant_id = args[1]
  //calculate conv rate based on merchnat
  var conv_rate float64 = 1;
  //calculate points based on txn
  CustomerDataObj.points = args[2] * conv_rate
  fmt.Printf("Input from user:%s\n", CustomerDataObj)
  customerTxsAsBytes, err := stub.GetState(customerTxStr)
  if err != nil {
   return nil, errors.New("Failed to get consumer Transactions")
  }
  json.Unmarshal(customerTxsAsBytes, &CustomerDataList)
  CustomerDataList = append(CustomerDataList, CustomerDataObj)
  jsonAsBytes, _ := json.Marshal(CustomerDataList)
  err = stub.PutState(customerTxStr, jsonAsBytes)
  if err != nil {
   return nil, err
  }
 return nil, errors.New("Received unknown function invocation: ")
}

func (t *SimpleChaincode) CustomerTransaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
 if args[2] == "KMT"{
  t.InitCustData(stub, args);
 }
 if args[2] == "SMC" {
  t.InitCustMerchantData(stub, args);
 }
 return nil, errors.New("Received unknown function invocation: ")
}
// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
 fmt.Println("query is running " + function)
 // Handle different functions
 if function == "dummy_query" {           //read a variable
  fmt.Println("hi there " + function)      //error
  return nil, nil;
 }
 fmt.Println("query did not find func: " + function)      //error
 return nil, errors.New("Received unknown function query: " + function)
}
