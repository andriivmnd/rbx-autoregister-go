package helpers

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
)

var Proxies = []*Proxy{{Client: &http.Client{}}}

type Proxy struct {
	Client *http.Client
	Ip     string
}

func LoadProxies() {
	pFile, err := os.Open("proxies.txt")
	if err != nil {
	}
	defer func(pFile *os.File) {
		err := pFile.Close()
		if err != nil {

		}
	}(pFile)

	scanner := bufio.NewScanner(pFile)
	for scanner.Scan() {
		parse, _ := url.Parse(fmt.Sprintf("http://%s", scanner.Text()))
		Proxies = append(Proxies, &Proxy{Client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(parse),
			},
		}, Ip: scanner.Text()})
	}

	Proxies = Proxies[1:]

}

func GetRandomProxy() *Proxy {
	return Proxies[rand.Intn(len(Proxies))]
}
