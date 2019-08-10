package main

import (
	"github.com/spf13/pflag"
)

type config struct {
	listener string
	logging  bool
	dir      string
	watch    bool
	noip     noipConfig
}

type noipConfig struct {
	user     string
	pass     string
	host     string
	interval int
}

var c = config{}

func init() {
	pflag.StringVarP(&c.listener, "listener", "l", "127.0.0.1:8989", "Listener that should be binded")
	pflag.BoolVar(&c.logging, "logging", false, "Write logs to STDOUT")
	pflag.StringVarP(&c.dir, "directory", "d", ".", "Directory that should be served")
	pflag.StringVar(&c.noip.user, "noip.user", "", "User to access no-ip")
	pflag.StringVar(&c.noip.pass, "noip.pass", "", "Password to access no-ip")
	pflag.StringVar(&c.noip.host, "noip.host", "", "Host to update via no-ip")
	pflag.IntVar(&c.noip.interval, "noip.interval", 30, "Interval to update no-ip")
	pflag.BoolVarP(&c.watch, "watch", "w", false, "Watch and refresh")
}
