package controllers

import (
	"encoding/json"
	"errors"

	"fmt"

	"log"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	websocket "github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"

	"github.com/satori/go.uuid"
	"option.bzza.com/helpers"
	"option.bzza.com/models"
	"option.bzza.com/models/redis"
	"option.bzza.com/system"
)

var sMapWsConns = new(sync.Map)

type hub struct {
	Conn *net.Conn
	Op   websocket.OpCode
	Ok   bool
	Freq time.Duration
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//server -> {"msgType":"serverName","serverName":"Binary Demo","serverType":"0","useOptions":"1","useOnlyOptions":"1","useNewOptions":"1","msgResult":"Success"}
type commandCode2012Return struct {
	MsgType        string `json:"msgType"`        //"serverName"
	ServerName     string `json:"serverName"`     //"Binary Demo"
	ServerType     string `json:"serverType"`     //"0"
	UseOptions     string `json:"useOptions"`     //"1"
	UseOnlyOptions string `json:"useOnlyOptions"` //"1"
	UseNewOptions  string `json:"useNewOptions"`  //"1"
	MsgResult      string `json:"msgResult"`      //"Success"
}

func (this *ws) commandCode2012Return() {
	s := commandCode2012Return{
		MsgType:        "serverName",
		ServerName:     "Binary Demo",
		ServerType:     "0",
		UseOptions:     "1",
		UseOnlyOptions: "1",
		UseNewOptions:  "1",
		MsgResult:      "Success",
	}

	this.wsSend(&s)
}

//{"msgType":"traderType","isLoggedInByInvestorPassword":"0","isLoggedInFundAccount":"0"}
type commandCode2038Return struct {
	MsgType                      string `json:"msgType"`                      //traderType
	IsLoggedInByInvestorPassword string `json:"isLoggedInByInvestorPassword"` //0
	IsLoggedInFundAccount        string `json:"isLoggedInFundAccount"`        //0
}

//{"msgType":"holidays","holidaysArray":[],"msgResult":"Success"}
type commandCode2038ReturnHolidays struct {
	MsgType       string   `json:"msgType"` //holidays
	HolidaysArray []string `json:"holidaysArray"`
	MsgResult     string   `json:"msgResult"` //Success
}

func (this *ws) commandCode2038Return() {
	s := commandCode2038Return{
		MsgType: "traderType",
		IsLoggedInByInvestorPassword: "0",
		IsLoggedInFundAccount:        "0",
	}

	this.wsSend(&s)

	s1 := commandCode2038ReturnHolidays{
		MsgType:       "holidays",
		HolidaysArray: []string{},
		MsgResult:     "Success",
	}

	this.wsSend(&s1)
}

//{"msgType":"symbolGroups","symbolsGroupsArray":[{"id":"1","name":"forex","description":"forex group","tradeMode":"1;0;0;0;0;2;0;0;23;59;3;0;0;23;59;4;0;0;23;59;5;0;0;23;59;6;0;0;22;59;7;0;0;0;0;","HolidaysCount":"0"},{"id":"2","name":"cfd","description":"cfd group","tradeMode":"1;0;0;0;0;2;15;30;21;59;3;15;30;21;59;4;15;30;21;59;5;15;30;21;59;6;15;30;21;59;7;0;0;0;0;","HolidaysCount":"0"},{"id":"3","name":"metals","description":"metals group","tradeMode":"1;0;0;0;0;2;0;0;23;59;3;0;0;23;59;4;0;0;23;59;5;0;0;23;59;6;0;0;22;59;7;0;0;0;0;","HolidaysCount":"0"},{"id":"4","name":"index","description":"index group","tradeMode":"1;0;0;0;0;2;0;0;23;59;3;0;0;23;59;4;0;0;23;59;5;0;0;23;59;6;0;0;22;59;7;0;0;0;0;","HolidaysCount":"0"},{"id":"5","name":"oil","description":"oil group","tradeMode":"1;0;0;0;0;2;0;0;23;59;3;0;0;23;59;4;0;0;23;59;5;0;0;23;59;6;0;0;23;59;7;0;0;0;0;","HolidaysCount":"0"}],"msgResult":"Success"}
type commandCode2036Return struct {
	MsgType            string                                      `json:"msgType"`            //symbolGroups
	SymbolsGroupsArray *[]commandCode2036Return_SymbolsGroupsArray `json:"symbolsGroupsArray"` //symbolsGroupsArray
	MsgResult          string                                      `json:"msgResult"`          //Success
}
type commandCode2036Return_SymbolsGroupsArray struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	TradeMode     string `json:"tradeMode"`
	HolidaysCount string `json:"holidaysCount"`
}

func (this *ws) commandCode2036Return() {
	sArr := []commandCode2036Return_SymbolsGroupsArray{
		{
			Id:            "1",
			Name:          "forex",
			Description:   "forex group",
			TradeMode:     "1;0;0;0;0;2;0;0;23;59;3;0;0;23;59;4;0;0;23;59;5;0;0;23;59;6;0;0;22;59;7;0;0;0;0;",
			HolidaysCount: "0",
		},
		{
			Id:            "2",
			Name:          "cfd",
			Description:   "cfd group",
			TradeMode:     "1;0;0;0;0;2;15;30;21;59;3;15;30;21;59;4;15;30;21;59;5;15;30;21;59;6;15;30;21;59;7;0;0;0;0;",
			HolidaysCount: "0",
		},
		{
			Id:            "3",
			Name:          "metals",
			Description:   "metals group",
			TradeMode:     "1;0;0;0;0;2;0;0;23;59;3;0;0;23;59;4;0;0;23;59;5;0;0;23;59;6;0;0;22;59;7;0;0;0;0;",
			HolidaysCount: "0",
		},
		{
			Id:            "4",
			Name:          "index",
			Description:   "index group",
			TradeMode:     "1;0;0;0;0;2;0;0;23;59;3;0;0;23;59;4;0;0;23;59;5;0;0;23;59;6;0;0;22;59;7;0;0;0;0;",
			HolidaysCount: "0",
		},
		{
			Id:            "5",
			Name:          "oil",
			Description:   "oil group",
			TradeMode:     "1;0;0;0;0;2;0;0;23;59;3;0;0;23;59;4;0;0;23;59;5;0;0;23;59;6;0;0;23;59;7;0;0;0;0;",
			HolidaysCount: "0",
		},
	}

	s := commandCode2036Return{
		MsgType:            "symbolGroups",
		SymbolsGroupsArray: &sArr,
		MsgResult:          "Success",
	}
	this.wsSend(&s)

}

type commandCode2023Return struct {
	MsgType      string          `json:"msgType"` //symbols
	SymbolsArray *[]symbolsArray `json:"symbolsArray"`
	MsgResult    string          `json:"msgResult"`
}

func (this *ws) wsSend(s interface{}) (err error) {
	b := new([]byte)
	var wf websocket.Frame
	switch reflect.TypeOf(s).String() {
	case "*[]byte":
		wf = websocket.NewFrame(this.Op, true, *s.(*[]byte))
	case "[]byte":
		wf = websocket.NewFrame(this.Op, true, (*b))
	default:
		*b, err = json.Marshal(s)
		if err != nil {
			log.Println(err)
			return
		}
		wf = websocket.NewFrame(this.Op, true, (*b))
	}

	err = websocket.WriteFrame(*this.Conn, wf)
	//wsutil.WriteClientMessage(*this.Conn, this.Op, *b)
	return
}

func (this *ws) commandCode2023Return() (err error) {
	bytesAddr := system.ReadFile("./symbols.json")
	sArr := []symbolsArray{}
	err = json.Unmarshal(*bytesAddr, &sArr)
	if err != nil {
		log.Println(err)
		return
	}

	this.wsSend(&commandCode2023Return{
		MsgType:      "symbols",
		SymbolsArray: &sArr,
		MsgResult:    "Success",
	})

	return
}

type symbolsArray struct {
	SymbolName             string `json:"symbolName"`      //AUDCAD
	QuoteSymbolName        string `json:"quoteSymbolName"` //USDCAD
	QuoteSymbolLocation    string `json:"quoteSymbolLocation"`
	QuoteSymbolCurrency    string `json:"quoteSymbolCurrency"`
	Description            string `json:"description"`
	SpreadBid              string `json:"spreadBid"`
	SpreadAsk              string `json:"spreadAsk"`
	SpreadType             string `json:"spreadType"`
	SwapLong               string `json:"swapLong"`
	SwapShort              string `json:"swapShort"`
	Delay                  string `json:"delay"`
	StopLevel              string `json:"stopLevel"`
	Digits                 string `json:"digits"`
	GapLevel               string `json:"gapLevel"`
	CalculationTypeID      string `json:"calculationTypeID"`
	MarginSymbolLocation   string `json:"marginSymbolLocation"`
	MarginSymbol           string `json:"marginSymbol"`
	MarginCurrency         string `json:"marginCurrency"`
	ContractSize           string `json:"contractSize"`
	Percentage             string `json:"percentage"`
	CoverageMarginStrategy string `json:"coverageMarginStrategy"`
	SymbolGroupID          string `json:"symbolGroupID"`
	Commission             string `json:"commission"`
	Tradeforbidden         string `json:"tradeforbidden"`
}

