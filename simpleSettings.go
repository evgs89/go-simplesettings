package simplesettings

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type settingsSection struct {
	Values map[string]string
}

func newSettingsSection() *settingsSection {
	ss := settingsSection{Values: make(map[string]string)}
	return &ss
}

func (ss *settingsSection) AddValue(name string, sv string) {
	ss.Values[name] = sv
}

func (ss *settingsSection) DeleteValue(name string) {
	delete(ss.Values, name)
}

// Settings is storage object for program settings
type Settings struct {
	Sections []string
	Values   map[string]*settingsSection
}

// NewSettings - Settings constructor.
func NewSettings() *Settings {
	s := Settings{
		Sections: []string{"#root"},
		Values:   make(map[string]*settingsSection),
	}
	return &s
}

// NewSettingsFromFile - alternative constructor for creating Settings object of given file.
func NewSettingsFromFile(filename string) *Settings {
	s := NewSettings()
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Can't open settings file")
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	currentSection := "#root" // special!!!
	currentSectionValues := newSettingsSection()
	for scanner.Scan() {
		st := scanner.Text()
		if (!strings.HasPrefix(st, "#")) && len(st) > 2 {
			if strings.HasPrefix(st, "[") {
				// it's section name
				s.Values[currentSection] = currentSectionValues
				currentSectionValues = newSettingsSection()
				currentSection = strings.Trim(st, "[ ]")
				s.Sections = append(s.Sections, currentSection)
			} else {
				// it's possibly value
				split := strings.SplitN(st, "=", 2)
				if len(split) > 1 {
					key, value := strings.Trim(split[0], " "), strings.Trim(split[1], " ")
					currentSectionValues.AddValue(key, value)
				} else {
					// nope, something else
					log.Fatal("Error parsing settings file")
				}
			}
		}
	}
	s.Values[currentSection] = currentSectionValues
	return s
}

func (s *Settings) getVal(section, param string) string {
	if section == "" {
		section = "#root"
	}
	var tmp *settingsSection
	if s.Values[section] == tmp {
		log.Fatalf("No such section in settings file: %v\n", section)
		return ""
	}
	sectionMap := *s.Values[section]
	settingsValue := sectionMap.Values[param]
	if settingsValue == "" {
		log.Fatalf("No such param in section %v of settings file: %v\n", section, param)
		return ""
	}
	if strings.HasPrefix(settingsValue, `"`) && strings.HasSuffix(settingsValue, `"`) {
		settingsValue = strings.Trim(settingsValue, `"`)
	}
	return settingsValue
}

// GetStr - Get string value from param of section in Settings object.
func (s *Settings) GetStr(section, param string) string {
	return s.getVal(section, param)
}

// GetInt - Get integer value from param of section in Settings object. Fails if can't convert value to integer!
func (s *Settings) GetInt(section, param string) int {
	val := s.getVal(section, param)
	intVal, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("Error parsing settings value %v - %v to int\n", section, param)
		return 0
	}
	return intVal
}
