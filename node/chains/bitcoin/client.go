package bitcoin

import (
	brpc "./rpc"
	"time"
	"github.com/Oneledger/protocol/node/log"

)


func GetBtcClient(rpcPort int) *brpc.Bitcoind {

	usr, pass := getCredential()
	cli, err :=  brpc.New("127.0.0.1", rpcPort , usr, pass, false)

	if err != nil{
		log.Error(err.Error())
		return nil
	}

	return cli
}

func getCredential() (usr string, pass string){
	//todo: getCredential from database which should be randomly generated when register or import if user already has bitcoin node
	usr = "oltest01"
	pass = "olpass01"
	return
}


func ScheduleBlockGeneration(cli brpc.Bitcoind, interval time.Duration ) chan bool {
	ticker := time.NewTicker(interval * time.Second)
	stop := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				cli.Generate(1)
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

