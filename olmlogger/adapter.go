package olmlogger

import (
	"github.com/CubicrootXYZ/gologger"
	"maunium.net/go/mautrix/crypto"
)

// New assembles a new logger for olm encryption.
func New(logger gologger.Logger) crypto.Logger {
	return &olmlogger{
		logger: logger,
	}
}

type olmlogger struct {
	logger gologger.Logger
}

func (logger *olmlogger) Error(message string, args ...interface{}) {
	logger.logger.Errorf(message, args...)
}

func (logger *olmlogger) Warn(message string, args ...interface{}) {
	logger.logger.Errorf(message, args...)
}

func (logger *olmlogger) Debug(message string, args ...interface{}) {
	logger.logger.Errorf(message, args...)
}

func (logger *olmlogger) Trace(message string, args ...interface{}) {
	logger.logger.Errorf(message, args...)
}
