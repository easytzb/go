package main

import (
	"gopkg.in/redis.v4"
	"log"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	myopenid := "ogw2owQfmfI1AG0yFncJ5J8BV_4M"

	redis := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := redis.Ping().Result()
	if err != nil {
		log.Println("redis 连接失败", err)
		return
	}
	defer redis.Close()

	latestTimeStr, _ := redis.HGet("userLatestTime", myopenid).Result()
	var msg string
	if "" != latestTimeStr {
		latestTime, err1 := strconv.ParseInt(latestTimeStr, 10, 64)
		if err1 != nil {
			log.Println("最后时间转64位整数失败", err1)
			return
		}

		diff := time.Now().Sub(time.Unix(latestTime, 0))
		if diff < 42*time.Hour {
			return
		}
		msg = "最后的交互距现在已有 " + time.Now().Sub(time.Unix(latestTime, 0)).String()
	} else {
		msg = "还没有记录最后交互时间，来一个"
	}

	err = exec.Command("/usr/bin/php", "/webser/www/tchat/cron/index.php", "sendMsg2Me", msg).Run()
	if err != nil {
		log.Println("微信信息发送失败:6", err)
		return
	}
}
