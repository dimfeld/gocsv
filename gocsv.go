package gocsv

import (
	"encoding/csv"
	"fmt"
	"io"
)

// A Reader may be constructed using either the New function, or directly by
// setting the Reader member to an existing csv.Reader and calling
// ReadHeader().
type Reader struct {
	*csv.Reader
	FieldIndex map[string]int
	Field      []string
}

type Writer struct {
	*csv.Writer
	AllowUnknown bool
	FieldIndex   map[string]int
	Field        []string
}

type Record map[string]string

func (c *Reader) makeRecord(values []string) Record {
	record := make(Record, len(values))
	for i, field := range c.Field {
		if i > len(values) {
			break
		}
		record[field] = values[i]
	}
	return record
}

func (c *Reader) Read() (Record, error) {
	fields, err := c.Reader.Read()
	if err != nil {
		return nil, err
	}
	return c.makeRecord(fields), nil
}

func (c *Reader) ReadAll() ([]Record, error) {
	allValues, err := c.Reader.ReadAll()
	if err != nil {
		return nil, err
	}
	records := make([]Record, len(allValues))
	for i, values := range allValues {
		records[i] = c.makeRecord(values)
	}
	return records, nil
}

func (c *Reader) ReadHeader() error {
	header, err := c.Reader.Read()
	if err != nil {
		return err
	}

	c.Field = header
	c.FieldIndex = make(map[string]int, len(c.Field))
	for i, field := range header {
		c.FieldIndex[field] = i
	}

	return nil
}

// NewReader constructs a Reader object, reading the first line from the
// supplied io.Reader and interpreting it as a header line.
func NewReader(r io.Reader) (c *Reader, err error) {
	c = &Reader{Reader: csv.NewReader(r)}
	err = c.ReadHeader()
	return
}

func NewWriter(w io.Writer, fields []string) *Writer {
	writer := &Writer{
		Writer:     csv.NewWriter(w),
		Field:      fields,
		FieldIndex: map[string]int{},
	}

	for i, f := range fields {
		writer.FieldIndex[f] = i
	}

	return writer
}

func (w *Writer) Write(values Record) error {
	record := make([]string, len(w.Field))
	for key, value := range values {
		pos, ok := w.FieldIndex[key]
		if !ok {
			if !w.AllowUnknown {
				return fmt.Errorf("Unknown field %s", key)
			}
			continue
		}

		record[pos] = value
	}

	return w.Writer.Write(record)
}

func (w *Writer) WriteHeader() error {
	return w.Writer.Write(w.Field)
}
