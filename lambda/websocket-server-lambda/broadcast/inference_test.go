package broadcast

import (
	"testing"
)

func TestInitialInferenceData(t *testing.T) {
	// Test DynamoDB query input parameters with valid information
	dbInput, err := InitialInferenceData("EURUSD")
	if err != nil {
		t.Errorf("DynamoDB query was invalid, failed with error: %e", err)
	}

	tableName := dbInput.TableName
	attributeNames := dbInput.ExpressionAttributeNames
	attributeValues := dbInput.ExpressionAttributeValues
	if *tableName != "TechnicalAnalysisTable" {
		t.Errorf("DynamoDB query was invalid, table name was %s, expected: TechnicalAnalysisTable", *tableName)
	}
	if *attributeNames["#0"] != "Date" {
		t.Errorf("DynamoDB query was invalid, attribute 0 was %s, expected: Date", *attributeNames["#0"])
	}
	if *attributeNames["#1"] != "EURUSDInference" {
		t.Errorf("DynamoDB query was invalid, attribute 1 was %s, expected: EURUSDInference", *attributeNames["#1"])
	}
	if *attributeNames["#2"] != "Time" {
		t.Errorf("DynamoDB query was invalid, attribute 2 was %s, expected: Time", *attributeNames["#2"])
	}
	if *attributeValues[":0"].S != "inference" {
		t.Errorf("DynamoDB query was invalid, expected attribute value was %s, expected: inference", *attributeValues[":0"].S)
	}
}

func TestGetInferenceStructField(t *testing.T) {
	// Test whether function does not panic if passed empty struct
	emptyInferenceStruct := inferenceTable{}
	gbpusdInference := getInferenceStructField("GBPUSD", &emptyInferenceStruct)
	if len(*gbpusdInference) != 0 {
		t.Errorf("Inference struct data was invalid, got: %d, wanted: %d", len(*gbpusdInference), 0)
	}

	// Test response when invalid symbol is passed
	response := getInferenceStructField("eurusd", &emptyInferenceStruct)
	if response != nil {
		t.Errorf("Inference struct data was invalid, got: %d, wanted: %d", len(*response), 0)
	}
}
