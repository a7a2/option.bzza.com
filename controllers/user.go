package controllers

import (
	//	"encoding/json"

	"strconv"
	"time"

	"net/http"
	//	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"option.bzza.com/models"
	"option.bzza.com/models/redis"

	"option.bzza.com/helpers"
)

//{"name":"jialiang","patronimic":"","surname":"kang","country":"中国","region":"demo","city":"佛山","eMail":"me1@a7a2.com","address":"demo","phone":"13516521122","index":"0","deposit":5000,"delay":"0","leverage":"100","realAccount":"0","groupID":"2","groupName":"Demo"}:
type postReg struct { // post reg
	Name        string `formam:"name" json:"name" form:"name" binding:"required"`
	Patronimic  string `formam:"patronimic" json:"patronimic" form:"patronimic" binding:"required"`
	Surname     string `formam:"surname" json:"surname" form:"surname" binding:"required"`
	Country     string `formam:"country" json:"country" form:"country" binding:"required"`
	Region      string `formam:"region" json:"region" form:"region" binding:"required"`
	City        string `formam:"city" json:"city" form:"city" binding:"required"`
	Email       string `formam:"eMail" json:"eMail" form:"eMail" binding:"required"`
	Address     string `formam:"address" json:"address" form:"address" binding:"required"`
	Phone       string `formam:"phone" json:"phone" form:"phone" binding:"required"`
	Index       string `formam:"index" json:"index" form:"index" binding:"required"` //number
	Deposit     int    `formam:"deposit" json:"deposit" form:"deposit" binding:"required"`
	Delay       string `formam:"delay" json:"delay" form:"delay" binding:"required"`
	Leverage    string `formam:"leverage" json:"leverage" form:"leverage" binding:"required"`          //杠杠
	RealAccount string `formam:"realAccount" json:"realAccount" form:"realAccount" binding:"required"` //0、1
	GroupID     string `formam:"groupID" json:"groupID" form:"groupID" binding:"required"`             //2期权
	GroupName   string `formam:"groupName" json:"groupName" form:"groupName" binding:"required"`       // Real、Demo
}

//regResult {"msgType":"authorizeData","login":"11019","password":"HTfLXF","investorPassword":"GKQIR4","msgResult":"Success"}
func PostReg(c *gin.Context) {
	//c.JSON(http.StatusOK, gin.H{"msgResult": "Success"})
	var reg postReg
	b, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	err = simpleDecodeUnmarshal(b, &reg)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	switch { //检查提交内容是否符合规则
	case models.ReEmail.MatchString(reg.Email) == false || (reg.RealAccount == "1" && reg.Deposit != 0):
		c.JSON(http.StatusOK, gin.H{"error": models.ErrEmail})
	}
	password := helpers.GetRandomString(6)
	investorPassword := helpers.GetRandomString(6)

	var phone int64
	var index, delay, leverage, realAccount, groupID int
	var deposit float64
	if phone, err = strconv.ParseInt(reg.Phone, 10, 64); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	if index, err = strconv.Atoi(reg.Index); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	if deposit, err = strconv.ParseFloat(strconv.Itoa(reg.Deposit), 64); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	if delay, err = strconv.Atoi(reg.Delay); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	if leverage, err = strconv.Atoi(reg.Leverage); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	if realAccount, err = strconv.Atoi(reg.RealAccount); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	if groupID, err = strconv.Atoi(reg.GroupID); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	varUuid := uuid.NewV4()
	tNow := time.Now()
	user := models.Users{
		Uuid:             varUuid,
		Password:         helpers.EncryptPwd2Db([]byte(password)),
		InvestorPassword: helpers.EncryptPwd2Db([]byte(investorPassword)),
		Name:             reg.Name,
		Patronimic:       reg.Patronimic,
		Surname:          reg.Surname,
		Country:          reg.Country,
		Region:           reg.Region,
		City:             reg.City,
		Email:            reg.Email,
		Address:          reg.Address,
		Phone:            phone,
		Index:            index,
		Balance:          deposit,
		Delay:            delay,
		Leverage:         leverage,
		RealAccount:      realAccount,
		GroupID:          groupID,
		GroupName:        reg.GroupName,
		CreatedAt:        tNow,

		RegistrationIP: c.ClientIP(),
	}

	m := models.Money{
		UserUuid:      varUuid,
		Uuid:          uuid.NewV4(),
		Equity:        deposit,
		Status:        0,
		OperationDate: tNow.Unix(),
		Comment:       "Deposit",
	}

	if err = m.Insert(); err != nil {
		m.RealDelete()
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	if err = user.Insert(); err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	if err == nil {
		r := regResult{
			MsgType:          "authorizeData",
			Login:            user.Login,
			Password:         password,
			InvestorPassword: investorPassword,
			MsgResult:        "Success",
		}
		c.JSON(http.StatusOK, &r)
		return
	}

}

//regResult {"msgType":"authorizeData","login":"11019","password":"HTfLXF","investorPassword":"GKQIR4","msgResult":"Success"}
type regResult struct {
	MsgType          string `json:"msgType"`
	Login            int    `json:"login"`
	Password         string `json:"password"`
	InvestorPassword string `json:"investorPassword"`
	MsgResult        string `json:"msgResult"`
}

type authResp struct {
	Auth  string `json:"auth"`
	Token string `json:"token"`
}

func GetLogin(c *gin.Context) { //auth //RETURN {"auth":"yes","token":"68E1E21377164081B75C10BE2C131CE1"}
	loginNum, password := c.Param("loginNum"), c.Param("password")
	loginInt, err := strconv.Atoi(loginNum)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	user := &models.Users{}
	err = user.UserAuth(loginInt, helpers.EncryptPwd2Db([]byte(password)))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	} else if user.Email != "" { //ok登陆成功
		token := helpers.MD5([]byte(loginNum + user.LastLoginDate.String() + user.Uuid.String()))
		//user.Password = "" //存入reis前去掉密码
		user.LastLoginDate = time.Now()
		redis.RedisClient.Set(token, user.Uuid.String(), time.Hour*24)
		resp := authResp{Auth: "yes", Token: token}
		c.JSON(http.StatusOK, resp)
	}
}
