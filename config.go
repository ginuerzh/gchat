// config
package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Server      string `short:"s" long:"server" default:"talk.google.com:443" description:"XMPP server [ip:port]"`
	Resource    string `short:"r" long:"resource" description:"resource"`
	EnableProxy bool
	Proxy       string `short:"p" long:"proxy" description:"Proxy server [ip:port]"`
	UseSysProxy bool   `short:"a" long:"sproxy" description:"Use system proxy"`
	NoTLS       bool   `long:"notls" description:"Do not use tls"`
	Username    string `long:"user" description:"Username"`
	Password    string `long:"pass" description:"Password"`
	AutoLogin   bool
	EnableDebug bool `short:"d" long:"debug" description:"Enable debug" json:"-"`
	NoGui       bool `long:"nogui" description:"Run on command mode" json:"-"`
}

func (config *Config) Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(config)
}

func (config *Config) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(config)
}
