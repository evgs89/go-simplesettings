package simplesettings

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

// Settings is storage object for program settings
type Settings struct {
	lock     sync.RWMutex
	Sections map[string]*SettingsSection
}

// NewSettings - Settings object constructor.
func NewSettings() *Settings {
	sections := make(map[string]*SettingsSection)
	sections[""] = newSettingsSection()
	s := Settings{Sections: sections}
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
	//currentSectionValues := NewSettingsSection()
	for scanner.Scan() {
		st := scanner.Text()
		if (!strings.HasPrefix(st, "#")) && (!strings.HasPrefix(st, ";")) && len(st) > 2 {
			if strings.HasPrefix(st, "[") {
				// it's section name
				currentSection = strings.Trim(st, "[ ]")
			} else {
				// it's possibly value
				if value, err := newValueFromString(st); err == nil {
					s.getOrCreateSection(currentSection, true).addValue(value)
				}
			}
		}
	}
	return s
}

func (s *Settings) getOrCreateSection(section string, create bool) *SettingsSection {
	s.lock.RLock()
	ss := s.Sections[section]
	if ss == nil {
		if !create {
			s.lock.RUnlock()
			log.Fatalf("no such section in settings: %v", section)
		}
		s.lock.RUnlock()
		s.lock.Lock()
		ss = newSettingsSection()
		s.Sections[section] = ss
		s.lock.Unlock()
	} else {
		s.lock.RUnlock()
	}
	return ss
}

// GetSection to pass it to sub-routine (for example, db settings section to db driver)
// Panics if section not found in settings object
func (s *Settings) GetSection(section string) *SettingsSection {
	return s.getOrCreateSection(section, false)
}

//AddSection adds new section to settings object and returns pointer to created section
// Panics if such section is already exists
func (s *Settings) AddSection(section string) *SettingsSection {
	if s.HasSection(section) {
		panic(fmt.Sprintf("Section %v already is in settings, overwriting is not allowed due to data safety", section))
	}
	return s.getOrCreateSection(section, true)
}

// SectionsList is a way to simply get slice of section names
func (s *Settings) SectionsList() []string {
	sectionsList := make([]string, 0, len(s.Sections))
	s.lock.RLock()
	for key := range s.Sections {
		sectionsList = append(sectionsList, key)
	}
	s.lock.RUnlock()
	return sectionsList
}

// HasSection returns true if settings object has section with given name
func (s *Settings) HasSection(section string) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.Sections[section] != nil
}

// Get string value from Settings object
func (s *Settings) Get(section, key string) string {
	return s.GetSection(section).Get(key)
}

// GetInt - get integer value from Settings object
func (s *Settings) GetInt(section, key string) int {
	return s.GetSection(section).GetInt(key)
}

// GetBool - get boolean value from Settings object
func (s *Settings) GetBool(section, key string) bool {
	return s.GetSection(section).GetBool(key)
}

// GetArray - get []string value from Settings object
func (s *Settings) GetArray(section, key string) []string {
	return s.GetSection(section).GetArray(key)
}

// Set string, int, bool, slice value to Settings object
func (s *Settings) Set(section, key string, value interface{}) error {
	ss := s.getOrCreateSection(section, true)
	return ss.Set(key, value)
}

// String prints all settings in INI-like style
func (s *Settings) String() string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	out := fmt.Sprintf("%v\n", s.GetSection(""))
	for section, values := range s.Sections {
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
