package gostruct

import (
	"fmt"
	"testing"
)

type address struct {
	Streat  string
	ZipCode string
	State   string
	City    []*city
	Emails  []string
	Maps    map[string]interface{}
	In      interface{}
}

type city struct {
	Name string
	InUS bool
	Park *park
}

type park struct {
	Name     string
	Location string
	Maps     map[string]interface{}
	Emails   []string
}

func TestGetStruct(t *testing.T) {
	s := &address{}
	s.Streat = "311 wind st"
	s.ZipCode = "77479"
	s.State = "Taxes"
	s.Emails = []string{"123@123.com", "456@456.com"}
	m := make(map[string]interface{})
	m["dd"] = "dd"
	m["cc"] = "cc"
	m["bb"] = "bb"

	s.Maps = m

	s.City = append(s.City, &city{Name: "Sugar Land", InUS: true, Park: &park{Name: "Name", Location: "location", Maps: m}})
	s.In = "string222"

	field, err := GetField(s, "Emails")
	if err != nil {
		panic(err)
	}

	fmt.Println("Emails:", field)

	field, err = GetField(s, "Emails[0]")
	if err != nil {
		panic(err)
	}

	fmt.Println("Emails[0]:", field)

	field3, err := GetField(s, `City[0].Park.Maps`)
	if err != nil {
		panic(err)
	}
	fmt.Println("City[0].Park.Maps:", field3)

}
