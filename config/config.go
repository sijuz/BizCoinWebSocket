package config

import "github.com/gorilla/websocket"

// Upgrader is any-event handler
var Upgrader = websocket.Upgrader{}

// Conf is app configurations
var Conf = Config{}
