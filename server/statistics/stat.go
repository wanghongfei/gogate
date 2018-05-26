package stat

import (
	"github.com/alecthomas/log4go"
)

// 流量记录器
type TraficStat struct {
	store		TraficInfoStore
	traficChan	chan *TraficInfo
}

func NewStraficStat(logDir string, queueSize int, traficStore TraficInfoStore) *TraficStat {
	return &TraficStat{
		traficChan: make(chan *TraficInfo, queueSize),
		store: traficStore,
	}
}

// 启动流量记录routine
func (ts *TraficStat) StartRecordTrafic() {
	go ts.traficLogRoutine()
}

// 记录流量
func (ts *TraficStat) RecordTrafic(info *TraficInfo) {
	// 验证
	if nil == info || info.ServiceId == "" || info.Timestamp == 0 {
		// 无效数据丢弃
		return
	}

	ts.traficChan <- info
}

func (ts *TraficStat) traficLogRoutine() {
	for trafic := range ts.traficChan {
		err := ts.store.Send(trafic)
		if nil != err {
			log4go.Error(err)
		}
	}

}

// 定义流量信息
type TraficInfo struct {
	ServiceId		string
	Success			bool
	Timestamp		int
}

// 定义流量数据存储方式
type TraficInfoStore interface {
	Send(info *TraficInfo) error
}


