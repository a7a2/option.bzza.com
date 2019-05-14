package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"time"

	"option.bzza.com/models/redis"
	"option.bzza.com/system"
)

type optionSettingsArrayString struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	SymbolName   string `json:"symbolName"`
	Type         string `json:"type"`
	IsUseSetting string `json:"isUseSetting"`
	OptionParams string `json:"optionParams"`
}

type optionSettingsArray struct {
	Id           string          `json:"id"`
	Name         string          `json:"name"`
	SymbolName   string          `json:"symbolName"`
	Type         string          `json:"type"`
	IsUseSetting string          `json:"isUseSetting"`
	OptionParams *[]optionParams `json:"optionParams"`
}

type optionParams struct {
	Id               string `json:"id"`
	Period           string `json:"period"`
	StopLine         string `json:"stopLine"`
	PayoutPercentage string `json:"payoutPercentage"`
	EarlyClosing     string `json:"earlyClosing"`
	Interval         string `json:"interval"`
}

var optionSettingsArrS []optionSettingsArrayString

var osaSyncMap *sync.Map

func LoopGetOptionSettingsFromRedis() {
	for {
		select {
		case <-time.After(time.Second * 5):
			//OptionSettingsArrayString ---star
			cc := &redis.Cache{
				Key:    "OptionSettingsArrayString",
				Field:  "0",
				It:     &optionSettingsArrS,
				Expire: time.Hour * 365 * 24 * 10,
				//Buf
				IsArr: false,
			}
			err := cc.RedisHGet2Struct()
			if err != nil {
				fmt.Println(err)
			}

			//OptionSettingsArray ---star
			ssmc := redis.RedisClient.HGetAll("OptionSettingsArray")
			if ssmc.Err() != nil {
				fmt.Println(ssmc.Err())
				return
			}
			for k, _ := range ssmc.Val() {
				var osa optionSettingsArray
				cc := &redis.Cache{
					Key:   "OptionSettingsArray",
					Field: k,
					It:    &osa,
					IsArr: false,
				}
				err = cc.RedisHGet2Struct()
				if err != nil {
					fmt.Println(err)
					continue
				}

				osaSyncMap.Store(cc.It.(*optionSettingsArray).Type+cc.It.(*optionSettingsArray).SymbolName, cc.It.(*optionSettingsArray))
			}
		}
	}
}

func OptionSettingsArrayInit() { //把OptionSettingsArray。json配置文件读出并结构化到redis
	for {
		select {
		case <-system.ChanOptionSettingsArrayInit:
			if system.Conf.Web.IsMasterServer != true { //主服务器才能修改json配置文件
				fmt.Println("主服务器才能修改json配置文件")
				return
			}

			//OptionSettingsArrayString ---start
			bytesAddr := system.ReadFile("./optionSettingsArray.json")
			err := json.Unmarshal(*bytesAddr, &optionSettingsArrS)
			if err != nil {
				log.Println(err)
				return
			}
			cc := &redis.Cache{
				Key:    "OptionSettingsArrayString",
				Field:  "0",
				It:     &optionSettingsArrS,
				Expire: time.Hour * 365 * 24 * 10,
				//Buf
				IsArr: false,
			}
			err = cc.Struct2RedisHSet()
			if err != nil {
				panic(err)
			}

			//OptionSettingsArray ----start
			var newOptionParams string
			opArr := []optionParams{}
			for _, v := range optionSettingsArrS {
				newOptionParams = strings.Replace(v.OptionParams, "\\", "", -1)
				err := json.Unmarshal([]byte(newOptionParams), &opArr)
				if err != nil {
					log.Println(err)
					return
				}

				o := optionSettingsArray{
					Id:           v.Id,
					Name:         v.Name,
					SymbolName:   v.SymbolName,
					Type:         v.Type,
					IsUseSetting: v.IsUseSetting,
					OptionParams: &opArr,
				}

				cc := &redis.Cache{
					Key:    "OptionSettingsArray",
					Field:  v.Type + v.SymbolName,
					It:     &o,
					Expire: time.Hour * 365 * 24 * 10,
					//Buf
					IsArr: false,
				}
				err = cc.Struct2RedisHSet()
				if err != nil {
					panic(err)
				}
			}
		}
	}

}
