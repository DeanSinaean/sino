package sino
import (
	"net/http"
	"net/url"
	"os"
	"io"
//	"fmt"
)

func GetByProxy(target_addr, proxy_addr string) (*http.Response, error) {
	request, _ := http.NewRequest("GET", target_addr, nil)
	proxy, err := url.Parse(proxy_addr)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}
	return client.Do(request)
}

//usage example
// var proxy_url string="http://username:password@ip:port/"
func download(file_name string,file_url string,proxy_url string) {
	resp, _ := GetByProxy(file_url, proxy_url)
	defer resp.Body.Close()
	file,_:=os.Create(file_name)
	io.Copy(file,resp.Body)
}
