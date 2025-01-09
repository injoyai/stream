package api

import (
	"github.com/injoyai/goutil/frame/mux"
)

var Server = mux.New()

func Run(port int) error {
	Server.SetPort(port)
	return Server.Run()
}
