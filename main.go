package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)


var letters = []rune("abcdefghijklmnopqrstuvwxyz")

var accountList string

type Data struct {
	Catchall string `json:"catchall"`
	Amount   int    `json:"amount"`
}

func main(){
	configFile, err := os.Open("config.json")
	if(err != nil) { fmt.Println(err) }
	config, _ := ioutil.ReadAll(configFile)

	var data Data
	json.Unmarshal(config, &data)

	proxyListCont, err := ioutil.ReadFile("proxies.txt")
	proxyListArr := strings.Split(string(proxyListCont), "\r\n")
	if(err != nil) { fmt.Println(err) }
	
	var wg sync.WaitGroup
	for i := 0; i < data.Amount; i++ {
		time.Sleep(100 * time.Millisecond)
		wg.Add(1)
		
		go genAcc(data.Catchall, proxyListArr, &wg)
	}
	wg.Wait()
	writeFile()
	
	time.Sleep(1 * time.Second)
	fmt.Println("Done!")
	fmt.Scanln()
}

func genAcc(catchall string, proxyListArr []string, wg *sync.WaitGroup){	
	defer wg.Done()

	rand.Seed(time.Now().UTC().UnixNano())
	proxySeed := rand.Intn(len(proxyListArr))

	proxy := strings.Split(proxyListArr[proxySeed], ":")
	
	var proxyURL string
	if(len(proxy) == 4){
		proxyURL = "http://" + proxy[2] + ":" + proxy[3] + "@" + proxy[0] + ":" + proxy[1]
	} else {
		proxyURL = "http://" + proxy[0] + ":" + proxy[1]
	}

	proxyParsed, err := url.Parse(proxyURL)
	if(err != nil){ fmt.Println(err) }

	tr := &http.Transport{
		MaxIdleConns:       1,
		IdleConnTimeout:    7 * time.Second,
		DisableCompression: true,
		Proxy:              http.ProxyURL(proxyParsed),
	}

	client := &http.Client{
		Transport:  tr,
	}

	emailVal := randSeq(7)
	time.Sleep(10 * time.Millisecond)
	passwordVal := randSeq(8)
	email := emailVal + catchall

	form := url.Values{}
	form.Add("email", email)
	form.Add("password", passwordVal)
	form.Add("confirmpassword", passwordVal)
	form.Add("gender", "no-thanks")
	form.Add("source", "SSENSE_EN_SIGNUP")

	postUrl := "https://www.ssense.com/en-us/account/register"

	req, _ := http.NewRequest(http.MethodPost, postUrl, strings.NewReader(form.Encode()))
	req.Header.Set("accept", "application/json")
	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("content-length", "115")
	req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Set("dnt", "1")
	req.Header.Set("origin", "https://www.ssense.com")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("referer", "https://www.ssense.com/en-us/account/login")
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="90", "Google Chrome";v="90"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	if(resp.StatusCode == 200){
		fmt.Println("[" + strconv.Itoa(resp.StatusCode) + "]" + " Successfully made account")

		accountList = accountList + email + ":" + passwordVal + "\n"
	} else {
		fmt.Println("[" + strconv.Itoa(resp.StatusCode) + "]" + " Account creation failed") 
	}
}

func writeFile(){
	file, err := os.OpenFile("accounts.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if(err != nil){ fmt.Println(err) }

	defer file.Close()

	if _, err := file.WriteString(accountList); err != nil {
		fmt.Println(err)
	}
}

func randSeq(n int) string {
    b := make([]rune, n)
	rand.Seed(time.Now().UTC().UnixNano())
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}