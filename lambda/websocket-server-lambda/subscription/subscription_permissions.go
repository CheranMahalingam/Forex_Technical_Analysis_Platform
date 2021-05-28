package subscription

// Uses websocket message key to check whether the user should be sent past 24 hours of exchange rates
func permissionToGetRates(message string) bool {
	if message == "EURUSD" || message == "GBPUSD" || message == "USDJPY" || message == "AUDCAD" {
		return true
	}
	return false
}

// Uses websocket message key to check whether the user should be sent latest ML inference
func permissionToGetInference(message string) bool {
	if message == "EURUSDInference" || message == "GBPUSDInference" || message == "USDJPYInference" || message == "AUDCADInference" {
		return true
	}
	return false
}

// Uses websocket message key to check whether the user should be sent past 24 hours of market news
func permissionToGetMarketNews(message string) bool {
	return message == "MarketNews"
}
