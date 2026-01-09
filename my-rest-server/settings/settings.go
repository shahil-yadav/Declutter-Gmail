package settings

import (
	"log"

	"github.com/go-ini/ini"
)

var cfg *ini.File

type Auth struct {
	RedirectUrl string // redirect-url setup from google-cloud-console
	Secret      string // used to create cookie secrets
	SessionName string // session name of cookie
}

var AuthSettings = &Auth{}

type Server struct {
	// can be debug or release
	RunMode  string
	HttpPort int
}

var ServerSetting = &Server{}

// parsed app.ini config
func Setup() {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("server", ServerSetting)
	mapTo("auth", AuthSettings)

	// ...
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
