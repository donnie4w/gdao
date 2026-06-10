// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package base

import (
	"github.com/donnie4w/go-logger/logger"
)

var _logging = logger.NewLogger().SetFormat(logger.FORMAT_LEVELFLAG)

var _logger = new(log)
var Logger = _logger

type log struct {
	IsVaild bool
}

func (l *log) SetLogger(on bool) {
	l.IsVaild = on
}

func (l *log) Debug(v ...interface{}) {
	if l.IsVaild {
		_logging.Debug(v...)
	}
}

func (l *log) Debugf(format string, v ...interface{}) {
	if l.IsVaild {
		_logging.Debugf(format, v...)
	}
}

func (l *log) Info(v ...interface{}) {
	if l.IsVaild {
		_logging.Info(v...)
	}
}

func (l *log) Infof(format string, v ...interface{}) {
	if l.IsVaild {
		_logging.Infof(format, v...)
	}
}

func (l *log) Warn(v ...interface{}) {
	if l.IsVaild {
		_logging.Warn(v...)
	}
}

func (l *log) Warnf(format string, v ...interface{}) {
	if l.IsVaild {
		_logging.Warnf(format, v...)
	}
}

func (l *log) Error(v ...interface{}) {
	if l.IsVaild {
		_logging.Error(v...)
	}
}

func (l *log) Errorf(format string, v ...interface{}) {
	if l.IsVaild {
		_logging.Errorf(format, v...)
	}
}
