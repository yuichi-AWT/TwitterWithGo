package main

import (
	"./twitter"
	"./file"
	"fmt"
	"bytes"
	"strings"
	"os"
)


func main(){
    var buf [512]byte;

	//ユーザ情報をファイルから読み込む
	f, err := file.Open("acountinfo.txt",0,0);
	if f == nil {
		fmt.Printf("Error: %s\n",err);
		os.Exit(1);
	}

	f.Read(&buf);
	str := bytes.NewBuffer(&buf).String();
	name := strings.Split(str, "\n", 10000)[0]
	passwd := strings.Split(str, "\n", 10000)[1]
	tweet := strings.Split(str, "\n", 10000)[2]

	//コマンドライン入力からtweetを受け取る
	/*
	in := file.Stdin;
	fmt.Printf("tweet >");
	in.Read(&buf);
	str = bytes.NewBuffer(&buf).String();
	tweet := strings.Split(str, "\n", 10000)[0];
	in.Close();
	*/

	//確認
	/*
	fmt.Printf("name   =>%s<=\n",name);
	fmt.Printf("passwd =>%s<=\n",passwd);
	*/
	fmt.Printf("tweet  =>%s<=\n",tweet);

	client := twitter.ClientBuilder(name,passwd);

	//ClientのtimeLineを取得してみる
	//HomeTimeLine => Clientが受け取ったtweetのリスト
	//FriendsTimeline => 友達だけのもの？かと思うが，Homeと同じ
	//userTimeline => 自分のツイートのみ取得
	// Replies,Mentions は不明，パニックを起こすのだが・・・
	/*
	timeline := client.HomeTimeline();
	*/
	timeline := client.HomeTimeline();
	tweetList := timeline.Tweet;
	for i := 0; i < 5 ; i++ {
		fmt.Printf("%s -> %s \n",tweetList[i].User.Screen_name,tweetList[i].Text);
	}



/*
	//FriendListを取得・表示
	users := client.Friends();
	userList := users.User;

	for i:= 0 ; i < len(userList); i++	{
		fmt.Printf("friend -> %s  display_name -> %s \n", userList[i].Name, userList[i].Screen_name);
	}
*/

/*
	//tweetをデコレーション
	tweet = "[test tweet]" + tweet

	//つぶやく
	client.Post(tweet);
*/
}
