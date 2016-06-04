package live

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func GetWorker() {
	const path = "http://127.0.0.1:8010/api/workers/"
	var (
		err      error
		response *http.Response
		data     []interface{}
	)

	response, err = http.Get(path)
	if err != nil {
		log.Println(err)
		// TODO check for error and if OWTF is down shutdown monitor gracefully
	}
	defer response.Body.Close()
	var dataByte []byte
	dataByte, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(dataByte, &data)
	if err != nil {
		log.Println("Error occured during decoding")
	}
	fmt.Println(data)
	fmt.Println(len(data))
}
