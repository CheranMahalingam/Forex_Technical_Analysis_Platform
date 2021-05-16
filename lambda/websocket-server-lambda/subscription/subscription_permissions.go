package subscription

func permissionToGetRates(message string) bool {
	if message == "EURUSD" || message == "GBPUSD" {
		return true
	}
	return false
}

func permissionToGetInference(message string) bool {
	if message == "EURUSDInference" || message == "GBPUSDInference" {
		return true
	}
	return false
}

func permissionToGetMarketNews(message string) bool {
	if message == "MarketNews" {
		return true
	}
	return false
}
