package main

import (
	"BizCoinWebSocket/config"
	"BizCoinWebSocket/funcs"
	"flag"
	"log"
	"net/http"
)

var addr = flag.String(
	"addr",
	config.Conf.HostPort,
	"http service address",
)

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/time_game", funcs.TimeGame)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
