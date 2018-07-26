package server

import (
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"github.com/alecthomas/log4go"
)

type graceListener struct {
	originalLn			net.Listener

	// 是否存处在关闭的过程中
	isShuttingDown		int64

	// 用来标记是否可以close此listener, 当可以关闭时, 此channel会关闭
	canShutdownNow		chan struct{}

	// 当前连接数
	connCount			int64

	// 最大关闭等待时间
	maxWait				time.Duration
}


func newGraceListener(maxWait time.Duration, ln net.Listener) *graceListener {
	return &graceListener{
		originalLn: ln,
		isShuttingDown: 0,
		connCount: 0,
		maxWait: maxWait,
		canShutdownNow: make(chan struct{}),
	}
}

func (gln *graceListener) Accept() (net.Conn, error) {
	conn, err := gln.originalLn.Accept()
	if nil != err {
		return nil, err
	}

	// 连接数+1
	atomic.AddInt64(&gln.connCount, 1)

	return &graceConnection{
		Conn: conn,
		gln: gln,
	}, nil
}

func (gln *graceListener) Close() error {

	log4go.Info("waiting for pending connection to close")
	err := gln.waitAll()
	if nil != err {
		log4go.Warn("%v, force close", err)

	} else {
		log4go.Warn("all connections are gracefully closed")
	}

	// 这里是为了等待异步日志输出完成
	time.Sleep(500 * time.Millisecond)

	err = gln.originalLn.Close()
	if nil != err {
		return err
	}

	return err
}

func (gln *graceListener) waitAll() error {
	// 标记处于存在关闭状态
	atomic.AddInt64(&gln.isShuttingDown, 1)

	if atomic.LoadInt64(&gln.connCount) == 0 {
		close(gln.canShutdownNow)
		return nil
	}

	select {
	case <-gln.canShutdownNow:
		return nil

	case <-time.After(gln.maxWait):
		return fmt.Errorf("cannot close all connections after %s", gln.maxWait)
	}

	return nil
}

func (gln *graceListener) closeConnection() {
	// 连接数-1
	atomic.AddInt64(&gln.connCount, -1)

	// 如果当前处理正在关闭状态, 且连接数为0
	if gln.isShuttingDown != 0 && gln.connCount == 0 {
		close(gln.canShutdownNow)
	}
}

func (gln *graceListener) Addr() net.Addr {
	return gln.originalLn.Addr()
}

type graceConnection struct {
	net.Conn

	// 保存Listener指针
	gln *graceListener
}

func (c *graceConnection) Close() error {
	err := c.Conn.Close()
	if nil != err {
		return nil
	}

	c.gln.closeConnection()
	return nil
}
