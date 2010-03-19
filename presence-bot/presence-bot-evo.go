package main

import (
	"twitter"
	"file"
	"bytes"
	"strings"
	"rand"
	"time"
)

func main(){
	s := GetGoroku()

	client := twitter.ClientBuilder("your account", "password")

	tweet := "【在席情報】" + s
	client.Post(tweet)
}

func GetGoroku() (s string){
	file, _ := file.Open("info.txt", 0, 0)
	if file == nil{
		return ""
	}

	var buf [100000]byte
	file.Read(&buf)

	bb := bytes.NewBuffer(&buf)
	str := bb.String()

	strs := strings.Split(str, "\n", 10000)

	rand.Seed(time.Nanoseconds())
	tmp := rand.Int31()
	i := tmp % (int32)(len(strs)-1)
	return strs[i]
}