//{"msgType":"lastQuote","lastQuotesArray":[{"symbol":"AUDSGD","bid":"1.043","ask":"1.0442","date":"1461365693"},{"symbol":"_YHOO","bid":"52.72","ask":"52.9","date":"1501871098"},{"symbol":"_DD","bid":"84.57","ask":"84.6","date":"1504216798"},{"symbol":"_ROSN","bid":"318.88","ask":"318.97","date":"1521218390"},{"symbol":"_SBER","bid":"256.79","ask":"256.9","date":"1521218398"},{"symbol":"_GAZP","bid":"140.02","ask":"140.25","date":"1521218399"},{"symbol":"_LKOH","bid":"3843.9","ask":"3845.6","date":"1521218399"},{"symbol":"RTS","bid":"1254.29","ask":"1254.3","date":"1521219060"},{"symbol":"MICEX","bid":"2294.6","ask":"2294.61","date":"1521219060"},{"symbol":"SB","bid":"12.64","ask":"12.72","date":"1521223116"},{"symbol":"EURRUB","bid":"70.649","ask":"70.779","date":"1521223197"},{"symbol":"KC","bid":"118.2","ask":"118.3","date":"1521224973"},{"symbol":"USDRUB","bid":"57.481","ask":"57.586","date":"1521226798"},{"symbol":"ZC","bid":"381.38","ask":"383.88","date":"1521227688"},{"symbol":"ZS","bid":"1049.13","ask":"1050.13","date":"1521227699"},{"symbol":"CT","bid":"82.59","ask":"82.93","date":"1521227999"},{"symbol":"_AA","bid":"47.22","ask":"47.3","date":"1521233689"},{"symbol":"_WU","bid":"20.05","ask":"20.08","date":"1521233970"},{"symbol":"_HPQ","bid":"23.49","ask":"23.52","date":"1521233971"},{"symbol":"_GE","bid":"14.28","ask":"14.31","date":"1521233975"},{"symbol":"_HAL","bid":"45.88","ask":"45.93","date":"1521233982"},{"symbol":"_PFE","bid":"36.75","ask":"36.78","date":"1521233988"},{"symbol":"_XRX","bid":"31.09","ask":"31.12","date":"1521233992"},{"symbol":"_BAC","bid":"32.14","ask":"32.17","date":"1521233993"},{"symbol":"_T","bid":"36.96","ask":"36.99","date":"1521233994"},{"symbol":"_KO","bid":"43.45","ask":"43.48","date":"1521233996"},{"symbol":"_INTC","bid":"51.14","ask":"51.15","date":"1521233996"},{"symbol":"_CSCO","bid":"44.98","ask":"45","date":"1521233996"},{"symbol":"_JPM","bid":"115.37","ask":"115.4","date":"1521233997"},{"symbol":"_PG","bid":"78.9","ask":"78.93","date":"1521233998"},{"symbol":"_IP","bid":"54.56","ask":"54.59","date":"1521233998"},{"symbol":"_VZ","bid":"48.48","ask":"48.51","date":"1521233998"},{"symbol":"_WMT","bid":"89.54","ask":"89.57","date":"1521233998"},{"symbol":"_HD","bid":"178.98","ask":"179.01","date":"1521233998"},{"symbol":"_AXP","bid":"95.58","ask":"95.61","date":"1521233998"},{"symbol":"_MO","bid":"63.15","ask":"63.18","date":"1521233998"},{"symbol":"_MCD","bid":"162.36","ask":"162.39","date":"1521233998"},{"symbol":"_JNJ","bid":"133.06","ask":"133.09","date":"1521233998"},{"symbol":"_MRK","bid":"55.63","ask":"55.66","date":"1521233998"},{"symbol":"_EBAY","bid":"42.44","ask":"42.47","date":"1521233999"},{"symbol":"_FB","bid":"184.99","ask":"185.02","date":"1521233999"},{"symbol":"_DIS","bid":"102.83","ask":"102.86","date":"1521233999"},{"symbol":"_YNDX","bid":"42.3","ask":"42.33","date":"1521233999"},{"symbol":"_PM","bid":"103.23","ask":"103.26","date":"1521233999"},{"symbol":"_UTX","bid":"128.35","ask":"128.4","date":"1521233999"},{"symbol":"_MSFT","bid":"94.44","ask":"94.47","date":"1521233999"},{"symbol":"_MCO","bid":"167.28","ask":"167.31","date":"1521233999"},{"symbol":"_MMM","bid":"237.23","ask":"237.27","date":"1521233999"},{"symbol":"_IBM","bid":"160.28","ask":"160.31","date":"1521233999"},{"symbol":"_XOM","bid":"74.92","ask":"74.95","date":"1521233999"},{"symbol":"_GOOG","bid":"1135.89","ask":"1136.19","date":"1521233999"},{"symbol":"_HON","bid":"151.66","ask":"151.69","date":"1521233999"},{"symbol":"_CVX","bid":"114.99","ask":"115.02","date":"1521233999"},{"symbol":"_AMZN","bid":"1570.9","ask":"1571.2","date":"1521233999"},{"symbol":"_AAPL","bid":"177.84","ask":"177.98","date":"1521233999"},{"symbol":"_CAT","bid":"156.4","ask":"156.43","date":"1521233999"},{"symbol":"_BA","bid":"329.6","ask":"329.82","date":"1521233999"},{"symbol":"JP225","bid":"21395","ask":"21415","date":"1521237244"},{"symbol":"US500","bid":"2755.03","ask":"2755.73","date":"1521237260"},{"symbol":"DE30","bid":"12396","ask":"12401","date":"1521237279"},{"symbol":"F40","bid":"5137","ask":"5142.5","date":"1521237282"},{"symbol":"WTI","bid":"62.21","ask":"62.28","date":"1521237286"},{"symbol":"UK100","bid":"7066.5","ask":"7070.5","date":"1521237286"},{"symbol":"USTEC","bid":"7041.9","ask":"7046.2","date":"1521237295"},{"symbol":"US30","bid":"24951","ask":"24961","date":"1521237299"},{"symbol":"DX","bid":"89.753","ask":"89.793","date":"1521237300"},{"symbol":"XBRUSD","bid":"66.12","ask":"66.15","date":"1521237377"},{"symbol":"NG","bid":"2.688","ask":"2.698","date":"1521237518"},{"symbol":"XAUUSD","bid":"1313.93","ask":"1314.24","date":"1521237570"},{"symbol":"USDSGD","bid":"1.3173","ask":"1.3181","date":"1521237571"},{"symbol":"NZDSGD","bid":"0.9505","ask":"0.9517","date":"1521237572"},{"symbol":"NZDCHF","bid":"0.68664","ask":"0.68767","date":"1521237578"},{"symbol":"GBPAUD","bid":"1.80704","ask":"1.80922","date":"1521237580"},{"symbol":"EURCAD","bid":"1.60953","ask":"1.60987","date":"1521237591"},{"symbol":"GBPCHF","bid":"1.32665","ask":"1.32827","date":"1521237591"},{"symbol":"EURCHF","bid":"1.16987","ask":"1.17035","date":"1521237591"},{"symbol":"USDNOK","bid":"7.70691","ask":"7.72433","date":"1521237591"},{"symbol":"AUDNZD","bid":"1.06793","ask":"1.06945","date":"1521237591"},{"symbol":"USDSEK","bid":"8.18625","ask":"8.20234","date":"1521237592"},{"symbol":"GBPSGD","bid":"1.8371","ask":"1.8383","date":"1521237592"},{"symbol":"GBPCAD","bid":"1.82562","ask":"1.82704","date":"1521237592"},{"symbol":"NZDCAD","bid":"0.94502","ask":"0.94571","date":"1521237592"},{"symbol":"AUDCAD","bid":"1.0097","ask":"1.0109","date":"1521237592"},{"symbol":"EURJPY","bid":"130.191","ask":"130.261","date":"1521237593"},{"symbol":"GBPJPY","bid":"147.721","ask":"147.818","date":"1521237593"},{"symbol":"USDJPY","bid":"105.952","ask":"105.953","date":"1521237595"},{"symbol":"AUDJPY","bid":"81.726","ask":"81.783","date":"1521237595"},{"symbol":"NZDJPY","bid":"76.46","ask":"76.52","date":"1521237595"},{"symbol":"EURSGD","bid":"1.6178","ask":"1.619","date":"1521237595"},{"symbol":"AUDUSD","bid":"0.77109","ask":"0.77174","date":"1521237596"},{"symbol":"EURAUD","bid":"1.59234","ask":"1.5938","date":"1521237596"},{"symbol":"EURGBP","bid":"0.88116","ask":"0.88131","date":"1521237596"},{"symbol":"USDCHF","bid":"0.95174","ask":"0.95235","date":"1521237596"},{"symbol":"USDDKK","bid":"6.05957","ask":"6.06203","date":"1521237596"},{"symbol":"USDMXN","bid":"18.67972","ask":"18.69992","date":"1521237596"},{"symbol":"AUDCHF","bid":"0.7338","ask":"0.73381","date":"1521237596"},{"symbol":"NZDUSD","bid":"0.72113","ask":"0.72177","date":"1521237596"},{"symbol":"EURUSD","bid":"1.22878","ask":"1.22905","date":"1521237597"},{"symbol":"XAGUSD","bid":"16.32","ask":"16.36","date":"1521237598"},{"symbol":"USDCAD","bid":"1.30942","ask":"1.30989","date":"1521237599"},{"symbol":"GBPUSD","bid":"1.39363","ask":"1.39423","date":"1521237599"}],"msgResult":"Success"}
type commandCode2029Return struct {
	MsgType         string             `json:"msgType"` //lastQuote
	LastQuotesArray *[]lastQuotesArray `json:"lastQuotesArray"`
	MsgResult       string             `json:"msgResult"`
}
type lastQuotesArray struct {
	Symbol string `json:"symbol"`
	Bid    string `json:"bid"`
	Ask    string `json:"ask"`
	Date   string `json:"date"`
}

func (this *ws) commandCode2029Return() {
	bytesAddr := system.ReadFile("./lastQuotesArray.json")
	sArr := []lastQuotesArray{}
	err := json.Unmarshal(*bytesAddr, &sArr)
	if err != nil {
		log.Println(err)
		return
	}

	this.wsSend(&commandCode2029Return{
		MsgType:         "lastQuote",
		LastQuotesArray: &sArr,
		MsgResult:       "Success",
	})

	var ok bool
	var it interface{}
	var hArr []*hub
	if it, ok = sMapWsConns.Load(this.User.Uuid); ok {
		hArr = it.([]*hub)
	}
	if len(hArr) > 9 { //每个用户只能拥有10个广播名额
		return
	}
	hArr = append(hArr, &hub{Conn: this.Conn, Op: this.Op, Ok: true, Freq: 5})
	sMapWsConns.Store(this.User.Uuid, hArr) //记录ws连接
}

