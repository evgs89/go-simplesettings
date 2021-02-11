package simplesettings

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type settingsSection map[string]*SettingsValue

func newSettingsSection() *settingsSection {
	ss := settingsSection{}
	return &ss
}

func (ss settingsSection) AddValue(sv *SettingsValue) {
	ss[sv.Name] = sv
}

func (ss settingsSection) DeleteValue(name string) {
	delete(ss, name)
}

func (ss settingsSection) String() string {
	out := ""
	for _, val := range ss {
		out += fmt.Sprintf("%v\n", val)
	}
	return out
}

type SettingsValue struct {
	Name        string
	StringValue string
}

// NewValue - SettingsValue constructor
func NewValue(name string, value interface{}) *SettingsValue {
	switch value.(type) {
	case string:
		return &SettingsValue{Name: name, StringValue: value.(string)}
	case int:
		return &SettingsValue{name, strconv.Itoa(value.(int))}
	case bool:
		if value.(bool) {
			return &SettingsValue{name, "TRUE"}
		} else {
			return &SettingsValue{name, "FALSE"}
		}
	case []string:
		return &SettingsValue{name, strings.Join(value.([]string), ", ")}
	default:
		panic(fmt.Sprintf("Can't save value as Settings element: %v", value))
		return nil
	}
}

func newValueFromString(st string) (*SettingsValue, error) {
	// it's possibly value
	split := strings.SplitN(st, "=", 2)
	if len(split) > 1 {
		name, value := strings.Trim(split[0], " "), strings.Trim(split[1], " ")
		return NewValue(name, value), nil
	} else {
		// nope, something else
		log.Printf("Error parsing string %v", st)
		return nil, nil
	}
}

// ParseString returns settings value as string
func (val *SettingsValue) ParseString() string {
	return val.StringValue
}

// ParseInt returns settings value as integer or panics
func (val *SettingsValue) ParseInt() int {
	if intVal, err := strconv.Atoi(val.StringValue); err != nil {
		panic(fmt.Sprintf("Error parsing param %v to int", val.Name))
		return 0
	} else {
		return intVal
	}
}

// ParseArray returns settings value as slice of strings
func (val *SettingsValue) ParseArray() []string {
	splitVal := strings.Split(val.StringValue, ",")
	var arrVal []string
	for _, v := range splitVal {
		arrVal = append(arrVal, strings.Trim(v, " "))
	}
	return arrVal
}

// ParseBool tries to parse boolean value, if can't, it will return true if string is not empty
func (val *SettingsValue) ParseBool() bool {
	if strings.ToLower(val.StringValue) == "false" || val.StringValue == "" || val.StringValue == "0" {
		return false
	}
	return true
}

func (val *SettingsValue) String() string {
	return fmt.Sprintf("%v = %v", val.Name, val.StringValue)
}
