package main

import (
	"BizCoinWebSocket/config"
	"BizCoinWebSocket/funcs"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var addr = flag.String(
	"addr",
	config.Conf.HostPort,
	"http service address",
)

func main() {
	fmt.Println("127.0.0.1:3000")
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/time_game", funcs.TimeGame)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
