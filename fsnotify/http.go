package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func get(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		log.Println("获取网页内容失败:1", url, err)
		return nil, err
	}
	defer r.Body.Close()

	return ioutil.ReadAll(r.Body)
}
