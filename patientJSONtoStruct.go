package main

import (
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"io"
	"os"
)

func patientJSONtoStruct() (*[]Patient, error) {
	jsonFile, err := os.Open("list_patients.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	var patients []Patient
	err = json.Unmarshal(byteValue, &patients)
	if err != nil {
		return nil, err
	}
	return &patients, nil
}

func appendToJSON(patient Patient, patients []Patient) error {

	patients = append(patients, patient)
	patientsJSON, err := json.Marshal(patients)

	if err != nil {
		return err
	}

	err = os.WriteFile("list_patients.json", patientsJSON, 0644)
	if err != nil {
		return err
	}
	return nil
}
func removeFromJSON(id string, patients []Patient) error {
	patients = lo.Filter(patients, func(p Patient, index int) bool {
		return p.Guid != id
	})
	patientsJSON, err := json.Marshal(patients)

	if err != nil {
		return err
	}

	err = os.WriteFile("list_patients.json", patientsJSON, 0644)
	if err != nil {
		return err
	}
	return nil
}

func updateJSON(id string, patient Patient, patients []Patient) error {

	patients = lo.Map(patients, func(p Patient, index int) Patient {
		if p.Guid == id {
			return patient
		}
		return p
	})
	patientsJSON, err := json.Marshal(patients)

	if err != nil {
		return err
	}

	err = os.WriteFile("list_patients.json", patientsJSON, 0644)
	if err != nil {
		return err
	}
	return nil
}