func SymbolBroadcast(qDetails *quoteDetails) {
	var err error
	var b []byte
	s := commandCode2029Return_quoteDetails{MsgType: "quote", QuoteDetails: qDetails, MsgResult: "Success"}
	if b, err = json.Marshal(&s); err != nil {
		return
	}
	str := string(b)
	sMapWsConns.Range(func(k, v interface{}) bool {
		for kk, vHub := range v.([]*hub) {
			if vHub.Ok {
				if err = websocket.WriteFrame(*vHub.Conn, websocket.NewTextFrame(str)); err != nil {
					v.([]*hub)[kk].Ok = false
					sMapWsConns.Store(k.(uuid.UUID), v.([]*hub)) //发送失败 可能是对方关闭页面 或网络故障 采用临时关闭该conn 直到客户端发回2025指令恢复
				}
			}
		}
		return true
	})
}

func loopDelConn() { //5分钟后删除已经关闭的websocket连接
	for {
		select {
		case <-time.After(time.Minute):
			var sumAll, sumInvalid int
			sMapWsConns.Range(func(k, v interface{}) bool {
				sumAll += 1
				for kk, vHub := range v.([]*hub) {
					if vHub.Ok == false {
						if vHub.Freq > 0 {
							v.([]*hub)[kk].Freq -= 1
							sMapWsConns.Store(k.(uuid.UUID), v.([]*hub))
							continue
						} else {
							if len(v.([]*hub)) == 1 { //该用于就只有一个ws连接
								sMapWsConns.Delete(k.(uuid.UUID))
							} else {
								slice := v.([]*hub)
								slice = append(slice[:kk], slice[kk+1:]...)
								sMapWsConns.Store(k.(uuid.UUID), slice)
							}
						}
						//break
					}
					sumInvalid += 1
				}
				return true
			})
			fmt.Println("websocket,总", sumAll, " ,在线:", sumAll-sumInvalid, ", 离开:", sumInvalid)
		}
	}
}

//{"msgType":"pong","msgResult":"Success"}
type commandCode2025Return struct {
	MsgType   string `json:"msgType"`   //pong
	MsgResult string `json:"msgResult"` //Success
}

func (this *ws) commandCode2025Return() {
	pong := commandCode2025Return{
		MsgType:   "pong",
		MsgResult: "Success",
	}
	err := this.wsSend(&pong)

	//恢复conn的状态
	if err == nil {
		it, ok := sMapWsConns.Load(this.User.Uuid)
		if !ok {
			return
		}
		for k, v := range it.([]*hub) {
			if v.Conn == this.Conn {
				if v.Ok == false {
					it.([]*hub)[k].Ok = true
					sMapWsConns.Store(this.User.Uuid, it.([]*hub))
				}
				return
			}
		}

	}
}

//{"msgType":"traderData","id":"11043","name":"jialiang","patronymic":"","surname":"kang","country":"中国","region":"demo","city":"佛山","mailIndex":"0","address":"demo","phoneNumber":"13516521123","email":"me2@a7a2.com","allowTrade":"1","registrationDate":"3/16/2018 8:34:12 AM","lastLoginDate":"3/18/2018 4:31:59 PM","delay":"0","leverage":"100","comment":"","stopout":"20","depositPercent":"0","minDeposit":"0","maxOrderCount":"30","registrationIP":"117.136.12.129","priceStrategy":"0","coverageStrategyID":"1","showbonuses":"0","equityMarginCall":"1","groupID":"2","msgResult":"Success"}
type commandCode2013Return_person struct {
	MsgType string `json:"msgType"` //traderData
	Person
	MsgResult string `json:"msgResult"`
}
type Person struct {
	Id                 string `json:"id"`         //login
	Name               string `json:"name"`       //name
	Patronymic         string `json:"patronymic"` //
	Surname            string `json:"surname"`
	Country            string `json:"country"`
	Region             string `json:"region"`
	City               string `json:"city"`
	MailIndex          string `json:"mailIndex"`
	Address            string `json:"address"`
	PhoneNumber        string `json:"phoneNumber"` //phone
	Email              string `json:"email"`
	AllowTrade         string `json:"allowTrade"`
	RegistrationDate   string `json:"registrationDate"` //CreatedAt
	LastLoginDate      string `json:"lastLoginDate"`
	Delay              string `json:"delay"`
	Leverage           string `json:"leverage"`
	Comment            string `json:"comment"`
	Stopout            string `json:"stopout"`
	DepositPercent     string `json:"depositPercent"`
	MinDeposit         string `json:"minDeposit"`
	MaxOrderCount      string `json:"maxOrderCount"`
	RegistrationIP     string `json:"registrationIP"`
	PriceStrategy      string `json:"priceStrategy"`
	CoverageStrategyID string `json:"coverageStrategyID"`
	Showbonuses        string `json:"showbonuses"`
	EquityMarginCall   string `json:"equityMarginCall"`
	GroupID            string `json:"groupID"`
}

func (this *ws) commandCode2013Return() {
	id, mailIndex, allowTrade, delay, leverage, stopout, maxOrderCount, priceStrategy, coverageStrategyID, equityMarginCall, groupID := strconv.Itoa(this.User.Login), strconv.Itoa(this.User.MailIndex), strconv.Itoa(this.User.AllowTrade), strconv.Itoa(this.User.Delay), strconv.Itoa(this.User.Leverage), strconv.Itoa(this.User.Stopout), strconv.Itoa(this.User.MaxOrderCount), strconv.Itoa(this.User.PriceStrategy), strconv.Itoa(this.User.CoverageStrategyID), strconv.Itoa(this.User.EquityMarginCall), strconv.Itoa(this.User.GroupID)
	depositPercent := strconv.Itoa(this.User.DepositPercent)
	//minDeposit := strconv.FormatFloat(this.User.MinDeposit, 'f', 2, 64)
	//showbonuses := strconv.FormatFloat(this.User.Showbonuses, 'f', 2, 64)
	phone := strconv.FormatInt(this.User.Phone, 10)
	p := Person{
		Id:                 id,
		Name:               this.User.Name,
		Patronymic:         this.User.Patronimic,
		Surname:            this.User.Surname,
		Country:            this.User.Country,
		Region:             this.User.Region,
		City:               this.User.City,
		MailIndex:          mailIndex,
		Address:            this.User.Address,
		PhoneNumber:        phone,
		Email:              this.User.Email,
		AllowTrade:         allowTrade,
		RegistrationDate:   this.User.CreatedAt.Format("1/2/2006 15:04:05 PM"),
		LastLoginDate:      this.User.LastLoginDate.Format("1/2/2006 15:04:05 PM"),
		Delay:              delay,
		Leverage:           leverage,
		Comment:            "",
		Stopout:            stopout,
		DepositPercent:     depositPercent,
		MinDeposit:         "0",
		MaxOrderCount:      maxOrderCount,
		RegistrationIP:     this.User.RegistrationIP,
		PriceStrategy:      priceStrategy,
		CoverageStrategyID: coverageStrategyID,
		Showbonuses:        "0",
		EquityMarginCall:   equityMarginCall,
		GroupID:            groupID,
	}

	this.wsSend(&commandCode2013Return_person{
		MsgType:   "traderData",
		Person:    p,
		MsgResult: "Success",
	})
}

type commandCode2067Return struct {
	MsgType        string            `json:"msgType"` //optionsSettingsName
	OptionTypeName *[]optionTypeName `json:"optionTypeName"`
	//"optionTypeName":[{"name":"Express","aliasName":"Express"}
	//{"name":"Classic","aliasName":"Classic"}
	//{"name":"One Touch","aliasName":"One Touch"},
	//{"name":"Range","aliasName":"Range"}],"msgResult":"Success"}
	MsgResult string `json:"msgResult"`
}
type optionTypeName struct {
	Name      string `json:"name"`
	AliasName string `json:"aliasName"`
}

//{"msgType":"positions","sumsDetails":{"sumPositions":"0","sumInput":"5000","sumOutput":"0","sumInputBonus":"0","sumOutputBonus":"0","sumHistory":"0","sumSpendBonus":"0"},"positionsArray":[],"msgResult":"Success"}
type commandCode2014Return struct {
	MsgType        string       `json:"msgType"`
	SumsDetails    *sumsDetails `json:"sumsDetails"`
	PositionsArray []string     `json:"positionsArray"`
	MsgResult      string       `json:"msgResult"`
}
type sumsDetails struct { //GetPositions
	SumPositions   string `json:"sumPositions"`
	SumInput       string `json:"sumInput"`
	SumOutput      string `json:"sumOutput"`
	SumInputBonus  string `json:"sumInputBonus"`
	SumOutputBonus string `json:"sumOutputBonus"`
	SumHistory     string `json:"sumHistory"`
	SumSpendBonus  string `json:"sumSpendBonus"`
}

type commandCode2061Return struct {
	MsgType             string                       `json:"msgType"`
	OptionSettingsArray *[]optionSettingsArrayString `json:"optionSettingsArray"`
	MsgResult           string                       `json:"msgResult"`
}

type commandCode2015Return struct {
	MsgType            string   `json:"msgType"`
	PendingOrdersArray []string `json:"pendingOrdersArray"`
	MsgResult          string   `json:"msgResult"`
}

func nameConn(conn net.Conn) string {
	return conn.LocalAddr().String() + " > " + conn.RemoteAddr().String()
}

func Wshandler(w http.ResponseWriter, r *http.Request) {
	conn, _, hs, err := websocket.UpgradeHTTP(r, w, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("%s: established websocket connection: %+v", nameConn(conn), hs)

	//fmt.Println(websocket.DefaultClientWriteBufferSize)

	user := new(models.Users)
	sl := strings.Split(r.RequestURI, "/") //get token
	if len(sl) != 3 {
		log.Println("len(sl) != 3")
		return
	}

	strCmd := redis.RedisClient.Get(sl[2])
	if strCmd.Err() != nil {
		log.Println(strCmd.Err())
		return
	}
	userUuidStr := strCmd.Val()
	if err = user.FindUsersByUuid(userUuidStr); err != nil {
		log.Println(strCmd.Err())
		return
	}

	redis.Store(user.Uuid.String(), user, time.Hour*24)

	this := &ws{
		Conn:   &conn, //*websocket.Conn
		User:   user,  //来源于redis
		PByte:  []byte{},
		MsgMap: make(map[string]interface{}),
		Token:  sl[2],
	}
	//sumSyncMap.Store(&conn, &sum{0.0, 0.0})
	go func() {
		defer conn.Close()

		var codeF64 float64
		var codeStr string
		var ok bool
		for {

			msgBytes, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				log.Println(err)
				break
			}
			err = json.Unmarshal(msgBytes, &this.MsgMap)
			if err != nil {
				continue
			}

			if codeF64, ok = this.MsgMap["commandCode"].(float64); ok {
				codeStr = strconv.FormatFloat(codeF64, 'f', 0, 64)
			} else if codeStr, ok = this.MsgMap["commandCode"].(string); !ok {
				log.Println("codeStr, ok = this.MsgMap[\"commandCode\"].(string)")
				continue
			}
			this.PByte = msgBytes
			this.Op = op

			if op == websocket.OpClose {
				continue
			}
			this.commandCode(codeStr)
		}
	}()
}

