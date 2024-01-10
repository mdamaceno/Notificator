package main

import "github.com/mdmaceno/notificator/config"

func main() {
	env := config.Envs()
	config.InitDB(env)
}
