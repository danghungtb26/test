package utilities

import (
	"sync"

	"git.ctisoftware.vn/back-end/utilities/data/provider/log"
)

var (
	_logger log.Logger
	m       sync.Mutex
)

func SetLogger(logger log.Logger) {
	m.Lock()
	defer m.Unlock()
	_logger = logger
}

func Logger() log.Logger {
	m.Lock()
	defer m.Unlock()
	if _logger == nil {
		_logger = log.NoopLogger
	}
	return _logger
}
