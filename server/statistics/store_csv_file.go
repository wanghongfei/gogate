package stat

import (
	"bytes"
	"errors"
	"fmt"
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

func NewCsvFileTraficInfoStore(logDir string) *CsvFileTraficInfoStore {
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

func (fs *CsvFileTraficInfoStore) Close() error {
	errMsg := ""

	for _, file := range fs.fileMap {
		closeErr := file.Close()
		if nil != closeErr {
			fmt.Sprintf("%s%s;", errMsg, closeErr.Error())
		}
	}

	if "" != errMsg {
		return errors.New(errMsg)
	}

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
	buf.WriteString(strconv.FormatInt(info.timestamp, 10))
	buf.WriteString(",")

	buf.WriteString(strconv.Itoa(info.SuccessCount))
	buf.WriteString(",")

	buf.WriteString(strconv.Itoa(info.FailedCount))
	buf.WriteString(",")

	buf.WriteString(info.ServiceId)
	buf.WriteString("\n")

	return &buf
}


