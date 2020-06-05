package main

import (
	"BizCoinWebSocket/additionally"
	"BizCoinWebSocket/config"
	"BizCoinWebSocket/funcs"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

func configure() error{
	conf, err := ioutil.ReadFile("config/config.json")
	if  err != nil {
		return err
	}

	var newConf config.Config
	err = json.Unmarshal(conf, &newConf)
	if err != nil {
		return err
 	}
	config.Conf = newConf

	return nil
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	if err := configure(); err != nil {
		log.Println("error in configure config:", err)
		additionally.SendError("if err := configure(); err != nil {", err)
		return
	}

	var addr = flag.String(
		"addr",
		config.Conf.HostPort,
		"http service address",
	)
	log.Println(config.Conf.HostPort)
	close := make(chan struct{})
	go funcs.ConnList(close)

	http.HandleFunc("/time_game", funcs.TimeGame)
	log.Fatal(http.ListenAndServe(*addr, nil))
	close <- struct{}{}
}
