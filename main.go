package main

import (
	"RobloxGen/helpers"
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/gookit/color"
)

var (
	modkernel32         = syscall.NewLazyDLL("kernel32.dll")
	procSetConsoleTitle = modkernel32.NewProc("SetConsoleTitleW")
)

var (
	accounts = 0
)

func setConsoleTitle(title string) {
	utf16Title := syscall.StringToUTF16Ptr(title)
	_, _, _ = procSetConsoleTitle.Call(uintptr(unsafe.Pointer(utf16Title)))
}

type ChallengeMetadata struct {
	CaptchaToken     string `json:"captchaToken"`
	UnifiedCaptchaID string `json:"unifiedCaptchaId"`
	DataExchangeBlob string `json:"dataExchangeBlob"`
	ActionType       string `json:"actionType"`
	RequestPath      string `json:"requestPath"`
	RequestMethod    string `json:"requestMethod"`
}

type CaptchaSolved struct {
	UnifiedCaptchaID string `json:"unifiedCaptchaId"`
	CaptchaToken     string `json:"captchaToken"`
	ActionType       string `json:"actionType"`
}

func logWithTime(msg string) {
	timeStr := time.Now().Format("15:04:05")
		("[<dark_gray>" + timeStr + "</>]" + msg)
}

func Birthday() string {
	timestampStr := "1989-02-01T08:00:00.000Z"
	layout := "2006-01-02T15:04:05.000Z"
	timestamp, err := time.Parse(layout, timestampStr)
	if err != nil {
		fmt.Println("Error parsing timestamp:", err)
		return ""
	}

	minTime := timestamp.AddDate(-1, 0, 0)
	maxTime := timestamp.AddDate(1, 0, 0)
	minYear := minTime.Year()
	maxYear := maxTime.Year()

	rand.Seed(time.Now().UnixNano())
	randomYear := rand.Intn(maxYear-minYear+1) + minYear

	randomTime := time.Date(randomYear, time.Month(rand.Intn(12)+1), rand.Intn(28)+1, rand.Intn(24), rand.Intn(60), rand.Intn(60), 0, time.UTC)

	return randomTime.Format(layout)
}

func checkUsername(p *helpers.Proxy, username string) string {
	for {
		username = username + helpers.GenerateRandomString(2)
		reqBody, _ := json.Marshal(map[string]interface{}{
			"username": username,
			"birthday": "2000-05-19T22:18:30.951Z",
			"context":  0,
		})

		req, err := http.NewRequest("POST", "https://auth.roblox.com/v1/usernames/validate", bytes.NewBuffer(reqBody))
		if err != nil {
			return username
		}
		req.Header.Set("authority", "auth.roblox.com")
		req.Header.Set("accept", "application/json, text/plain, */*")
		req.Header.Set("accept-language", "en-US,en;q=0.7")
		req.Header.Set("content-type", "application/json;charset=UTF-8")
		req.Header.Set("cookie", "rbx-ip2=; GuestData=UserID=-1339924708; RBXImageCache=timg=2Wzsq9ymwCnZlyeCIi_sQjBsd0WnGHEOUIsp_pdZAEs53wrsiHByEi3UjIO3UC9LCTaSk_J2jJ1mZbd6vGAU_uP9IVjdZwTuM7LnP_yxV76d480affDjFDPLgj0z3VObJldOL0B-JpA309qjBp_rtOIvdrFeSlW0bDeLqa_UNlFq0Bgcmu50FewCx0XX7r8TRd3BsNufV94O48eNMMspBg; RBXEventTrackerV2=CreateDate=8/23/2023 9:56:10 PM&rbxid=&browserid=186064096094")
		req.Header.Set("origin", "https://www.roblox.com")
		req.Header.Set("referer", "https://www.roblox.com/")
		req.Header.Set("sec-ch-ua", `"Chromium";v="116", "Not)A;Brand";v="24", "Brave";v="116"`)
		req.Header.Set("sec-ch-ua-mobile", "?0")
		req.Header.Set("sec-ch-ua-platform", `"Windows"`)
		req.Header.Set("sec-fetch-dest", "empty")
		req.Header.Set("sec-fetch-mode", "cors")
		req.Header.Set("sec-fetch-site", "same-site")
		req.Header.Set("sec-gpc", "1")
		req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36")
		res, err := p.Client.Do(req)
		if err != nil {
			return username
		}

		defer res.Body.Close()

		var result map[string]interface{}
		json.NewDecoder(res.Body).Decode(&result)

		if result["code"] != nil && result["code"].(float64) == 0 {
			return username
		}
		return username
	}
}

