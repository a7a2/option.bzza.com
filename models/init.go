package models

import (
	"fmt"

	"strings"

	"regexp"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/lib/pq"

	"option.bzza.com/models/redis"
	"option.bzza.com/system"
	//"option.bzza.com/system/getSslCer"
)

func init() {
	//getSslCer.GetSslCer()
	system.YamlInit()
	redis.InitRedis()
	initDB()
	initRegexp()
	InitReplaceConfigJs()

	go system.FileChangeWatcher()

}

var DB *gorm.DB

func initDB() {
	var err error
	DB, err = gorm.Open(system.Conf.Db.DriverName, system.Conf.Db.DataSourceName)
	if err != nil {
		panic(err)
	}

	DB.LogMode(true)
	DB.CreateTable(&Users{})
	DB.CreateTable(&Trade{})
	DB.CreateTable(&Money{}) //Money
	//DB.AutoMigrate(&Users{}) //AutoMigrate will ONLY create tables, missing columns and missing indexes, and WON'T change existing column's type or delete unused columns to protect your data.
	//DB.Model(&Users{}).AddUniqueIndex("member_unique", "eMail", "phone", "login")
}

var ReEmail, RePortConfigJs *regexp.Regexp

func initRegexp() {
	ReEmail = regexp.MustCompile(`[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+`)
	RePortConfigJs = regexp.MustCompile(`(Port: '[0-9]{0,5}'\,)`)
	//fmt.Println(ReEmail.MatchString("me@8aa2.com"))
}

func InitReplaceConfigJs() {
	bytes := new([]byte)
	bytes = system.ReadFile(system.Conf.Web.JsConfigPath)
	sl := RePortConfigJs.FindAllString(string(*bytes), -1)
	var oldStr string
	for _, v := range sl {
		oldStr = fmt.Sprintf("%s%s", oldStr, v)
	}

	newStr := "Port: '" + system.Conf.Web.Port + "',"

	if newStr != oldStr {
		configStr := strings.Replace(string(*bytes), oldStr, newStr, 1)
		err := system.WriteFile(system.Conf.Web.JsConfigPath, []byte(configStr))
		if err != nil {
			panic(err)
		}
	}
}

type XCaDeviceid struct {
	Symbol                  string
	Cacloudmarketinstanceid string
}
