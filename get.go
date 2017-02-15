package gostruct

import (
	"github.com/Sirupsen/logrus"
	"reflect"
	"strconv"
	"strings"
	"fmt"
)

func GetField(s interface{}, fieldName string) (interface{}, error) {
	field, err := getNestField(s, fieldName, fieldName)
	if err != nil {
		return nil, err
	}
	return field.Interface(), nil
}

func getNestField(s interface{}, fullName, fieldName string) (reflect.Value, error) {

	//Make sure types are strcut or ptr
	t := reflect.TypeOf(s).Kind()

	if t == reflect.Struct || t == reflect.Ptr {

	} else {
		return reflect.Value{}, fmt.Errorf("Struct must be struct interface")
	}

	val := getReflectValue(s)

	if i := strings.Index(fieldName, "."); i > -1 {
		currFieldName := fieldName[0:i]
		fname := getFieldName(currFieldName)
		field := val.FieldByName(fname)
		if !field.IsValid() {
			//We should ignore the error since there might be empty field.
			return reflect.Value{}, fmt.Errorf("No such field: %s in obj", fname)
		}
		nextFieldName := fieldName[i+1 : len(fieldName)]
		fieldValue := reflect.ValueOf(field.Interface())
		switch fieldValue.Kind() {
		case reflect.Slice:
			index, err := getFieldSliceIndex(currFieldName)
			if err != nil {
				return reflect.Value{}, err
			}

			if index != -1 {
				field = fieldValue.Index(index)
			}

			return getNestField(field.Interface(), fullName, nextFieldName)

		case reflect.Map:
			key, err := getFieldMapKey(currFieldName)

			if err != nil {
				return reflect.Value{}, err
			}
			if key != "" {
				field = field.MapIndex(reflect.ValueOf(key))
			}
			return getNestField(field.Interface(), fullName, nextFieldName)

		}
		return getNestField(field.Interface(), fullName, nextFieldName)

	}

	if !val.IsValid() {
		return reflect.Value{}, fmt.Errorf("Nil pointer: %s in obj", fullName)
	}

	field := val.FieldByName(getFieldName(fieldName))

	if !field.IsValid() {
		return field, fmt.Errorf("No such field: %s in obj", fieldName)
	}

	switch field.Kind() {
	case reflect.Slice:
		index, err := getFieldSliceIndex(fieldName)
		if err != nil {
			return reflect.Value{}, err
		}
		logrus.Debugf("Field slice index %d", index)

		if index != -1 {
			field = field.Index(index)
		}
	case reflect.Map:
		key, err := getFieldMapKey(fieldName)
		logrus.Debugf("Field map key %s", key)

		if err != nil {
			return reflect.Value{}, err
		}
		if key != "" {
			field = field.MapIndex(reflect.ValueOf(key))
		}
	case reflect.Struct:
		//TODO
	}

	return field, nil

}

func getFieldName(fieldName string) string {
	if strings.Index(fieldName, "[") >= 0 {
		return fieldName[0:strings.Index(fieldName, "[")]
	}

	return fieldName
}

func getFieldSliceIndex(fieldName string) (int, error) {
	if strings.Index(fieldName, "[") >= 0 {
		index := fieldName[strings.Index(fieldName, "[")+1 : strings.Index(fieldName, "]")]
		logrus.Debugf("sssssss %d", index)
		i, err := strconv.Atoi(index)
		if err != nil {
			return -2, nil
		}
		return i, nil
	}

	return -1, nil
}

func getFieldMapKey(fieldName string) (string, error) {
	if strings.Index(fieldName, "[") >= 0 {
		key := fieldName[strings.Index(fieldName, "[")+1 : strings.Index(fieldName, "]")]
		return strconv.Unquote(key)
	}

	return "", nil
}

func getReflectValue(in interface{}) reflect.Value {
	var value reflect.Value
	if reflect.TypeOf(in).Kind() == reflect.Ptr {
		value = reflect.ValueOf(in).Elem()
	} else {
		value = reflect.ValueOf(in)
	}
	return value
}
