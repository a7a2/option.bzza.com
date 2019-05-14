package helpers

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	//	"errors"
	//	"os"
	"strings"

	"log"
	"net/http"

	"crypto/hmac"
	"crypto/sha256"
	"io"

	"math/rand"

	"strconv"
	"time"

	//	"github.com/pzduniak/argon2"
	"golang.org/x/crypto/scrypt"
)

func EncryptPwd2Db(password []byte) string { //注册密码由这里开始加密
	b, err := scrypt.Key(password, salt, 16384, 8, 1, 32)
	if err != nil {
		panic(err)
	}
	return MD5(b)
}

var salt []byte = []byte("a89fdefdaa501f0a16af38e6fc708f1c") //用于argon2加密的salt
//func argon2(text string) (strHash string) { //argon2加密
//	hashByte, _ := argon2.Key([]byte(text), salt, 10, 4, 16384, 32, argon2.Argon2i)
//	strHash = hex.EncodeToString(hashByte)
//	return
//}

// 生成32位MD5
func MD5(b []byte) string {
	ctx := md5.New()
	ctx.Write(b)
	return hex.EncodeToString(ctx.Sum(nil))
}

//生成随机字符串
func GetRandomString(i64 int64) string {
	str := "123456789abcdefghijklmnpqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ" //移除0、o、O
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	intB, err := strconv.Atoi(strconv.FormatInt(int64(len(bytes)), 10))
	if err != nil {
		return ""
	}

	for i := int64(0); i < i64; i++ {
		result = append(result, bytes[r.Intn(intB)])
	}
	return string(result)
}

//get HmacSHA256
func GetHmacCode(s string, key []byte) string {
	h := hmac.New(sha256.New, key)
	io.WriteString(h, s)
	return string(h.Sum(nil))
}

//base64 Encode
func Base64Encode(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

func HttpRequst(Url, metod, reqData string) ([]byte, error) {
	reader := strings.NewReader(reqData) //reqData= "a=bhh"

	//client := &http.Client{}
	req := new(http.Request)
	var err error
	req, err = http.NewRequest(metod, Url, reader)
	if err != nil {
		log.Println("http.NewRequest:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8") //application/json; charset=utf-8
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36")
	req.Header.Set("Referer", "http://binary.utip.org/")
	req.Header.Set("Origin", "http://binary.utip.org")
	req.Header.Set("Host", "binary.utip.org:8085")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,th;q=0.7")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Add("dnt", "1")
	req.Header.Add("cache-control", "no-cache")
	//	req.Header.Add("connection", "keep-alive")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("client.Do:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var bytes []byte
	bytes, err = ioutil.ReadAll(resp.Body)
	return bytes, err
}
