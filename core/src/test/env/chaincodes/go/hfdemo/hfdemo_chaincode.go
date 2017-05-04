package hfdemo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type TJUStudentInfoChainCode struct {
}

type Student struct {
	Id string `json:"id"`
	Name string `json:"name"`
	University string `json:"university"`
	Degree string `json:"degree"`
	Start string `json:"start"`
	End string `json:"end"`
	EducationQualifications []string `json:"educationQualifications"`
	InternInfos []InternInfo `json:"interninfos"`
}

type InternInfo struct {
	StudentId string `json:"studentid"`
	Name string `json:"name"`
	WorkingId string `json:"workingid"`
	Company string `json:"company"`
	Department string `json:"department"`
	Position string `json:"position"`
	Start string `json:"start"`
	End string `json:"end"`
}

func main() {
	err := shim.Start(new(TJUStudentInfoChainCode))
	if err != nil {
		fmt.Printf("Error starting TJUStudentInfoChainCode: %s", err)
	}
}

// Init initializes chaincode
func (t *TJUStudentInfoChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke - entry point for invocations
func (t *TJUStudentInfoChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running:" + function)

	// function handler
	// TODO: batch operations
	if function == "insert"  {
		return t.insert(stub, args)
	} else if function == "insertBatch" {
		return t.insertBatch(stub, args)
	} else if function == "remove" {
		return t.remove(stub, args)
	} else if function == "removeBatch" {
		return t.removeBatch(stub, args)
	} else if function == "update" {
		return t.update(stub, args)
	} else if function == "updateBatch" {
		return t.updateBatch(stub, args)
	}


	fmt.Println("invoke did not find func:" + function) //error
	return shim.Error("Received unknown function invocation")
}

func (t *TJUStudentInfoChainCode) insert(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1{
		return shim.Error("Incorrect number of arguments.Excepting 1")
	}
	info := args[0]

	// parse info from json
	var jsonInfo Student
	err = json.Unmarshal([]byte(info), &jsonInfo)
	if err != nil {
		return shim.Error(err.Error())
	}

	// check if already exists
	id := jsonInfo.Id
	tryFetch, err := stub.GetState(id)
	if err != nil {
		return shim.Error("Fail to get by " + id + " " + err.Error())
	}
	if tryFetch != nil {
		fmt.Println("Id " + id + " already exists.Content:" + tryFetch)
		return shim.Error("Id " + id + " already exists.")
	}

	fmt.Println("insert id:" + id + " content:" + info)
	err = stub.PutState(id, info)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("- insert success")
	return shim.Success(nil)
}

func (t *TJUStudentInfoChainCode) insertBatch(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1{
		return shim.Error("Incorrect number of arguments.Excepting 1")
	}
	infos := args[0]

	// parse infos from json
	var jsonInfo []Student
	err = json.Unmarshal([]byte(infos), &jsonInfo)
	if err != nil {
		return shim.Error(err.Error())
	}

	for idx,val := range jsonInfo {
		fmt.Printf("Index:%d,Value:%s", idx, val)
		id := val.Id
		tryFetch, err := stub.GetState(id)
		if err != nil {
			return shim.Error(err.Error())
		}
		if tryFetch != nil {
			fmt.Println("Id " + id + " already exists.Content:" + tryFetch)
			return shim.Error("Id " + id + " already exists.")
		}

		fmt.Println("insert id:" + id + " content:" + val)
		err = stub.PutState(id, val)
		if err != nil {
			return shim.Error(err.Error())
		}
		fmt.Println("- insert success")
	}

	return shim.Success(nil)
}

func (t *TJUStudentInfoChainCode) remove(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// remove by id
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments.Excepting 1.")
	}
	id := args[0]

	fmt.Println("remove id:" + id)
	err = stub.DelState(id)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("- remove success")
	return shim.Success(nil)
}

func (t *TJUStudentInfoChainCode) removeBatch(stub shim.ChaincodeStubInterface, args []string) pb.Response  {
	// remove by id
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments.Excepting 1.")
	}
	ids := args[0]

	var jsonIds []string
	err := json.Unmarshal([]byte(ids), &jsonIds)
	if err != nil {
		return shim.Error(err.Error())
	}

	for idx,val := range jsonIds {
		fmt.Printf("remove index:%d,id:%s\n",idx,val)
		err = stub.DelState(val)
		if err != nil {
			return shim.Error(err.Error())
		}
		fmt.Println("- remove success")
	}

	return shim.Success(nil)
}

func (t *TJUStudentInfoChainCode) update(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// update by id
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments.Excepting 1")
	}
	info := args[0]

	// parse info from json
	var jsonInfo Student
	err = json.Unmarshal([]byte(info), &jsonInfo)
	if err != nil {
		return shim.Error(err.Error())
	}

	// check if already exists
	id := jsonInfo.Id
	tryFetch, err := stub.GetState(id)
	if err != nil {
		return shim.Error("Fail to get by " + id + " " + err.Error())
	} else if tryFetch == nil {
		fmt.Println("Id " + id + " not exists.")
		return shim.Error("Id " + id + " not exists.")
	}

	fmt.Println("update id:" + id + " content:" + info)
	err = stub.PutState(id, info)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("- update success")
	return shim.Success(nil)
}

func (t *TJUStudentInfoChainCode) updateBatch(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1{
		return shim.Error("Incorrect number of arguments.Excepting 1")
	}
	infos := args[0]

	// parse infos from json
	var jsonInfo []Student
	err = json.Unmarshal([]byte(infos), &jsonInfo)
	if err != nil {
		return shim.Error(err.Error())
	}

	for idx,val := range jsonInfo {
		fmt.Printf("Index:%d,Value:%s", idx, val)
		id := val.Id
		tryFetch, err := stub.GetState(id)
		if err != nil {
			return shim.Error(err.Error())
		}
		if tryFetch == nil {
			fmt.Println("Id " + id + " not exists.")
			return shim.Error("Id " + id + " not exists.")
		}

		fmt.Println("update id:" + id + " content:" + val)
		err = stub.PutState(id, val)
		if err != nil {
			return shim.Error(err.Error())
		}
		fmt.Println("- insert success")
	}

	return shim.Success(nil)
}

// Query - query of the chaincode
func (t *TJUStudentInfoChainCode) Query(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("query is running:" + function)

	// function handler
	if function == "query" {
		return t.query(stub, args)
	} else if function == "queryBatch" {
		return t.queryBatch(stub, args)
	}

	fmt.Println("query did not find func:" + function) // error
	return shim.Error("Received unknown function query")
}

func (t *TJUStudentInfoChainCode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// query by id
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments.Excepting 1")
	}
	id := args[0]

	fmt.Printf("- query on:\n%s\n", id)
	res, err := stub.GetState(id)
	if err != nil {
		return shim.Error(err.Error())
	}
	if res == nil {
		return shim.Error("Fail to get state for " + id)
	}
	fmt.Printf("query responses:%s\n", string(res))

	return shim.Success(res)
}

func (t *TJUStudentInfoChainCode) queryBatch(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// query by student id
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments.Expected 2")
	}
	startId := args[0]
	endId := args[1]
	if startId > endId {
		return shim.Error("Incorrect start & end for " + startId + " & " + endId)
	}

	fmt.Printf("- query batch on: \n%s-%s\n", startId, endId)
	res, err := stub.GetStateByRange(startId, endId)
	if err != nil {
		return shim.Error(err.Error())
	}
	if res == nil {
		return shim.Error("Fail to get state from " + startId + " to " + endId)
	}
	fmt.Printf("query responses:%s\n", string(res))

	return shim.Success(res)
}

