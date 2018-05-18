package mobileLog

import (
	//"fmt"
	"log"
	"os"
	//"time"
)

type ImpLog struct {
	mf *os.File

	ml *log.Logger
}

func NewImpLog(logPath string) (*ImpLog, error) {
	mlog := &ImpLog{mf: nil}
	err := mlog.InitMobileLog(logPath)
	if err != nil {
		return nil, err
	}
	return mlog, nil
}

func (mlog *ImpLog) InitMobileLog(logPath string) error {

	f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	if err != nil {
		return err
	}
	mlog.mf = f

	//	mlog.mf.WriteString("InitMobileLog open===")

	mlog.ml = log.New(mlog.mf, "ipfs-m:", log.Ltime)

	//log.SetOutput(mlog.mf)

	mlog.ml.Print("InitMobileLog open===system\n")

	return nil
}

func (mlog *ImpLog) CloseMobileLog() {
	if mlog.mf != nil {
		mlog.ml.Print("InitMobileLog close===")
		mlog.mf.Close()
	}
}

func (mlog *ImpLog) Print(v ...interface{}) {
	//now := time.Now()
	//year, _, _ := now.Date()
	//hour, min, sec := t.Clock()
	//mlog.mf.WriteString("ipfs:[" + "]" + fmt.Sprint(v...) + "\n")
	mlog.ml.Print(v)
}
