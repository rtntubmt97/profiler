package intervalListeners

import (
	"testing"
	"time"
)

func TestConsoleLog(t *testing.T) {
	consoleLog := NewConsolLog()
	consoleLog.Listen(genProfiles(10), time.Now(), 0)

	time.Sleep(time.Second * 2)
	consoleLog.Listen(genProfiles(10), time.Now(), 0)
}

func TestFileLog(t *testing.T) {
	fileLog, file := NewFileLog("/home/tumd/golang-repositories/profiler/test/out.txt")

	fileLog.Listen(genProfiles(0), time.Now(), 0)
	time.Sleep(time.Second * 2)
	fileLog.Listen(genProfiles(10), time.Now(), 0)
	time.Sleep(time.Second * 2)

	file.Sync()
}
