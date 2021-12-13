package db

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"time"
)

// NullInt wraps sql.NullInt32 and implements some interfaces to make life easier.
type NullInt sql.NullInt32

// MarshalJSON implements the JSON marshal interface for NullInt.
func (i *NullInt) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		i.Int32 = 0
		i.Valid = true
	}

	return json.Marshal(i.Int32)
}

// UnmarshalJSON implements the JSON unmarshal interface for NullInt.
func (i *NullInt) UnmarshalJSON(data []byte) error {
	var num *int32

	if err := json.Unmarshal(data, &num); err != nil {
		return err
	}

	if num != nil {
		i.Int32 = *num
		i.Valid = true
	} else {
		i.Valid = false
	}

	return nil
}

// Scan implements the Scanner interface for NullInt.
func (i *NullInt) Scan(value interface{}) error {
	var num sql.NullInt32
	if err := num.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*i = NullInt{num.Int32, false}
	} else {
		*i = NullInt{num.Int32, true}
	}

	return nil
}

// NullString wraps sql.NullString and implements some interfaces to make life easier.
type NullString sql.NullString

// MarshalJSON implements the JSON marshal interface for NullString.
func (s *NullString) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		s.String = ""
		s.Valid = true
	}

	return json.Marshal(s.String)
}

// UnmarshalJSON implements the JSON unmarshal interface for NullString.
func (s *NullString) UnmarshalJSON(data []byte) error {
	var str *string

	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	if str != nil {
		s.String = *str
		s.Valid = true
	} else {
		s.Valid = false
	}

	return nil
}

// Scan implements the Scanner interface for NullString.
func (s *NullString) Scan(value interface{}) error {
	var str sql.NullString
	if err := str.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*s = NullString{str.String, false}
	} else {
		*s = NullString{str.String, true}
	}

	return nil
}

// NullTime wraps sql.NullTime and implements some interfaces to make life easier.
type NullTime sql.NullTime

// MarshalJSON implements the JSON marshal interface for NullTime.
func (t *NullTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(t.Time)
}

// UnmarshalJSON implements the JSON unmarshal interface for NullTime.
func (t *NullTime) UnmarshalJSON(data []byte) error {
	var ts *time.Time

	if err := json.Unmarshal(data, &ts); err != nil {
		return err
	}

	if ts != nil {
		t.Time = *ts
		t.Valid = true
	} else {
		t.Valid = false
	}

	return nil
}

// Scan implements the Scanner interface for NullTime.
func (t *NullTime) Scan(value interface{}) error {
	var nt sql.NullTime
	if err := nt.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*t = NullTime{nt.Time, false}
	} else {
		*t = NullTime{nt.Time, true}
	}

	return nil
}
