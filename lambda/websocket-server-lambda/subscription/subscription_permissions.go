package subscription

func permissionToGetRates(message string) bool {
	if message == "EURUSD" || message == "GBPUSD" || message == "USDJPY" || message == "AUDCAD" {
		return true
	}
	return false
}

func permissionToGetInference(message string) bool {
	if message == "EURUSDInference" || message == "GBPUSDInference" || message == "USDJPYInference" || message == "AUDCADInference" {
		return true
	}
	return false
}

func permissionToGetMarketNews(message string) bool {
	return message == "MarketNews"
}
