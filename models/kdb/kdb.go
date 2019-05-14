package kdb

import (
	"strconv"
	"strings"
	//	"errors"
	"fmt"

	"github.com/sv/kdbgo"
	"option.bzza.com/system"
)

var conn *kdb.KDBConn

func init() {
	var err error
	var kdbServer1Slice []string
	if kdbServer1Slice = strings.Split(system.Conf.Kdb.Server1, ":"); len(kdbServer1Slice) != 2 {
		panic("config.yaml kdb server1 error")
	}
	var port int
	if port, err = strconv.Atoi(kdbServer1Slice[1]); err != nil {
		panic("config.yaml kdb server1 error")
	}

	if conn, err = kdb.DialKDB(kdbServer1Slice[0], port, ""); err != nil {
		panic("Failed to connect kdb+:" + err.Error())
	}

	var dataK *kdb.K
	if dataK, err = conn.Call("tables[]"); err != nil {
		panic(err)
	}

	ok := false
	if dataK.Len() >= 1 {
		for _, v := range dataK.Data.([]string) {
			if v == "data" {
				ok = true
			}
		}
	}

	if !ok { //table data is exist
		if _, err = conn.Call("data:([]date:`datetime$();sym:`symbol$();bid:`float$();ask:`float$())"); err != nil {
			panic("kdb+ Query failed init:" + err.Error())
		}
	}

}

func Call(str string) (data *kdb.K, err error) {
	//data:([[]date:`datetime$()];sym:`symbol$();bid:`float$();ask:`float$())
	//`data insert (2013.07.01T10:03:54.348z;`IBM;20.83f;20.88f)
	if data, err = conn.Call(str); err != nil {
		fmt.Println("kdb+ Query failed:", err, " ", str)
	}
	return
}

type Data struct {
	Symbol string
	Bid    string
	Ask    string
	Date   string
}

func (d *Data) InsertIfNotExist() (data *kdb.K, err error) {
	var kdbStr string
	//	kdbStr = fmt.Sprintf("select [1]bid from data where date=%s,bid=%s,ask=%s,sym=`%s", d.Date, d.Bid, d.Ask, d.Symbol)
	//	data, err = Call(kdbStr)
	//	if err != nil {
	//		return
	//	} else if data.Len() == 1 { //已经存在
	//		return
	//	}

	kdbStr = fmt.Sprintf("`data insert (%s;`%s;%vf;%vf)", d.Date, d.Symbol, d.Bid, d.Ask)
	if data, err = Call(kdbStr); err != nil {
		return
	}
	return
}
