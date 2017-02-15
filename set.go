package gostruct

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"strconv"
	"reflect"
	"strings"
)

func SetField(s interface{}, fieldName string, value interface{}) error {
	return setField(reflect.ValueOf(s), fieldName, fieldName, value)
}

func setField(v reflect.Value, fieldName, currName string, value interface{}) error {
	logrus.Debugf("First kind %s", v.Kind())

	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("Struct must be a pointer")
	}

	if v.IsNil() {
		logrus.Debugf("Field %s is nil and set a new ", fieldName)
		v.Set(reflect.New(v.Type().Elem()))
	}

	v = reflect.Indirect(v)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if fieldName == "" {

		switch inputv := value.(type) {
		case string:
			err := setStringValue(v, inputv)
			if err != nil {
				return err
			}
		case []byte:
			err := setStringValue(v, string(inputv))
			if err != nil {
				return err
			}
		default:
			valv := reflect.ValueOf(value)
			for valv.Kind() == reflect.Ptr {
				valv = valv.Elem()
			}
			v.Set(valv)
		}

		return nil
	}

	switch v.Kind() {
	case reflect.Struct, reflect.Ptr:
		currName, nextFieldName := getCurrAndNextFieldName(currName)
		currFieldName := getFieldName(currName)
		logrus.Debugf("Curr name %s and next field Name %s", currFieldName, nextFieldName)

		if v.Kind() == reflect.Struct {
			v = v.FieldByName(currFieldName)
		} else {
			v = v.Elem().FieldByName(currFieldName)
		}
		if !v.IsValid() {
			return fmt.Errorf("No such field: %s in obj", currFieldName)
		}
		//if nextFieldName == "" {
		//	nextFieldName = fieldName
		//}
		logrus.Debugf("Field type %s", v.Kind())
		if v.Kind() == reflect.Ptr {
			err := setField(v, fieldName, nextFieldName, value)
			if err != nil {
				return err
			}
		} else {
			err := setField(v.Addr(), fieldName, nextFieldName, value)
			if err != nil {
				return err
			}
		}
	case reflect.Slice:
		av := v
		currName, nextFieldName := getCurrAndNextFieldName(fieldName)
		logrus.Debugf("Curr name %s and next field Name %s", currName, nextFieldName)

		currFieldName := getFieldName(currName)
		index, err := getFieldSliceIndex(currName)
		if err != nil {
			return err
		}

		elementType := v.Type().Elem()
		var newslice reflect.Value
		logrus.Debugf("Sclice index %d", index)
		if index != -1 && index != -2 {
			//Set to specific
			arrayElement := av.Index(index)
			if !arrayElement.IsValid() {
				arrayElement.Set(reflect.New(elementType).Elem())
			}
			err = SetField(arrayElement.Addr().Interface(), currFieldName, value)
			if err != nil {
				return err
			}

		} else if index == -2 {
			valueOf := reflect.ValueOf(value)
			if v.Type() != valueOf.Type() {
				return fmt.Errorf("Provided value type (%v) didn't match obj field type (%v)\n", valueOf.Type(), v.Type())
			}
			v.Set(valueOf)
		} else {
			//Append to the slice
			arrayyElement := reflect.New(elementType).Elem()
			err = SetField(arrayyElement.Addr().Interface(), currFieldName, value)
			if err != nil {
				logrus.Errorf("Set Field Error %+v", err)
				return err
			}
			newslice = av
			logrus.Debugf("arrayyElement %+v %+v", arrayyElement)
			newslice = reflect.Append(newslice, arrayyElement)
			v.Set(newslice)

		}
	case reflect.Map:
		av := v
		currName, nextFieldName := getCurrAndNextFieldName(fieldName)
		logrus.Debugf("Curr name %s and next field Name %s", currName, nextFieldName)

		key, err := getFieldMapKey(currName)
		if err != nil {
			return err
		}

		logrus.Debugf("Map key %s", key)
		if key != "" {
			//Set to specific
			av.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
		} else if key == "" {
			valueOf := reflect.ValueOf(value)
			if v.Type() != valueOf.Type() {
				return fmt.Errorf("Provided value type (%v) didn't match obj field type (%v)\n", valueOf.Type(), v.Type())
			}
			v.Set(valueOf)
		}
	default:
		valueOf := reflect.ValueOf(value)
		if v.Type() != valueOf.Type() {
			return fmt.Errorf("Provided value type (%v) didn't match obj field type (%v)\n", valueOf.Type(), v.Type())
		}
		v.Set(valueOf)

	}
	return nil
}

func getCurrAndNextFieldName(name string) (string, string) {
	currName := name
	nextFieldName := ""
	if i := strings.Index(name, "."); i > -1 {
		currName = name[0:i]
		nextFieldName = name[i+1 : len(name)]
	}
	return currName, nextFieldName
}

func setStringValue(v reflect.Value, value string) (err error) {
	s := value

	// if type is []byte
	if v.Kind() == reflect.Slice && v.Type().Elem().Kind() == reflect.Uint8 {
		v.SetBytes([]byte(s))
		return
	}

	switch v.Kind() {
	case reflect.String:
		v.SetString(s)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var n int64
		n, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			return
		}
		if v.OverflowInt(n) {
			err = fmt.Errorf("overflow int64 for %d.", n)
			return
		}
		v.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		var n uint64
		n, err = strconv.ParseUint(s, 10, 64)
		if err != nil {
			return
		}
		if v.OverflowUint(n) {
			err = fmt.Errorf("overflow uint64 for %d.", n)
			return
		}
		v.SetUint(n)
	case reflect.Float32, reflect.Float64:
		var n float64
		n, err = strconv.ParseFloat(s, v.Type().Bits())
		if err != nil {
			return
		}
		if v.OverflowFloat(n) {
			err = fmt.Errorf("overflow float64 for %d.", n)
			return
		}
		v.SetFloat(n)
	case reflect.Bool:
		var n bool
		n, err = strconv.ParseBool(s)
		if err != nil {
			return
		}
		v.SetBool(n)
	default:
		err = fmt.Errorf("value %+v can only been set to primary type but was %+v", value, v)
	}

	return
}
