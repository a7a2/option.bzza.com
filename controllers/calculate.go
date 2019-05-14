package controllers

import (
	"fmt"
	"math"
	//	"math"
	"strconv"
	"strings"
	"time"

	"option.bzza.com/models"

	"option.bzza.com/models/redis"
)

var redisLockClearUpKey string = "redisLockClearUpKey"

func calLoop() { //结算 开奖
	for {
		select {
		case <-time.After(time.Second * 5):
			statusCmd := redis.RedisClient.SetNX(redisLockClearUpKey, "lock", time.Second)
			if statusCmd.Val() {
				tArr, err := calRedisHgetAll()
				if err != nil {
					redis.RedisClient.Del(redisLockClearUpKey)
					continue
				}
				calStepOne(&tArr)
			}
			redis.RedisClient.Del(redisLockClearUpKey)
		}
	}
}

func calStepClearDirection0123(v *models.Trade, ex int64) bool {
	var f64 float64
	if f64, _ = checkPriceNow(v.SymbolName, v.OpenPrice, time.Unix(ex, 0)); f64 == 0.0 {
		return false
	}

	switch v.Direction {
	case 0: //升
		if f64 > v.OpenPrice { //赢了 实时报价大于买入价格
			v.Profit = v.PayoutPercentage / 100 * v.InvestmentSum
		} else {
			v.Profit = 0 - v.InvestmentSum
		}
	case 1: //跌
		if f64 < v.OpenPrice {
			v.Profit = v.PayoutPercentage / 100 * v.InvestmentSum
		} else {
			v.Profit = 0 - v.InvestmentSum
		}
	case 2: //升
		if f64 > v.OpenPrice {
			v.Profit = v.PayoutPercentage / 100 * v.InvestmentSum
		} else {
			v.Profit = 0 - v.InvestmentSum
		}
	case 3: //跌
		if f64 < v.OpenPrice {
			v.Profit = v.PayoutPercentage / 100 * v.InvestmentSum
		} else {
			v.Profit = 0 - v.InvestmentSum
		}
	}
	v.ClosePrice = f64
	return true
}

func decimalToNum(f64 float64, n int) float64 { //1.2345 to 12345
	return f64 * math.Pow10(n)
}

var mapMarkSymbolMaxNum = make(map[string]int) //用于记录交易对的报价的小数点后最大有多少位

func getDecimalNum(f64 float64) int { //查找小数点后有几位1.2345 返回4
	s := strconv.FormatFloat(f64, 'f', -1, 64)
	sArr := strings.Split(s, ".")
	if len(sArr) == 2 {
		return len(sArr[1])
	}
	return 0
}

func returnMax(symbol string, nums ...int) { //返回int最大值 到map[symbol]
	for _, num := range nums {
		if num > mapMarkSymbolMaxNum[symbol] {
			mapMarkSymbolMaxNum[symbol] = num
		}
	}
	return
}

