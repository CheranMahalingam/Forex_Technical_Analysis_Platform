package broadcast

// Parses ExchangRateTable to construct a websocket payload containing ohlc data
func createCallbackSymbolMessage(symbolRate *[]exchangeRateTable, symbol string) *[]CallbackMessageSymbol {
	var newCallbackMessageData []CallbackMessageSymbol

	for _, rate := range *symbolRate {
		// Gets ohlc data for a particular currency pair
		symbolField := getSymbolStructField(symbol, &rate)

		// Ensures all fields for ohlc data are non-zero
		if checkIsRateValid(symbolField) {
			newData := CallbackMessageSymbol{
				Timestamp: rate.Date + " " + rate.Timestamp,
				Open:      symbolField.Open,
				High:      symbolField.High,
				Low:       symbolField.Low,
				Close:     symbolField.Close,
				Volume:    symbolField.Volume,
			}
			newCallbackMessageData = append(newCallbackMessageData, newData)
		}
	}
	return &newCallbackMessageData
}

// Parses Inference to construct a websocket payload containing exchange rate predictions
func createCallbackInferenceMessage(inferenceList *[]inferenceTable, symbol string) *CallbackMessageInference {
	// Gets inferences for correct symbol
	inferenceField := getInferenceStructField(symbol, &(*inferenceList)[0])
	return &CallbackMessageInference{
		Inference: *inferenceField,
		Date:      (*inferenceList)[0].Latest,
	}
}
