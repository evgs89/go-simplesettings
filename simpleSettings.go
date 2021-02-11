package simplesettings

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

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

// Get SettingsValue from settings
func (s *Settings) getVal(section, key string) *SettingsValue {
	ss := (*s)[section]
	if ss == nil {
		log.Fatalf("no such section in settings: %v", section)
		return nil
	}
	return (*(*s)[section])[key]
}

// Get string value from Settings object
func (s *Settings) Get(section, key string) string {
	return s.getVal(section, key).ParseString()
}

// GetInt - get integer value from Settings object
func (s *Settings) GetInt(section, key string) int {
	return s.getVal(section, key).ParseInt()
}

// GetBool - get boolean value from Settings object
func (s *Settings) GetBool(section, key string) bool {
	return s.getVal(section, key).ParseBool()
}

// GetArray - get []string value from Settings object
func (s *Settings) GetArray(section, key string) []string {
	return s.getVal(section, key).ParseArray()
}

// Set string, int, bool, slice value to settings object
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

// String prints all settings in INI-like style
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

// SaveToFile saves Settings object to INI-like text file
func (s *Settings) SaveToFile(filename string) error {
	data := []byte(s.String())
	err := ioutil.WriteFile(filename, data, 0644)
	return err
}
