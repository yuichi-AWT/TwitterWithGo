package twitter

import (
	"json"
	"bufio"
	"io"
	"net"
	"bytes"
	"encoding/base64"
	"strings"
	"http"
)

// 以下3つはjsonデータのパース用
type User struct {
	Name        string
	Screen_name string
}
type Tweet struct {
	User User
	Text string
}
type Tweets struct {
	Tweet []Tweet
}

// ユーザアカウント情報
type Auth struct {
	Name	string
	Passwd	string
}
// クライアント
// こいつにいくつかインタフェースを実装
type Client struct{
	Auth Auth
}

func ClientBuilder(name string, passwd string) (cc Client){
	var c Client
	c.Auth.Name = name
	c.Auth.Passwd = passwd

	return c
}

// Auth構造体に関連付けたメソッド
// BASIC認証用にName/PasswdをBASE64エンコード
func (a *Auth) Base64enc() (s string){
	authSrc := a.Name + ":" + a.Passwd
	bb := &bytes.Buffer{}
	encoder := base64.NewEncoder(base64.StdEncoding, bb)
	encoder.Write(strings.Bytes(authSrc))
	encoder.Close()
	authEnc := bb.String()

	return authEnc
}

// ツイートをポストするメソッド
func (c *Client)Post(stat string) {
	// 送信本文をURLエンコード
	sStatURL := http.URLEscape(stat)

	if con, err := net.Dial("tcp", "", "twitter.com:80"); err == nil {
		io.WriteString(con, "POST /statuses/update.json?status=" + sStatURL + " HTTP/1.1\r\n")
		io.WriteString(con, "Host: twitter.com\r\n")
		io.WriteString(con, "Authorization: Basic " + c.Auth.Base64enc() + "\r\n")
		io.WriteString(con, "\r\n")

		con.Close()
	}
}

// フレンドのTLを取得するメソッド
func (c *Client)Get(url string) (t Tweets) {
	var tweets Tweets

	if con, err := net.Dial("tcp", "", "twitter.com:80"); err == nil {
		io.WriteString(con, "GET " + url + " HTTP/1.1\r\n")
		io.WriteString(con, "Host: twitter.com\r\n")
		io.WriteString(con, "Authorization: Basic "+ c.Auth.Base64enc() +"\r\n")
		io.WriteString(con, "\r\n")

		reader := bufio.NewReader(con)

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err.String() == "EOF" {
					json.Unmarshal("{ \"tweet\": "+string(line)+"}", &tweets)
				}
				break
			}
		}
		con.Close()
	}

	return tweets
}

func (c *Client)HomeTimeline() (t Tweets){
	return c.Get("/statuses/home_timeline.json")
}
func (c *Client)FriendsTimeline() (t Tweets){
	return c.Get("/statuses/friends_timeline.json")
}
func (c *Client)UserTimeline() (t Tweets){
	return c.Get("/statuses/user_timeline.json")
}
func (c *Client)Replies() (t Tweets){
	return c.Get("/statuses/replies.json")
}
func (c *Client)Mentions() (t Tweets){
	return c.Get("/statuses/mentions.json")
}

type Users struct{
	User []User
}

func (c *Client)Friends() (u Users) {
	var users Users

	if con, err := net.Dial("tcp", "", "twitter.com:80"); err == nil {
		io.WriteString(con, "GET " + "/statuses/friends.json" + " HTTP/1.1\r\n")
		io.WriteString(con, "Host: twitter.com\r\n")
		io.WriteString(con, "Authorization: Basic "+ c.Auth.Base64enc() +"\r\n")
		io.WriteString(con, "\r\n")

		reader := bufio.NewReader(con)

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err.String() == "EOF" {
					json.Unmarshal("{ \"user\": "+string(line)+"}", &users)
				}
				break
			}
		}
		con.Close()
	}

	return users
}

