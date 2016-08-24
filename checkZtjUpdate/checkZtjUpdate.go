package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
)

func main() {
	//r, err := http.Get("http://m.qidian.com/book/showbook.aspx?bookid=3347595")
	r, err := http.Get("http://www.qidian.com/Book/3347595.aspx")
	if err != nil {
		fmt.Println("获取网页内容失败:1", err.Error())
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		fmt.Println("获取网页内容失败:2", err.Error())
		return
	}

	//fmt.Println(string(body[:]))
	//return

	//date := regexp.MustCompile(`\d{4}(/|-)\d{1,2}(/|-)\d{1,2} \d{1,2}:\d{1,2}:\d{1,2}`).FindString(string(body[:]))
	date := regexp.MustCompile(`dateModified">(\d{4}-\d\d-\d\d \d\d:\d\d)<\/span>`).FindSubmatch(body)
	if len(date) == 0 {
		fmt.Println("解析更新日期失败:3")
		return
	} else {
		//fmt.Println(date)
	}

	fr, err := os.Open("/webser/logs/ztj.date")
	if err != nil {
		fmt.Println("日期记录文件打开失败:4", err.Error())
		return
	}
	defer fr.Close()

	oldDate := make([]byte, 16)
	count, _ := fr.Read(oldDate)
	//if err != nil {
	//	fmt.Println("原更新时间读取失败:5", err.Error())
	//	return
	//}

	if count == 0 || string(oldDate[:]) != string(date[1][:]) {
		//新的更新时间写入文件
		fw, er := os.Create("/webser/logs/ztj.date")
		if er != nil {
			fmt.Println("日期记录文件打开失败:7", err.Error())
			return
		}
		defer fw.Close()

		if _, err := fw.Write(date[1]); err != nil {
			fmt.Println("更新时间写入失败:6", err.Error())
			return
		}

		//调用PHP脚本发送微信通知
		cmd := exec.Command("/usr/bin/php", "/webser/www/tchat/cron/index.php", "sendMsg2Me", "maoni updated!!!")
		err = cmd.Run()
		if err != nil {
			fmt.Println("微信信息发送失败:6", err.Error())
		}

	}

	return
}
