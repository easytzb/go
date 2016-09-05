package main

import (
	"log"
	"regexp"
)

//func xs45() ([]byte, error) {
func xs45() {
	url := "http://www.45xs.com/books/31/31565/"
	html, err := get(url)
	if err != nil {
		return
	}

	matchUrl := regexp.MustCompile(`"infot">.+?<a href="(.+?)" target`).FindSubmatch(html)
	if len(matchUrl) == 0 {
		log.Println("解析新URL失败")
		return
	} else {
		//log.Println(newUrl)
	}

	newUrl := url + string(matchUrl[1])
}
