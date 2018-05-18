package mobileLog

var (
	logname = "session.log"

	mlog MLog
)

type MLog interface {
	InitMobileLog(logPath string) error

	CloseMobileLog()

	Print(v ...interface{})
}

func NewMobileLog(logPath string) error {
	l, err := NewImpLog(logPath)
	if err == nil {
		mlog = l
	}
	return err
}

func Print(v ...interface{}) {
	mlog.Print(v)
}

func CloseMobileLog() {
	if mlog != nil {
		mlog.CloseMobileLog()
	}
}
