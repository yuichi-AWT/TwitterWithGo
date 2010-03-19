package twitter
 
import (
	"http"
	"strings"
	"fmt"
	"bufio"
)

type BasicAuthParams struct {
	UserID string
	Password string
}

type Client struct {
	Auth BasicAuthParams
}

func New(username string , password string) (c *Client) {
	c = new(Client)
	c.Auth = BasicAuthParams{ UserID: username, Password: password }
	return
}

func (c *Client)Send(method string, params map[string]string) (status_code int) {
	url := makeAPIUrl(c.Auth, method, params)
	response ,_ := http.Post(url, "application/xml", strings.NewReader(""))
//	fmt.Printf("%s\n%s\n", url, response.Status)
	status_code = response.StatusCode
	return
}

func (c *Client)Get(method string, params map[string]string) (line string) {
	url := makeAPIUrl(c.Auth, method, params)
	response ,_ ,_ := http.Get(url)
	reader := bufio.NewReader(response.Body)
	line, _ = reader.ReadString(0)
	fmt.Printf("%s\n", url)
	return
}

func makeAPIUrl(auth BasicAuthParams, method string, params map[string]string) (url string) {
	url = "http://" + auth.UserID + ":" + auth.Password + "@twitter.com/"
	url = url + method
	if len(params) > 0 {
		url += "?"
		for key, val := range params {
			url += key + "=" + val
		}
	}
	return
}
