// Package repository error.go contains errors
package repository

import "fmt"

// ErrNil means that u've givet nil entity for a create method
var ErrNil = fmt.Errorf("entity that u've given is nil")

var ErrExist = fmt.Errorf("such username already exist")
