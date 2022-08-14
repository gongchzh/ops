package ctl

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

type Transport struct {
	Transport http.RoundTripper
}

func (t *Transport) transport() http.RoundTripper {
	if nil != t.Transport {
		return t.Transport
	}
	return http.DefaultTransport
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0; .NET4.0C; .NET4.0E; .NET CLR 2.0.50727; .NET CLR 3.0.4506.2152; .NET CLR 3.5.30729)")
	return t.transport().RoundTrip(req)
}

func NewHttpClient() *http.Client {
	t := &Transport{}
	jar, _ := cookiejar.New(nil)
	client := http.DefaultClient
	client.Transport = t
	client.Jar = jar
	return client
}

type HtmlRes struct {
	Text string
}

func (res *HtmlRes) ParseText() string {
	return strings.ReplaceAll(res.Text, "<br>", "\n")
}

func (res *HtmlRes) Red(s ...interface{}) {

	res.Text += `<a11 style="font-style:normal;font-size:20px;color:red" >` + fmt.Sprint(s...) + `</a11><br>`
}

func (res *HtmlRes) Black(s ...interface{}) {
	res.Text += `<a11 style="font-style:normal;font-size:20px;color:black" >` + fmt.Sprint(s...) + `</a11><br>`
}

func (res *HtmlRes) Purple(s ...interface{}) {
	res.Text += `<a11 style="font-style:normal;font-size:20px;color:purple" >` + fmt.Sprint(s...) + `</a11><br>`
}

func (res *HtmlRes) Orange(s ...interface{}) {
	res.Text += `<a11 style="font-style:normal;font-size:20px;color:orange" >` + fmt.Sprint(s...) + `</a11><br>`
}

func (res *HtmlRes) Blue(s ...interface{}) {
	res.Text += `<a11 style="font-style:normal;font-size:20px;color:blue" >` + fmt.Sprint(s...) + `</a11><br>`
}

func (res *HtmlRes) Green(s ...interface{}) {
	res.Text += `<a11 style="font-style:normal;font-size:20px;color:green" >` + fmt.Sprint(s...) + `</a11><br>`
}

func (res *HtmlRes) Yellow(s ...interface{}) {
	res.Text += `<a11 style="font-style:normal;font-size:20px;color:yellow" >` + fmt.Sprint(s...) + `</a11><br>`
}

func (res *HtmlRes) Write() []byte {
	res.Text = `<!doctype html>
<html>
<body>` + res.Text + `</body>
	</html>`

	return []byte(res.Text)
}
