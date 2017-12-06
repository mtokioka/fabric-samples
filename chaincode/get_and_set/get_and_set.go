package main

import (
  "bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// 放り込む適当なデータ型定義
type SampleData struct {
	Name   string `json:"name"`
	Owner  string `json:"owner"`
}

// Instantiate/Upgrade 時に実行されるが、初期化処理は明示的に別ファンクション（initLedger）に分離した方が良い
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

// メソッド呼び出しの振り分け
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()
	if function == "get" {
		return s.get(APIstub, args)
	} else if function == "set" {
		return s.set(APIstub, args)
	} else if function == "get_all" {
		return s.getAll(APIstub)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

// 指定した key のデータを取ってくるだけ
func (s *SmartContract) get(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	dataAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(dataAsBytes)
}

// 全データを取ってくるだけ

func (s *SmartContract) getAll(APIstub shim.ChaincodeStubInterface) sc.Response {

  startKey := "D:0"
  endKey := "D:99"

  resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
  if err != nil {
    return shim.Error(err.Error())
  }
  defer resultsIterator.Close()

  // buffer is a JSON array containing QueryResults
  var buffer bytes.Buffer
  buffer.WriteString("[")

  bArrayMemberAlreadyWritten := false
  for resultsIterator.HasNext() {
    queryResponse, err := resultsIterator.Next()
    if err != nil {
      return shim.Error(err.Error())
    }
    // Add a comma before array members, suppress it for the first array member
    if bArrayMemberAlreadyWritten == true {
      buffer.WriteString(",")
    }
    buffer.WriteString("{\"Key\":")
    buffer.WriteString("\"")
    buffer.WriteString(queryResponse.Key)
    buffer.WriteString("\"")

    buffer.WriteString(", \"Record\":")
    // Record is a JSON object, so we write as-is
    buffer.WriteString(string(queryResponse.Value))
    buffer.WriteString("}")
    bArrayMemberAlreadyWritten = true
  }
  buffer.WriteString("]")

  fmt.Printf("- queryAllCars:\n%s\n", buffer.String())

  return shim.Success(buffer.Bytes())
}


// 初期化処理
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	data := []SampleData{
          SampleData{Name: "data0", Owner: "owner0"},
          SampleData{Name: "data1", Owner: "owner1"},
          SampleData{Name: "data2", Owner: "owner2"},
          SampleData{Name: "data3", Owner: "owner3"},
	}

	i := 0
	for i < len(data) {
		fmt.Println("i is ", i)
		dataAsBytes, _ := json.Marshal(data[i])
    APIstub.PutState("D:"+strconv.Itoa(i), dataAsBytes)
		fmt.Println("Added", data[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) set(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	var data = SampleData{Name: args[1], Owner: args[2]}

	dataAsBytes, _ := json.Marshal(data)
	APIstub.PutState(args[0], dataAsBytes)

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
