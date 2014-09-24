package gocsv

import (
	"bytes"
	"reflect"
	"testing"
)

var testFile = `a,b ,c 
1,2 ,3
4,5,6 
`

var expectedField = []string{"a", "b ", "c "}

var expectedRecords = []Record{
	Record{"a": "1", "b ": "2 ", "c ": "3"},
	Record{"a": "4", "b ": "5", "c ": "6 "},
}

var errorFile = `a,b,c
3,4,5,6
3,4,5`

func getExpectedData(trimTrailing bool) (fields []string, records []Record) {
	if !trimTrailing {
		return expectedField, expectedRecords
	}

	fields = make([]string, len(expectedField))
	for i, field := range expectedField {
		fields[i] = string(bytes.TrimRight([]byte(field), " "))
	}

	records = make([]Record, len(expectedRecords))
	for i, expectedRecord := range expectedRecords {
		record := make(Record)

		for field, value := range expectedRecord {
			trimmedField := string(bytes.TrimRight([]byte(field), " "))
			trimmedValue := string(bytes.TrimRight([]byte(value), " "))
			record[trimmedField] = trimmedValue
		}
		records[i] = record
	}

	return
}

func TestRead(t *testing.T) {
	testRead(t, false)
	testRead(t, true)
}

func testRead(t *testing.T, trimTrailing bool) {
	t.Log("Testing with trimTrailing", trimTrailing)
	reader := bytes.NewBufferString(testFile)
	csv, err := NewTrimmingReader(reader, false, trimTrailing)
	if err != nil {
		t.Fatalf("Failed to create reader: %s", err)
	}

	expectedField, expectedRecords := getExpectedData(trimTrailing)

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
	testReadAll(t, false)
	testReadAll(t, true)
}

func testReadAll(t *testing.T, trimTrailing bool) {
	t.Log("Testing with trimTrailing", trimTrailing)
	reader := bytes.NewBufferString(testFile)
	csv, err := NewTrimmingReader(reader, false, trimTrailing)
	if err != nil {
		t.Fatalf("Failed to create reader: %s", err)
	}

	expectedField, expectedRecords := getExpectedData(trimTrailing)

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

func TestWrite(t *testing.T) {
	buf := &bytes.Buffer{}
	w := NewWriter(buf, []string{"a", "b", "c"})

	err := w.Write(Record{"a": "1", "b": "2", "c": "3"})
	w.Flush()
	expected := "1,2,3\n"
	if err != nil {
		t.Error("Unexpected error", err)
	}
	if string(buf.Bytes()) != expected {
		t.Errorf("Expected %s, saw %s", expected, string(buf.Bytes()))
	}
	buf.Reset()

	err = w.Write(Record{"a": "1", "b": "2", "c": "3"})
	if err != nil {
		t.Error("Unexpected error", err)
	}
	err = w.Write(Record{"a": "2", "b": "3", "c": "4"})
	w.Flush()
	expected = "1,2,3\n2,3,4\n"
	if err != nil {
		t.Error("Unexpected error", err)
	}
	if string(buf.Bytes()) != expected {
		t.Errorf("Expected %s, saw %s", expected, string(buf.Bytes()))
	}
	buf.Reset()

	err = w.WriteHeader()
	w.Flush()
	expected = "a,b,c\n"
	if err != nil {
		t.Error("Unexpected error", err)
	}
	if string(buf.Bytes()) != expected {
		t.Errorf("Expected %s, saw %s", expected, string(buf.Bytes()))
	}
	buf.Reset()

	err = w.Write(Record{"a": "2", "b": "3", "d": "4"})
	if err == nil {
		t.Error("Expected error on unknown field")
	}
	w.Flush()
	if buf.Len() != 0 {
		t.Error("Write with unknown field still wrote data:",
			string(buf.Bytes()))
	}
	buf.Reset()

	w.AllowUnknown = true
	err = w.Write(Record{"a": "1", "b": "2", "d": "3"})
	if err != nil {
		t.Error("Unexpected error", err)
	}
	err = w.Write(Record{"a": "2", "b": "3", "c": "4"})
	w.Flush()
	expected = "1,2,\"\"\n2,3,4\n"
	if err != nil {
		t.Error("Unexpected error", err)
	}
	if string(buf.Bytes()) != expected {
		t.Errorf("Expected %s, saw %s", expected, string(buf.Bytes()))
	}
}
