package db

import (
	"testing"
	"time"
)

// TestMarshalNullInt makes sure that marshalling a valid
// NullInt works as expected.
func TestMarshalNullInt(t *testing.T) {
	// Given
	test := NullInt{
		Int32: 5,
		Valid: true,
	}

	// When
	result, err := test.MarshalJSON()

	// Then
	if err != nil {
		t.Errorf("error marshalling valid null int: %s\n", err.Error())
	}

	if string(result) != "5" {
		t.Fatalf("result does not match expected: got: %d, expected: 5\n", result)
	}
}

// TestMarshalNullInt_Null ensures that marshalling a NullInt
// without a valid value returns "0".
func TestMarshalNullInt_Null(t *testing.T) {
	// Given
	test := NullInt{
		Valid: false,
	}

	// When
	result, err := test.MarshalJSON()

	// Then
	if err != nil {
		t.Errorf("error marshalling a null NullInt: %s\n", err.Error())
	}

	if string(result) != "0" {
		t.Fatalf("result does not match expected: got: %d, expected: 0\n", result)
	}
}

// TestUnmarshalNullInt ensures that unmarshalling a NullInt from a valid
// input works as expected.
func TestUnmarshalNullInt(t *testing.T) {
	// Given
	var result NullInt
	test := "5"

	// When
	err := result.UnmarshalJSON([]byte(test))

	// Then
	if err != nil {
		t.Errorf("error unmarshalling NullInt: %s\n,", err.Error())
	}

	if !result.Valid {
		t.Fatal("unmarshalled NullInt is not valid")
	}

	if result.Int32 != 5 {
		t.Fatalf("result does not match expected: got %d, expected: 5\n", result.Int32)
	}
}

// TestUnmarshalNullInt_Null ensures that unmarshalling a NullInt from a valid
// null input works as expected.
func TestUnmarshalNullInt_Null(t *testing.T) {
	// Given
	var result NullInt

	// When
	err := result.UnmarshalJSON([]byte("null"))

	// Then
	if err != nil {
		t.Errorf("error unmarshalling NullInt: %s\n,", err.Error())
	}

	if result.Valid {
		t.Fatal("unmarshalled NullInt from nil value is valid")
	}

	if result.Int32 != 0 {
		t.Fatalf("result does not match expected: got %d, expected: 0\n", result.Int32)
	}
}

// TestScanNullInt ensures that scanning a valid NullInt from a database
// works as expected.
func TestScanNullInt(t *testing.T) {
	// Given
	var result NullInt

	// When
	err := result.Scan(5)

	// Then
	if err != nil {
		t.Errorf("error calling Scan on NullInt: %s\n", err.Error())
	}

	if !result.Valid {
		t.Fatal("scanned NullInt is not valid")
	}

	if result.Int32 != 5 {
		t.Fatalf("result does not match expected: got %d, expected: 5\n", result.Int32)
	}
}

// TestScanNullInt_Null ensures that scanning an invalid NullInt from a
// database works as expected.
func TestScanNullInt_Null(t *testing.T) {
	// Given
	var result NullInt

	// When
	err := result.Scan(nil)

	// Then
	if err != nil {
		t.Errorf("error calling Scan on NullInt: %s\n", err.Error())
	}

	if result.Valid {
		t.Fatal("scanned NullInt is valid, but shouldn't be")
	}

	if result.Int32 != 0 {
		t.Fatalf("result does not match expected: got %d, expected: 0\n", result.Int32)
	}
}

// TestMarshalNullString makes sure that marshalling a valid
// NullString works as expected.
func TestMarshalNullString(t *testing.T) {
	// Given
	test := NullString{
		String: "test",
		Valid:  true,
	}

	// When
	result, err := test.MarshalJSON()

	// Then
	if err != nil {
		t.Errorf("error marshalling valid NullString: %s\n", err.Error())
	}

	if string(result) != "\"test\"" {
		t.Fatalf("result does not match expected: got: %s, expected: \"test\"\n", result)
	}
}

// TestMarshalNullString_Null ensures that marshalling a NullString
// without a valid value works as expected.
func TestMarshalNullString_Null(t *testing.T) {
	// Given
	test := NullString{
		Valid: false,
	}

	// When
	result, err := test.MarshalJSON()

	// Then
	if err != nil {
		t.Errorf("error marshalling a null NullString: %s\n", err.Error())
	}

	if string(result) != "\"\"" {
		t.Fatalf("result does not match expected: got: %s, expected: \"\"\n", result)
	}
}

// TestUnmarshalNullString ensures that unmarshalling a NullString from a valid
// input works as expected.
func TestUnmarshalNullString(t *testing.T) {
	// Given
	var result NullString
	test := "\"test\""

	// When
	err := result.UnmarshalJSON([]byte(test))

	// Then
	if err != nil {
		t.Errorf("error unmarshalling NullString: %s\n,", err.Error())
	}

	if !result.Valid {
		t.Fatal("unmarshalled NullString is not valid")
	}

	if result.String != "test" {
		t.Fatalf("result does not match expected: got %s, expected: test\n", result.String)
	}
}

