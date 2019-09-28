package chat

import (
	"github.com/btcsuite/btclog"
)

var log btclog.Logger

func init() { DisableLog() }

//DisableLog ok
func DisableLog() { log = btclog.Disabled }

//UseLogger ok
func UseLogger(logger btclog.Logger) {
	log = logger
}
