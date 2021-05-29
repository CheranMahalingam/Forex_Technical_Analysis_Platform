package subscription

import (
	"testing"
)

func TestKeyConditionDayData(t *testing.T) {
	// Test key condition for data on Thursday
	currentDay := 4
	currentHour := 22
	// No way of accessing internal values so we check the length of key conditions
	keyCond := keyConditionDayData(currentDay, currentHour)
	if len(*keyCond) != 2 {
		t.Errorf("Invalid key conditions, got length: %d, expected length 2", len(*keyCond))
	}

	// Test key condition for data on Friday after market closes
	currentDay = 5
	currentHour = 23
	keyCond = keyConditionDayData(currentDay, currentHour)
	if len(*keyCond) != 2 {
		t.Errorf("Invalid key conditions, got length: %d, expected length 2", len(*keyCond))
	}

	// Test key condition for data on Saturday after market closing
	currentDay = 6
	keyCond = keyConditionDayData(currentDay, currentHour)
	if len(*keyCond) != 2 {
		t.Errorf("Invalid key conditions, got length: %d, expected length 2", len(*keyCond))
	}

	// Test key condition for data on Sunday before market opens
	currentDay = 0
	currentHour = 22
	keyCond = keyConditionDayData(currentDay, currentHour)
	if len(*keyCond) != 2 {
		t.Errorf("Invalid key conditions, got length: %d, expected length 2", len(*keyCond))
	}

	// Test key condition for data on Sunday after market opens
	currentHour = 23
	keyCond = keyConditionDayData(currentDay, currentHour)
	if len(*keyCond) != 3 {
		t.Errorf("Invalid key conditions, got length: %d, expected length 3", len(*keyCond))
	}

	// Test key condition for data on Monday
	currentDay = 1
	keyCond = keyConditionDayData(currentDay, currentHour)
	if len(*keyCond) != 3 {
		t.Errorf("Invalid key conditions, got length: %d, expected length 3", len(*keyCond))
	}
}
