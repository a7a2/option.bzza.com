package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

const (
	ErrUnkown      = "unkown"
	ErrEmail       = "email error"
	ErrEmailExists = "email already exists"
	ErrCommit      = "Commit error"
)

type GetSymbols struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Data    Symbols `json:"data"`
}
type Symbols struct {
	Title  string  `json:"title"`
	Symbol string  `json:"symbol"`
	Bid    float64 `json:"bid"`
	Ask    float64 `json:"ask"`
	Unix   int64   `json:"unix"`
}

type Users struct { // db user
	Uuid             uuid.UUID `gorm:"primary_key;column:uuid"`
	Login            int       `gorm:"AUTO_INCREMENT;unique_index;default:10001;type:integer;column:login"`
	Password         string    `gorm:"default:null;type:character(32);column:password"`
	InvestorPassword string    `gorm:"default:null;type:character(32);column:investorPassword"`
	//Money              float64   `json:"money" gorm:"column:money;type:numeric(11,2);default:0.00"`
	Name               string    `gorm:"default:null;type:character varying(20);column:name"`
	Patronimic         string    `gorm:"default:null;type:character varying(20);column:patronimic"`
	Surname            string    `gorm:"default:null;type:character varying(30);column:surname"`
	Country            string    `gorm:"default:null;type:character varying(30);column:country"`
	Region             string    `gorm:"default:null;type:character varying(30);column:region"`
	City               string    `gorm:"default:null;type:character varying(30);column:city"`
	Email              string    `gorm:"unique_index;type:character varying(100);default:null;column:eMail"` //邮箱
	Address            string    `gorm:"default:null;type:character varying(150);column:address"`
	DepositCurrency    string    `gorm:"default:'usd';type:character varying(20);column:depositCurrency"` //usd
	Phone              int64     `gorm:"unique_index;type:bigint;default:null;column:phone"`
	Index              int       `gorm:"column:index;default:0"`
	Balance            float64   `gorm:"column:balance;type:numeric(12,2);default:0.00"`
	SumOptionDeals     float64   `gorm:"column:sumOptionDeals;type:numeric(12,2);default:0.00"` //累计开单投资金额
	Delay              int       `gorm:"column:delay;default:0"`
	Leverage           int       `gorm:"column:leverage";default:100`                                //杠杠
	RealAccount        int       `gorm:"column:realAccount;default:0"`                               //0、1
	GroupID            int       `gorm:"column:groupID;default:2"`                                   //2期权
	GroupName          string    `gorm:"column:groupName;type:character varying(10);default:'Demo'"` // Real、Demo Test Admin
	CreatedAt          time.Time `json:"createdAt" gorm:"type:timestamp(6) with time zone;column:createdAt"`
	UpdatedAt          time.Time `json:"updatedAt" gorm:"type:timestamp(6) with time zone;column:updatedAt;default:null"`
	LastLoginDate      time.Time `json:"lastLoginDate" gorm:"type:timestamp(6) with time zone;column:lastLoginDate"`
	MailIndex          int       `json:"mailIndex" gorm:"column:mailIndex;default:0"`
	AllowTrade         int       `json:"allowTrade" gorm:"column:allowTrade;default:1"`
	Stopout            int       `json:"stopout" gorm:"column:stopout;default:20"`
	DepositPercent     int       `json:"depositPercent" gorm:"column:depositPercent;default:0"`
	MinDeposit         float64   `json:"minDeposit" gorm:"column:minDeposit;type:numeric(6,2);default:0.00"`
	MaxOrderCount      int       `json:"maxOrderCount" gorm:"column:maxOrderCount;default:30"`
	RegistrationIP     string    `json:"registrationIP" gorm:"column:registrationIP;type:inet;default:null"`
	PriceStrategy      int       `json:"priceStrategy" gorm:"column:priceStrategy;default:0"`
	CoverageStrategyID int       `json:"coverageStrategyID" gorm:"column:coverageStrategyID;default:1"`
	Showbonuses        float64   `json:"showbonuses" gorm:"column:showbonuses;type:numeric(10,2);default:0.00"`
	EquityMarginCall   int       `json:"equityMarginCall" gorm:"column:equityMarginCall;default:1"`
}

