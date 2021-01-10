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
	IntVal := s.Get("", "IntVal").ParseInt()  // 123
	BoolVal := s.Get("section1", "BoolVal").ParseBool() // true
	StrVal := s.Get("section1", "StrVal").ParseStr() // "abc"
	ArrVal := s.Get("section2", "ArrVal").ParseArr() // []string{"a", "b", "c"}
	
	// save
	err := s.Set("section2", "NewArrVal", []string{"aa", "bb", "cc"})
	err = s.Set("section2", "NewBoolVal", false)
	
	// write to disk
	err = s.SaveToFile("modified_settings.ini")
}
```