func calStepClearDirection4567(v *models.Trade, openDateInt64, exInt64 int64) (close bool) {
	var highPrice, lowPrice, lastPrice float64
	var ok bool
	if highPrice, lowPrice, lastPrice, ok = checkHighLowPrice(v.SymbolName, time.Unix(openDateInt64, 0), time.Unix(exInt64, 0)); ok == false {
		return
	}

	returnMax(v.SymbolName, getDecimalNum(highPrice), getDecimalNum(lowPrice), getDecimalNum(v.OpenPrice))
	if mapMarkSymbolMaxNum[v.SymbolName] <= 0 {
		close = false
		return
	}

	openPriceInt := decimalToNum(v.OpenPrice, mapMarkSymbolMaxNum[v.SymbolName])
	highPriceInt := decimalToNum(highPrice, mapMarkSymbolMaxNum[v.SymbolName]) //时间段内最高价int版
	lowPriceInt := decimalToNum(lowPrice, mapMarkSymbolMaxNum[v.SymbolName])   //时间段内最低价int版
	upOpenPriceInt := openPriceInt + float64(v.Interval)                       //4
	downOpenPriceInt := openPriceInt - float64(v.Interval)                     //5
	fmt.Println(exInt64)
	fmt.Println("v.OpenPrice=", v.OpenPrice)
	fmt.Println("highPrice=", highPrice, " 	lowPrice=", lowPrice)
	fmt.Println("highPriceInt=", highPriceInt, " 	lowPriceInt=", lowPriceInt)
	fmt.Println("upOpenPriceInt=", upOpenPriceInt, " 	downOpenPriceInt=", downOpenPriceInt)
	close = false
	switch v.Direction { //4升 、5降，只要指定时间到价格触及目标价格即赢。6范围内、7范围外
	case 4:
		if upOpenPriceInt <= highPriceInt {
			close = true
			v.ClosePrice = highPrice
			v.Profit = v.PayoutPercentage / 100 * v.InvestmentSum
		} else if time.Now().Unix() > exInt64 {
			close = true
			v.ClosePrice = highPrice
			v.Profit = 0 - v.InvestmentSum
		}
	case 5:
		if downOpenPriceInt >= lowPriceInt {
			close = true
			v.ClosePrice = lowPrice
			v.Profit = v.PayoutPercentage / 100 * v.InvestmentSum
		} else if time.Now().Unix() > exInt64 {
			close = true
			v.ClosePrice = lowPrice
			v.Profit = 0 - v.InvestmentSum
		}
	case 7: //范围内
		//实际最高价 大于 指定范围内上线，已经不在范围内，判输.实际最低价 小于 指定范围内下线 ，不再范围内，判输。可以结算
		if upOpenPriceInt < highPriceInt {
			fmt.Println(upOpenPriceInt-highPriceInt, " ", downOpenPriceInt-lowPriceInt)
			fmt.Println(time.Now().Unix(), " ", exInt64)
			close = true
			v.ClosePrice = highPrice
			v.Profit = 0 - v.InvestmentSum
		} else if downOpenPriceInt > lowPriceInt {
			close = true
			v.ClosePrice = lowPrice
			v.Profit = 0 - v.InvestmentSum
		} else if time.Now().Unix() > exInt64 {
			close = true
			v.ClosePrice = lastPrice
			v.Profit = v.PayoutPercentage / 100 * v.InvestmentSum
		}
	case 6:
		//实际最高价 大于 指定范围内上线，在范围外，判赢.实际最低价 小于 指定范围内下线 ，判赢。可以结算
		if upOpenPriceInt < highPriceInt {
			v.Profit = v.PayoutPercentage / 100 * v.InvestmentSum
			v.ClosePrice = highPrice
			close = true
		} else if downOpenPriceInt > lowPriceInt {
			v.Profit = v.PayoutPercentage / 100 * v.InvestmentSum
			v.ClosePrice = lowPrice
			close = true
		} else if time.Now().Unix() > exInt64 {
			close = true
			v.ClosePrice = lastPrice
			v.Profit = 0 - v.InvestmentSum
		}
	}

	return
}

func calStepOne(tArr *[]models.Trade) {
	for _, v := range *tArr {
		ex := v.OpenDate + int64(v.Period)

		switch {
		case v.Close:
			continue
		case (v.Direction == 0 || v.Direction == 1 || v.Direction == 2 || v.Direction == 3) && ex < time.Now().Unix():
			if calStepClearDirection0123(&v, ex) == false {
				continue
			}
		case (v.Direction == 4 || v.Direction == 5 || v.Direction == 6 || v.Direction == 7):
			if calStepClearDirection4567(&v, v.OpenDate, ex) == false {
				continue
			}
		default:
			continue
		}

		v.Close = true

		if time.Now().Unix() > ex {
			v.CloseDate = ex
		} else {
			v.CloseDate = time.Now().Unix()
		}

		u := &models.Users{
			Uuid: v.UserUuid,
		}
		if err := u.FindUsersByUuid(v.UserUuid.String()); err != nil {
			return
		}

		fmt.Println(v.Profit, " ", fmt.Sprintf("sumOptionDeals+%v", v.Profit))
		redis.RedisClient.HDel("valuation_clear", v.Uuid.String()) //结算完删除redis上的单
		//v.Update() //交易更新到库
		u.SumOptionDeals += v.Profit
		if v.Profit > 0.0 {
			u.Balance = u.Balance + v.Profit + v.InvestmentSum
			v.Balance = v.Balance + v.Profit + v.InvestmentSum
			//u.UpdateBalance(v.Profit) //用户更新到库
			//v.Profit -= v.InvestmentSum
		}

		models.Clear(&v, u)
		unicast(u) //结算后推送到用户端

	}
}

func calRedisHgetAll() (t []models.Trade, err error) {
	ssmc := redis.RedisClient.HGetAll("valuation_clear")
	if ssmc.Err() != nil {
		err = ssmc.Err()
		return
	}

	for k, _ := range ssmc.Val() {
		cc := &redis.Cache{
			Key:   "valuation_clear",
			Field: k,
			It:    &models.Trade{},
			//Expire:time.Second*10,
			//Buf
			IsArr: false,
		}
		err = cc.RedisHGet2Struct()
		if err != nil {
			fmt.Println(err)
			continue
		}

		t = append(t, *cc.It.(*models.Trade))
	}

	return

}
