package caskdb

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"unicode/utf8"
)

type ValTypeID uint32
type UglyRune rune

const (
	IntTypeID ValTypeID = iota
	Int8TypeID
	Int16TypeID
	Int32TypeID
	Int64TypeID
	UIntTypeID
	UInt8TypeID
	UInt16TypeID
	UInt32TypeID
	UInt64TypeID
	Float32TypeID
	Float64TypeID
	StrTypeID
	RuneTypeID
	BoolTrueTypeID
	BoolFalseTypeID
	BytesTypeID
)

func (r *Record) EncodeValue(val interface{}) error {
	// encode value into bytes and record its type in header
	buf := new(bytes.Buffer)

	switch val.(type) {
	case string:
		_, err := buf.WriteString(val.(string))
		r.Value = buf.Bytes()
		r.Header.ValueType = StrTypeID
		return err
	case UglyRune:
		uglyRune, _ := val.(UglyRune)
		_, err := buf.WriteRune(rune(uglyRune))
		r.Value = buf.Bytes()
		r.Header.ValueType = RuneTypeID
		return err
	case int:
		encoder := gob.NewEncoder(buf)
		err := encoder.Encode(val.(int))
		r.Value = buf.Bytes()
		r.Header.ValueType = IntTypeID
		return err
	case int8:
		encoder := gob.NewEncoder(buf)
		err := encoder.Encode(val.(int8))
		r.Value = buf.Bytes()
		r.Header.ValueType = Int8TypeID
		return err
	case int16:
		encoder := gob.NewEncoder(buf)
		err := encoder.Encode(val.(int16))
		r.Value = buf.Bytes()
		r.Header.ValueType = Int16TypeID
		return err
	case int32:
		encoder := gob.NewEncoder(buf)
		err := encoder.Encode(val.(int32))
		r.Value = buf.Bytes()
		r.Header.ValueType = Int32TypeID
		return err
	case int64:
		encoder := gob.NewEncoder(buf)
		err := encoder.Encode(val.(int64))
		r.Value = buf.Bytes()
		r.Header.ValueType = Int64TypeID
		return err
	case uint:
		encoder := gob.NewEncoder(buf)
		err := encoder.Encode(val.(uint))
		r.Value = buf.Bytes()
		r.Header.ValueType = UIntTypeID
		return err
	case uint8:
		encoder := gob.NewEncoder(buf)
		err := encoder.Encode(val.(uint8))
		r.Value = buf.Bytes()
		r.Header.ValueType = UInt8TypeID
		return err
	case uint16:
		encoder := gob.NewEncoder(buf)
		err := encoder.Encode(val.(uint16))
		r.Value = buf.Bytes()
		r.Header.ValueType = UInt16TypeID
		return err
	case uint32:
		encoder := gob.NewEncoder(buf)
		err := encoder.Encode(val.(uint32))
		r.Value = buf.Bytes()
		r.Header.ValueType = UInt32TypeID
		return err
	case uint64:
		encoder := gob.NewEncoder(buf)
		err := encoder.Encode(val.(uint64))
		r.Value = buf.Bytes()
		r.Header.ValueType = UInt64TypeID
		return err
	case float32:
		encoder := gob.NewEncoder(buf)
		err := encoder.Encode(val.(float32))
		r.Value = buf.Bytes()
		r.Header.ValueType = Float32TypeID
		return err
	case float64:
		encoder := gob.NewEncoder(buf)
		err := encoder.Encode(val.(float64))
		r.Value = buf.Bytes()
		r.Header.ValueType = Float64TypeID
		return err
	case bool:
		input := val.(bool)
		if input {
			//true
			TrueID := BoolTrueTypeID
			err := binary.Write(buf, binary.LittleEndian, &TrueID)
			r.Value = buf.Bytes()
			r.Header.ValueType = BoolTrueTypeID
			return err
		}
		FalseID := BoolFalseTypeID
		err := binary.Write(buf, binary.LittleEndian, &FalseID)
		r.Value = buf.Bytes()
		r.Header.ValueType = BoolFalseTypeID
		return err
	case []byte:
		_, err := buf.Write(val.([]byte))
		r.Value = val.([]byte)
		r.Header.ValueType = BytesTypeID
		return err
	default:
		return ErrInvalidValue
	}
}

func (r *Record) DecodeValue() (interface{}, error) {
	// decode the value depending on the ValueType that was stored
	switch r.Header.ValueType {
	case IntTypeID:
		var value int
		decoder := gob.NewDecoder(bytes.NewReader(r.Value))
		err := decoder.Decode(&value)
		return value, err
	case Int8TypeID:
		var value int8
		decoder := gob.NewDecoder(bytes.NewReader(r.Value))
		err := decoder.Decode(&value)
		return value, err
	case Int16TypeID:
		var value int16
		decoder := gob.NewDecoder(bytes.NewReader(r.Value))
		err := decoder.Decode(&value)
		return value, err
	case Int64TypeID:
		var value int64
		decoder := gob.NewDecoder(bytes.NewReader(r.Value))
		err := decoder.Decode(&value)
		return value, err
	case UIntTypeID:
		var value uint
		decoder := gob.NewDecoder(bytes.NewReader(r.Value))
		err := decoder.Decode(&value)
		return value, err
	case UInt8TypeID:
		var value uint8
		decoder := gob.NewDecoder(bytes.NewReader(r.Value))
		err := decoder.Decode(&value)
		return value, err
	case UInt16TypeID:
		var value uint16
		decoder := gob.NewDecoder(bytes.NewReader(r.Value))
		err := decoder.Decode(&value)
		return value, err
	case UInt32TypeID:
		var value uint32
		decoder := gob.NewDecoder(bytes.NewReader(r.Value))
		err := decoder.Decode(&value)
		return value, err
	case UInt64TypeID:
		var value uint64
		decoder := gob.NewDecoder(bytes.NewReader(r.Value))
		err := decoder.Decode(&value)
		return value, err
	case Float32TypeID:
		var value float32
		decoder := gob.NewDecoder(bytes.NewReader(r.Value))
		err := decoder.Decode(&value)
		return value, err
	case Float64TypeID:
		var value float64
		decoder := gob.NewDecoder(bytes.NewReader(r.Value))
		err := decoder.Decode(&value)
		return value, err
	case StrTypeID:
		return string(r.Value), nil
	case RuneTypeID:
		result, _ := utf8.DecodeRune(r.Value)
		if result == utf8.RuneError {
			return "", ErrInvalidValue
		}
		return UglyRune(result), nil
	case Int32TypeID:
		var value int32
		decoder := gob.NewDecoder(bytes.NewReader(r.Value))
		err := decoder.Decode(&value)
		if err != nil {
			return "", ErrInvalidValue
		}
		return value, nil
	case BoolTrueTypeID:
		var BoolTrueTypeID ValTypeID
		err := binary.Read(bytes.NewBuffer(r.Value), binary.LittleEndian, &BoolTrueTypeID)
		return true, err
	case BoolFalseTypeID:
		var BoolFalseTypeID ValTypeID
		err := binary.Read(bytes.NewBuffer(r.Value), binary.LittleEndian, &BoolFalseTypeID)
		return false, err
	case BytesTypeID:
		// let user do marshalling/unmarshalling for structs?
		return r.Value, nil
	default:
		return nil, ErrInvalidValue
	}
}
