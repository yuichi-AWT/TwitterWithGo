package main

import (
	"io"
	"net"
	"bufio"
	"regexp"
	"strings"
	"syscall"
)

func main(){
	_, strs := GetGoroku()

	fd, _ := syscall.Open("goroku.txt", 1, 644)

	for _, s := range strs {
		syscall.Write(fd, strings.Bytes(s))	
	}
}

func GetGoroku() (i int, s []string){
	goroku := make(map[int] string)

	if con, err := net.Dial("tcp", "", "XX.XX.XX.XX:YY"); err == nil {
		io.WriteString(con, "GET " + "/cgi-bin/show.cgi" + " HTTP/1.1\r\n")
		io.WriteString(con, "Host: XX.XX.XX.XX\r\n")
		io.WriteString(con, "Authorization: Basic XXXXXXXXXXXXXXXXX\r\n")
		io.WriteString(con, "\r\n")

		reader := bufio.NewReader(con)

		i := 0
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			if ret, _ := regexp.MatchString("<td colspan=\"5\">", line); ret == true {
				line, _ := reader.ReadString('\n')
				re, _ := regexp.Compile("[<br />|&nbsp;]")
				line = re.ReplaceAllString(line, "")
				goroku[i] = line
				i++
			}
		}
		con.Close()
	}

	length := len(goroku)
	ret := make([]string, length)
	for i, s := range goroku{
		ret[i] = s
	}

	return length, ret
}
