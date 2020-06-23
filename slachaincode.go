package main

import (
	"fmt"
	"bytes"
	"time"
	"strconv"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Logger for logging
var logger = shim.NewLogger("slachaincode")
//smartcontract struct
type SlaData struct {  

}
// Sla Contract data
type SlaContract struct {  
	Id string `json:"id"`
	CreatedBy string `json:"createdBy"`
	CreatedTimeStamp string `json:"createdTimeStamp"`
	UpdatedBy string `json:"updatedBy"`
	UpdatedTimeStamp string `json:"updatedTimeStamp"`
	Country string `json:"country"`
	FstsReqID string `json:"fstsReqID"`
	IsDeleted string `json:"isDeleted"`
	Last_update string `json:"last_update"`
	Size string `json:"size"`
	Status string `json:"status"`
	YOY string `json:"yOY"`
	MsgPerDay string `json:"msgPerDay"`
	MsgPerSec string `json:"msgPerSec"`
	Project_name string `json:"project_name"`
	Serviceid string `json:"serviceid"`
	FstsReqID_FK string `json:"fstsReqID_FK"`
	Approval1 string `json:"approval1"`
	Approval1TimeStamp string `json:"approval1TimeStamp"`
	Approval2 string `json:"approval2"`
	Approval2TimeStamp string `json:"approval2TimeStamp"`
	Frequency string `json:"frequency"`
	Slacontractobj string `json:"slacontractobj"`
	Version string `json:"Version"`
	Receiver  string `json:"receiver"`
	Sender  string `json:"sender"`
	CreateUserName string `json:"createUserName"`
	Approval1UserName string `json:"approval1UserName"`
	Approval2UserName  string `json:"approval2UserName"`
}

//NFR struct

type Nfr struct{
	NfrID string `json:"nfrID"`
	ContractID string `json:"contractID"`
	ContractVersion string `json:"contractVersion"`
	Projectid string `json:"projectid"`
	Projectname string `json:"projectname"`
	Serviceid string `json:"serviceid"`
	Servicename string `json:"servicename"`
	Country string `json:"country"`
	Sender string `json:"sender"`
	Receiver string `json:"receiver"`
	Isbreach string `json:"isbreach"`
	Category string `json:"category"`
	Expected string `json:"expected"`
	Actual string `json:"actual"`
	Breachpercent string `json:"breachpercent"`
	Breachaction string `json:"breachaction"`
	Domain string `json:"domain"`
	Date string `json:"date"`
	Nfrobj string `json:"nfrobj"`
	DataPushDate string `json:"datapushdate"`
}


func (t *SlaData) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("slachaincode Init")
	return shim.Success(nil)
}

func (t *SlaData) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "addSlaData" {
		// Make payment of X units from A to B
		return t.addSlaData(stub, args)
	} else if function == "addNfrData" {
		// Make payment of X units from A to B
		return t.addNfrData(stub, args)
	} else if function == "getSlaData" {
		// get an data from its state
		return t.getSlaData(stub, args)
	} else if function == "getSlaDataHistory" {
		// get an data from its state
		return t.getSlaDataHistory(stub, args)
	} else if function == "queryNfr" {
		// get an data from its state
		return t.queryNfr(stub, args)
	}
	

	return shim.Error("Invalid invoke function name. Expecting \"addSlaData\" \"getSlaData\" \"getSlaDataHistory\"")
}

// Transaction add sla contract record
func (t *SlaData) addSlaData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, value string    // data holding
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	key = args[0]
	value = args[1]
	slaContract:= SlaContract{}
	errslaContract:= json.Unmarshal([]byte(value), &slaContract)
	if errslaContract != nil {
		logger.Errorf("addSlaContractData : UnMarshalling Error : " + string(err.Error()))
		return shim.Error("Invalid SLAContract data, error in unmarshaling.")
	}

	if existingSlaContractData, _ := stub.GetState(args[0]); len(existingSlaContractData) <= 0 {
		logger.Infof("addSlaContractData: ContractId not exist")
		slaContract.Version="0";
	} else {
		existingSlaContract := SlaContract{}
		errexistingSlaContract := json.Unmarshal(existingSlaContractData, &existingSlaContract)
		if errexistingSlaContract != nil {
			logger.Errorf("addSlaContractData : UnMarshalling Error : " + string(err.Error()))
			return shim.Error("Invalid SLAContract data, error in unmarshaling..")
		}
		existingVersion,_:=strconv.Atoi(existingSlaContract.Version)
		slaContract.Version=strconv.Itoa(existingVersion+1);
	}
	// Write the state back to the ledger
	slaContractByte, err := json.Marshal(slaContract)
	if err != nil {
		logger.Errorf("addSlaContractData : Marshalling Error : " + string(err.Error()))
		return shim.Error("Invalid SLAData data, error in marshaling")
	}
	//Inserting DataBlock to BlockChain
	err = stub.PutState(key, slaContractByte)
	if err != nil {
		logger.Errorf("addNfrData : PutState Failed Error : " + string(err.Error()))
		return shim.Error("Invalid nfr data, error in pushing into ledger")
	}
	return shim.Success(nil)
}