func GetCSRF(p *helpers.Proxy) string {
	var token = ""

	req, err := http.NewRequest("POST", "https://auth.roblox.com/v2/signup", nil)
	if err != nil {
		return token
	}

	res, err := p.Client.Do(req)
	if err != nil {
		return token
	}

	token = res.Header.Get("x-csrf-token")

	if token == "" {
		fmt.Println("Didn't Get CSRF-Token from the req")
		return token
	}
	return token

}

func GetCaptchaInfo(csrf string, username string, password string, p *helpers.Proxy, birthday string) (string, string) {
	var id string
	var metadata string

	data := []byte(fmt.Sprintf(`{"username":"%s","password":"%s","birthday":"%s","gender":2,"isTosAgreementBoxChecked":true,"agreementIds":["adf95b84-cd26-4a2e-9960-68183ebd6393","91b2d276-92ca-485f-b50d-c3952804cfd6"]}`, username, password, birthday))

	req, err := http.NewRequest("POST", "https://auth.roblox.com/v2/signup", bytes.NewBuffer(data))
	if err != nil {
		return id, metadata
	}

	req.Header.Set("authority", "auth.roblox.com")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "en-US,en;q=0.7")
	req.Header.Set("content-type", "application/json;charset=UTF-8")
	req.Header.Set("cookie", "rbx-ip2=; GuestData=UserID=-1339924708; RBXImageCache=timg=2Wzsq9ymwCnZlyeCIi_sQjBsd0WnGHEOUIsp_pdZAEs53wrsiHByEi3UjIO3UC9LCTaSk_J2jJ1mZbd6vGAU_uP9IVjdZwTuM7LnP_yxV76d480affDjFDPLgj0z3VObJldOL0B-JpA309qjBp_rtOIvdrFeSlW0bDeLqa_UNlFq0Bgcmu50FewCx0XX7r8TRd3BsNufV94O48eNMMspBg; RBXEventTrackerV2=CreateDate=8/23/2023 9:56:10 PM&rbxid=&browserid=186064096094")
	req.Header.Set("origin", "https://www.roblox.com")
	req.Header.Set("referer", "https://www.roblox.com/")
	req.Header.Set("sec-ch-ua", `"Chromium";v="116", "Not)A;Brand";v="24", "Brave";v="116"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("sec-gpc", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36")
	req.Header.Set("x-csrf-token", csrf)

	res, err := p.Client.Do(req)

	if err != nil {
		return id, metadata
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	id = res.Header.Get("rblx-challenge-id")
	metadata = res.Header.Get("rblx-challenge-metadata")

	return id, metadata

}

func AccountGen(username string, password string, p *helpers.Proxy, birthday string) string {
	var (
		chalMetadataJSON ChallengeMetadata
		req              *http.Request
		res              *http.Response
		err              error
	)
	file, err := os.OpenFile("cookies.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return ""
	}
	csrf := GetCSRF(p)
	//logWithTime(" Got Csrf Token --> " + csrf)

	id, metadata := GetCaptchaInfo(csrf, username, password, p, birthday)
	chalMetadataDecoded, err := b64.StdEncoding.DecodeString(metadata)
	if err != nil {
		return ""
	}

	err = json.Unmarshal(chalMetadataDecoded, &chalMetadataJSON)
	if err != nil {
		return ""
	}

	//logWithTime(" Captcha Info --> " + id)
	solution := helpers.CapBypassSolver(chalMetadataJSON.DataExchangeBlob)

	if solution == "" {
		return ""
	}
	payload, err := json.Marshal(CaptchaSolved{
		UnifiedCaptchaID: chalMetadataJSON.UnifiedCaptchaID,
		CaptchaToken:     solution,
		ActionType:       chalMetadataJSON.ActionType,
	})

	if err != nil {
		fmt.Println(err)
		return ""
	}

	payloadB64 := b64.StdEncoding.EncodeToString(payload)
	authPayload := []byte(fmt.Sprintf(`{"username":"%s","password":"%s","birthday":"%s","gender":2,"isTosAgreementBoxChecked":true,"agreementIds":["adf95b84-cd26-4a2e-9960-68183ebd6393","91b2d276-92ca-485f-b50d-c3952804cfd6"]}`, username, password, birthday))
	req, err = http.NewRequest("POST", "https://auth.roblox.com/v2/signup", bytes.NewBuffer(authPayload))

	if err != nil {
		return ""
	}

	req.Header.Set("authority", "auth.roblox.com")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "en-US,en;q=0.7")
	req.Header.Set("content-type", "application/json;charset=UTF-8")
	req.Header.Set("cookie", "rbx-ip2=; GuestData=UserID=-1339924708; RBXImageCache=timg=2Wzsq9ymwCnZlyeCIi_sQjBsd0WnGHEOUIsp_pdZAEs53wrsiHByEi3UjIO3UC9LCTaSk_J2jJ1mZbd6vGAU_uP9IVjdZwTuM7LnP_yxV76d480affDjFDPLgj0z3VObJldOL0B-JpA309qjBp_rtOIvdrFeSlW0bDeLqa_UNlFq0Bgcmu50FewCx0XX7r8TRd3BsNufV94O48eNMMspBg; RBXEventTrackerV2=CreateDate=8/23/2023 9:56:10 PM&rbxid=&browserid=186064096094")
	req.Header.Set("origin", "https://www.roblox.com")
	req.Header.Set("rblx-challenge-id", id)
	req.Header.Set("rblx-challenge-metadata", payloadB64)
	req.Header.Set("rblx-challenge-type", "captcha")
	req.Header.Set("referer", "https://www.roblox.com/")
	req.Header.Set("sec-ch-ua", `"Chromium";v="116", "Not)A;Brand";v="24", "Brave";v="116"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("sec-gpc", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36")
	req.Header.Set("x-csrf-token", csrf)

	res, err = p.Client.Do(req)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	//body, _ := ioutil.ReadAll(res.Body)
	//fmt.Printf("%s\n", body)

	if res.StatusCode == 200 {
		for _, cookie := range res.Cookies() {
			if cookie.Name == ".ROBLOSECURITY" {
				cookieString := fmt.Sprintf("%s:%s:%s", username, password, cookie.Value)
				logWithTime("<green> SUCCESS</> Made A Account --> <lightBlue>" + username + "</>")
				accounts += 1
				file.WriteString(cookieString + "\n")
				file.Close()
				return cookieString
			}
		}
	}

	return "Hi"
}

func round(value float64) int {
	if value < 0 {
		return int(value - 0.5)
	}
	return int(value + 0.5)
}

func calculateCPM(created int, startTime time.Time) int {
	currentTime := time.Now()
	elapsedTime := currentTime.Sub(startTime)
	//elapsedMinutes := elapsedTime.Minutes()

	return int(math.Round(float64(created) / (elapsedTime.Seconds() / 60)))
}

func formatTimeDuration(duration time.Duration) string {
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func updateTitleInBackground(timestart time.Time) {

	for {
		CPM := calculateCPM(accounts, timestart)
		currentTime := time.Now()
		elapsedTime := currentTime.Sub(timestart)

		title := fmt.Sprintf("Slotth Roblox Gen | CPM: %d | Accounts Made: %d | Time: %s",
			CPM, accounts, formatTimeDuration(elapsedTime))

		if runtime.GOOS == "windows" {
			setConsoleTitle(title)
		} else {
			fmt.Printf("\033]0;%s\a", title)
		}

		time.Sleep(1 * time.Second)
	}
}

func main() {
	helpers.LoadProxies()
	helpers.LoadUsernames()
	timestart := time.Now()

	//var work int
	var workers int
	fmt.Print("Enter threads amount (example: 50): ")
	_, err := fmt.Scanf("%d", &workers)
	if err != nil {
		fmt.Println("input err:", err)
		return
	}

	var wg sync.WaitGroup
	ch := make(chan struct{}, workers)

	go updateTitleInBackground(timestart)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				<-ch
				proxy := helpers.GetRandomProxy()
				birthday := Birthday()
				username := helpers.GetRealisticUsername()
				usernameNew := checkUsername(proxy, username)
				password := helpers.GenerateRandomString(10)
				AccountGen(usernameNew, password, proxy, birthday)
				ch <- struct{}{}
			}
		}()
		ch <- struct{}{}
	}

	wg.Wait()
}
