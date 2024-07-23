// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package base

import (
	"fmt"
	"github.com/donnie4w/simplelog/logging"
)

var _logging = logging.NewLogger().SetFormat(logging.FORMAT_LEVELFLAG)

var logger = new(log)
var Logger = logger

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
		_logging.Debug(fmt.Sprintf(format, v...))
	}
}

func (l *log) Info(v ...interface{}) {
	if l.IsVaild {
		_logging.Info(v...)
	}
}

func (l *log) Infof(format string, v ...interface{}) {
	if l.IsVaild {
		_logging.Info(fmt.Sprintf(format, v...))
	}
}

func (l *log) Warn(v ...interface{}) {
	if l.IsVaild {
		_logging.Warn(v...)
	}
}

func (l *log) Warnf(format string, v ...interface{}) {
	if l.IsVaild {
		_logging.Warn(fmt.Sprintf(format, v...))
	}
}

func (l *log) Error(v ...interface{}) {
	if l.IsVaild {
		_logging.Error(v...)
	}
}

func (l *log) Errorf(format string, v ...interface{}) {
	if l.IsVaild {
		_logging.Error(fmt.Sprintf(format, v...))
	}
}
