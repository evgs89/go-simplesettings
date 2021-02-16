# SimpleSettings

## Overview

It's a simple package to read from and save to *.INI-like files. It parses almost all variants of INI and saves current settnings to file. 

*Warning* Saving settings to existing file will delete all comments!

## Usage example

_settings.ini_
```ini
IntVal = 123

# SomeComment

; SomeOtherComment

[ section1 ]
BoolVal = TRUE
StrVal = abc

[ section2 ]
ArrVal = a, b, c


```

_main.go_
```go
package main

import "github.com/evgs89/go-simplesettings"

func main() {
	s := simplesettings.NewSettingsFromFile("settings.ini")
	// read
	StrVal := s.Get("section1", "StrVal")       // "abc"
	IntVal := s.GetInt("", "IntVal")            // 123
	BoolVal := s.GetBool("section1", "BoolVal") // true
	ArrVal := s.GetArray("section2", "ArrVal")  // []string{"a", "b", "c"}

	// save
	err := s.Set("section2", "NewArrVal", []string{"aa", "bb", "cc"}) // nil
	err = s.Set("section2", "NewBoolVal", false)                      // nil
    
    // You can use SettingsSection object directly and pass it to some sub-routine
    section1 := (*s)["section1"]
    s1BoolVal := section1.GetBool("BoolVar")
    err = section1.Set("NewStringVal", "SomeString") 	

    newSection := simplesettings.NewSettingsSection()
    err = newSection.Set("NewIntVal", 123)
    (*s)["newSection"] = newSection

	// write to disk
	err = s.SaveToFile("modified_settings.ini") // nil
}
```
