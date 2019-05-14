package controllers

import (
	"sync"

	"option.bzza.com/system"
)

func init() {
	osaSyncMap = new(sync.Map)
	go OptionSettingsArrayInit()
	go LoopGetOptionSettingsFromRedis()
	go WsGetSymbolsFromA7a2()
	go loopDelConn() //删除websocket连接
	go calLoop()     //结算
	system.ChanOptionSettingsArrayInit <- true
}
