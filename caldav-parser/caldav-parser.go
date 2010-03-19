package main

import(
	"fmt"
	"xml"
	"os"
	"file"
	"regexp"
)

type Event struct{
	name string
	start string
	end string
}

func main(){
	multi, err := ReadLog()
	if err != nil{
		print(err.String())
		return
	}

	events := make([]Event, len(multi.Response))

	for i, res := range multi.Response{
		reg, _ := regexp.Compile("SUMMARY:[^\n]+\n")
		ss := reg.MatchStrings(res.Propstat.Prop.CalendarData)
		reg, _ = regexp.Compile("SUMMARY:")
		events[i].name = reg.ReplaceAllString(ss[0], "")

		reg, _ = regexp.Compile("DTSTART;[^\n]+\n")
		ss = reg.MatchStrings(res.Propstat.Prop.CalendarData)

		reg, _ = regexp.Compile("DTEND;[^\n]+\n")
		ss = reg.MatchStrings(res.Propstat.Prop.CalendarData)
	}

	for i, e := range events{
		fmt.Printf("%d:%s", i, e.name)
	}
}

type Prop struct{
	CalendarData string	
}

type Propstat struct{
	Prop Prop
}

type Response struct{
	Href string
	Propstat Propstat
}
type Multistatus struct{
	Response []Response
}

func ReadLog()(m Multistatus, e os.Error){
	var ret Multistatus

	file, _ := file.Open("caldav.log", 0, 0)
	if file == nil{
		return ret, os.NewError("can not open file.")
	}

	err := xml.Unmarshal(file, &ret)
	if err != nil{
		return ret, err
	}


	return ret, nil
}
