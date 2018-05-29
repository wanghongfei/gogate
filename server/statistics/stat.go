package stat

import (
	"time"

	"github.com/alecthomas/log4go"
)

// 流量记录器
type TraficStat struct {
	store			TraficInfoStore

	// trafic信息缓冲channel
	bufferChan		chan *TraficInfo
	writeChan		chan map[string]*TraficInfo
	// 发送间隔, 秒
	writeInterval	int
}

// 创建统计器
// queueSize: 最大可以放多少条流量对象不会block
// interval: 每多少秒调用一次存储对象发送这期间积累的数据
// traficStore: 数据保存逻辑的实现, 如CsvFileTrafficStore
func NewTrafficStat(queueSize, interval int, traficStore TraficInfoStore) *TraficStat {
	if interval < 1 {
		interval = 1
	}

	return &TraficStat{
		bufferChan: make(chan *TraficInfo, queueSize),
		writeChan: make(chan map[string]*TraficInfo, interval + 1),
		writeInterval: interval,

		store: traficStore,
	}
}

// 启动流量记录routine
func (ts *TraficStat) StartRecordTrafic() {
	// 启动统计routine
	go ts.traficAggregateRoutine()
	// 启动写日志任务routine
	go ts.traficLogRoutine()
}

// 记录流量
func (ts *TraficStat) RecordTrafic(info *TraficInfo) {
	// 验证
	if nil == info || info.SuccessCount < 0 || info.FailedCount < 0 {
		// 无效数据丢弃
		return
	}

	ts.bufferChan <- info
}

// 每ts.writeInterval秒累计一次此时间段内的流量信息, 封装成写任务扔到writeChan中
func (ts *TraficStat) traficAggregateRoutine() {
	ticker := time.NewTicker(time.Second * time.Duration(ts.writeInterval))

	for {
		<- ticker.C

		// 取出当前channel全部元素
		size := len(ts.bufferChan)
		if 0 == size {
			// 上一个时间周期内没有元素
			// skip
			continue
		}

		// 统计在此时间周期里的数据之和
		sumMap := make(map[string]*TraficInfo)
		for ix := 0; ix < size; ix++ {
			elem := <- ts.bufferChan

			targetInfo, exist := sumMap[elem.ServiceId]
			if !exist {
				targetInfo = new(TraficInfo)
				targetInfo.timestamp = time.Now().UnixNano() / 1e6
				targetInfo.ServiceId = elem.ServiceId
				sumMap[elem.ServiceId] = targetInfo
			}

			targetInfo.FailedCount += elem.FailedCount
			targetInfo.SuccessCount += elem.SuccessCount
		}

		ts.writeChan <- sumMap
	}
}

func (ts *TraficStat) traficLogRoutine() {
	for servMap := range ts.writeChan {
		for _, traffic := range servMap {
			err := ts.store.Send(traffic)
			if nil != err {
				log4go.Error(err)
			}
		}

	}

}

// 定义流量信息
type TraficInfo struct {
	ServiceId		string
	SuccessCount	int
	FailedCount		int

	// unix毫秒数
	timestamp		int64
}

// 定义流量数据存储方式
type TraficInfoStore interface {
	// 发送trafic数据
	Send(info *TraficInfo) error
	// 清理资源
	Close() error
}


