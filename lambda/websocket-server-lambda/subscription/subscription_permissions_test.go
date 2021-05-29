package subscription

import "testing"

func TestPermissionToGetRates(t *testing.T) {
	// Test if valid symbol has permission
	hasPermission := permissionToGetRates("EURUSD")
	if hasPermission != true {
		t.Errorf("Invalid permissions, got: %t, expected true", hasPermission)
	}

	// Test if valid symbol with invalid casing fails
	hasPermission = permissionToGetRates("eURUSD")
	if hasPermission != false {
		t.Errorf("Invalid permissions, got: %t, expected false", hasPermission)
	}
}

func TestPermissionToGetInference(t *testing.T) {
	// Test if valid inference string has permission
	hasPermission := permissionToGetInference("AUDCADInference")
	if hasPermission != true {
		t.Errorf("Invalid permissions, got: %t, expected true", hasPermission)
	}

	// Test if valid inference string with improper casing has permission
	hasPermission = permissionToGetInference("aUDCADInference")
	if hasPermission != false {
		t.Errorf("Invalid permissions, got: %t, expected false", hasPermission)
	}
}

func TestPermissionToGetMarketNews(t *testing.T) {
	// Test if valid string has permission
	hasPermission := permissionToGetMarketNews("MarketNews")
	if hasPermission != true {
		t.Errorf("Invalid permissions, got: %t, expected true", hasPermission)
	}
}
