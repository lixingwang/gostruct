# gostruct

A go tool to set/get struct fields value from expression string which similar jsonpath.

#Example

## Get
Get value from struct using expression
```
    City[0].Park.Address
    City[0].Park.Name
    Emails[0]
    field, err := GetField(s, `City[0].Park.Maps`)
	  field, err = GetField(s, "Emails[0]")
```

## Set
```
	err := SetField(s, "Person.Streat", "hello ST")
	err := SetField(s, `Maps["lix"]`, "wangzai")
	err := SetField(s, "Emails[0]", "555@125553.com")
 ```

