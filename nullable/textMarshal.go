package nullable

import (
	"encoding"
	"errors"
	"fmt"
	"strconv"
)

func (n Nullable[T]) MarshalText() ([]byte, error) {
	if !n.HasValue {
		return []byte{}, nil
	}

	value := any(n.Value)
	txt, ok := value.(encoding.TextMarshaler)
	if ok {
		return txt.MarshalText()
	}

	switch any(n.Value).(type) {
	case float32, float64:
		return marshalTextFloat(n)
	case bool:
		return marshalTextBool(any(n).(Nullable[bool]))
	case int, int8, int16, int32, int64:
		return marshalTextInt(n)
	}

	var ref T
	return []byte{}, fmt.Errorf("type %T cannot be marshalled to text", ref)
}

func marshalTextInt[T any](f Nullable[T]) ([]byte, error) {
	if !f.HasValue {
		return []byte{}, nil
	}

	var value int64
	switch any(f.Value).(type) {
	case int:
		value = int64(any(f.Value).(int))
	case int8:
		value = int64(any(f.Value).(int8))
	case int16:
		value = int64(any(f.Value).(int16))
	case int32:
		value = int64(any(f.Value).(int32))
	case int64:
		value = any(f.Value).(int64)
	}

	return []byte(strconv.FormatInt(value, 10)), nil
}

func marshalTextFloat[T any](f Nullable[T]) ([]byte, error) {
	if !f.HasValue {
		return []byte{}, nil
	}

	var value float64
	switch any(f.Value).(type) {
	case float32:
		value = float64(any(f.Value).(float32))
	case float64:
		value = any(f.Value).(float64)
	}

	return []byte(strconv.FormatFloat(value, 'f', -1, 64)), nil
}

func marshalTextBool(b Nullable[bool]) ([]byte, error) {
	if !b.HasValue {
		return []byte{}, nil
	}
	if !b.Value {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

func (n *Nullable[T]) UnmarshalText(text []byte) error {
	value := any(&n.Value)
	txt, ok := value.(encoding.TextUnmarshaler)
	if ok {
		return txt.UnmarshalText(text)
	}

	switch any(n.Value).(type) {
	case bool:
		return unmarshalTextBool(text, any(n).(*Nullable[bool]))
	case float32, float64:
		return unmarshalTextFloat(text, n)
	case int, int8, int16, int32, int64:
		return unmarshalTextInt(text, n)
	}

	var ref T
	return fmt.Errorf("type %T unmarshal", ref)
}

func unmarshalTextBool(text []byte, b *Nullable[bool]) error {
	str := string(text)
	switch str {
	case "", "null":
		b.HasValue = false
		return nil
	case "true":
		b.Value = true
	case "false":
		b.Value = false
	default:
		return errors.New("null: invalid input for UnmarshalText:" + str)
	}
	b.HasValue = true
	return nil
}

func unmarshalTextFloat[T any](text []byte, f *Nullable[T]) error {
	str := string(text)
	if str == "" || str == "null" {
		f.HasValue = false
		return nil
	}

	var size int
	v := any(f.Value)
	switch v.(type) {
	case float32:
		size = 32
	case float64:
		size = 64
	}

	n, err := strconv.ParseFloat(string(text), size)
	if err != nil {
		return fmt.Errorf("null: couldn't unmarshal text: %w", err)
	}

	switch v.(type) {
	case float32:
		f.Value = any(float32(n)).(T)
	case float64:
		f.Value = any(n).(T)
	}

	f.HasValue = true
	return err
}

func unmarshalTextInt[T any](text []byte, f *Nullable[T]) error {
	str := string(text)
	if str == "" || str == "null" {
		f.HasValue = false
		return nil
	}

	var size int
	v := any(f.Value)
	switch v.(type) {
	case int8:
		size = 8
	case int16:
		size = 16
	case int32, int:
		size = 32
	case int64:
		size = 64
	}

	n, err := strconv.ParseInt(str, 10, size)
	if err != nil {
		return fmt.Errorf("null: couldn't unmarshal text: %w", err)
	}

	switch v.(type) {
	case int8:
		f.Value = any(int8(n)).(T)
	case int16:
		f.Value = any(int16(n)).(T)
	case int32:
		f.Value = any(int32(n)).(T)
	case int:
		f.Value = any(int(n)).(T)
	case int64:
		f.Value = any(n).(T)
	}

	f.HasValue = true
	return err
}