type ws struct {
	Conn   *net.Conn
	Err    error
	MsgMap map[string]interface{}
	PByte  []byte
	User   *models.Users
	Op     websocket.OpCode
	Token  string
}
type sum struct {
	OptionDealsF64 float64
	SpendBonusF64  float64
}

//var sumSyncMap = new(sync.Map)

type commandCode2030Return struct {
	MsgType              string `json:"msgType"`
	Id                   string `json:"id"`
	GroupName            string `json:"groupName"`
	Description          string `json:"description"`
	RegistrationDate     string `json:"registrationDate"`
	DepositCurrency      string `json:"depositCurrency"`
	PositionStep         string `json:"positionStep"`
	AllowDemo            string `json:"allowDemo"`
	MinOptionSum         string `json:"minOptionSum"`
	MaxOptionSum         string `json:"maxOptionSum"`
	OptionSumLevel1      string `json:"optionSumLevel1"`
	OptionSumLevel2      string `json:"optionSumLevel2"`
	OptionSumLevel3      string `json:"optionSumLevel3"`
	OptionEqualityPrices string `json:"optionEqualityPrices"`
	MsgResult            string `json:"msgResult"`
}

func (this *ws) commandCode2030Return() {
	//r := map[interface{}]interface{}{"msgType": "traderGroup", "id": "2", "groupName": "general (demo)", "description": "dg", "registrationDate": "1/24/2013", "depositCurrency": "usd", "positionStep": "2", "allowDemo": "-1", "minOptionSum": "1", "maxOptionSum": "10000", "optionSumLevel1": "10", "optionSumLevel2": "25", "optionSumLevel3": "100", "optionEqualityPrices": "-1", "msgResult": "Success"}

	r := commandCode2030Return{
		MsgType:              "traderGroup",
		Id:                   this.User.Uuid.String(),
		GroupName:            this.User.GroupName, //"general (demo)"
		Description:          this.User.GroupName,
		RegistrationDate:     this.User.CreatedAt.Format("1/2/2006"),
		DepositCurrency:      this.User.DepositCurrency,
		PositionStep:         "2",
		AllowDemo:            "-1",
		MinOptionSum:         "1",
		MaxOptionSum:         "100000",
		OptionSumLevel1:      "10",
		OptionSumLevel2:      "25",
		OptionSumLevel3:      "100",
		OptionEqualityPrices: "-1",
		MsgResult:            "Success",
	}
	err := this.wsSend(&r)
	if err != nil {
		fmt.Println(err)
	}
}

func (this *ws) commandCode2067Return() {
	//			r := map[string]interface{}{"msgType":"optionsSettingsName","optionTypeName":[{"name":"Express","aliasName":"Express"},{"name":"Classic","aliasName":"Classic"},{"name":"One Touch","aliasName":"One Touch"},{"name":"Range","aliasName":"Range"}],"msgResult":"Success"}
	r := []optionTypeName{
		{
			Name:      "Express",
			AliasName: "Express",
		},
		{
			Name:      "Classic",
			AliasName: "Classic",
		},
		{
			Name:      "One Touch",
			AliasName: "One Touch",
		},
		{
			Name:      "Range",
			AliasName: "Range",
		},
	}
	this.wsSend(&commandCode2067Return{
		MsgType:        "optionsSettingsName",
		OptionTypeName: &r,
		MsgResult:      "Success",
	})
}

func (this *ws) commandCode2014Return() {
	//{"msgType":"positions","sumsDetails":{"sumPositions":"0","sumInput":"5000","sumOutput":"0","sumInputBonus":"0","sumOutputBonus":"0","sumHistory":"0","sumSpendBonus":"0"},"positionsArray":[],"msgResult":"Success"}
	this.User.FindUsersByUuid(this.User.Uuid.String())
	if err := this.User.FindUsersByUuid(this.User.Uuid.String()); err == nil {
		r := sumsDetails{
			SumPositions:   "0",
			SumInput:       strconv.FormatFloat(this.User.Balance-this.User.SumOptionDeals, 'f', 2, 64),
			SumOutput:      "0",
			SumInputBonus:  "0",
			SumOutputBonus: "0",
			SumHistory:     "0",
			SumSpendBonus:  "0",
		}
		this.wsSend(&commandCode2014Return{
			MsgType:        "positions",
			SumsDetails:    &r,
			PositionsArray: []string{},
			MsgResult:      "Success",
		})
	}

}

func (this *ws) commandCode2061Return() {
	this.wsSend(&commandCode2061Return{
		MsgType:             "optionsSettings",
		OptionSettingsArray: &optionSettingsArrS,
		MsgResult:           "Success",
	})
}

func (this *ws) commandCode2015Return() {
	this.wsSend(&commandCode2015Return{
		MsgType:            "pendingOrders",
		PendingOrdersArray: []string{},
		MsgResult:          "Success",
	})
}

type commandCode2062Return struct {
	MsgType              string                `json:"msgType"`
	SumsDetails          *sumsDetailsPositions `json:"sumsDetails"`
	OptionPositionsArray []string              `json:"optionPositionsArray"`
	MsgResult            string                `json:"msgResult"`
}

func unicast(u *models.Users) {
	sMapWsConns.Range(func(k, v interface{}) bool {
		for _, vHub := range v.([]*hub) {
			if vHub.Ok {
				cc2062Return(u, vHub.Op, vHub.Conn)
			}
		}
		return true
	})

}

func cc2062Return(u *models.Users, op websocket.OpCode, conn *net.Conn) { //GetOptionPositions
	var tArr []models.Trade
	var err error
	if tArr, err = models.ListTrade("close=false and delete=false and \"userUuid\"='" + u.Uuid.String() + "'"); err != nil {
		fmt.Println(err)
		return
	}

	sliceOp := make([]optionPositionsArray_return, 0)
	for _, v := range tArr {
		expiryDate := v.OpenDate + int64(v.Period)
		op := optionPositionsArray_return{
			Id:               v.Uuid,
			OptionType:       strconv.Itoa(v.OptionType),
			Direction:        strconv.Itoa(v.Direction),
			SettingName:      v.SymbolName,
			SymbolName:       v.SymbolName,
			StopLine:         strconv.Itoa(v.StopLine),
			PayoutPercentage: strconv.FormatFloat(v.PayoutPercentage, 'f', -1, 64),
			EarlyClosing:     strconv.Itoa(v.EarlyClosing),
			Interval:         strconv.Itoa(v.Interval),
			OpenPrice:        strconv.FormatFloat(v.OpenPrice, 'f', -1, 64),
			OpenDate:         strconv.FormatInt(v.OpenDate, 10),
			ExpiryDate:       strconv.FormatInt(expiryDate, 10),
			InvestmentSum:    strconv.FormatFloat(v.InvestmentSum, 'f', -1, 64),
		}
		sliceOp = append(sliceOp, op)
	}

	sumsDetails := sumsDetailsPositions{
		SumOptionPositions: "0",
		SumOptionDeals:     strconv.FormatFloat(u.SumOptionDeals, 'f', 2, 64),
		SumSpendBonus:      "0",
	}

	var b []byte
	var wf websocket.Frame
	var cc2062R interface{}
	if len(sliceOp) == 0 {
		cc2062R = &commandCode2062Return{
			MsgType:              "optionsPositions",
			SumsDetails:          &sumsDetails,
			OptionPositionsArray: []string{},
			MsgResult:            "Success",
		}
	} else {
		cc2062R = &optionsPositionsSumsDetails_return{
			MsgType:              "optionsPositions",
			SumsDetails:          &sumsDetails,
			OptionPositionsArray: &sliceOp,
			MsgResult:            "Success",
		}
	}
	if b, err = json.Marshal(cc2062R); err != nil {
		log.Println(err)
		return
	}

	wf = websocket.NewFrame(op, true, b)
	websocket.WriteFrame(*conn, wf)
}

func (this *ws) commandCode2062Return() {
	if err := this.User.FindUsersByUuid(this.User.Uuid.String()); err != nil {
		return
	}
	unicast(this.User)
	//cc2062Return(this.User, this.Op, this.Conn, 0.0)
	//{"msgType":"optionsPositions","sumsDetails":{"sumOptionPositions":"0","sumOptionDeals":"0","sumSpendBonus":"0"},"optionPositionsArray":[],"msgResult":"Success"}
}

func (this *ws) commandCode2003Return() {
	var cp2003 changePassword2003
	this.Err = json.Unmarshal(this.PByte, &cp2003)
	if this.Err != nil {
		this.wsSend(map[string]interface{}{"error": "2003:" + this.Err.Error()})
		return
	}

	newPassword := helpers.EncryptPwd2Db([]byte(cp2003.NewPassword))
	isInvestorPassword := cp2003.IsInvestorPassword

	if this.Err = this.User.UserAuth(this.User.Login, helpers.EncryptPwd2Db([]byte(cp2003.OldPassword))); this.Err != nil { //验证旧密码ok
		this.wsSend(map[string]interface{}{"error": ":密码错误:"})
		return
	}
	if isInvestorPassword { //修改投资密码
		this.User.InvestorPassword = newPassword
	} else {
		this.User.Password = newPassword
	}
	if this.Err = this.User.Update(); this.Err != nil {
		this.wsSend(map[string]interface{}{"error": "2003:" + this.Err.Error()})
		return
	}
	this.wsSend(&result{
		MsgType:   "changePassword",
		MsgResult: "Success",
	})
}

