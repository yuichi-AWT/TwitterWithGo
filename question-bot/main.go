package main

import (
	"twitter";
	"fmt";
	"xml";
	"regexp";
	"strings";
	"time";
	"strconv";
)

type User struct {
	Name string;
	ScreenName string;
}

type Status struct {
	CreatedAt string;
	Id string;
	Text string;
	User User;
}

type Statuses struct {
	Status []Status;
}

type Result struct {
	Statuses []Statuses;
}

func main() {
	client := twitter.New("your acount", "password");
	since_id := ""
	for true {
		fmt.Println(since_id)
		since_id = quester(client, since_id)
		time.Sleep(600000000000)
	}
}

func quester(client *twitter.Client, since_id string)(string){
	params := map[string]string{}
	if since_id != ""{
		params["since_id"] = since_id
	}
	jsn := client.Get("statuses/home_timeline.xml", params)
	var statuses Statuses;
	xml.Unmarshal(strings.NewReader(jsn), &statuses);

	if len(statuses.Status) == 0{
		return since_id
	}
	topTweet := statuses.Status[0]
	fmt.Println(topTweet)
	topTime, _ := time.Parse(time.RubyDate, topTweet.CreatedAt);
	for i := len(statuses.Status)-1 ; i >= 0; i-- {
		status := statuses.Status[i]
		if b ,_ := regexp.MatchString("?|？$", status.Text); b {
			fmt.Printf("user name:%s\n", status.User.ScreenName)
			fmt.Printf("text:%s\n", status.Text);
			tm, _ := time.Parse(time.RubyDate, status.CreatedAt);
			fmt.Printf("time:%d:%d:%d:%d:%d:%d\n", tm.Year, tm.Month, tm.Day, tm.Hour, tm.Minute, tm.Second);

			if tm.Seconds() + 10*60 > topTime.Seconds() {
				fmt.Println("end")
				return since_id
			}

			num, _ := count(statuses.Status[0:i], status.User.ScreenName)
			num_s := strconv.Itoa(num)
			message := "【集計】" + status.User.ScreenName +"さんに反応した人の数→" + num_s
			fmt.Println(message)
			client.Send("statuses/update.xml", map[string]string{"status": message})
		}
	}
	return topTweet.Id
}

func count(statuses []Status, user_id string) (num int, user_name []string){
	num = 0
	user_name = make([]string, len(statuses))
	for i := len(statuses)-1 ; i >= 0; i-- {
		status := statuses[i]
		if b ,_ := regexp.MatchString("@" + user_id, status.Text); b{
			fmt.Println(status.Text)
			user_name[num] = status.User.Name
			num ++
		}
	}
	return
}
