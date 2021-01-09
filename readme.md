# SimpleSettings

## Overview

It's a simple package to read from and save to *.INI-like files. It parses almost all variants of INI and saves current settnings to file. 

*Warning* Saving settings to existing file will delete all comments!

## Usage example

_settings.ini_
```ini
IntVal = 123

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
	IntVal := s.Get("", "IntVal")  // 123
	BoolVal := s.Get("section1", "BoolVal") // true
	StrVal := s.Get("section1", "StrVal") // "abc"
	ArrVal := s.Get("section2", "ArrVal") // []string{"a", "b", "c"}
	
	// save
	err := s.Set("section2", "NewArrVal", []string{"aa", "bb", "cc"})
	
	// write to disk
	err = s.SaveToFile("modified_settings.ini")
}
```