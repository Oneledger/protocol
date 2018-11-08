package monitor

import (
	"os"
	"time"
)

func (monitor Monitor) CheckStatus(status_ch chan Status) {
	i := 0
	for {
		time.Sleep(time.Second)
		i = i + 1
		if i >= monitor.TickerThreshold {
			status_ch <- Status{"Reach the ticker threshold, might have a dead loop", STATUS_DEADLOOP}
			return
		}
	}
}

func (monitor Monitor) CheckUnique() (Status, bool) {
	if _, err := os.Stat(monitor.PidFilePath); !os.IsNotExist(err) {
		switch monitor.RunningMode {
		case AGGRESIVE_MODE:
			os.Remove(monitor.PidFilePath)
			return Status{"ovm.pid file exists, there is another ovm running or exit abnormally, but we can still run a new thread", STATUS_WARNING}, false
		case CONSERVATIVE_MODE:
			return Status{"ovm.pid file exists, there is another ovm running or exit abnormally", STATUS_ALREADY_RUNNING}, true
		default:
			return Status{"ovm.pid file exists, there is another ovm running or exit abnormally, but we can still run a new thread", STATUS_WARNING}, false
		}
	}
	return Status{"OK", STATUS_OK}, false
}

func (monitor Monitor) GetPidFilePath() string {
	return monitor.PidFilePath
}

func CreateMonitor(tickerThreshold int, runningMode RunningMode, pidPath string) Monitor {
	return Monitor{tickerThreshold, runningMode, pidPath}
}
