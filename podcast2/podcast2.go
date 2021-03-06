package main

import (
	"github.com/mmcdole/gofeed"
	"gopkg.in/redis.v5"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func main() {

	atom := map[string]string{
		"IPN": "http://ipn.li/feed",
		"PAGE SEVEN 纪录片": "http://nj.lizhi.fm/rss/29470.xml",
		"Teahour":        "http://teahour.fm/feed.xml",
		"一席":             "http://rss.kaolafm.com/MZ_RSS/rss/1100000046138/album.xml",
		"反派影评":           "http://www.ximalaya.com/album/4127591.xml",
		"吴晓波频道":          "http://www.ximalaya.com/album/269179.xml",
		"新闻酸菜馆":          "http://since1989.org/feed/wasai",
		"旅行麻辣烫":          "http://trip.since1989.org/feed/podcast",
		"极客电台":           "http://geek.wasai.org/feed/podcast",
		"枝桠":             "http://rss.kaolafm.com/MZ_RSS/rss/1100000083390/album.xml",
		"电影不无聊":          "http://rss.kaolafm.com/MZ_RSS/rss/1100000046058/album.xml",
		"电影沙龙":           "http://www.ximalaya.com/album/273057.xml",
		"罗辑思维":           "http://podcast.taobility.com/category/logicthinking/feed/",
		"大内密谈":           "http://nj.lizhi.fm/rss/14275.xml",
		"轻阅读":            "http://feed.cri.cn/rss/31295a33-de9b-409a-82ce-7905207f6c00",
	}

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

	fp, client := gofeed.NewParser(), &http.Client{}
	for t, url := range atom {
		req, err := http.NewRequest("GET", url, nil)

		if strings.Contains(url, "nj.lizhi.fm") {
			req.Header.Add("Accept-Encoding", "identity")
		}

		res, err := client.Do(req)
		if err != nil {
			log.Println(t, "获取feed内容失败:1", url, err)
			continue
		}
		if res != nil {
			defer func() {
				ce := res.Body.Close()
				if ce != nil {
					err = ce
				}
			}()
		}

		if res.StatusCode < 200 || res.StatusCode > 300 {
			log.Println(t, "获取feed内容状态码不是200:2", res.StatusCode)
			continue
		}

		body, err2 := ioutil.ReadAll(res.Body)
		if err2 != nil {
			log.Println(t, "获取feed内容body部分失败:3", err2)
			continue
		}

		feed, err := fp.ParseString(string(body))
		if err != nil || feed.Items[0].Published == "" {
			log.Println(t, "解析feed错误:4", err)
			continue
		}

		itemsLen := len(feed.Items)
		if itemsLen == 0 {
			log.Println(t, "items长度为0:5", err)
			continue
		}

		oldPublished, err2 := redis.HGet("atom1", t).Int64()
		if err2 != nil {
			//默认10天前
			log.Println(t, "最后更新时间获取失败", err2)
			oldPublished = 0
		}

		//正则匹配出item
		match := regexp.MustCompile(`(<item>(?:.|\s)+?<\/item>)`).FindAllString(string(body), -1)
		if match == nil || len(match) == 0 {
			log.Println(t, "正则匹配出错")
			continue
		}

		loopTimes := itemsLen
		if itemsLen > 2 {
			loopTimes = 2
		}
		for i := 0; i < loopTimes; i++ {
			newPublished := feed.Items[i].PublishedParsed.Unix()
			if newPublished <= oldPublished {
				//还没有更新
				break
			}

			err = exec.Command("/usr/bin/php", "/webser/www/tchat/cron/index.php", "sendMsg2Me", t+" "+feed.Items[i].Title).Run()
			if err != nil {
				log.Println(t, "微信信息发送失败:6", err)
				continue
			}

			_, err = redis.LPush("atomQue", match[i]).Result()
			if err != nil {
				log.Println(t, "item写入失败", err)
				continue
			}
		}
		_, err = redis.HSet("atom1", t, strconv.FormatInt(feed.Items[0].PublishedParsed.Unix(), 10)).Result()
		if err != nil {
			log.Println(t, "更新日期写入失败", err)
			continue
		}
	}
}
