package simplesettings

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type settingsSection map[string]*settingsValue

func newSettingsSection() *settingsSection {
	ss := settingsSection{}
	return &ss
}

func (ss settingsSection) AddValue(sv *settingsValue) {
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

type settingsValue struct {
	Name        string
	StringValue string
}

func NewValue(name string, value interface{}) *settingsValue {
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
		return NewValue(name, value), nil
	} else {
		// nope, something else
		log.Printf("Error parsing string %v", st)
		return nil, nil
	}
}

func (val *settingsValue) ParseString() string {
	return val.StringValue
}

func (val *settingsValue) ParseInt() int {
	if intVal, err := strconv.Atoi(val.StringValue); err != nil {
		panic(fmt.Sprintf("Error parsing param %v to int", val.Name))
		return 0
	} else {
		return intVal
	}
}

func (val *settingsValue) ParseArray() []string {
	splitVal := strings.Split(val.StringValue, ",")
	var arrVal []string
	for _, v := range splitVal {
		arrVal = append(arrVal, strings.Trim(v, " "))
	}
	return arrVal
}

func (val *settingsValue) ParseBool() bool {
	if strings.ToLower(val.StringValue) == "false" || val.StringValue == "" || val.StringValue == "0" {
		return false
	}
	return true
}

func (val *settingsValue) String() string {
	return fmt.Sprintf("%v = %v", val.Name, val.StringValue)
}

// Settings is storage object for program settings
type Settings map[string]*settingsSection

// NewSettings - Settings constructor.
func NewSettings() *Settings {
	s := Settings{"": newSettingsSection()}
	return &s
}

// NewSettingsFromFile - alternative constructor for creating Settings object of given INI-file.
func NewSettingsFromFile(filename string) *Settings {
	s := NewSettings()
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Can't open settings file")
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	currentSection := ""
	currentSectionValues := newSettingsSection()
	for scanner.Scan() {
		st := scanner.Text()
		if (!strings.HasPrefix(st, "#")) && (!strings.HasPrefix(st, ";")) && len(st) > 2 {
			if strings.HasPrefix(st, "[") {
				// it's section name
				(*s)[currentSection] = currentSectionValues
				currentSectionValues = newSettingsSection()
				currentSection = strings.Trim(st, "[ ]")
			} else {
				// it's possibly value
				if value, err := newValueFromString(st); err == nil {
					currentSectionValues.AddValue(value)
				}
			}
		}
	}
	(*s)[currentSection] = currentSectionValues
	return s
}

func (s *Settings) Get(section, key string) *settingsValue {
	ss := (*s)[section]
	if ss == nil {
		log.Fatalf("no such section in settings: %v", section)
		return nil
	}
	return (*(*s)[section])[key]
}

func (s *Settings) Set(section, key string, value interface{}) error {
	val := NewValue(key, value)
	if val == nil {
		return fmt.Errorf("can't save as value: %v", val)
	}
	ss := (*s)[section]
	if ss == nil {
		ss = &settingsSection{}
		(*s)[section] = ss
	}
	ss.AddValue(val)
	return nil
}

func (s *Settings) String() string {
	out := fmt.Sprintf("%v\n", (*s)[""])
	for section, values := range *s {
		if section == "" {
			continue
		}
		out += fmt.Sprintf("[ %v ]\n%v\n", section, values)
	}
	return out
}

func (s *Settings) SaveToFile(filename string) error {
	data := []byte(s.String())
	err := ioutil.WriteFile(filename, data, 0644)
	return err
}
