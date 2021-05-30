package finance

import (
	"testing"
	"time"
)

func TestCreateExpression(t *testing.T) {
	currentDate := time.Now()
	expression, err := createExpression(currentDate.Format("2006-01-02"))
	if err != nil {
		t.Errorf("Function failed with unexpected error: %e", err)
	}

	if *expression.Names()["#0"] != "Date" {
		t.Errorf("Filter expression key was invalid, got: %s, expected: Date", *expression.Names()["#0"])
	}
	if *expression.Values()[":0"].S != currentDate.Format("2006-01-02") {
		t.Errorf("Filter expression key condition Date was invalid, got: %s, expected: %s", *expression.Values()[":0"].S, currentDate.Format("2006-01-02"))
	}
}