func (t *SlaData) addNfrData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var nfrList []Nfr
	var dataPushDate string
	if len(args) < 2 {
		return shim.Error("Invalid number of arguments provided for transaction")
	}
	dataPushDate= args[0] 
	err = json.Unmarshal([]byte(args[1]), &nfrList)
	if err != nil {
		logger.Errorf("addNfrData : UnMarshalling Error : " + string(err.Error()))
		return shim.Error("Invalid nfr data, error in unmarshaling.")
	}
	for i:=0; i<len(nfrList); i++ {
		var nfrData Nfr
		nfrData = nfrList[i]
		nfrData.DataPushDate = dataPushDate;
		if slaContractData, _ := stub.GetState(nfrData.ContractID); len(slaContractData) <= 0 {
				logger.Infof("addNfrData: ContractId not exist")
				return shim.Error("ContractID not exist")
		} else {
			existingSlaContract := SlaContract{}
			errexistingSlaContract := json.Unmarshal(slaContractData, &existingSlaContract)
			if errexistingSlaContract != nil {
				logger.Errorf("addNfrData : UnMarshalling Error : " + string(err.Error()))
				return shim.Error("Invalid nfr data, error in unmarshaling..")
			}
			nfrData.ContractVersion= existingSlaContract.Version
		}
		nfrData.Nfrobj="1";
		nfrAsByte, err := json.Marshal(nfrData)
		if err != nil {
			logger.Errorf("addNfrData : Marshalling Error : " + string(err.Error()))
			return shim.Error("Invalid nfr data, error in marshaling")
		}
		//Inserting DataBlock to BlockChain
		err = stub.PutState(nfrData.NfrID, nfrAsByte)
		if err != nil {
			logger.Errorf("addNfrData : PutState Failed Error : " + string(err.Error()))
			return shim.Error("Invalid nfr data, error in pushing into ledger")
		}
	}
	resultData := map[string]interface{} {
		"trxnID": stub.GetTxID(),
		"message": "Bulk NFRData added into ledger successfully",
	}

	respJSON, _ := json.Marshal(resultData)
	return shim.Success(respJSON)
}

// query callback representing the query of a chaincode
func (t *SlaData) getSlaData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key string // uniqueId
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting uniqeId to query")
	}

	key = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"no data for  " + key + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"key\":\"" + key + "\",\"value\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

func (t *SlaData) queryNfr(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var jsonResp string
	if len(args) != 1 {
		logger.Errorf("queryNFR:Invalid number of arguments are provided for transaction")
		jsonResp = "{\"Error\":\"Invalid number of arguments are provided for transaction\"}"
		return shim.Error(jsonResp)
	}
	var records []Nfr
	queryString := args[0]
	logger.Infof("Query Selector : " + string(queryString))
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		logger.Errorf("queryNfr:GetQueryResult is Failed with error :" + string(err.Error()))
		jsonResp = "{\"Error\":\"GetQueryResult is Failed with error- \"" + string(err.Error()) + "\"}"
		return shim.Error(jsonResp)
	}
	for resultsIterator.HasNext() {
		record := Nfr{}
		recordBytes, _ := resultsIterator.Next()
		if (string(recordBytes.Value)) == "" {
			continue
		}
		err = json.Unmarshal(recordBytes.Value, &record)
		if err != nil {
			logger.Errorf("queryNfr:Unable to unmarshal nfr retrieved :" + string(err.Error()))
			jsonResp = "{\"Error\":\"GetQueryResult is Failed with error- \"" + string(err.Error()) + "\"}"
			return shim.Error(jsonResp)
		}
		records = append(records, record)
	}
	resultData := map[string]interface{}{
		"status":    "true",
		"Nfr": records,
	}
	respJson, _ := json.Marshal(resultData)
	return shim.Success(respJson)
}

func (t *SlaData) retriveNfrRecords(stub shim.ChaincodeStubInterface, criteria string, indexs ...string) []Nfr {

	var finalSelector string
	records := make([]Nfr, 0)

	if len(indexs) == 0 {
		finalSelector = fmt.Sprintf("{\"selector\":%s }", criteria)

	} else {
		finalSelector = fmt.Sprintf("{\"selector\":%s , \"use_index\" :\"%s\" }", criteria, indexs[0])
	}

	logger.Infof("Query Selector : %s", finalSelector)
	resultsIterator, _ := stub.GetQueryResult(finalSelector)
	for resultsIterator.HasNext() {
		record := Nfr{}
		recordBytes, _ := resultsIterator.Next()
		err := json.Unmarshal(recordBytes.Value, &record)
		if err != nil {
			logger.Infof("Unable to unmarshal Nfr retrived:: %v", err)
		}
		records = append(records, record)
	}
	return records
}

// query callback representing the query of a chaincode
func (t *SlaData) getSlaDataHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key string // uniqueId
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting uniqeId to query")
	}

	key = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetHistoryForKey(key)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"no data for  " + key + "\"}"
		return shim.Error(jsonResp)
	}
	defer Avalbytes.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for Avalbytes.HasNext() {
		response, err := Avalbytes.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")

		buffer.WriteString(string(response.Value))

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getSlaDataHistory returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}



func main() {
	err := shim.Start(new(SlaData))
	if err != nil {
		fmt.Printf("Error starting slachaincode: %s", err)
	}
}
