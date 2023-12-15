package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const (
	timeOffsetHours  = 24
	maxMessageLength = 55
)

type Config struct {
	SlackApiURL      string `json:"slackApiURL"`
	PostMessage      string `json:"postMessage"`
	MonitoredChannel string `json:"monitoredChannel"`
	SendTo           string `json:"sendTo"`
}

func LoadConfig(filename string) (Config, error) {
	var config Config

	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}

var token string

type Response struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	ReplyCount int    `json:"reply_count"`
	Text       string `json:"text"`
	User       string `json:"user"`
}

func Env_load() {
	// 環境変数をロードしてグローバル変数 token に設定
	err := godotenv.Load("env/.env")
	if err != nil {
		fmt.Println("Error loading .env files:")
	}

	token = os.Getenv("TOKEN")
}

func main() {

	Env_load()

	messages, err := getMessagesFromSlack()
	if err != nil {
		fmt.Println("Error getting messages:", err)
		return
	}

	messageList := filterMessages(messages)

	if len(messageList) > 0 {
		err := postMessageToSlack(messageList)
		if err != nil {
			fmt.Println("Error posting message:", err)
			return
		}
	} else {
		fmt.Println("no messages")
	}
}

func getMessagesFromSlack() ([]Message, error) {

	config, err := LoadConfig("settings/setting.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	referenceTime := time.Now().Add(-timeOffsetHours * time.Hour).Unix()

	params := url.Values{}
	params.Add("channel", config.MonitoredChannel)
	params.Add("oldest", fmt.Sprintf("%v", referenceTime))

	request, err := http.NewRequest("GET", config.SlackApiURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var result Response
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Messages, nil
}

func filterMessages(messages []Message) []string {
	var messageList []string
	for _, message := range messages {
		if message.ReplyCount == 0 {
			if message.Text != "" {
				messageList = append(messageList, " > "+message.Text)
			}
		}
	}
	return messageList
}

func postMessageToSlack(messageList []string) error {

	config, err := LoadConfig("setting.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	outputContext := "*<自動送信>  返信の無い質問があります* \n"
	for i, message := range messageList {
		trimmedMessage := message
		if len(message) > 50 {
			trimmedMessage = "(" + strconv.Itoa(i+1) + ") > " + message[:maxMessageLength] + "...(以下略)"
		}
		outputContext += trimmedMessage + "\n"
	}

	data := url.Values{}
	data.Set("token", token)
	data.Set("channel", config.SendTo)
	data.Set("text", outputContext)

	req, err := http.NewRequest("POST", config.PostMessage, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("Message sent. Status Code:", resp.StatusCode)
	return nil
}
