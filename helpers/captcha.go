package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gookit/color"
)

func logWithTime(msg string) {
	timeStr := time.Now().Format("15:04:05")
	color.Println("[<dark_gray>" + timeStr + "</>]" + msg)
}

type SolutionResponse struct {
	Solution struct {
		Token string `json:"token"`
	} `json:"solution"`
	ErrorID int `json:"errorId"`
}

func CapBypassSolver(blob string) string {
	// logWithTime(" Sending Request to (capbypass)")

	payload := map[string]interface{}{
		"api_key":               "SLOTTH-BAAF46F254936A3A9D9F968BE68FA053",
		"site_key":              "A2A14B1D-1AF3-C791-9BBC-EE33CC7A0A6F",
		"site_iframe_url":       "https://www.roblox.com",
		"site_template":         "roblox_register",
		"arkose_api_url":        "https://roblox-api.arkoselabs.com",
		"force_old_games":       false,
		"force_audio":           false,
		"ignore_silent_pass":    false,
		"debug_challenge":       true,
		"optional_custom_proxy": "socks5://jew:jew@174.114.197.180:1080",
		"data": map[string]interface{}{
			"blob": blob,
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Json stuff had a issue %s", err)
		return ""
	}

	req, err := http.NewRequest("POST", "https://nigger.zone/cap/solve/arkose", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Printf("Couldn't Send request to CapBypass %s", err)
		return ""
	}

	req.Header.Set("Content-Type", "application/json")

	logWithTime("<magenta> INFO</> Sent a Request to CapBypass (" + blob[:23] + ")")

	client := &http.Client{Timeout: time.Duration(180) * time.Second}

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed To send a request to capbypass: %s", err)
		return ""
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	body, _ := ioutil.ReadAll(res.Body)

	fmt.Printf("%s\n", body)

	return string(body)
}
