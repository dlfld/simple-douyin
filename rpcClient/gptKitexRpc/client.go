package gptKitexRpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// const ipPort = "http://43.130.60.218:8888"
const ipPort = "http://127.0.0.1:8888"

type ChatReqBody struct {
	Content string `json:"content"`
}

func SendChatMsg(msg string) {
	url := ipPort + "/chatCompletion"
	body := ChatReqBody{
		Content: msg,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return
	}
	fmt.Println(jsonBody)
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println("Response:", string(data))

}
