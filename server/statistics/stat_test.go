package stat

import (
	"fmt"
	"testing"
	"time"
)

func TestNewCsvFileTraficInfoStore(t *testing.T) {
	info := &TraficInfo{
		ServiceId: "user-service",
		SuccessCount: 10,
		FailedCount: 1,
		timestamp: time.Now().UnixNano() / 10e6,
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

func TestNewStraficStat(t *testing.T) {
	stat := NewStraficStat(10, 1, NewCsvFileTraficInfoStore("/tmp"))
	stat.StartRecordTrafic()

	ticker := time.NewTicker(time.Millisecond * 400)
	count := 0
	for {
		<- ticker.C

		info := &TraficInfo{
			ServiceId: "dog-service",
			SuccessCount: 1,
			FailedCount: 0,
		}

		stat.RecordTrafic(info)
		fmt.Println("put")

		count ++
		if count > 10 {
			break
		}
	}
}