type Trade struct {
	UserUuid         uuid.UUID `gorm:"column:userUuid;index"`
	Uuid             uuid.UUID `gorm:"primary_key;column:uuid"`
	Balance          float64   `gorm:"default:0.0;type:numeric(10,2);column:balance"`
	SymbolName       string    `gorm:"default:null;type:character varying(20);column:symbolName"`
	OptionType       int       `gorm:"column:optionType;default:0"`   //0是express
	OpenPrice        float64   `gorm:"default:null;column:openPrice"` //开单 价格
	ClosePrice       float64   `gorm:"default:null;column:closePrice"`
	Profit           float64   `gorm:"default:0.0;type:numeric(10,2);column:profit"`         //输赢实际金额
	Direction        int       `gorm:"column:direction;default:0"`                           //趋势 买升还是跌
	InvestmentSum    float64   `gorm:"default:null;type:numeric(10,2);column:investmentSum"` //投资金额
	StopLine         int       `gorm:"column:stopLine;default:0"`
	PayoutPercentage float64   `gorm:"default:null;type:numeric(10,2);column:payoutPercentag"` //盈利 实际 得到赔率
	EarlyClosing     int       `gorm:"column:earlyClosing;default:0"`
	Interval         int       `gorm:"column:interval;default:0"`
	Period           int       `gorm:"column:period;default:60"` //结算 的时间周期（秒）
	OptionParamId    int       `gorm:"column:optionParamId;default:0"`
	OptionSettingsId int       `gorm:"column:optionSettingsId;default:0"`
	OpenDate         int64     `gorm:"default:null;column:openDate;index"` //开单时间
	CloseDate        int64     `gorm:"default:null;column:closeDate;index"`
	Close            bool      `gorm:"column:close;default:false;index"` //true为结算了
	Delete           bool      `gorm:"column:delete;default:false;index"`
}

//{"msgType":"deposits","depositsArray":[{"depositID":"12147","equity":"5000","status":"0","operationDate":"1522667655","comment":"First deposit"}],"msgResult":"Success"}
type Money struct {
	UserUuid      uuid.UUID `gorm:"column:userUuid;index"`
	Uuid          uuid.UUID `gorm:"primary_key;column:uuid"`
	Equity        float64   `gorm:"default:null;type:numeric(10,2);column:equity"`          //金额
	Status        int       `gorm:"column:status;default:0"`                                //0存款
	OperationDate int64     `gorm:"default:null;column:operationDate;index"`                //操作时间
	Comment       string    `gorm:"default:null;type:character varying(20);column:comment"` //注释
}

func Clear(t *Trade, u *Users) {
	g := DB.Begin()
	if g.Save(t).Error != nil {
		g.Rollback()
	}

	if g.Save(u).Error != nil {
		g.Rollback()
	}
	g.Commit()

}

func (user *Users) Insert() error {
	//db.Set("gorm:query_option", "FOR UPDATE").First(&user, 10)
	var iCount int
	g := DB.Model(&Users{}).Count(&iCount)
	if g.Error != nil {
		iCount = 0
	}

	user.Login = iCount + 10001
	return DB.Create(user).Error

}

// update user
func (user *Users) Update() (err error) {
	err = DB.Save(user).Error
	return
}

func (user *Users) UpdateBalance(f64 float64) error {
	return DB.Model(user).UpdateColumns(map[string]interface{}{
		"deposit": gorm.Expr(fmt.Sprintf("deposit + %f", f64)),
	}).Error
}

func (user *Users) FindUsersByUuid(uuidStr string) error {
	g := DB.First(&user, "uuid = ?", uuidStr)
	return g.Error
}

func (user *Users) UserAuth(login int, password string) error {
	g := DB.First(&user, "login = ? and password = ?", login, password)
	return g.Error
}

func (t *Trade) Insert() error {
	return DB.Create(t).Error
}

func (t *Trade) Update() error {
	return DB.Save(t).Error
}

func (t *Trade) FindTradeByUuid(uuidStr string) error {
	g := DB.First(t, "uuid = ?", uuidStr)
	return g.Error
}

func (m *Money) Insert() error {
	return DB.Create(m).Error
}

func (m *Money) RealDelete() error {
	return DB.Delete(m).Error
}

func ListTrade(where string) (tArr []Trade, err error) {
	err = DB.Where(where).Find(&tArr).Error
	return
}

func ListMoney(where string) (mArr []Money, err error) {
	err = DB.Where(where).Find(&mArr).Error
	return
}
