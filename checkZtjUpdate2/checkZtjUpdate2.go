package main

import (
	"gopkg.in/redis.v5"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	//r, err := http.Get("http://m.qidian.com/book/showbook.aspx?bookid=3347595")
	url := "http://zetianjixiaoshuo.com/"
	r, err := http.Get(url)
	if err != nil {
		log.Println("获取网页内容失败:1", url, err)
		return
	}

	//log.Println(r.Header.Get("Last-Modified"), r.Body)
	//return

	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		log.Println("获取网页内容失败:2", err)
		return
	}

	//log.Println(string(body[:]))
	//return

	//date := regexp.MustCompile(`\d{4}(/|-)\d{1,2}(/|-)\d{1,2} \d{1,2}:\d{1,2}:\d{1,2}`).FindString(string(body[:]))
	match := regexp.MustCompile(`最新章节：<a href="(.+?)" title`).FindSubmatch(body)
	if len(match) == 0 {
		log.Println("解析最新章节地址失败:3")
		return
	} else {
		//log.Println(date)
	}

	redis := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err = redis.Ping().Result()
	if err != nil {
		log.Println("redis 连接失败", err)
		return
	}
	defer redis.Close()

	oldUrl, _ := redis.HGet("ztj", "url").Result()
	if oldUrl == string(match[1]) {
		log.Println("还没有更新")
		return
	}

	newUrl := url + string(match[1])
	nr, err2 := http.Get(newUrl)
	if err2 != nil {
		log.Println("获取网页内容失败:1", newUrl, err2)
		return
	}

	//log.Println(r.Header.Get("Last-Modified"), r.Body)
	//return

	nbody, err3 := ioutil.ReadAll(nr.Body)
	nr.Body.Close()
	if err3 != nil {
		log.Println("获取网页内容失败:2", err3)
		return
	}

	nmatch := regexp.MustCompile(`<p>(.|\s)+?</p>`).FindAllString(string(nbody), -1)
	if len(nmatch) == 0 {
		log.Println("解析最新内容失败:3")
		return
	}

	msg := regexp.MustCompile(`\n|<p>|</p>`).ReplaceAllString(strings.Join(nmatch, "__LINE__"), "")

	//调用PHP脚本发送微信通知
	err = exec.Command("/usr/bin/php", "/webser/www/tchat/cron/index.php", "sendMsg2Me", msg).Run()
	if err != nil {
		log.Println("微信信息发送失败:6", err)
		return
	}

	//新的更新时间写入文件
	_, err = redis.HSet("ztj", "url", string(match[1])).Result()
	if err != nil {
		log.Println("新url写入失败	:7", err)
		return
	}

	return
}
