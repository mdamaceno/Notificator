package config

import "os"

type AppConf struct {
	APP_ID       string
	DB_HOST      string
	DB_USER      string
	DB_PASS      string
	DB_NAME      string
	DB_PORT      string
	DB_TEST_HOST string
	DB_TEST_USER string
	DB_TEST_PASS string
	DB_TEST_NAME string
	DB_TEST_PORT string
}

func Envs() *AppConf {
	return &AppConf{
		APP_ID:       os.Getenv("APP_ID"),
		DB_HOST:      os.Getenv("DB_HOST"),
		DB_USER:      os.Getenv("DB_USER"),
		DB_PASS:      os.Getenv("DB_PASS"),
		DB_NAME:      os.Getenv("DB_NAME"),
		DB_PORT:      os.Getenv("DB_PORT"),
		DB_TEST_HOST: os.Getenv("DB_TEST_HOST"),
		DB_TEST_USER: os.Getenv("DB_TEST_USER"),
		DB_TEST_PASS: os.Getenv("DB_TEST_PASS"),
		DB_TEST_NAME: os.Getenv("DB_TEST_NAME"),
		DB_TEST_PORT: os.Getenv("DB_TEST_PORT"),
	}
}
