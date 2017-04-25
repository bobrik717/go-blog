package models

import (
	"crypto/rand"
	"fmt"
	"html/template"
)

func GenerateId() string {
	b := make([]byte,16)
	rand.Read(b)
	return fmt.Sprintf("%x",b)
}

func Unescape(x string) interface {} {
	return template.HTML(x)
}