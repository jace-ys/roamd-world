package config

import (
	"github.com/alecthomas/kong"
)

type StaticBase struct {
	Config    string `help:"Path to YAML configuration file." env:"CONFIG" type:"path" default:"config/default.yaml"`
	Port      int    `help:"Port binding for the gRPC server." env:"PORT" default:"50051"`
	AdminPort int    `help:"Port binding for the HTTP admin server." env:"ADMIN_PORT" default:"8081"`
	LogLevel  string `help:"Minimum log level." env:"LOG_LEVEL" enum:"debug,info,warn,error,fatal" default:"debug"`
}

func MustParseStatic(static any) {
	ctx := kong.Parse(static, kong.Name("service"))
	if ctx.Error != nil {
		panic(ctx.Error)
	}
}
