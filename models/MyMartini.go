package models

import (
	"github.com/codegangsta/martini"
	"os"
)

type MyMartiniClassic struct {
	martini.ClassicMartini
}

func (m *MyMartiniClassic) RunCustom(port string) {
	if len(port) == 0 {
		port = os.Getenv("PORT")
	}

	if len(port) == 0 {
		port = "3000"
	}

	host := os.Getenv("HOST")

	m.RunOnAddr(host + ":" + port)
}
