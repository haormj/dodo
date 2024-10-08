package util

import (
	"reflect"

	"github.com/spf13/cast"
)

type StructField struct {
	Field      reflect.StructField
	FieldName  string
	FieldValue reflect.Value
	IsExist    bool
}

func (sf *StructField) Bool() (bool, bool) {
	if sf.IsExist {
		return cast.ToBool(sf.FieldValue.Interface()), sf.IsExist
	} else {
		return false, sf.IsExist
	}
}

func (sf *StructField) BoolE() (bool, bool, error) {
	if sf.IsExist {
		b, err := cast.ToBoolE(sf.FieldValue.Interface())
		return b, sf.IsExist, err
	} else {
		return false, sf.IsExist, nil
	}
}

func (sf *StructField) Int() (int, bool) {
	if sf.IsExist {
		return cast.ToInt(sf.FieldValue.Interface()), sf.IsExist
	} else {
		return 0, sf.IsExist
	}
}

func (sf *StructField) IntE() (int, bool, error) {
	if sf.IsExist {
		i, err := cast.ToIntE(sf.FieldValue.Interface())
		return i, sf.IsExist, err
	} else {
		return 0, sf.IsExist, nil
	}
}

func (sf *StructField) Float64() (float64, bool) {
	if sf.IsExist {
		return cast.ToFloat64(sf.FieldValue.Interface()), sf.IsExist
	} else {
		return 0, sf.IsExist
	}
}

func (sf *StructField) Float64E() (float64, bool, error) {
	if sf.IsExist {
		f, err := cast.ToFloat64E(sf.FieldValue.Interface())
		return f, sf.IsExist, err
	} else {
		return 0, sf.IsExist, nil
	}
}

func (sf *StructField) String() (string, bool) {
	if sf.IsExist {
		return cast.ToString(sf.FieldValue.Interface()), sf.IsExist
	} else {
		return "", sf.IsExist
	}
}

func (sf *StructField) StringE() (string, bool, error) {
	if sf.IsExist {
		str, err := cast.ToStringE(sf.FieldValue.Interface())
		return str, sf.IsExist, err
	} else {
		return "", sf.IsExist, nil
	}
}

func (sf *StructField) Interface() (interface{}, bool) {
	return sf.FieldValue.Interface(), sf.IsExist
}

// SetStructField v is ptr to struct
func SetStructField(v interface{}, fieldValue interface{}, fieldNames ...string) bool {
	sf := GetStructField(v, fieldNames...)
	if !sf.IsExist {
		return false
	}
	if !sf.FieldValue.CanSet() {
		return false
	}
	fieldValueType := reflect.TypeOf(fieldValue)
	fieldValueValue := reflect.ValueOf(fieldValue)
	if sf.Field.Type != fieldValueType {
		if sf.FieldValue.Kind() == reflect.Ptr &&
			sf.Field.Type.Elem() == fieldValueType {
			v := reflect.New(fieldValueType)
			v.Elem().Set(fieldValueValue)
			fieldValueValue = v
		} else {
			return false
		}
	}
	sf.FieldValue.Set(fieldValueValue)
	return true
}

// GetStructField get struct field
func GetStructField(v interface{}, fieldNames ...string) *StructField {
	var sf *StructField
	fieldValue := reflect.ValueOf(v)
	for _, fieldName := range fieldNames {
		sf = getStructField(fieldValue, fieldName)
		if !sf.IsExist {
			break
		}
		fieldValue = sf.FieldValue
	}
	return sf
}

func getStructField(v reflect.Value, fieldName string) *StructField {
	sf := &StructField{}
	vValue := v
	for {
		if vValue.Kind() != reflect.Ptr {
			break
		}
		vValue = reflect.Indirect(vValue)
	}
	if vValue.Kind() != reflect.Struct {
		return sf
	}
	vType := vValue.Type()
	for i := 0; i < vType.NumField(); i++ {
		field := vType.Field(i)
		if field.Name == fieldName {
			fValue := vValue.Field(i)
			sf := &StructField{
				Field:      field,
				FieldName:  fieldName,
				FieldValue: fValue,
				IsExist:    true,
			}
			return sf
		}
	}
	return sf
}
