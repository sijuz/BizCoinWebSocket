package additionally

import (
	"BizCoinWebSocket/config"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func SendError(content string, errMes error){

	data, err := json.Marshal(map[string]string{
		"key": "jhFjk71SH5JbQu18hZLpQjb1691",
		"error": "GOLANG:"+content + " " + errMes.Error(),
		"host": config.Conf.MyHostPort,
	})
	if  err != nil {
		return
	}

	toSend := bytes.NewReader(data)

	recv, err := http.Post(
		"http://"+config.Conf.BugReportHostPort,
		"application/json",
		toSend,
		)
	if err != nil {
		log.Println(recv.Body, err)
		return
	}

	log.Println(recv.Body)

	return
}
