package service

import "spm/config"

var services []func()

var cfg = config.Read("service")
var key = []byte(cfg["service_key"])

// register service
func register(s func()) {
	services = append(services, s)
}

// start service
func Start() {
	for _, service := range services {
		go service()
	}
}
