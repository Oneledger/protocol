package bitcoin

import (
	brpc "./rpc"
	"time"

)

//var btcClient, err = btc.New("127.0.0.1",18831, "oltest01", "olpass01", true)


func ScheduleBlockGeneration(cli brpc.Bitcoind, interval time.Duration ) chan bool {
	ticker := time.NewTicker(interval * time.Second)
	stop := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				//TODO add generate block rpc call wrapped
				//call generate block for num
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()
	return stop
}

func StopBlockGeneration(stop chan bool) {
	close(stop)
}

