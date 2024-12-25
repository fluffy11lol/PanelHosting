package models

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestTariffInitialization(t *testing.T) {
	expected := Tariff{
		ID:    "123",
		Name:  "Basic Plan",
		SSD:   256,
		CPU:   4,
		RAM:   8,
		Price: 500,
	}

	actual := Tariff{
		ID:    "123",
		Name:  "Basic Plan",
		SSD:   256,
		CPU:   4,
		RAM:   8,
		Price: 500,
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Tariff initialization failed.\nExpected: %+v\nGot: %+v", expected, actual)
	}
}

func TestTariffJSONSerialization(t *testing.T) {
	tariff := Tariff{
		ID:    "123",
		Name:  "Basic Plan",
		SSD:   256,
		CPU:   4,
		RAM:   8,
		Price: 500,
	}

	expectedJSON := `{"id":"123","name":"Basic Plan","ssd":256,"cpu":4,"ram":8,"price":500}`

	serialized, err := json.Marshal(tariff)
	if err != nil {
		t.Fatalf("Failed to serialize Tariff: %v", err)
	}

	if string(serialized) != expectedJSON {
		t.Errorf("JSON serialization failed.\nExpected: %s\nGot: %s", expectedJSON, string(serialized))
	}
}

func TestTariffJSONDeserialization(t *testing.T) {
	jsonData := `{"id":"123","name":"Basic Plan","ssd":256,"cpu":4,"ram":8,"price":500}`
	expected := Tariff{
		ID:    "123",
		Name:  "Basic Plan",
		SSD:   256,
		CPU:   4,
		RAM:   8,
		Price: 500,
	}

	var actual Tariff
	err := json.Unmarshal([]byte(jsonData), &actual)
	if err != nil {
		t.Fatalf("Failed to deserialize JSON into Tariff: %v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("JSON deserialization failed.\nExpected: %+v\nGot: %+v", expected, actual)
	}
}
