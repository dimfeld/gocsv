package gocsv

import (
	"bytes"
	"reflect"
	"testing"
)

var testFile = `a,b,c
1,2,3
4,5,6
`

var expectedField = []string{"a", "b", "c"}

var expectedRecords = []Record{
	Record{"a": "1", "b": "2", "c": "3"},
	Record{"a": "4", "b": "5", "c": "6"},
}

var errorFile = `a,b,c
3,4,5,6
3,4,5`

func TestRead(t *testing.T) {
	reader := bytes.NewBufferString(testFile)
	csv, err := NewReader(reader)
	if err != nil {
		t.Fatalf("Failed to create reader: %s", err)
	}

	if !reflect.DeepEqual(csv.Field, expectedField) {
		t.Fatalf("Expected header %v, saw %v", expectedField, csv.Field)
	}

	for i, expected := range expectedRecords {
		record, err := csv.Read()
		if err != nil {
			t.Errorf("Error reading record %d: %s", i, err)
		}

		if !reflect.DeepEqual(expected, record) {
			t.Errorf("Record %d expected %v, saw %v",
				i, expected, record)
		}
	}
}

func TestReadAll(t *testing.T) {
	reader := bytes.NewBufferString(testFile)
	csv, err := NewReader(reader)
	if err != nil {
		t.Fatalf("Failed to create reader: %s", err)
	}

	if !reflect.DeepEqual(csv.Field, expectedField) {
		t.Fatalf("Expected header %v, saw %v", expectedField, csv.Field)
	}

	records, err := csv.ReadAll()
	if err != nil {
		t.Errorf("Error reading records: %s", err)
	}
	for i, expected := range expectedRecords {
		if !reflect.DeepEqual(expected, records[i]) {
			t.Errorf("Record %d expected %v, saw %v",
				i, expected, records[i])
		}
	}
}

func TestError(t *testing.T) {
	reader := bytes.NewBufferString("")
	_, err := NewReader(reader)
	if err == nil {
		t.Errorf("Expected error with empty content")
	}

	reader = bytes.NewBufferString(errorFile)
	c, _ := NewReader(reader)
	_, err = c.Read()
	if err == nil {
		t.Errorf("Expected error from Read with too-long line")
	}

	reader = bytes.NewBufferString(errorFile)
	c, _ = NewReader(reader)
	_, err = c.ReadAll()
	if err == nil {
		t.Errorf("Expected error from ReadAll with too-long line")
	}
}
