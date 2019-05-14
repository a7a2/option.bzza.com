package controllers

import (
	"flag"
	"fmt"
	"io"

	"os"

	"log"

	"net/url"

	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"option.bzza.com/helpers"
	"option.bzza.com/system"
)

type Sign struct {
	Key          string
	Metod        string
	Accept       string
	Content_MD5  string
	Content_Type string
	Date         string
	Headers      string
	Url          string
}

var httpORhttps string

func httpGetSymbolsFromAliyun() {
	httpORhttps = system.Conf.Datasource.Aliyun.HttpORhttps

	u, _ := url.Parse(httpORhttps)
	timestampStr := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	s := Sign{
		Key:          system.Conf.Datasource.Aliyun.AppKey,
		Metod:        "GET",
		Accept:       "*/*",
		Content_Type: "application/json; charset=utf-8",
		Date:         timestampStr,
		Headers:      "X-Ca-Key,X-Ca-Timestamp",
		Url:          u.Path,
	}
	signature := s.Sign()

	client := &http.Client{}
	req, err := http.NewRequest(s.Metod, httpORhttps, nil)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Set("Content-Type", s.Content_Type)
	req.Header.Set("Accept", s.Accept)
	req.Header.Set("X-Ca-Key", s.Key)
	req.Header.Set("X-Ca-Timestamp", s.Date)
	req.Header.Set("X-Ca-Signature-Headers", s.Headers)
	req.Header.Set("X-Ca-Signature", signature)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("client.Do:", err)
		return
	}
	defer resp.Body.Close()

	//fmt.Println(resp.Header.Get("X-Ca-Error-Message"))

	var bytes []byte
	stdout := os.Stdout
	_, err = io.Copy(stdout, resp.Body)
	n, err := stdout.Read(bytes)
	fmt.Println(n, ":", string(bytes))
}

func LoopHttpGetSymbolsFromAliyun() {
	for {
		select {
		case <-time.After(time.Second * 3):
			//httpGetSymbolsFromAliyun()
			go WsGetSymbolsFromAliyun()
		}
	}

}

var addr = flag.String("addr", "realtime.market.alicloudapi.com:8080", "aliyun api ws")

func WsGetSymbolsFromAliyun() {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	////Url := "ws://realtime.market.alicloudapi.com:8080/w2/reg/finance"
	//	u, _ := url.Parse(Url)

	timestampStr := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	s := Sign{
		Key:          system.Conf.Datasource.Aliyun.AppKey,
		Metod:        "GET",
		Accept:       "*/*",
		Content_Type: "application/json; charset=utf-8",
		Date:         timestampStr,
		Headers:      "X-Ca-Key,X-Ca-Timestamp",
		Url:          u.Path,
	}
	signature := s.Sign()

	wsHeaders := http.Header{
		"Origin":                   []string{"http://realtime.market.alicloudapi.com/w2/reg/finance"},
		"Sec-WebSocket-Extensions": {"permessage-deflate; client_max_window_bits, x-webkit-deflate-frame"},
		"x-ca-deviceid":            []string{fmt.Sprintf("%s@%s", uuid.New().String(), s.Key)},
		"Accept":                   []string{s.Accept},
		"Content-Type":             []string{s.Content_Type},
		"X-Ca-Key":                 []string{s.Key},
		"X-Ca-Timestamp":           []string{s.Date},
		"X-Ca-Signature-Headers":   []string{s.Headers},
		"X-Ca-Signature":           []string{signature},
	}

	var dialer *websocket.Dialer
	conn, resp, err := dialer.Dial(u.String(), wsHeaders)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	fmt.Println("resp:", resp.Header)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}
		fmt.Printf("received: %s\n", message)
	}
}

func (s *Sign) Sign() string {
	headers := "X-Ca-Key:" + s.Key + "\n" + "X-Ca-Timestamp:" + s.Date + "\n"
	stringToSign := s.Metod + "\n" + s.Accept + "\n" + s.Content_MD5 + "\n" + s.Content_Type + "\n" + "\n" + headers + s.Url
	//fmt.Println(stringToSign)
	stringToSign = helpers.GetHmacCode(stringToSign, []byte(system.Conf.Datasource.Aliyun.AppSecret))
	stringToSign = helpers.Base64Encode([]byte(stringToSign))
	return stringToSign
}
