package stat

import (
	"bytes"
	"os"
	"strconv"
)

// 文件流量存储器
type CsvFileTraficInfoStore struct {
	// 流量日志文件所在目录
	logDir		string

	// serviceId(string) -> *File
	fileMap		map[string]*os.File
}

func NewFileTraficInfoStore(logDir string) *CsvFileTraficInfoStore {
	return &CsvFileTraficInfoStore{
		logDir: logDir,
		fileMap: make(map[string]*os.File),
	}
}

func (fs *CsvFileTraficInfoStore) Send(info *TraficInfo) error {
	buf := fs.ToCsv(info)
	f, err := fs.getFile(info.ServiceId)
	if nil != err {
		return err
	}

	buf.WriteTo(f)

	return nil
}

// 从map中取出日志文件, 如果没有则打开
func (fs *CsvFileTraficInfoStore) getFile(servId string) (*os.File, error) {
	logFile, exist := fs.fileMap[servId]
	if !exist {
		// 不存在则创建
		f, err := os.OpenFile(fs.logDir + "/" + servId + ".log", os.O_CREATE | os.O_APPEND | os.O_RDWR, 0644)
		if nil != err {
			return nil, err
		}

		logFile = f
		fs.fileMap[servId] = f
	}

	return logFile, nil
}

func (fs *CsvFileTraficInfoStore) ToCsv(info *TraficInfo) *bytes.Buffer {
	var buf bytes.Buffer
	buf.WriteString(info.ServiceId)
	if info.Success {
		buf.WriteString("1")
	} else {
		buf.WriteString("0")
	}
	buf.WriteString(strconv.Itoa(info.Timestamp))
	buf.WriteString("\n")

	return &buf
}


