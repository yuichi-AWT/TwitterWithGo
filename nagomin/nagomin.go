package main

import (
	"fmt"
	"twitter"
	"time"
	"file"
	"os"
	"rand"
	"strings"
	"bufio"
	"container/vector"
)

const (
	MIN_TIME   = 1800 //botツイートする際の最低間隔(秒)
	TWEET_TIME = 3600 //botがツイートする条件(秒)
	BOT_NAME   = "nagomi_bot"
	PASSWORD   = "saturati"
)

func main() {
	//実行間隔を600秒に設定
	var wait int64 = 600000000000

	//データベースの読み込み処理
	f, _ := file.Open("nagomin.txt", os.O_RDONLY, 644)
	buf := make([]byte, 65536)
	f.Read(buf)
	f.Close()
	var nagomin vector.StringVector
	nagomin = strings.Split(string(buf), "\n", 0)
	nagomin = nagomin[0 : len(nagomin)-1]

	//ボットの処理の開始
	for {
		tmp := nagomi_bot(nagomin)
		time.Sleep(wait)
		nagomin = tmp
	}
}

func nagomi_bot(nagomin vector.StringVector) vector.StringVector {
	rand.Seed(time.Nanoseconds())

	// client情報を作成
	c := twitter.ClientBuilder(BOT_NAME, PASSWORD)

	// Timelineを取得
	tweets := c.HomeTimeline(0, 0, 0, 0)
	//fmt.Printf("tweet = %v\n",tweets)

	// botへのリプライを取得
	direct_messages := c.DirectMessages(0, 0, 0, 0)

	//DMの内容をデータベースに登録する処理の開始
	for i := 0; i < len(direct_messages); i++ {
		f, _ := file.Open("nagomin.txt", os.O_WRONLY|os.O_APPEND, 644)
		bufy := bufio.NewWriter(f)

		//DMの内容がすでにデータベースに登録されてないかチェック
		flag := 0
		for j := 0; j < len(nagomin); j++ {
			//登録されている場合の処理
			if string(nagomin[j]) == direct_messages[i].Text {
				flag = 1
				break
			}
		}
		//データベースに新規登録処理
		if flag == 0 {
			fmt.Printf("%v\n", direct_messages[i].Text)
			bufy.WriteString(direct_messages[i].Text + "\n")
			//データベースに追加
			bufy.Flush()
			//すでに読み込んでいるデータベースを更新
			nagomin.Push(direct_messages[i].Text)
		}
		f.Close()
	}

	//現在の時刻の取得
	current_t := time.LocalTime()

	tweet := tweets[0]
	//最新ツイートの時刻を取得
	last_tweet_t, _ := time.Parse(time.RubyDate, tweet.Created_at)
	timespan := (current_t.Seconds() - last_tweet_t.Seconds())

	//ツイートする最低間隔を設定(botが連続でツイートすることを防止する)
	if timespan < MIN_TIME {
		return nagomin
	}

	//APIからのつぶやきを最新ツイートに含めないようにする処理
	i := 0
	for i = 0; tweets[i].Source == "<a href=\"http://apiwiki.twitter.com/\" rel=\"nofollow\">API</a>"; i++ {
	}

	tweet = tweets[i]
	last_tweet_t, _ = time.Parse(time.RubyDate, tweet.Created_at)

	//現在時刻と最新ツイートの時刻との差分を取得
	timespan = (current_t.Seconds() - last_tweet_t.Seconds())

	//ツイート処理
	if timespan > TWEET_TIME {
		str := nagomin[rand.Int31n(int32(len(nagomin)))]
		fmt.Printf("%d秒経過したので，ツイートします．\n", timespan)
		c.Post(fmt.Sprintf("【nagomin】%s", str))
	}
	return nagomin
}