type result struct {
	MsgType   string `json:"msgType"` //optionsPositions
	MsgResult string `json:"msgResult"`
}

type changePassword2003 struct {
	CommandCode        string `json:"commandCode"`
	TraderID           string `json:"traderID"`
	OldPassword        string `json:"oldPassword"`
	NewPassword        string `json:"newPassword"`
	IsInvestorPassword bool   `json:"isInvestorPassword"`
}

type sumsDetailsPositions struct {
	SumOptionPositions string `json:"sumOptionPositions"` //付出金
	SumOptionDeals     string `json:"sumOptionDeals"`     //统计总结算后的输赢金额
	SumSpendBonus      string `json:"sumSpendBonus"`
}

type commandCode2029Return_quoteDetails struct {
	MsgType      string        `json:"msgType"` //quote
	QuoteDetails *quoteDetails `json:"quoteDetails"`
	MsgResult    string        `json:"msgResult"`
}
type quoteDetails struct {
	Symbol string `json:"symbol"`
	Bid    string `json:"bid"`
	Ask    string `json:"ask"`
	Date   string `json:"date"`
}

func (this *ws) commandCode(codeStr string) {
	switch codeStr {
	case "2012": //服务器信息
		this.commandCode2012Return()
	case "2038": //交易类型信息 、节假日
		this.commandCode2038Return()
	case "2036": //各个类型的具体参数
		this.commandCode2036Return()
	case "2023": // GetSymbols
		this.commandCode2023Return()
	case "2025": //ping
		this.commandCode2025Return()
	case "2029": //GetLastQuotes
		this.commandCode2029Return()
	case "2013": //GetClientData
		this.commandCode2013Return()
	case "2030": //GetTraderGroup
		this.commandCode2030Return()
	case "2067": // GetOptionsAliasTypeName:
		this.commandCode2067Return()
	case "2014":
		//chanCommandCode2014 <- true
		this.commandCode2014Return()
	case "2061":
		this.commandCode2061Return()
	case "2015":
		this.commandCode2015Return()
	case "2062": //GetOptionPositions
		this.commandCode2062Return()
	case "2003": //ChangePassword
		this.commandCode2003Return()
	case "2060": //OpenOptionPosition
		if err := this.commandCode2060Return(); err != nil {
			this.wsSend(map[string]interface{}{"error": models.ErrCommit})
		}
	case "2066": //CloseOptionPosition
		this.commandCode2066Return()
	case "2011": //GetOperationHistory
		this.commandCode2011Return()
	case "2016": //GetDepositOperations
		this.commandCode2016Return()
	case "2063": //GetOptionDeals
		this.commandCode2063Return()
	case "2024": //GetExecutedOrders
		this.commandCode2024Return()
	}
}

type result2024 struct {
	MsgType             string   `json:"msgType"`
	ExecutedOrdersArray []string `json:"executedOrdersArray"`
	MsgResult           string   `json:"msgResult"`
}

//{"msgType":"executedOrders","executedOrdersArray":[],"msgResult":"Success"}
func (this *ws) commandCode2024Return() {
	this.wsSend(&result2024{
		MsgType:             "executedOrders",
		ExecutedOrdersArray: []string{},
		MsgResult:           "Success",
	})
}

