package stat

import (
	"testing"
	"time"
)

func TestNewCsvFileTraficInfoStore(t *testing.T) {
	info := &TraficInfo{
		ServiceId: "user-service",
		Success: true,
		Timestamp: time.Now().UnixNano() / 10e6,
	}

	cf := NewCsvFileTraficInfoStore("/tmp")


	err := cf.Send(info)
	if nil != err {
		t.Error(err)
		return
	}

	err = cf.Close()
	if nil != err {
		t.Error(err)
	}
}