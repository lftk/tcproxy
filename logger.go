package main

import (
	"sync/atomic"

	"github.com/golang/glog"
)

var _id int64

type logger int64

func trafficLogger() logger {
	return logger(atomic.AddInt64(&_id, 1))
}

func (l logger) TraffIn(b []byte) {
	glog.Infof("[IN#%d]%s", int64(l), string(b))
}

func (l logger) TraffOut(b []byte) {
	glog.Infof("[OUT#%d]%s", int64(l), string(b))
}
