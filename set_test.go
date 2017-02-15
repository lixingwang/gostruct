package gostruct

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestSetSlice(t *testing.T) {
	s := &address{}
	s.Streat = "311 wind st"
	s.ZipCode = "77477"

	err := SetField(s, "Emails[]", []string{"123@123.com", "456@456.com"})
	if err != nil {
		panic(err)
	}
	err = SetField(s, "Emails[0]", "555@125553.com")
	err = SetField(s, "Emails[1]", "555@125553.com")

	if err != nil {
		panic(err)
	}
	v, _ := json.Marshal(s)
	fmt.Println(string(v))

}

func TestSetMap(t *testing.T) {
	s := &address{}
	s.Streat = "311 wind st"
	s.ZipCode = "77477"
	m := make(map[string]interface{})
	m["dd"] = "dd"
	m["cc"] = "cc"
	m["bb"] = "bb"
	err := SetField(s, "Maps", m)
	if err != nil {
		panic(err)
	}

	err = SetField(s, `Maps["lix"]`, "wangzai")
	if err != nil {
		panic(err)
	}
	v, _ := json.Marshal(s)
	fmt.Println(string(v))

}

func TestSetNestMap(t *testing.T) {
	s := &address{}
	s.Streat = "311 wind st"
	s.ZipCode = "77477"
	m := make(map[string]interface{})
	m["dd"] = "dd"
	m["cc"] = "cc"
	m["bb"] = "bb"

	err := SetField(s, "Streat", "hello ST")
	err = SetField(s, "Maps", m)

	if err != nil {
		panic(err)
	}

	v, _ := json.Marshal(s)
	fmt.Println("%s", string(v))

}
