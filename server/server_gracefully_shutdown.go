package server

import (
	"net"
	"sync/atomic"
	"time"

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
	err := gln.originalLn.Close()
	if nil != err {
		return err
	}

	gln.markShutdown()

	return err
}

func (gln *graceListener) markShutdown() {
	// 标记处于正在关闭状态
	atomic.AddInt64(&gln.isShuttingDown, 1)

	if atomic.LoadInt64(&gln.connCount) == 0 {
		close(gln.canShutdownNow)
	}
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