//{"msgType":"optionDeals","optionDealsArray":[{"id":"107128","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"1","openPrice":"1.40666","openDate":"1522667718","closeDate":"1522667778","closePrice":"1.40678","investmentSum":"10","profit":"-10","payout":"0","balance":"4990","isEarlyClosing":"0","spendbonus":"0"},{"id":"107129","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"1","openPrice":"1.40674","openDate":"1522667740","closeDate":"1522667800","closePrice":"1.40654","investmentSum":"10","profit":"7","payout":"17","balance":"4997","isEarlyClosing":"0","spendbonus":"0"},{"id":"107130","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"1","openPrice":"1.40672","openDate":"1522667750","closeDate":"1522667810","closePrice":"1.40656","investmentSum":"10","profit":"7","payout":"17","balance":"5004","isEarlyClosing":"0","spendbonus":"0"},{"id":"107131","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"1","openPrice":"1.40669","openDate":"1522667735","closeDate":"1522667855","closePrice":"1.40656","investmentSum":"10","profit":"7","payout":"17","balance":"5011","isEarlyClosing":"0","spendbonus":"0"},{"id":"107132","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"0","openPrice":"1.40656","openDate":"1522667805","closeDate":"1522667865","closePrice":"1.40653","investmentSum":"10","profit":"-10","payout":"0","balance":"5001","isEarlyClosing":"0","spendbonus":"0"},{"id":"107133","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"0","openPrice":"1.40655","openDate":"1522667809","closeDate":"1522667869","closePrice":"1.40655","investmentSum":"10","profit":"-10","payout":"0","balance":"4991","isEarlyClosing":"0","spendbonus":"0"},{"id":"107134","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"0","openPrice":"1.40655","openDate":"1522667812","closeDate":"1522667872","closePrice":"1.40659","investmentSum":"10","profit":"7","payout":"17","balance":"4998","isEarlyClosing":"0","spendbonus":"0"},{"id":"107135","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"0","openPrice":"1.40654","openDate":"1522667816","closeDate":"1522667876","closePrice":"1.40653","investmentSum":"10","profit":"-10","payout":"0","balance":"4988","isEarlyClosing":"0","spendbonus":"0"},{"id":"107136","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"0","openPrice":"1.40654","openDate":"1522667819","closeDate":"1522667879","closePrice":"1.40655","investmentSum":"10","profit":"7","payout":"17","balance":"4995","isEarlyClosing":"0","spendbonus":"0"},{"id":"107137","settingName":"EURAUD","optionType":"0","symbolName":"EURAUD","direction":"0","openPrice":"1.60484","openDate":"1522675368","closeDate":"1522675548","closePrice":"1.60489","investmentSum":"99","profit":"69.3","payout":"168.3","balance":"5064.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107138","settingName":"EURAUD","optionType":"2","symbolName":"EURAUD","direction":"4","openPrice":"1.6048","openDate":"1522675910","closeDate":"1522675970","closePrice":"1.60514","investmentSum":"10","profit":"-10","payout":"0","balance":"5054.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107139","settingName":"EURAUD","optionType":"2","symbolName":"EURAUD","direction":"5","openPrice":"1.60492","openDate":"1522675948","closeDate":"1522676008","closePrice":"1.60507","investmentSum":"10","profit":"-10","payout":"0","balance":"5044.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107140","settingName":"EURAUD","optionType":"2","symbolName":"EURAUD","direction":"5","openPrice":"1.60495","openDate":"1522675967","closeDate":"1522676027","closePrice":"1.60497","investmentSum":"10","profit":"-10","payout":"0","balance":"5034.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107141","settingName":"EURAUD","optionType":"2","symbolName":"EURAUD","direction":"4","openPrice":"1.60506","openDate":"1522675977","closeDate":"1522676037","closePrice":"1.60495","investmentSum":"10","profit":"-10","payout":"0","balance":"5024.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107142","settingName":"EURAUD","optionType":"1","symbolName":"EURAUD","direction":"2","openPrice":"1.60485","openDate":"1522675569","closeDate":"1522676469","closePrice":"1.60525","investmentSum":"100","profit":"70","payout":"170","balance":"5094.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107143","settingName":"EURUSD","optionType":"0","symbolName":"EURUSD","direction":"0","openPrice":"1.23167","openDate":"1522745282","closeDate":"1522745342","closePrice":"1.2319","investmentSum":"10","profit":"7.5","payout":"17.5","balance":"5101.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107144","settingName":"AUDJPY","optionType":"2","symbolName":"AUDJPY","direction":"5","openPrice":"81.666","openDate":"1522745349","closeDate":"1522745409","closePrice":"81.659","investmentSum":"10","profit":"-10","payout":"0","balance":"5091.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107314","settingName":"GBPCAD","optionType":"0","symbolName":"GBPCAD","direction":"1","openPrice":"1.78951","openDate":"1523005128","closeDate":"1523005188","closePrice":"1.78948","investmentSum":"10","profit":"7","payout":"17","balance":"5098.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107315","settingName":"GBPCAD","optionType":"0","symbolName":"GBPCAD","direction":"0","openPrice":"1.78872","openDate":"1523005742","closeDate":"1523005802","closePrice":"1.78921","investmentSum":"10","profit":"7","payout":"17","balance":"5105.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107316","settingName":"GBPCAD","optionType":"0","symbolName":"GBPCAD","direction":"0","openPrice":"1.78871","openDate":"1523005754","closeDate":"1523005814","closePrice":"1.78925","investmentSum":"10","profit":"7","payout":"17","balance":"5112.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107317","settingName":"GBPCAD","optionType":"0","symbolName":"GBPCAD","direction":"0","openPrice":"1.78969","openDate":"1523006001","closeDate":"1523006061","closePrice":"1.7896","investmentSum":"10","profit":"-10","payout":"0","balance":"5102.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107318","settingName":"GBPCAD","optionType":"0","symbolName":"GBPCAD","direction":"0","openPrice":"1.78882","openDate":"1523008621","closeDate":"1523008681","closePrice":"1.78892","investmentSum":"10","profit":"7","payout":"17","balance":"5109.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107319","settingName":"GBPCAD","optionType":"0","symbolName":"GBPCAD","direction":"1","openPrice":"1.78887","openDate":"1523008668","closeDate":"1523008728","closePrice":"1.78874","investmentSum":"10","profit":"7","payout":"17","balance":"5116.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107320","settingName":"GBPCAD","optionType":"0","symbolName":"GBPCAD","direction":"0","openPrice":"1.78877","openDate":"1523008752","closeDate":"1523008812","closePrice":"1.78891","investmentSum":"10","profit":"7","payout":"17","balance":"5123.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107321","settingName":"GBPCAD","optionType":"0","symbolName":"GBPCAD","direction":"1","openPrice":"1.78913","openDate":"1523008838","closeDate":"1523008898","closePrice":"1.78907","investmentSum":"10","profit":"7","payout":"17","balance":"5130.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107322","settingName":"GBPCAD","optionType":"0","symbolName":"GBPCAD","direction":"0","openPrice":"1.79004","openDate":"1523013268","closeDate":"1523013328","closePrice":"1.79003","investmentSum":"10","profit":"-10","payout":"0","balance":"5120.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107323","settingName":"GBPCAD","optionType":"0","symbolName":"GBPCAD","direction":"0","openPrice":"1.79027","openDate":"1523013398","closeDate":"1523013458","closePrice":"1.78996","investmentSum":"10","profit":"-10","payout":"0","balance":"5110.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107324","settingName":"GBPCAD","optionType":"0","symbolName":"GBPCAD","direction":"0","openPrice":"1.78992","openDate":"1523013676","closeDate":"1523013736","closePrice":"1.78971","investmentSum":"10","profit":"-10","payout":"0","balance":"5100.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107325","settingName":"GBPCAD","optionType":"0","symbolName":"GBPCAD","direction":"0","openPrice":"1.78984","openDate":"1523013685","closeDate":"1523013745","closePrice":"1.78957","investmentSum":"10","profit":"-10","payout":"0","balance":"5090.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107326","settingName":"GBPCAD","optionType":"0","symbolName":"GBPCAD","direction":"0","openPrice":"1.78988","openDate":"1523013689","closeDate":"1523013749","closePrice":"1.78964","investmentSum":"10","profit":"-10","payout":"0","balance":"5080.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107327","settingName":"GBPCAD","optionType":"0","symbolName":"GBPCAD","direction":"0","openPrice":"1.79034","openDate":"1523013966","closeDate":"1523014026","closePrice":"1.79047","investmentSum":"10","profit":"7","payout":"17","balance":"5087.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107342","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"0","openPrice":"1.40967","openDate":"1523246075","closeDate":"1523246135","closePrice":"1.40965","investmentSum":"10","profit":"-10","payout":"0","balance":"5077.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107343","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"0","openPrice":"1.40945","openDate":"1523246335","closeDate":"1523246395","closePrice":"1.40931","investmentSum":"10","profit":"-10","payout":"0","balance":"5067.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107344","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"1","openPrice":"1.40929","openDate":"1523246481","closeDate":"1523246541","closePrice":"1.4094","investmentSum":"10","profit":"-10","payout":"0","balance":"5057.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107345","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"1","openPrice":"1.40938","openDate":"1523246500","closeDate":"1523246560","closePrice":"1.40941","investmentSum":"10","profit":"-10","payout":"0","balance":"5047.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107346","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"1","openPrice":"1.40941","openDate":"1523246570","closeDate":"1523246630","closePrice":"1.40929","investmentSum":"10","profit":"7","payout":"17","balance":"5054.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107347","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"0","openPrice":"1.40938","openDate":"1523246591","closeDate":"1523246651","closePrice":"1.40933","investmentSum":"10","profit":"-10","payout":"0","balance":"5044.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107348","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"1","openPrice":"1.40937","openDate":"1523246994","closeDate":"1523247054","closePrice":"1.40927","investmentSum":"10","profit":"7","payout":"17","balance":"5051.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107349","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"0","openPrice":"1.40937","openDate":"1523247008","closeDate":"1523247068","closePrice":"1.40925","investmentSum":"10","profit":"-10","payout":"0","balance":"5041.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107350","settingName":"GBPUSD","optionType":"2","symbolName":"GBPUSD","direction":"5","openPrice":"1.40937","openDate":"1523247019","closeDate":"1523247079","closePrice":"1.40925","investmentSum":"10","profit":"-10","payout":"0","balance":"5031.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107351","settingName":"GBPUSD","optionType":"3","symbolName":"GBPUSD","direction":"7","openPrice":"1.40933","openDate":"1523247041","closeDate":"1523247101","closePrice":"1.40927","investmentSum":"10","profit":"6.5","payout":"16.5","balance":"5038.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107352","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"0","openPrice":"1.40919","openDate":"1523247154","closeDate":"1523247214","closePrice":"1.40918","investmentSum":"10","profit":"-10","payout":"0","balance":"5028.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107353","settingName":"EURCHF","optionType":"0","symbolName":"EURCHF","direction":"1","openPrice":"1.17753","openDate":"1523247159","closeDate":"1523247219","closePrice":"1.1776","investmentSum":"10","profit":"-10","payout":"0","balance":"5018.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107354","settingName":"EURAUD","optionType":"0","symbolName":"EURAUD","direction":"1","openPrice":"1.59617","openDate":"1523247190","closeDate":"1523247250","closePrice":"1.59612","investmentSum":"10","profit":"7","payout":"17","balance":"5025.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107355","settingName":"USDCHF","optionType":"0","symbolName":"USDCHF","direction":"0","openPrice":"0.95964","openDate":"1523247213","closeDate":"1523247273","closePrice":"0.95963","investmentSum":"10","profit":"-10","payout":"0","balance":"5015.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107356","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"1","openPrice":"1.40916","openDate":"1523247254","closeDate":"1523247314","closePrice":"1.40918","investmentSum":"10","profit":"-10","payout":"0","balance":"5005.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107357","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"0","openPrice":"1.40917","openDate":"1523247283","closeDate":"1523247343","closePrice":"1.40918","investmentSum":"10","profit":"7","payout":"17","balance":"5012.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107358","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"0","openPrice":"1.40918","openDate":"1523247347","closeDate":"1523247407","closePrice":"1.40907","investmentSum":"10","profit":"-10","payout":"0","balance":"5002.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107359","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"0","openPrice":"1.40904","openDate":"1523247358","closeDate":"1523247478","closePrice":"1.40906","investmentSum":"10","profit":"7","payout":"17","balance":"5009.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107360","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"0","openPrice":"1.40906","openDate":"1523247371","closeDate":"1523247551","closePrice":"1.40908","investmentSum":"10","profit":"7","payout":"17","balance":"5016.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107361","settingName":"AUDUSD","optionType":"0","symbolName":"AUDUSD","direction":"0","openPrice":"0.76861","openDate":"1523247494","closeDate":"1523247554","closePrice":"0.76862","investmentSum":"10","profit":"7","payout":"17","balance":"5023.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107362","settingName":"AUDUSD","optionType":"0","symbolName":"AUDUSD","direction":"1","openPrice":"0.76863","openDate":"1523247548","closeDate":"1523247668","closePrice":"0.76852","investmentSum":"10","profit":"7","payout":"17","balance":"5030.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107363","settingName":"GBPUSD","optionType":"1","symbolName":"GBPUSD","direction":"2","openPrice":"1.40935","openDate":"1523246946","closeDate":"1523247846","closePrice":"1.40877","investmentSum":"10","profit":"-10","payout":"0","balance":"5020.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107364","settingName":"GBPUSD","optionType":"1","symbolName":"GBPUSD","direction":"2","openPrice":"1.40934","openDate":"1523247035","closeDate":"1523247935","closePrice":"1.40879","investmentSum":"10","profit":"-10","payout":"0","balance":"5010.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107365","settingName":"GBPUSD","optionType":"3","symbolName":"GBPUSD","direction":"7","openPrice":"1.40879","openDate":"1523247948","closeDate":"1523248008","closePrice":"1.40882","investmentSum":"10","profit":"6.5","payout":"16.5","balance":"5016.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107366","settingName":"GBPUSD","optionType":"3","symbolName":"GBPUSD","direction":"6","openPrice":"1.40879","openDate":"1523247953","closeDate":"1523248013","closePrice":"1.4088","investmentSum":"10","profit":"-10","payout":"0","balance":"5006.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107367","settingName":"GBPUSD","optionType":"2","symbolName":"GBPUSD","direction":"4","openPrice":"1.40881","openDate":"1523247958","closeDate":"1523248018","closePrice":"1.40886","investmentSum":"10","profit":"-10","payout":"0","balance":"4996.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107368","settingName":"GBPUSD","optionType":"2","symbolName":"GBPUSD","direction":"5","openPrice":"1.40881","openDate":"1523247963","closeDate":"1523248023","closePrice":"1.40886","investmentSum":"10","profit":"-10","payout":"0","balance":"4986.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107369","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"1","openPrice":"1.40881","openDate":"1523248106","closeDate":"1523248166","closePrice":"1.40881","investmentSum":"10","profit":"-10","payout":"0","balance":"4976.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107370","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"1","openPrice":"1.40887","openDate":"1523248133","closeDate":"1523248193","closePrice":"1.40882","investmentSum":"10","profit":"7","payout":"17","balance":"4983.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107371","settingName":"GBPUSD","optionType":"1","symbolName":"GBPUSD","direction":"2","openPrice":"1.40907","openDate":"1523247388","closeDate":"1523248288","closePrice":"1.40899","investmentSum":"10","profit":"-10","payout":"0","balance":"4973.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107372","settingName":"AUDUSD","optionType":"1","symbolName":"AUDUSD","direction":"3","openPrice":"0.76862","openDate":"1523247563","closeDate":"1523248463","closePrice":"0.76836","investmentSum":"10","profit":"7","payout":"17","balance":"4980.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107373","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"1","openPrice":"1.40917","openDate":"1523248560","closeDate":"1523248620","closePrice":"1.40911","investmentSum":"10","profit":"7","payout":"17","balance":"4987.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107374","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"1","openPrice":"1.40912","openDate":"1523248573","closeDate":"1523248633","closePrice":"1.40919","investmentSum":"10","profit":"-10","payout":"0","balance":"4977.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107375","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"1","openPrice":"1.40918","openDate":"1523248628","closeDate":"1523248688","closePrice":"1.40902","investmentSum":"10","profit":"7","payout":"17","balance":"4984.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107376","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"1","openPrice":"1.40919","openDate":"1523248633","closeDate":"1523248693","closePrice":"1.40904","investmentSum":"10","profit":"7","payout":"17","balance":"4991.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107377","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"1","openPrice":"1.40903","openDate":"1523248717","closeDate":"1523248777","closePrice":"1.409","investmentSum":"10","profit":"7","payout":"17","balance":"4998.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107378","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"0","openPrice":"1.40902","openDate":"1523248725","closeDate":"1523248785","closePrice":"1.40902","investmentSum":"10","profit":"-10","payout":"0","balance":"4988.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107379","settingName":"GBPUSD","optionType":"2","symbolName":"GBPUSD","direction":"4","openPrice":"1.409","openDate":"1523248732","closeDate":"1523248792","closePrice":"1.40902","investmentSum":"10","profit":"-10","payout":"0","balance":"4978.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107380","settingName":"GBPUSD","optionType":"3","symbolName":"GBPUSD","direction":"7","openPrice":"1.40911","openDate":"1523248742","closeDate":"1523248802","closePrice":"1.40908","investmentSum":"10","profit":"6.5","payout":"16.5","balance":"4985.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107381","settingName":"GBPUSD","optionType":"3","symbolName":"GBPUSD","direction":"7","openPrice":"1.40905","openDate":"1523248845","closeDate":"1523248905","closePrice":"1.40908","investmentSum":"10","profit":"6.5","payout":"16.5","balance":"4991.8","isEarlyClosing":"0","spendbonus":"0"},{"id":"107382","settingName":"GBPUSD","optionType":"3","symbolName":"GBPUSD","direction":"7","openPrice":"1.40903","openDate":"1523248880","closeDate":"1523248940","closePrice":"1.40898","investmentSum":"10","profit":"6.5","payout":"16.5","balance":"4998.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107383","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"0","openPrice":"1.41295","openDate":"1523284607","closeDate":"1523284667","closePrice":"1.41265","investmentSum":"10","profit":"-10","payout":"0","balance":"4988.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107384","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"0","openPrice":"1.41297","openDate":"1523284611","closeDate":"1523284671","closePrice":"1.41263","investmentSum":"10","profit":"-10","payout":"0","balance":"4978.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107385","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"0","openPrice":"1.41291","openDate":"1523285019","closeDate":"1523285079","closePrice":"1.41328","investmentSum":"10","profit":"7","payout":"17","balance":"4985.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107386","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"0","openPrice":"1.41338","openDate":"1523285111","closeDate":"1523285171","closePrice":"1.41337","investmentSum":"10","profit":"-10","payout":"0","balance":"4975.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107387","settingName":"GBPUSD","optionType":"2","symbolName":"GBPUSD","direction":"4","openPrice":"1.41349","openDate":"1523285130","closeDate":"1523285190","closePrice":"1.41336","investmentSum":"10","profit":"-10","payout":"0","balance":"4965.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107388","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"0","openPrice":"1.41338","openDate":"1523285201","closeDate":"1523285261","closePrice":"1.41336","investmentSum":"10","profit":"-10","payout":"0","balance":"4955.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107389","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"0","openPrice":"1.41338","openDate":"1523285211","closeDate":"1523285271","closePrice":"1.41332","investmentSum":"10","profit":"-10","payout":"0","balance":"4945.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107390","settingName":"GBPUSD","optionType":"2","symbolName":"GBPUSD","direction":"4","openPrice":"1.41327","openDate":"1523285282","closeDate":"1523285342","closePrice":"1.41329","investmentSum":"2000","profit":"-2000","payout":"0","balance":"2945.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107391","settingName":"GBPUSD","optionType":"0","symbolName":"GBPUSD","direction":"0","openPrice":"1.41416","openDate":"1523285935","closeDate":"1523285995","closePrice":"1.41428","investmentSum":"10","profit":"7","payout":"17","balance":"2952.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107520","settingName":"GBPUSD","optionType":"3","symbolName":"GBPUSD","direction":"6","openPrice":"1.43051","openDate":"1523879503","closeDate":"1523879563","closePrice":"1.43053","investmentSum":"1011","profit":"-1011","payout":"0","balance":"1941.3","isEarlyClosing":"0","spendbonus":"0"},{"id":"107521","settingName":"GBPUSD","optionType":"3","symbolName":"GBPUSD","direction":"6","openPrice":"1.43053","openDate":"1523879513","closeDate":"1523879573","closePrice":"1.43053","investmentSum":"1011","profit":"-1011","payout":"0","balance":"930.3","isEarlyClosing":"0","spendbonus":"0"}],"msgResult":"Success"}
type result2063 struct {
	MsgType          string              `json:"msgType"` //optionDeals
	OptionDealsArray *[]optionDealsArray `json:"optionDealsArray"`
	MsgResult        string              `json:"msgResult"`
}

