package simplesettings

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type settingsValue struct {
	Name        string
	StringValue string
}

// NewValue - settingsValue constructor
func newValue(name string, value interface{}) *settingsValue {
	switch value.(type) {
	case string:
		return &settingsValue{Name: name, StringValue: value.(string)}
	case int:
		return &settingsValue{name, strconv.Itoa(value.(int))}
	case bool:
		if value.(bool) {
			return &settingsValue{name, "TRUE"}
		} else {
			return &settingsValue{name, "FALSE"}
		}
	case []string:
		return &settingsValue{name, strings.Join(value.([]string), ", ")}
	default:
		panic(fmt.Sprintf("Can't save value as Settings element: %v", value))
		return nil
	}
}

func newValueFromString(st string) (*settingsValue, error) {
	// it's possibly value
	split := strings.SplitN(st, "=", 2)
	if len(split) > 1 {
		name, value := strings.Trim(split[0], " "), strings.Trim(split[1], " ")
		return newValue(name, value), nil
	} else {
		// nope, something else
		log.Printf("Error parsing string %v", st)
		return nil, nil
	}
}

// ParseString returns settings value as string
func (val *settingsValue) ParseString() string {
	return val.StringValue
}

// ParseInt returns settings value as integer or panics
func (val *settingsValue) ParseInt() int {
	if intVal, err := strconv.Atoi(val.StringValue); err != nil {
		panic(fmt.Sprintf("Error parsing param %v to int", val.Name))
		return 0
	} else {
		return intVal
	}
}

// ParseArray returns settings value as slice of strings
func (val *settingsValue) ParseArray() []string {
	splitVal := strings.Split(val.StringValue, ",")
	var arrVal []string
	for _, v := range splitVal {
		arrVal = append(arrVal, strings.Trim(v, " "))
	}
	return arrVal
}

// ParseBool tries to parse boolean value, if can't, it will return true if string is not empty
func (val *settingsValue) ParseBool() bool {
	if strings.ToLower(val.StringValue) == "false" || val.StringValue == "" || val.StringValue == "0" {
		return false
	}
	return true
}

func (val *settingsValue) String() string {
	return fmt.Sprintf("%v = %v", val.Name, val.StringValue)
}
