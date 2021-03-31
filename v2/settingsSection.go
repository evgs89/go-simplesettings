package simplesettings

import (
	"fmt"
	"sync"
)

// SettingsSection is a structure to hold key-value pairs and process them as settings values
type SettingsSection struct {
	lock   sync.RWMutex
	Values map[string]*settingsValue
}

func newSettingsSection() *SettingsSection {
	ss := &SettingsSection{}
	ss.Values = make(map[string]*settingsValue)
	return ss
}

func (ss *SettingsSection) addValue(sv *settingsValue) {
	ss.lock.Lock()
	defer ss.lock.Unlock()
	ss.Values[sv.Name] = sv
}

// DeleteValue with given name from this section
func (ss *SettingsSection) DeleteValue(name string) {
	ss.lock.Lock()
	defer ss.lock.Unlock()
	delete(ss.Values, name)
}

func (ss *SettingsSection) String() string {
	out := ""
	ss.lock.RLock()
	for _, val := range ss.Values {
		out += fmt.Sprintf("%v\n", val)
	}
	ss.lock.RUnlock()
	return out
}

func (ss *SettingsSection) getVal(key string) settingsValue {
	ss.lock.RLock()
	defer ss.lock.RUnlock()
	val := ss.Values[key]
	if val == nil {
		panic(fmt.Sprintf("Key %v not found in this settings section", key))
	}
	return *val
}

// Get string value from SettingsSection object
func (ss *SettingsSection) Get(key string) string {
	return ss.getVal(key).ParseString()
}

// GetInt - get integer value from SettingsSection object
func (ss *SettingsSection) GetInt(key string) int {
	return ss.getVal(key).ParseInt()
}

// GetBool - get boolean value from SettingsSection object
func (ss *SettingsSection) GetBool(key string) bool {
	return ss.getVal(key).ParseBool()
}

// GetArray - get []string value from SettingsSection object
func (ss *SettingsSection) GetArray(key string) []string {
	return ss.getVal(key).ParseArray()
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
