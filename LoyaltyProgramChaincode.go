package main

import (
	"fmt"
	"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// LoyaltyProgramChaincode
type LoyaltyProgramChaincode struct {
}

var stateMerchant = "_merchantState" //name for the key/value that will store a list of all merchants
var stateCustomer = "_stateCustomer" //name for the key/value that will store a list of all customers
var stateTransaction = "_stateTransaction"
var stateCustMerchantStr = "_stateCustMerchantStr"

type MerchantData struct{
	MERCHANT_NAME string `json:"MERCHANT_NAME"`
	POINT_PER_RS string `json:"POINT_PER_RS"`
	JOINING_BONUS string `json:"JOINING_BONUS"`
}

type CustomerData struct{
 customer_first_name string `json:"CUSTOMER_FIRST_NAME"`
 customer_last_name string `json:"CUSTOMER_LAST_NAME"`
 mobile_no string `json:"MOBILE_NO"`
}

type CustomerMerchantData struct{
 customer_id string `json:"CUSTOMER_ID"`
 merchant_id string `json:"MERCHANT_ID"`
 points string `json:"POINTS"`
}

func (t *LoyaltyProgramChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// Initialize the chaincode
	fmt.Println("Init for Merchant Data")
  var emptyMerchantTxs []MerchantData
  jsonAsBytes0, _ := json.Marshal(emptyMerchantTxs)
	var err error
  err = stub.PutState(stateMerchant, jsonAsBytes0)

	fmt.Println("Init for Customer Data")
  var emptyCustTxs []CustomerData
  jsonAsBytes1, _ := json.Marshal(emptyCustTxs)
  err = stub.PutState(stateCustomer, jsonAsBytes1)

	fmt.Println("Init for Customer For a Merchant")
  var emptyCustMerchantTxs []CustomerMerchantData
  jsonAsBytes2, _ := json.Marshal(emptyCustMerchantTxs)
  err = stub.PutState(stateCustMerchantStr, jsonAsBytes2)
	if err != nil {
		return nil, err
	}
	return nil, nil
}


func (t *LoyaltyProgramChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "registerMerchant" {
		return t.registerMerchant(stub, args)
	}
	if function == "registerCustomer" {
		return t.registerCustomer(stub, args)
	}
	if function == "performTransaction" {
		return t.performTransaction(stub, args)
	}
	/*if function == "transferPoints" {
		return t.transferPoints(stub, args)
	}  */

	return nil, nil
}


func (t *LoyaltyProgramChaincode) registerMerchant(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Printf("In register merchant method")
	var MerchantDataObj MerchantData
	var MerchantDataList []MerchantData
	MerchantDataObj.MERCHANT_NAME = args[0]
	MerchantDataObj.POINT_PER_RS = args[1]
	MerchantDataObj.JOINING_BONUS = args[2]
	merchantTxsAsBytes, err := stub.GetState(stateMerchant)
  json.Unmarshal(merchantTxsAsBytes, &MerchantDataList)
  MerchantDataList = append(MerchantDataList, MerchantDataObj)
	jsonAsBytes, _ := json.Marshal(MerchantDataList)
	err = stub.PutState(stateMerchant, jsonAsBytes)

//testing wht is in ledger
	merchantTxnAsBytes, err := stub.GetState(stateMerchant)

//	fmt.Printf("In register merchant method" + merchantTxnAsBytes)

	if err != nil {
		return nil, err
	}
	return merchantTxnAsBytes, nil
}

func (t *LoyaltyProgramChaincode) registerCustomer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var CustomerDataObj CustomerData
  var CustomerDataList []CustomerData

  // Initialize the chaincode
  CustomerDataObj.customer_first_name = args[0]
  CustomerDataObj.customer_last_name = args[1]
	CustomerDataObj.mobile_no = args[2]
  customerTxsAsBytes, err := stub.GetState(stateCustomer)

  json.Unmarshal(customerTxsAsBytes, &CustomerDataList)
  CustomerDataList = append(CustomerDataList, CustomerDataObj)
  jsonAsBytes, _ := json.Marshal(CustomerDataList)
  err = stub.PutState(stateCustomer, jsonAsBytes)

// save loyalty with bonus points
var CustomerTxnDataObj CustomerMerchantData
CustomerTxnDataObj.customer_id = args[2] // customer mobile no is used as unique id for customer
CustomerTxnDataObj.merchant_id = args[3]
if args[3] == "KMT"{
	CustomerTxnDataObj.points = "100"; // as of now hardcoded but todo fetch from merchant joiing bonus value
} else if args[3] == "SMC" {
	CustomerTxnDataObj.points = "	150";
}

t.saveLoyalty(stub, CustomerTxnDataObj)
if err != nil {
	return nil, err
}
 return nil, nil

}

