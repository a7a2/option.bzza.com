package controllers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	//	"reflect"
	//	"reflect"
	"strconv"
	"time"

	"github.com/aglyzov/ws-machine"
	kdbgo "github.com/sv/kdbgo"
	"option.bzza.com/models"
	kdb "option.bzza.com/models/kdb"
	"option.bzza.com/system"
)

var wsORwss string

func WsGetSymbolsFromA7a2() {
	wsORwss = system.Conf.Datasource.A7a2.WsORwss
	wsm := machine.New(wsORwss, http.Header{
		"Authorization:APPCODE": {"aa685fe9b31643e4aaf139fa3df9bc5c"},
		"Content-Type":          {"application/json; charset=utf-8"},
		"gateway_channel":       {"websocket"},
	},
	)

	//	var err error
	//	var msgStr string
	var bytes []byte

	for {
		select {
		case st := <-wsm.Status:
			fmt.Println("STATE:", st.State)
			if st.Error != nil {
				fmt.Println("ERROR:", st.Error)
			}
			if st.State == machine.CONNECTED {
				//fmt.Println("SUBSCRIBE: live_trades")
				wsm.Output <- []byte(``)
			}
		case bytes = <-wsm.Input:
			unmarshalSymbolsAndBroadcast(&bytes) //不能用go,单发

		}
	}
}

func unmarshalSymbolsAndBroadcast(bytes *[]byte) {
	var getSymbols models.GetSymbols
	err := json.Unmarshal(*bytes, &getSymbols)
	if err != nil {
		fmt.Println("unmarshalSymbols:json.Unmarshal:", err, "  ", string(*bytes))
		return
	}

	if getSymbols.Data.Bid < 0.0 || getSymbols.Data.Ask < 0.0 {
		return
	}
	qDetails := &quoteDetails{
		Symbol: getSymbols.Data.Symbol,
		Bid:    strconv.FormatFloat(getSymbols.Data.Bid, 'f', -1, 64),
		Ask:    strconv.FormatFloat(getSymbols.Data.Ask, 'f', -1, 64),
		Date:   strconv.FormatInt(getSymbols.Data.Unix, 10),
	}

	d := kdb.Data{
		Symbol: getSymbols.Data.Symbol,
		Bid:    qDetails.Bid,
		Ask:    qDetails.Ask,
		Date:   time.Unix(getSymbols.Data.Unix, 0).Format("2006.01.02T15:04:05z"),
	}
	go d.InsertIfNotExist()
	SymbolBroadcast(qDetails)

	//kdbStr := fmt.Sprintf("`data insert (%s;`%s;%vf;%vf)", time.Unix(getSymbols.Data.Unix, 0).Format("2006.01.02T15:04:05z"), getSymbols.Data.Symbol, getSymbols.Data.Bid, getSymbols.Data.Ask)
	//	fmt.Println(kdbStr)

	//go kdb.Call(kdbStr)

}

func checkPriceNow(symbol string, openPrice float64, t time.Time) (float64, bool) { //匹配最近30秒的3条数据,相等即返回true，有数据返回最近一条float64
	//for i := 0; i <= 1; i++ {
	kdbStr := fmt.Sprintf("select [3;>date]bid from data where sym=`%s,date<=%s,date>%s", symbol, t.Format("2006.01.02T15:04:05z"), t.Add(-time.Second*30).Format("2006.01.02T15:04:05z"))
	//fmt.Println(kdbStr)
	data, err := kdb.Call(kdbStr)
	if err != nil {
		return 0.0, false
	}
	var f64 float64 = 0.0

	if _, ok := data.Data.(kdbgo.Table); !ok {
		return 0.0, false
	}

	for _, v := range data.Data.(kdbgo.Table).Data {
		vv, ok := v.Data.([]float64)
		if ok {
			for _, price := range vv {
				if math.IsInf(price, 0) || math.IsNaN(price) {
					continue
				} else if f64 == 0.0 {
					f64 = price
				}

				if price == openPrice {
					return price, true
				}
			}
		}
	}

	return f64, false
}

func checkHighLowPrice(symbol string, tStart, tEnd time.Time) (highPrice, lowPrice, lastPrice float64, ok bool) { //寻找指定时间段内max、low值
	kdbStr := fmt.Sprintf("select high:max bid,low:min bid,close:last bid from data where sym=`%s,date>%s,date<%s", symbol, tStart.Format("2006.01.02T15:04:05z"), tEnd.Format("2006.01.02T15:04:05z"))
	fmt.Println(kdbStr)
	data, err := kdb.Call(kdbStr)
	if err != nil {
		return
	}

	if _, ok = data.Data.(kdbgo.Table); !ok {
		return
	}

	var vv []float64
	for _, v := range data.Data.(kdbgo.Table).Data {
		vv, ok = v.Data.([]float64)
		if ok {
			for k, price := range vv {
				fmt.Println(k, " ", price)
				if math.IsInf(price, 0) || math.IsNaN(price) {
					continue
				}

				switch {
				case highPrice <= 0.0:
					highPrice = price
				case lowPrice <= 0.0:
					lowPrice = price
				case lastPrice <= 0.0:
					lastPrice = price
				}

			}
		}
	}

	if highPrice == 0 || lowPrice == 0 || lastPrice == 0 {
		return 0.0, 0.0, 0.0, false
	}
	return
}
