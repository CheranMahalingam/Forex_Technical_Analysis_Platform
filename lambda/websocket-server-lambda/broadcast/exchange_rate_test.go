package broadcast

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func TestInitialOhlcData(t *testing.T) {
	// Test DynamoDB query input parameters with valid information
	currentDate := time.Now()
	keyCond := expression.Key("Date").Equal(expression.Value(currentDate.Format("2006-01-02")))
	dbInput, err := InitialOhlcData("EURUSD", keyCond)
	if err != nil {
		t.Errorf("DynamoDB query was invalid, failed with error: %e", err)
	}

	tableName := dbInput.TableName
	attributeNames := dbInput.ExpressionAttributeNames
	attributeValues := dbInput.ExpressionAttributeValues
	if *tableName != "SymbolRateTable" {
		t.Errorf("DynamoDB query was invalid, table name was %s, expected: SymbolRateTable", *tableName)
	}
	if *attributeNames["#0"] != "Date" {
		t.Errorf("DynamoDB query was invalid, attribute 0 was %s, expected: Date", *attributeNames["#0"])
	}
	if *attributeNames["#1"] != "EURUSD" {
		t.Errorf("DynamoDB query was invalid, attribute 1 was %s, expected: Date", *attributeNames["#1"])
	}
	if *attributeNames["#2"] != "Timestamp" {
		t.Errorf("DynamoDB query was invalid, attribute 2 was %s, expected: Date", *attributeNames["#2"])
	}
	if *attributeValues[":0"].S != currentDate.Format("2006-01-02") {
		t.Errorf("DynamoDB query was invalid, expected attribute value was %s, expected: %s", *attributeValues[":0"].S, currentDate.Format("2006-01-02"))
	}
}

func TestGetSymbolStructField(t *testing.T) {
	// Test whether function does not panic if passed empty struct
	emptySymbolStruct := exchangeRateTable{}
	eurusd := getSymbolStructField("EURUSD", &emptySymbolStruct)
	if eurusd.Open != 0 {
		t.Errorf("Symbol struct data was invalid, got: %f, wanted: %d", eurusd.Open, 0)
	}

	// Test response when invalid symbol is passed
	response := getSymbolStructField("eurusd", &emptySymbolStruct)
	if response != nil {
		t.Errorf("Symbol struct data was invalid, got: %v, wanted: %v", response, nil)
	}
}

func TestCheckIsRateValid(t *testing.T) {
	// Test whether invalid if the default struct is passed
	invalidOhlcData := ohlcData{}
	isValid := checkIsRateValid(&invalidOhlcData)
	if isValid != false {
		t.Errorf("OHLC data was invalid, got: %t, wanted: %t", isValid, true)
	}

	// Test whether invalid when one field is zero
	invalidOhlcData = ohlcData{Open: 0, High: 3, Low: 2, Close: 1, Volume: 5}
	isValid = checkIsRateValid(&invalidOhlcData)
	if isValid != false {
		t.Errorf("OHLC data was invalid, got: %t, wanted: %t", isValid, true)
	}

	// Test whether valid when all fields are non zero
	invalidOhlcData = ohlcData{Open: 1, High: 3, Low: 2, Close: 1, Volume: 5}
	isValid = checkIsRateValid(&invalidOhlcData)
	if isValid != true {
		t.Errorf("OHLC data was valid, got: %t, wanted: %t", isValid, true)
	}
}
