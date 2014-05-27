package gocsv

import (
	"encoding/csv"
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
	if err != nil {
		return
	}
	return
}