type optionDealsArray struct {
	Id             string `json:"id"`
	SettingName    string `json:"settingName"`
	OptionType     string `json:"optionType"`
	SymbolName     string `json:"symbolName"`
	Direction      string `json:"direction"`
	OpenPrice      string `json:"openPrice"`
	OpenDate       string `json:"openDate"`
	CloseDate      string `json:"closeDate"`
	ClosePrice     string `json:"closePrice"`
	InvestmentSum  string `json:"investmentSum"`
	Profit         string `json:"profit"`
	Payout         string `json:"payout"`
	Balance        string `json:"balance"`
	IsEarlyClosing string `json:"isEarlyClosing"`
	Spendbonus     string `json:"spendbonus"`
}

func (this *ws) commandCode2063Return() (err error) {
	var q getDealsQuery
	if err = json.Unmarshal(this.PByte, &q); err != nil {
		return
	}

	where := fmt.Sprintf("\"userUuid\"='%s' and \"closeDate\">%v and \"closeDate\"<%v", this.User.Uuid.String(), q.StartDate, q.FinishDate)
	var tArr []models.Trade
	if tArr, err = models.ListTrade(where); err != nil {
		return
	}
	var oArr []optionDealsArray
	for _, v := range tArr {
		o := optionDealsArray{
			Id:             v.Uuid.String(),
			SettingName:    v.SymbolName,
			OptionType:     strconv.Itoa(v.OptionType),
			SymbolName:     v.SymbolName,
			Direction:      strconv.Itoa(v.Direction),
			OpenPrice:      strconv.FormatFloat(v.OpenPrice, 'f', -1, 64),
			OpenDate:       strconv.FormatInt(v.OpenDate, 10),
			CloseDate:      strconv.FormatInt(v.CloseDate, 10),
			ClosePrice:     strconv.FormatFloat(v.ClosePrice, 'f', -1, 64),
			InvestmentSum:  strconv.FormatFloat(v.InvestmentSum, 'f', -1, 64),
			Profit:         strconv.FormatFloat(v.Profit, 'f', -1, 64),
			Payout:         "0",
			Balance:        strconv.FormatFloat(v.Balance, 'f', 2, 64),
			IsEarlyClosing: strconv.Itoa(v.EarlyClosing),
			Spendbonus:     "0",
		}
		oArr = append(oArr, o)
	}

	this.wsSend(&result2063{
		MsgType:          "optionDeals",
		OptionDealsArray: &oArr,
		MsgResult:        "Success",
	})
	return nil
}

func (this *ws) commandCode2016Return() (err error) {
	where := fmt.Sprintf("\"userUuid\"='%s'", this.User.Uuid.String())
	mArr, err := models.ListMoney(where)
	if err != nil {
		return
	}
	var cc2016Arr []commandCode2016Return
	for _, v := range mArr {
		cc2016Arr = append(cc2016Arr, commandCode2016Return{
			DepositID:     v.Uuid.String(),
			Equity:        strconv.FormatFloat(v.Equity, 'f', 2, 64),
			Status:        strconv.Itoa(v.Status),
			OperationDate: strconv.FormatInt(v.OperationDate, 10),
			Comment:       v.Comment,
		})
	}
	this.wsSend(&result2016{
		MsgType:       "deposits",
		MsgResult:     "Success",
		DepositsArray: &cc2016Arr,
	})
	return nil
}

type result2016 struct {
	MsgType       string                   `json:"msgType"` //optionsPositions
	DepositsArray *[]commandCode2016Return `json:"depositsArray"`
	MsgResult     string                   `json:"msgResult"`
}

//{"msgType":"deposits","depositsArray":[{"depositID":"12147","equity":"5000","status":"0","operationDate":"1522667655","comment":"First deposit"}],"msgResult":"Success"}
type commandCode2016Return struct {
	DepositID     string `json:"depositID"`
	Equity        string `json:"equity"`
	Status        string `json:"status"`
	OperationDate string `json:"operationDate"`
	Comment       string `json:"comment"`
}

//{"commandCode":"2011","startDate":"0","finishDate":"1523662994"}
type getDealsQuery struct { //
	CommandCode string `json:"commandCode"`
	StartDate   string `json:"startDate"`
	FinishDate  string `json:"finishDate"`
}

//{"msgType":"deals","dealsArray":[],"msgResult":"Success"}
type commandCode2011Return struct {
	MsgType    string   `json:"msgType"` //deals
	DealsArray []string `json:"dealsArray"`
	MsgResult  string   `json:"msgResult"`
}

func (this *ws) commandCode2011Return() (err error) {
	var q getDealsQuery
	if err = json.Unmarshal(this.PByte, &q); err != nil {
		return
	}

	this.wsSend(&commandCode2011Return{
		MsgType:    "deals",
		DealsArray: []string{},
		MsgResult:  "Success",
	})
	return nil
}

