package main

import "github.com/rai-project/raid/cmd"

var AppSecret string

func init() {
	cmd.AppSecret = AppSecret
}