func (t *LoyaltyProgramChaincode) performTransaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var CustomerTxnDataObj CustomerMerchantData

  CustomerTxnDataObj.customer_id = args[0]
  CustomerTxnDataObj.merchant_id = args[1]
	txnAmt, err := strconv.Atoi(args[2]); // txn amt for transaction at KMT/SMC

	var conv_rate_fr_KMT int = 1/25; // as of now hardcoded but todo fetch from merchant point per rs value
	var conv_rate_fr_SMC int = 1/20; // as of now hardcoded but todo fetch from merchant point per rs value
  //calculate points based on txn
	if args[1] == "KMT"{
		CustomerTxnDataObj.points = strconv.Itoa(txnAmt * conv_rate_fr_KMT)
	} else if args[1] == "SMC" {
		CustomerTxnDataObj.points = strconv.Itoa(txnAmt * conv_rate_fr_SMC)
	}

// saving loyalty points after transaction
	t.saveLoyalty(stub, CustomerTxnDataObj)
	if err != nil {
		return nil, err
	}
	return nil, nil

}

func (t *LoyaltyProgramChaincode) saveLoyalty(stub shim.ChaincodeStubInterface, custMerchantData CustomerMerchantData) ([]byte, error) {
	var CustomerTxnDataObj CustomerMerchantData
  var CustomerTxnList []CustomerMerchantData

  // Initialize the chaincode
  CustomerTxnDataObj.customer_id = custMerchantData.customer_id
  CustomerTxnDataObj.merchant_id = custMerchantData.merchant_id
  customerTxnAsBytes, err := stub.GetState(stateCustMerchantStr)

  json.Unmarshal(customerTxnAsBytes, &CustomerTxnList)


	length := len(CustomerTxnList)

	// iterate
	for i := 0; i < length; i++ {
		obj := CustomerTxnList[i]
		if (custMerchantData.customer_id == obj.customer_id && custMerchantData.merchant_id == obj.merchant_id){
			obj.points = obj.points + CustomerTxnDataObj.points
			CustomerTxnList = append(CustomerTxnList, obj)
		}
	}

  jsonAsBytes, _ := json.Marshal(CustomerTxnList)
  err = stub.PutState(stateCustMerchantStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	return nil, nil

}

//transfer points based on UPS
/*
func (t *LoyaltyProgramChaincode) transferPoints(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {


	var CustomerTxnDataObj1 CustomerMerchantData
	var CustomerTxnDataObj2 CustomerMerchantData

	mobile_no_user_1 := args[0]
	mobile_no_user_2 := args[1]
	transfer_points := args[2]
	source_merchant := args[3]
	dest_merchant := args[4]

  str1 := `{"CUSTOMER_ID": "` + mobile_no_user_1 + `", "MERCHANT_ID": "` + source_merchant + `"}`
	res1, err = t.GetLoyaltyPoints(stub shim.ChaincodeStubInterface, str1 rgs []string)
	json.Unmarshal(res1, &CustomerTxnDataObj1)

	str2 := `{"CUSTOMER_ID": "` + mobile_no_user_2 + `", "MERCHANT_ID": "` + dest_merchant + `"}`
	res2, err = t.GetLoyaltyPoints(stub shim.ChaincodeStubInterface, str2 rgs []string)
	json.Unmarshal(res1, &CustomerTxnDataObj2)

  CustomerTxnDataObj1.POINTS = strconv.Itoa(CustomerTxnDataObj1.POINTS - transfer_points)
	jsonAsBytes1, _ := json.Marshal(CustomerTxnDataObj1)
  err1 = stub.PutState(stateCustomer, jsonAsBytes)

	if err1 != nil {
		return nil, err1
	}

  KMT_UPS_RT := 0.10
	SMC_UPS_RT := 0.15
	if (source_merchant == "KMT" && dest_merchant == "SMC"){
		transfer_points = transfer_points + ((transfer_points * KMT_UPS_RT)/SMC_UPS_RT)
		CustomerTxnDataObj2.POINTS = strconv.Itoa(CustomerTxnDataObj2.POINTS + transfer_points)
	}


	jsonAsBytes2, _ := json.Marshal(CustomerTxnDataObj2)
	err2 = stub.PutState(stateCustomer, jsonAsBytes)
	return nil, nil
}
*/

// Query callback representing the query of a chaincode
func (t *LoyaltyProgramChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// get loyalty points

	res, err := t.GetLoyaltyPoints(stub, args)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (t *LoyaltyProgramChaincode) GetLoyaltyPoints(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

  var CustomerTxnList []CustomerMerchantData
	var CustomerTxnListMatched []CustomerMerchantData
	var customerId = args[0]
	var merchantId = args[1]

	customerTxnAsBytes, err := stub.GetState(stateCustMerchantStr)
  json.Unmarshal(customerTxnAsBytes, &CustomerTxnList)

	length := len(CustomerTxnList)

	// iterate
	for i := 0; i < length; i++ {
		obj := CustomerTxnList[i]
		if (customerId == obj.customer_id && merchantId == obj.merchant_id){
			CustomerTxnListMatched = append(CustomerTxnListMatched,obj)
		}
	}

		res, err := json.Marshal(CustomerTxnListMatched)
		if err != nil {
			return nil, err
		}
		return res, nil

}


func main() {
	err := shim.Start(new(LoyaltyProgramChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
