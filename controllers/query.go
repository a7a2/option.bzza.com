package controllers

import (
	//	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"net/url"

	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	kdbgo "github.com/sv/kdbgo"

	"option.bzza.com/models"
	kdb "option.bzza.com/models/kdb"
	"option.bzza.com/models/redis"
)

//{"token":"FC61884137334D44BE8598C05949BBBF","commandCode":"2017","symbolName":"EURUSD","period":"54","startDate":"1522416100","finishDate":"1522416111"}:
type queryResp struct {
	MsgType      string          `json:"msgType"`
	SymbolName   string          `json:"symbolName"`
	Period       string          `json:"period"`
	Description  string          `json:"description"`
	ArchiveArray *[]archiveArray `json:"archiveArray"`
	MsgResult    string          `json:"msgResult"`
}

type archiveArray struct {
	Open   string `json:"open"`
	Close  string `json:"close"`
	High   string `json:"high"`
	Low    string `json:"low"`
	Volume string `json:"volume"`
	Date   string `json:"date"`
}

type queryReqDate struct {
	Token       string `json:"token" formam:"token"`
	CommandCode string `json:"commandCode" formam:"commandCode"`
	SymbolName  string `json:"symbolName" formam:"symbolName"`
	Period      string `json:"period" formam:"period"`
	StartDate   string `json:"startDate" formam:"startDate"`
	FinishDate  string `json:"finishDate" formam:"finishDate"`
	CandleCount string `json:"candleCount" formam:"candleCount"`
}

type queryReqCandleCount struct {
	Token       string `json:"token"`
	CommandCode string `json:"commandCode"`
	SymbolName  string `json:"symbolName"`
	Period      string `json:"period"`
	CandleCount string `json:"candleCount"`
}

type queryMetod struct {
	It        interface{}
	C         *gin.Context
	QueryResp *queryResp
	Body      *[]byte
}

func (q *queryMetod) redisGetQueryRespAddr(bodyMd5 string) (err error) {
	//先从redis取得输出的缓存
	return nil
	c := redis.Cache{
		Key:   "query",
		Field: bodyMd5,
		It:    q.It.(*queryResp),
		IsArr: false,
	}
	if err = c.RedisHGet2Struct(); err == nil {
		q.C.JSON(http.StatusOK, c.It)
		return
	}
	return
}

func (q *queryMetod) redisSetQueryRespAddr(bodyMd5 string) (err error) {
	return
}

func (q *queryMetod) queryReqCandleCount(mapQuery map[string]string) (err error) {
	symbolName := mapQuery["symbolName"]

	startTime := time.Now().In(time.Local)
	endTime := time.Now().Add(time.Minute * -5).In(time.Local)
	setpTime := endTime

	q.selectKdb(mapQuery, startTime, setpTime, symbolName)
	return
}

func (q *queryMetod) selectKdb(mapQuery map[string]string, startTime, setpTime time.Time, symbolName string) {
	var td time.Duration
	switch mapQuery["period"] {
	case "50": //1m
		td = time.Minute
	case "51": //15m
		td = 15 * time.Minute
	case "52": //1h
		td = time.Hour
	case "53": //1d
		td = time.Hour * 24
	case "54": //5s
		td = 5 * time.Second
	default:
		td = time.Minute * 10
	}

	var err error
	var ex time.Duration

	for startTime.Sub(setpTime) > 0 {
		setpTime = setpTime.Add(td)
		kdbStr := fmt.Sprintf("select open:first bid,close:last bid,high:max bid,low:min bid from data where sym=`%s, date>%s, date<%s", symbolName, setpTime.Add(-td).Format("2006.01.02T15:04:05z"), setpTime.Format("2006.01.02T15:04:05z"))

		if time.Now().Sub(setpTime) > td {
			ex = time.Hour * 24 * 30
		} else {
			ex = td
			fmt.Println(td, " kdbStr:", kdbStr)
		}
		dataK := new(kdbgo.K)
		//var buf bytes.Buffer
		//		cc := redis.Cache{
		//			Key:    "selectKdb",
		//			Field:  kdbStr,
		//			It:     dataK,
		//			IsArr:  false,
		//			Expire: ex,
		//			Buf:    buf,
		//		}

		if err = redis.Get(kdbStr, dataK); err != nil {
			dataK, err = kdb.Call(kdbStr)
			if err != nil {
				continue
			}
			redis.Store(kdbStr, dataK, ex)
		}

		if _, ok := dataK.Data.(kdbgo.Table); ok {
			q.getOCHL(dataK.Data.(kdbgo.Table), setpTime)
		}

	}
}

