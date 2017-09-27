/*
 * The smart contract for Employee Office Use-Case
 * @ Author : ashutosh.phoujdar@oracle.com
 */
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the Employee structure, with 4 properties.  Structure tags are used by encoding/json library

type Employee struct {
	EmployeeID    string `json:"empId"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	DateOfJoining string `json:"dateOfJoining"`
	OfficeID      string `json:"officeId"`
}

// Define the Office structure, with 4 properties.  Structure tags are used by encoding/json library
type Office struct {
	OfficeID     string `json:"officeId"`
	BuildingName string `json:"buildingName"`
	StreetName   string `json:"streetName"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
}

/*
 * The Init method is called when the Smart Contract "EmployeeOffice" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "EmployeeOffice"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "createEmployee" {
		return s.createEmployee(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createOffice" {
		return s.createOffice(APIstub, args)
	} else if function == "assignOffice" {
		return s.assignOffice(APIstub, args)
	} else if function == "updateEmployee" {
		return s.updateEmployee(APIstub, args)
	} else if function == "queryAllEmployees" {
		return s.queryAllEmployees(APIstub)
	} else if function == "queryAllOffice" {
		return s.queryAllOffices(APIstub)
	} else if function == "queryEmployeeOfficeName" {
		return s.queryEmployeeOfficeName(APIstub, args)
	} else if function == "queryEmployeesInOffice" {
		return s.queryEmployeesInOffice(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

/*
Method : initLedger
Parameters : NONE
Description : This method is used to initialize the Blocks in Hyperledger
*/
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	employees := []Employee{
		Employee{EmployeeID: "EMP1001", FirstName: "Vineet", LastName: "Timble", DateOfJoining: "01/01/2000", OfficeID: "3"},
		Employee{EmployeeID: "EMP1002", FirstName: "Amit", LastName: "Saxena", DateOfJoining: "01/01/2006", OfficeID: "2"},
		Employee{EmployeeID: "EMP1003", FirstName: "Ashutosh", LastName: "Phoujdar", DateOfJoining: "01/01/2014", OfficeID: "2"},
		Employee{EmployeeID: "EMP1004", FirstName: "Niraj", LastName: "Pandey", DateOfJoining: "01/01/2012", OfficeID: "1"},
		Employee{EmployeeID: "EMP1005", FirstName: "Dinesh", LastName: "Juturu", DateOfJoining: "01/01/2015", OfficeID: "2"},
		Employee{EmployeeID: "EMP1006", FirstName: "Rajesh", LastName: "Annaji", DateOfJoining: "01/01/2012", OfficeID: "2"},
	}
	i := 0
	for i < len(employees) {
		fmt.Println("i is ", i)
		empAsBytes, _ := json.Marshal(employees[i])
		APIstub.PutState("EMPLOYEE"+strconv.Itoa(i), empAsBytes)
		fmt.Println("Added", employees[i])
		i = i + 1
	}

	offices := []Office{
		Office{OfficeID: "OFF1", BuildingName: "Nirlon Compound", StreetName: "Off Western Express Highway", City: "Mumbai", State: "Maharashtra", Country: "India"},
		Office{OfficeID: "OFF2", BuildingName: "Global Axis", StreetName: "Road No 9", City: "Bangalore", State: "Karnataka", Country: "India"},
		Office{OfficeID: "OFF3", BuildingName: "Ambrosia", StreetName: "Bavdhan Khurd", City: "Pune", State: "Maharashtra", Country: "India"},
	}

	j := 0
	for j < len(offices) {
		fmt.Println("j is ", j)
		officeAsBytes, _ := json.Marshal(offices[j])
		APIstub.PutState("OFFICE"+strconv.Itoa(j), officeAsBytes)
		fmt.Println("Added", offices[j])
		j = j + 1
	}
	return shim.Success(nil)
}

/*
Method : createEmployee
Parameters : EmployeeID, FirstName, LastName, DateOdJoining, OfficeID
Description : This method is used to add new employee block
*/
func (s *SmartContract) createEmployee(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var emp = Employee{EmployeeID: args[1], FirstName: args[2], LastName: args[3], DateOfJoining: args[4], OfficeID: args[5]}

	empAsBytes, _ := json.Marshal(emp)
	APIstub.PutState(args[0], empAsBytes)

	return shim.Success(nil)
}

/*
Method : updateEmployee
Parameters : EmployeeID, FirstName, LastName, DateOdJoining, OfficeID
Description : This method is used to find an existing employee and update his/her details in block
*/
func (s *SmartContract) updateEmployee(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	empAsBytes, _ := APIstub.GetState(args[1])
	emp := Employee{}
	json.Unmarshal(empAsBytes, &emp)

	emp.FirstName = args[2]
	emp.LastName = args[3]
	emp.DateOfJoining = args[4]
	emp.OfficeID = args[5]

	empAsBytes, _ = json.Marshal(emp)
	APIstub.PutState(args[0], empAsBytes)

	return shim.Success(nil)
}

/*
Method : assignOffice
Parameters : EmployeeID, OfficeID
Description : This method is used to Find an Employee and assign office location id
*/
func (s *SmartContract) assignOffice(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	empAsBytes, _ := APIstub.GetState(args[1])
	emp := Employee{}
	json.Unmarshal(empAsBytes, &emp)
	emp.OfficeID = args[2]

	empAsBytes, _ = json.Marshal(emp)
	APIstub.PutState(args[0], empAsBytes)

	return shim.Success(nil)
}

/*
Method : createOffice
Parameters : OfficeID, BuildingName, StreeName, City, State, Country
Description : This method is used to Add new office location
*/
func (s *SmartContract) createOffice(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	var office = Office{OfficeID: args[1], BuildingName: args[2], StreetName: args[3], City: args[4], State: args[5], Country: args[6]}
	officeAsBytes, _ := json.Marshal(office)
	APIstub.PutState(args[0], officeAsBytes)

	return shim.Success(nil)
}

/*
Method : queryAllEmployees
Parameters : NONE
Description : This method is used to query all employee blocks
*/
func (s *SmartContract) queryAllEmployees(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "EMP1001"
	endKey := "EMP1006"
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

	fmt.Printf("- queryAllEmployees:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*
Method : queryAllOffices
Parameters : NONE
Description : This method is used to query all office blocks
*/
func (s *SmartContract) queryAllOffices(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "OFF1"
	endKey := "OFF3"
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

	fmt.Printf("- queryAllOffices:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*
Method : queryEmployeeOfficeName
Parameters : EmployeeID
Description: This method is used to query office name of given employee
*/
func (s *SmartContract) queryEmployeeOfficeName(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	empAsBytes, _ := APIstub.GetState(args[1])
	emp := Employee{}
	json.Unmarshal(empAsBytes, &emp)

	officeAsBytes, _ := APIstub.GetState(emp.OfficeID)
	office := Office{}
	json.Unmarshal(officeAsBytes, &office)

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"office\",\"officeId\":\"%s\"}}", office.OfficeID)

	resultsIterator, err := APIstub.GetQueryResult(queryString)

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

	fmt.Printf("- queryEmployeeOfficeName:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*
Method : queryEmployeesInOffice
Parameters : OfficeID
Description : This method is used to query employees in office
*/
func (s *SmartContract) queryEmployeesInOffice(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"employee\",\"officeId\":\"%s\"}}", args[1])

	resultsIterator, err := APIstub.GetQueryResult(queryString)

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

	fmt.Printf("- queryEmployeeOfficeName:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
