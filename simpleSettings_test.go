package simplesettings

import (
	"os"
	"testing"
)

const testSettingsFilename = "test_settings.ini"

func TestNewSettingsFromFile(t *testing.T) {
	t.Parallel()
	s := NewSettingsFromFile(testSettingsFilename)
	if s.getVal("", "ListValue") == nil {
		t.Error("Fail to read param from root")
	}
	if s.getVal("SECTION 1", "BoolValue3") == nil {
		t.Error("Fail to read empty param")
	}
}

func TestSettings_Get(t *testing.T) {
	t.Parallel()
	s := NewSettingsFromFile(testSettingsFilename)
	// test string
	tvs := s.Get("SECTION 1", "StringValue1")
	if tvs != "SomeText" {
		t.Error("Fail to read string")
	}
	// test int
	tvi := s.GetInt("SECTION 2", "IntValue")
	if tvi != 2398 {
		t.Error("Fail to read int")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Error("ParseInt didn't Fail on wrong value")
		}
	}()
	tvi = s.GetInt("SECTION 1", "StringValue1")
	// test bool
	tvb := s.GetBool("SECTION 1", "BoolValue1")
	if !tvb {
		t.Error("Fail to read bool")
	}
	tvb = s.GetBool("SECTION 1", "BoolValue2")
	if tvb {
		t.Error("Fail to read bool")
	}
	tvb = s.GetBool("SECTION 1", "BoolValue3")
	if tvb {
		t.Error("Fail to read bool")
	}
	tvb = s.GetBool("SECTION 1", "BoolValue2")
	if tvb {
		t.Error("Fail to read bool")
	}
	tvb = s.GetBool("SECTION 1", "StringValue1")
	if !tvb {
		t.Error("Fail to read bool")
	}
	// test array
	tva := s.GetArray("", "ListValue")
	if len(tva) != 3 {
		t.Error("Fail to read array")
	}
}

func assertEqual(t *testing.T, v1, v2 interface{}) {
	if v2 != v1 {
		t.Errorf("%v must be equal to %v", v2, v1)
	}
}

func TestSettings_Set(t *testing.T) {
	t.Parallel()
	s := NewSettings()
	_ = s.Set("", "Val1", 123)
	assertEqual(t, s.Get("", "Val1"), "123")
	_ = s.Set("section1", "Val2", true)
	assertEqual(t, s.Get("section1", "Val2"), "TRUE")
	_ = s.Set("section1", "Val3", "abc")
	assertEqual(t, s.Get("section1", "Val3"), "abc")
	_ = s.Set("section2", "Val4", []string{"a", "b", "c"})
	assertEqual(t, s.Get("section2", "Val4"), "a, b, c")
	_ = s.SaveToFile("generated.ini")
	_ = os.Remove("generated.ini")
}
