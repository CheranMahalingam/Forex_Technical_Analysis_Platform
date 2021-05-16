package subscription

func createCallbackMessage(symbolRate *[]ExchangeRateTable, symbol string) *[]CallbackMessageData {
	var newCallbackMessageData []CallbackMessageData
	for _, rate := range *symbolRate {
		symbolField := getSymbolStructField(symbol, &rate)
		if checkIsRateValid(symbolField) {
			newData := CallbackMessageData{
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

func createCallbackInferenceMessage(inferenceList *[]Inference, symbol string) *CallbackMessageInference {
	inferenceField := getInferenceStructField(symbol, (*inferenceList)[0])
	return &CallbackMessageInference{
		Inference: *inferenceField,
		Date:      (*inferenceList)[0].Latest,
	}
}