func (q *queryMetod) getOCHL(tbl kdbgo.Table, setpTime time.Time) {
	ncols := len(tbl.Data)
	if len(tbl.Data) == 0 {
		return
	}
	nrows := int(tbl.Data[0].Len())

	archive := archiveArray{}
	throwAway := false
	for i := 0; i < nrows; i++ {
		for j := 0; j < ncols; j++ {
			f64, ok := tbl.Data[j].Index(i).(float64)
			if !ok || math.IsInf(f64, 0) || math.IsNaN(f64) {
				throwAway = true
				i = nrows + 2 //jump out
				//fmt.Println("jump out")
				break
				//fmt.Printf("%v %v\t", reflect.TypeOf(tbl.Data[j].Index(i)), tbl.Data[j].Index(i).(float32))
			} else {
				switch j {
				case 0:
					archive.Open = strconv.FormatFloat(f64, 'f', -1, 64)
				case 1:
					archive.Close = strconv.FormatFloat(f64, 'f', -1, 64)
				case 2:
					archive.High = strconv.FormatFloat(f64, 'f', -1, 64)
				case 3:
					archive.Low = strconv.FormatFloat(f64, 'f', -1, 64)
				}
			}
		}
	}

	if throwAway == false {
		archive.Date = strconv.FormatInt(setpTime.Unix(), 10)
		archive.Volume = strconv.Itoa(rand.Intn(1000))
		*q.It.(*queryResp).ArchiveArray = append(*q.It.(*queryResp).ArchiveArray, archive)
	}

	return
}

func (q *queryMetod) getQuoteArchiveStartEnd(mapQuery map[string]string) {
	symbolName, startDateStr, endDateStr := mapQuery["symbolName"], mapQuery["startDate"], mapQuery["finishDate"]
	startDateI64, _ := strconv.ParseInt(startDateStr, 10, 64)
	endDateI64, _ := strconv.ParseInt(endDateStr, 10, 64)

	startTime := time.Unix(startDateI64, 0).In(time.Local)
	if startDateI64 > time.Now().Unix() {
		startTime = time.Now()
	}

	endTime := time.Unix(endDateI64, 0).In(time.Local)
	setpTime := endTime

	q.selectKdb(mapQuery, startTime, setpTime, symbolName)

	return
}

func Query(c *gin.Context) { //{"token":"F8F3F817FE8742CEB21DF370888979F5","commandCode":"2018","symbolName":"EURUSD","period":"54","candleCount":"500"}:
	b, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	//var token string //, period, symbolName, startDate, finishDate, commandCode, candleCount
	mapQuery := make(map[string]string)
	if err = simpleDecodeUnmarshal(b, &mapQuery); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	//验证token
	user := &models.Users{}

	strCmd := redis.RedisClient.Get(mapQuery["token"])
	if strCmd.Err() != nil {
		fmt.Println(strCmd.Err())
		return
	}
	token := strCmd.Val()
	if err = user.FindUsersByUuid(token); err != nil {
		fmt.Println(strCmd.Err())
		return
	}

	aArr := new([]archiveArray)
	q := queryResp{
		MsgType:      "quoteArchive",
		SymbolName:   mapQuery["symbolName"],
		Period:       mapQuery["period"],
		Description:  mapQuery["symbolName"],
		MsgResult:    "Success",
		ArchiveArray: aArr,
	}
	qMetod := &queryMetod{}
	switch {
	case mapQuery["startDate"] != "":
		qMetod = &queryMetod{
			It: &q,
			C:  c,
			//Body: &b,
		}
		qMetod.getQuoteArchiveStartEnd(mapQuery)

	case mapQuery["candleCount"] != "":
		qMetod = &queryMetod{
			It: &q,
			C:  c,
			//Body: &b,
		}
		qMetod.queryReqCandleCount(mapQuery)
	}

	if *qMetod.It.(*queryResp).ArchiveArray == nil { //Archive is not found
		qMetod.It.(*queryResp).MsgResult = "Archive is not found"
	}

	c.JSON(http.StatusOK, qMetod.It)
	return
}

func simpleDecodeUnmarshal(b []byte, it interface{}) (err error) {
	strQuery, err := url.QueryUnescape(string(b))
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(strQuery), it)
	return
}
