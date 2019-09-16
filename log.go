package main

import (
	"io"
	"os"

	"github.com/btcsuite/btclog"
	"github.com/yarikbratashchuk/cafebot/chat"
	"github.com/yarikbratashchuk/cafebot/order"
)

var log btclog.Logger

func fatalf(format string, params ...interface{}) {
	log.Criticalf(format, params...)
	os.Exit(1)
}

func setupLog(dest io.Writer, loglevel string) {
	logBackend := btclog.NewBackend(dest)
	lvl, _ := btclog.LevelFromString(loglevel)

	orderLog := logBackend.Logger("ORDR")
	chatLog := logBackend.Logger("CHAT")
	log = logBackend.Logger("SRVR")

	orderLog.SetLevel(lvl)
	chatLog.SetLevel(lvl)
	log.SetLevel(lvl)

	chat.UseLogger(chatLog)
	order.UseLogger(orderLog)
}
