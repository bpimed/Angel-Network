package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Patient struct {
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
}

func (s *SmartContract) Set(ctx contractapi.TransactionContextInterface, patientID string, name string, lastname string) error {

	patient := Patient{
		Name:     name,
		Lastname: lastname,
	}

	patientAsBytes, err := json.Marshal(patient)
	if err != nil {
		fmt.Printf("Marshal error: %s", err.Error())
		return err
	}

	return ctx.GetStub().PutState(patientID, patientAsBytes)
}

func (s *SmartContract) Query(ctx contractapi.TransactionContextInterface, patientID string) (*Patient, error) {

	patientAsBytes, err := ctx.GetStub().GetState(patientID)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if patientAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", patientID)
	}

	patient := new(Patient)

	err = json.Unmarshal(patientAsBytes, patient)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal error. %s", err.Error())
	}

	return patient, nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error creating PatientControl chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting PatientControl chaincode: %s", err.Error())
	}
}