// TestUnmarshalNullString_Null ensures that unmarshalling a NullString from a valid
// null input works as expected.
func TestUnmarshalNullString_Null(t *testing.T) {
	// Given
	var result NullString

	// When
	err := result.UnmarshalJSON([]byte("null"))

	// Then
	if err != nil {
		t.Errorf("error unmarshalling NullString: %s\n,", err.Error())
	}

	if result.Valid {
		t.Fatal("unmarshalled NullString from nil value is valid")
	}

	if result.String != "" {
		t.Fatalf("result does not match expected: got %s, expected: \"\"\n", result.String)
	}
}

// TestScanNullString ensures that scanning a valid string from a database
// works as expected.
func TestScanNullString(t *testing.T) {
	// Given
	var result NullString

	// When
	err := result.Scan("test")

	// Then
	if err != nil {
		t.Errorf("error calling Scan on NullString: %s\n", err.Error())
	}

	if !result.Valid {
		t.Fatal("scanned NullString is not valid")
	}

	if result.String != "test" {
		t.Fatalf("result does not match expected: got %s, expected: test\n", result.String)
	}
}

// TestScanNullString_Null ensures that scanning a null string from a
// database works as expected.
func TestScanNullString_Null(t *testing.T) {
	// Given
	var result NullString

	// When
	err := result.Scan(nil)

	// Then
	if err != nil {
		t.Errorf("error calling Scan on NullString: %s\n", err.Error())
	}

	if result.Valid {
		t.Fatal("scanned NullString is valid, but shouldn't be")
	}

	if result.String != "" {
		t.Fatalf("result does not match expected: got %s, expected: \"\"\n", result.String)
	}
}

// TestMarshalNullTime makes sure that marshalling a valid
// NullTime works as expected.
func TestMarshalNullTime(t *testing.T) {
	// Given
	now := time.Now()
	test := NullTime{
		Time:  now,
		Valid: true,
	}
	marshalledNow, _ := now.MarshalJSON()

	// When
	result, err := test.MarshalJSON()

	// Then
	if err != nil {
		t.Errorf("error marshalling valid NullTime: %s\n", err.Error())
	}

	if string(result) != string(marshalledNow) {
		t.Fatalf("result does not match expected: got: %s, expected: %s\n", result, marshalledNow)
	}
}

// TestMarshalNullTime_Null ensures that marshalling a NullTime
// without a valid value works as expected.
func TestMarshalNullTime_Null(t *testing.T) {
	// Given
	test := NullTime{
		Valid: false,
	}

	// When
	result, err := test.MarshalJSON()

	// Then
	if err != nil {
		t.Errorf("error marshalling a null NullTime: %s\n", err.Error())
	}

	if string(result) != "null" {
		t.Fatalf("result does not match expected: got: %s, expected: null\n", result)
	}
}

// TestUnmarshalNullTime ensures that unmarshalling a NullTime from a valid
// input works as expected.
func TestUnmarshalNullTime(t *testing.T) {
	// Given
	now := time.Now()
	var result NullTime
	marshalledNow, _ := now.MarshalJSON()

	// When
	err := result.UnmarshalJSON(marshalledNow)

	// Then
	if err != nil {
		t.Errorf("error unmarshalling NullTime: %s\n,", err.Error())
	}

	if !result.Valid {
		t.Fatal("unmarshalled NullTime is not valid")
	}

	if !result.Time.Equal(now) {
		t.Fatalf("result does not match expected: got %s, expected: %s\n", result.Time, now)
	}
}

// TestUnmarshalNullTime_Null ensures that unmarshalling a NullTime from a valid
// null input works as expected.
func TestUnmarshalNullTime_Null(t *testing.T) {
	// Given
	var result NullTime

	// When
	err := result.UnmarshalJSON([]byte("null"))

	// Then
	if err != nil {
		t.Errorf("error unmarshalling NullTime: %s\n,", err.Error())
	}

	if result.Valid {
		t.Fatal("unmarshalled NullTime from nil value is valid")
	}

	if !result.Time.IsZero() {
		t.Fatal("unmarshalled NullTime is not zero value")
	}
}

// TestScanNullTime ensures that scanning a valid time from a database
// works as expected.
func TestScanNullTime(t *testing.T) {
	// Given
	now := time.Now()
	var result NullTime

	// When
	err := result.Scan(now)

	// Then
	if err != nil {
		t.Errorf("error calling Scan on NullTime: %s\n", err.Error())
	}

	if !result.Valid {
		t.Fatal("scanned NullTime is not valid")
	}

	if !result.Time.Equal(now) {
		t.Fatalf("result does not match expected: got %s, expected: %s\n", result.Time, now)
	}
}

// TestScanNullTime_Null ensures that scanning a null time from a
// database works as expected.
func TestScanNullTime_Null(t *testing.T) {
	// Given
	var result NullTime

	// When
	err := result.Scan(nil)

	// Then
	if err != nil {
		t.Errorf("error calling Scan on NullTime: %s\n", err.Error())
	}

	if result.Valid {
		t.Fatal("scanned NullTime is valid, but shouldn't be")
	}

	if !result.Time.IsZero() {
		t.Fatal("scanned NullTime is not zero value")
	}
}
