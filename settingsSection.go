package simplesettings

import (
	"fmt"
)

type SettingsSection map[string]*settingsValue

func NewSettingsSection() *SettingsSection {
	ss := SettingsSection{}
	return &ss
}

func (ss SettingsSection) addValue(sv *settingsValue) {
	ss[sv.Name] = sv
}

func (ss SettingsSection) DeleteValue(name string) {
	delete(ss, name)
}

func (ss SettingsSection) String() string {
	out := ""
	for _, val := range ss {
		out += fmt.Sprintf("%v\n", val)
	}
	return out
}

// Get string value from SettingsSection object
func (ss *SettingsSection) Get(key string) string {
	return (*ss)[key].ParseString()
}

// GetInt - get integer value from SettingsSection object
func (ss *SettingsSection) GetInt(key string) int {
	return (*ss)[key].ParseInt()
}

// GetBool - get boolean value from SettingsSection object
func (ss *SettingsSection) GetBool(key string) bool {
	return (*ss)[key].ParseBool()
}

// GetArray - get []string value from SettingsSection object
func (ss *SettingsSection) GetArray(key string) []string {
	return (*ss)[key].ParseArray()
}

// Set string, int, bool, slice value to SettingsSection object
func (ss *SettingsSection) Set(key string, value interface{}) error {
	val := newValue(key, value)
	if val == nil {
		return fmt.Errorf("can't save as value: %v", val)
	}
	ss.addValue(val)
	return nil
}