func (this *ws) commandCode2066Return() { //删除交易单
	msgResult := "Success"
	positionID, ok := this.MsgMap["positionID"].(string)
	t := new(models.Trade)
	err := t.FindTradeByUuid(positionID)
	switch {
	case !ok:
		msgResult = "false"
	case err != nil:
		msgResult = err.Error()
	case t.Close:
		msgResult = "false"
	case t.Delete:
		msgResult = "false"
	case t.EarlyClosing <= 0: //没有设置可提前删除时间的 不能删
		msgResult = "false"
	case t.OpenDate+int64(t.Period-t.EarlyClosing) < time.Now().Unix(): //现在的时间大于可删除时间 禁止删除
		fmt.Println(t.OpenDate+int64(t.Period-t.EarlyClosing), " ", time.Now().Unix())
		msgResult = "Option position can not be closed after stop time"
	case this.User.Uuid != t.UserUuid: //非该交易单的用户提交的，权力不够
		msgResult = "Permission denied"
	}

	if err = this.User.FindUsersByUuid(this.User.Uuid.String()); err != nil {
		return
	}
	var money float64
	if msgResult == "Success" { //符合删除条件
		switch t.OptionType {
		case 2:
			money = t.InvestmentSum * 0.6
		case 3:
			money = t.InvestmentSum * 0.5
		}
		this.User.SumOptionDeals += money
		this.User.Balance += money
		this.User.Update()
		t.CloseDate = time.Now().Unix()
		t.Close = true
		t.Delete = true
		t.Balance = this.User.Balance + money
		t.Profit = money - t.InvestmentSum
		t.Update()
		redis.RedisClient.HDel("valuation_clear", t.Uuid.String()) //删除redis上的单
		unicast(this.User)                                         //移除前端已经删除的交易单
	}
	this.wsSend(&commandCode2066Return{
		MsgType:          "closeOptionPosition",
		OptionPositionID: positionID,
		MsgResult:        msgResult,
	})

}

type commandCode2066Return struct {
	MsgType          string `json:"msgType"` //closeOptionPosition
	OptionPositionID string `json:"optionPositionID"`
	MsgResult        string `json:"msgResult"` //Success
}

//{"commandCode":"2060","settingName":"EURUSD","optionType":"0","openPrice":"1.23162","direction":"0","investmentSum":"10","symbolName":"EURUSD","stopLine":"0","payoutPercentage":"75","earlyClosing":"0","interval":"0","period":"60","optionParamId":"745","optionSettingsId":"1"}
type commandCode2060Query struct { //OpenOptionPosition query
	CommandCode      string `json:"commandCode"`
	SettingName      string `json:"settingName"`
	OptionType       string `json:"optionType"`
	OpenPrice        string `json:"openPrice"`
	Direction        string `json:"direction"`
	InvestmentSum    string `json:"investmentSum"`
	SymbolName       string `json:"symbolName"`
	StopLine         string `json:"stopLine"`
	PayoutPercentage string `json:"payoutPercentage"`
	EarlyClosing     string `json:"earlyClosing"`
	Interval         string `json:"interval"`
	Period           string `json:"period"`
	OptionParamId    string `json:"optionParamId"`
	OptionSettingsId string `json:"optionSettingsId"`
}

func findOptionTypeKey(optionType int) (key string) {
	switch optionType {
	case 0:
		key = "Express"
	case 1:
		key = "Classic"
	case 2:
		key = "One Touch"
	case 3:
		key = "Range"
	}
	return
}

func (this *ws) check2060(t *models.Trade) (optionParamsCheckOk bool) {
	optionParamsCheckOk = false
	//提交的数据检查--start
	f64, ok := checkPriceNow(t.SymbolName, t.OpenPrice, time.Now())
	fmt.Println("~~~~~~~~~checkPriceNow f64:", f64, " ", t.OpenPrice, " ", ok)
	switch {
	case t.OptionType != 0 && t.OptionType != 1 && t.OptionType != 2 && t.OptionType != 3:
		optionParamsCheckOk = false
		return
	case ok == false || f64 == 0.0:
		optionParamsCheckOk = false
		return
	case t.InvestmentSum <= 0.0:
		optionParamsCheckOk = false
		return
	}

	switch t.OptionType {
	case 0:
		if t.Direction != 0 && t.Direction != 1 {
			optionParamsCheckOk = false
			return
		}
	case 1:
		if t.Direction != 2 && t.Direction != 3 {
			optionParamsCheckOk = false
			return
		}
	case 2:
		if t.Direction != 4 && t.Direction != 5 {
			optionParamsCheckOk = false
			return
		}
	case 3:
		if t.Direction != 6 && t.Direction != 7 {
			optionParamsCheckOk = false
			return
		}
	}

	it, ok := osaSyncMap.Load(fmt.Sprintf("%v%s", findOptionTypeKey(t.OptionType), t.SymbolName))
	if !ok {
		optionParamsCheckOk = false
		return
	}

	periodStr := strconv.Itoa(t.Period)
	payoutPercentageStr := strconv.FormatFloat(t.PayoutPercentage, 'f', 0, 64)
	earlyClosingStr := strconv.FormatFloat(float64(t.EarlyClosing), 'f', 0, 64)
	stopLineStr := strconv.Itoa(t.StopLine)
	intervalStr := strconv.Itoa(t.Interval)
	optionParamIdStr := strconv.Itoa(t.OptionParamId)

	//fmt.Println(periodStr, " ", payoutPercentageStr, " ", earlyClosingStr, " ", stopLineStr, " ", intervalStr)
	for _, v := range *it.(*optionSettingsArray).OptionParams {
		//fmt.Println(v.Period, " ", v.PayoutPercentage, " ", v.EarlyClosing, " ", v.StopLine, " ", v.Interval)
		if v.Id == optionParamIdStr && v.Period == periodStr && v.PayoutPercentage == payoutPercentageStr && v.EarlyClosing == earlyClosingStr && v.StopLine == stopLineStr && v.Interval == intervalStr {
			optionParamsCheckOk = true
			break
		}
	}
	fmt.Println("optionParamsCheckOk=", optionParamsCheckOk)
	//提交的数据检查--end
	return
}

var redisLockDepositkey string = "RedisLockDepositkey"

func (this *ws) commandCode2060Return() (err error) { //OpenOptionPosition
	if this.User.AllowTrade != 1 {
		err = errors.New("Not allow trade!")
		return
	}

	var cc2060 commandCode2060Query
	if err = json.Unmarshal(this.PByte, &cc2060); err != nil {
		return
	}

	var stopLine, earlyClosing, interval, period, optionParamId, optionSettingsId, optionType, direction int
	var openPrice, investmentSum, payoutPercentage float64

	if direction, err = strconv.Atoi(cc2060.Direction); err != nil {
		return
	}
	if optionType, err = strconv.Atoi(cc2060.OptionType); err != nil {
		return
	}
	if stopLine, err = strconv.Atoi(cc2060.StopLine); err != nil {
		return
	}
	if earlyClosing, err = strconv.Atoi(cc2060.EarlyClosing); err != nil {
		return
	}
	if interval, err = strconv.Atoi(cc2060.Interval); err != nil {
		return
	}
	if period, err = strconv.Atoi(cc2060.Period); err != nil {
		return
	}
	if optionParamId, err = strconv.Atoi(cc2060.OptionParamId); err != nil {
		return
	}
	if optionSettingsId, err = strconv.Atoi(cc2060.OptionSettingsId); err != nil {
		return
	}

	if openPrice, err = strconv.ParseFloat(cc2060.OpenPrice, 64); err != nil {
		return
	}

	if investmentSum, err = strconv.ParseFloat(cc2060.InvestmentSum, 64); err != nil {
		return
	}
	if payoutPercentage, err = strconv.ParseFloat(cc2060.PayoutPercentage, 64); err != nil {
		return
	}

	if err = this.User.FindUsersByUuid(this.User.Uuid.String()); err != nil {
		return
	}
	t := &models.Trade{
		UserUuid:         this.User.Uuid,
		Uuid:             uuid.NewV4(),
		SymbolName:       cc2060.SymbolName,
		OptionType:       optionType,
		OpenPrice:        openPrice,
		Direction:        direction,
		InvestmentSum:    investmentSum,
		StopLine:         stopLine,
		PayoutPercentage: payoutPercentage,
		EarlyClosing:     earlyClosing,
		Interval:         interval,
		Period:           period,
		OptionParamId:    optionParamId,
		OptionSettingsId: optionSettingsId,
		OpenDate:         time.Now().Unix(),
		Balance:          this.User.Balance - investmentSum,
	}
	if ok := this.check2060(t); !ok {
		this.wsSend(map[string]string{"msgType": "openOptionPosition", "msgResult": "Trade forbidden"})
		return
	}

	if this.User.Balance = this.User.Balance - investmentSum; this.User.Balance < 0 { //消费后 扣钱
		this.User.Balance = this.User.Balance + investmentSum
		redis.RedisClient.Del(redisLockDepositkey + this.User.Uuid.String())
		this.wsSend(map[string]string{"msgType": "openOptionPosition", "msgResult": "false"})
		return
	}
	if err = this.User.Update(); err != nil { //更新数据到数据库
		redis.RedisClient.Del(redisLockDepositkey + this.User.Uuid.String())
		return
	}

	redis.Store(this.User.Uuid.String(), this.User, time.Hour*24) //为多服务器并发考虑 cache最新更新的user数据

	//redis cache 代结算的单----start
	c := redis.Cache{
		Key:    "valuation_clear", //结算专用key
		Field:  t.Uuid.String(),
		It:     &t, // for cache,address only
		Expire: time.Second * time.Duration(period) * 2,
		//Buf    bytes.Buffer
		IsArr: false,
	}
	if err = c.Struct2RedisHSet(); err != nil {
		return
	}
	//redis cache 代结算的单----end

	if err = t.Insert(); err != nil {
		return
	}
	this.wsSend(map[string]string{"msgType": "openOptionPosition", "msgResult": "Success"})

	this.commandCode2062Return()

	return
}

type optionsPositionsSumsDetails_return struct {
	MsgType              string                         `json:"msgType"` //optionsPositions
	SumsDetails          *sumsDetailsPositions          `json:"sumsDetails"`
	OptionPositionsArray *[]optionPositionsArray_return `json:"optionPositionsArray"`
	MsgResult            string                         `json:"msgResult"`
}

type optionPositionsArray_return struct {
	Id               uuid.UUID `json:"id"`
	OptionType       string    `json:"optionType"`
	Direction        string    `json:"direction"`
	SettingName      string    `json:"settingName"`
	SymbolName       string    `json:"symbolName"`
	StopLine         string    `json:"stopLine"`
	PayoutPercentage string    `json:"payoutPercentage"`
	EarlyClosing     string    `json:"earlyClosing"`
	Interval         string    `json:"interval"`
	OpenPrice        string    `json:"openPrice"`
	OpenDate         string    `json:"openDate"`
	ExpiryDate       string    `json:"expiryDate"`
	InvestmentSum    string    `json:"investmentSum"`
}
