package unit

import "testing"

func TestValidConversion(t *testing.T) {
	convert := GetMeasurementConverter()
	_, err := convert(1, "km", "m")
	if err != nil {
		t.Fatalf("Expected err to be nil, got %+v", err)
	}
}

func TestInvalidConversionShouldReturnError(t *testing.T) {
	convert := GetMeasurementConverter()
	_, err := convert(1, "km", "hr")
	if err == nil {
		t.Fatalf("Expected err not to be nil")
	}
}
