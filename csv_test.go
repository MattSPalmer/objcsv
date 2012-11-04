package objcsv

import (
	"bytes"
	"encoding/csv"
	"io"
	"reflect"
	"testing"
)

type dumbStruct struct {
	FName string
	LName string
	Age   int
	Worth float64
}

func TestWriteHeader(t *testing.T) {
	d := dumbStruct{"Hello", "World", 42, 10.01}
	buf := new(bytes.Buffer)
	cw := csv.NewWriter(buf)
	err := writeHeader(cw, d)

	if err != nil {
		t.Fatalf("Unexpected error from writeHeader: %v", err)
	}

	cw.Flush()
	actual := buf.String()
	expected := "FName,LName,Age,Worth\n"
	if actual != expected {
		t.Errorf("Header is wrong, expected %q, got %q", expected, actual)
	}
}

func TestWriteRow(t *testing.T) {
	d := dumbStruct{"Matt", "Palmer", 23, 40000.00}
	buf := new(bytes.Buffer)
	cw := csv.NewWriter(buf)
	err := writeRow(cw, d)

	if err != nil {
		t.Fatalf("Unexpected error from writeRow: %v", err)
	}

	cw.Flush()
	actual := buf.String()
	expected := "Matt,Palmer,23,40000.00\n"
	if actual != expected {
		t.Errorf("Row is wrong, expected %q, got %q", expected, actual)
	}
}

func TestReadCSV(t *testing.T) {
	data := `FName,LName,Age,Worth
Matt,Palmer,23,40000.00
Robert,Sesek,23,60000.00
Barack,Obama,51,80000000.00
Mitt,Romney,59,9950000000.00`

	buf := new(bytes.Buffer)

	_, err := buf.WriteString(data)
	if err != nil {
		t.Fatalf("Error in testing data: %v", err)
	}

	slice := make([]dumbStruct, 0)
	err = ReadCSV(buf, &slice)
	if err != nil {
		t.Fatalf("Unexpected error from ReadCSV: %v", err)
	}

	expected := []dumbStruct{
		{"Matt", "Palmer", 23, 40000.00},
		{"Robert", "Sesek", 23, 60000.00},
		{"Barack", "Obama", 51, 80000000.00},
		{"Mitt", "Romney", 59, 9950000000.00},
	}

	if !reflect.DeepEqual(slice, expected) {
		t.Errorf("Expected %v but got %v", expected, slice)
	}

}

func TestWriteCSV(t *testing.T) {
	d := []dumbStruct{
		{"Matt", "Palmer", 23, 40000.00},
		{"Robert", "Sesek", 23, 60000.00},
		{"Barack", "Obama", 51, 80000000.00},
		{"Mitt", "Romney", 59, 9950000000.00},
	}
	buf := new(bytes.Buffer)

	if err := WriteCSV(buf, d); err != nil {
		t.Fatalf("Unexpected error from WriteCSV: %v", err)
	}

	expected := []string{
		"FName,LName,Age,Worth\n",
		"Matt,Palmer,23,40000.00\n",
		"Robert,Sesek,23,60000.00\n",
		"Barack,Obama,51,80000000.00\n",
		"Mitt,Romney,59,9950000000.00\n",
	}
	for i := 0; i < 5; i++ {
		actual, err := buf.ReadString('\n')
		if err != nil {
			t.Fatalf("Unexpected output at line %d: %v", i, err)
		}
		if actual != expected[i] {
			t.Errorf("Row %d is wrong, expected %q, got %q", i, expected[i], actual)
		}
	}
	s, err := buf.ReadString('\n')
	if s != "" || err != io.EOF {
		t.Errorf("Got unexpected data %q where expected EOF", s)
	}
}
