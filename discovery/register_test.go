package discovery

import (
	"testing"
	"time"
)

func TestStartRegister(t *testing.T) {
	StartRegister()
	time.Sleep(time.Second * 60)
}