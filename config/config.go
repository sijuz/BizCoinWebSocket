package config

import "github.com/gorilla/websocket"

// Upgrader is any-event handler
var Upgrader = websocket.Upgrader{}

// Conf is app configurations
var Conf = Config{
	HostPort:          "127.0.0.1:3000",
	DbHostPort:        "35.228.208.75:3306",
	Workers:           4,
	DbName:            "bizgame",
	DbUser:            "root",
	DbPassword:        "PxK7jy7D4cjBNBE0",
	VkAppID:           "asdf",
	VkAppSecret:       "asdf",
	VkAppServiceToken: "asdf",
	UseVkSignChecker:  false,
	MinPrice: 100,
	MaxPrice: 1000,
	MinLoss: 700,
	MaxLoss: 1500,
	MinProfit: 200,
	MaxProfit: 1200,
}

// DBSettings is db params
//var DBSettings map[string]string = map[string]string{
//	"host": "35.228.208.75",
//	"port": "3306",
//	"name": "bizgame",
//	"user": "root",
//	"pswd": "PxK7jy7D4cjBNBE0",
//}
