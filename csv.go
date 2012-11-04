package objcsv

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

func csvError(msg string, args ...interface{}) error {
	return fmt.Errorf("CSV: "+msg, args...)
}

// ReadCSV reads from reader r into a pointer to a slice of the destination struct type.
func ReadCSV(r io.Reader, ptrToSlice interface{}) (error) {
  pv := reflect.ValueOf(ptrToSlice)
	if pv.IsNil() {
		return csvError("expected %v to have a non-nil value", ptrToSlice)
	}
	if pv.Kind() != reflect.Ptr || pv.Elem().Kind() != reflect.Slice {
		return csvError("expected a pointer to a slice, got %v instead", pv.Kind())
	}

	t := pv.Elem().Type().Elem()
	cr := csv.NewReader(r)
	didHeader := false
	for {
		record, err := cr.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if !didHeader {
			didHeader = true
			continue
		}

		rv, err := recordToValue(record, t)
		if err != nil {
			return err
		}
    pv.Elem().Set(reflect.Append(pv.Elem(), rv.Elem()))
	}
	return nil
}

func recordToValue(ss []string, t reflect.Type) (reflect.Value, error) {
	v := reflect.New(t)
	if len(ss) != t.NumField() {
		return v, csvError("recordToValue: can't decode CSV record to struct (field mismatch)")
	}
	for i, s := range ss {
    if s == "" {
      continue
    }
		f := t.Field(i)
		ft := f.Type.Kind()
		fv := v.Elem().Field(i)
		typeError := func(msg string, e error) error {
      return csvError("%v: type %v, field %v: %v", msg, ft, f.Name, e)
		}

		switch ft {
		case reflect.String:
			fv.SetString(s)
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			theInt, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return v, typeError("int", err)
			}
			fv.SetInt(theInt)
		case reflect.Float64:
			theFloat, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return v, typeError("float64", err)
			}
			fv.SetFloat(theFloat)
		case reflect.Bool:
			theBool, err := strconv.ParseBool(s)
			if err != nil {
				return v, typeError("bool", err)
			}
			fv.SetBool(theBool)
		default:
			return v, csvError("recordToValue: unsupported type %v in field %v", ft, f.Name)
		}

	}
	return v, nil
}

// WriteCSV writes to writer w from a slice of data records.
func WriteCSV(w io.Writer, slice interface{}) error {
	// ensure that d is a slice containing at least some data.
	v, err := reflectSlice(slice)
	if err != nil {
		return err
	}
	if v.Len() == 0 {
		return csvError("data is empty")
	}

	cw := csv.NewWriter(w)
	if err := writeHeader(cw, v.Index(0).Interface()); err != nil {
		return err
	}
	for i := 0; i < v.Len(); i++ {
		e := v.Index(i)
		if err := writeRow(cw, e.Interface()); err != nil {
			return err
		}
	}
	cw.Flush()
	return nil
}

func writeHeader(w *csv.Writer, d interface{}) error {
	v := reflect.ValueOf(d)
	header := make([]string, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		header[i] = v.Type().Field(i).Name
	}
	if err := w.Write(header); err != nil {
		return err
	}
	return nil
}

func writeRow(w *csv.Writer, d interface{}) error {
	v := reflect.ValueOf(d)
	row := make([]string, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		switch t := v.Field(i).Interface(); t.(type) {
		case int:
			row[i] = fmt.Sprintf("%d", t)
		case float64:
			row[i] = fmt.Sprintf("%.2f", t)
		case string:
			row[i] = fmt.Sprintf("%s", t)
		default:
			row[i] = fmt.Sprintf("%v", t)
		}
	}
	if err := w.Write(row); err != nil {
		return err
	}
	return nil
}

func reflectSlice(i interface{}) (reflect.Value, error) {
	v := reflect.ValueOf(i)
	if v.IsNil() {
		return v, csvError("expected %v to have a non-nil value", i)
	}
	if v.Kind() != reflect.Slice {
		return v, csvError("expected a slice, got %v instead", v.Kind())
	}
	return v, nil
}
