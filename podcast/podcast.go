package main

import (
	"github.com/mmcdole/gofeed"
	"gopkg.in/redis.v4"
	"log"
	"os/exec"
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

	fp := gofeed.NewParser()
	for t, url := range atom {
		feed, _ := fp.ParseURL(url)
		if feed.Items[0].Published == "" {
			continue
		}

		oldPublished, _ := redis.HGet("atom", t).Result()
		if feed.Items[0].Published == oldPublished {
			continue
		}

		err := exec.Command("/usr/bin/php", "/webser/www/tchat/cron/index.php", "sendMsg2Me", t+" "+feed.Items[0].Title).Run()
		if err != nil {
			log.Println(t, "微信信息发送失败:6", err)
			continue
		}

		_, err = redis.HSet("atom", t, feed.Items[0].Published).Result()
		if err != nil {
			log.Println(t, "更新日期写入失败", err)
			continue
		}
	}
}
