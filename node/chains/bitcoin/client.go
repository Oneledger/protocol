package bitcoin

import (
	btc "github.com/Oneledger/prototype-api/btcrpc"
)

var btcClient, err = btc.New("127.0.0.1",8331, "olbtc", "olbtcpw01", true)



