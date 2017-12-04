package logger

import (
	"testing"
	"time"
)

var log = NewJiaLog()

func TestDebug(t *testing.T) {
	log.SetEnableConsole(true)
	log.SetEnableLog(true)
	log.SetLogPath("./" + time.Now().Format("2006-01-02"))
	log.Debug("hello", "nihao")
	log.Info("good", "boy")
	time.Sleep(2 * time.Second)

}
