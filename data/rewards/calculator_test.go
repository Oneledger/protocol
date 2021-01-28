package rewards

/*
import "time"


import (
	"fmt"
	"github.com/Oneledger/protocol/storage"
	"github.com/stretchr/testify/assert"
	db "github.com/tendermint/tm-db"
	"testing"
	"time"
)

const (
	prefix = "rewcml"
)

var (
	calculator                   *RewardCalculator
	estimatedsecondspercycleCalc = int64(1728)
	blockspeedcalculatecycleCalc = int64(4)
	options                      *Options
	genesisTime                  time.Time
	rewardYears                  RewardYears
)

func init() {
	//Create Test DB
	newDB := db.NewDB("test", db.MemDBBackend, "")
	cs := storage.NewState(storage.NewChainState("chainstate", newDB))
	genesisTime = time.Now()
	options = &Options{}
	options.BlockSpeedCalculateCycle = blockspeedcalculatecycleCalc
	options.EstimatedSecondsPerCycle = estimatedsecondspercycleCalc

	rewardYears = RewardYears{Years: []RewardYear{}}

	calculator = NewRewardCalculator(cs, prefix)
	calculator.SetOptions(options)
	calculator.Init(genesisTime)
}

func TestRewardCalculator_secondsPerCycleLatest(t *testing.T) {
	for i := 0; i < 20; i++ {
		time := time.Now()
		calculator.Reset(int64(i+1), time, rewardYears)
		height, endtime := calculator.secondsPerCycleLatest()
		fmt.Println("height: ", height, "endtime: ", endtime)
	}

}

func TestRewardCalculator_SaveTimeStamp(t *testing.T) {
	timestamp := time.Now()
	err := calculator.SaveTimeStamp(100, timestamp)
	assert.Equal(t, nil, err)

	query, err := calculator.GetTimeStamp(100)
	assert.Equal(t, nil, err)
	assert.Equal(t, timestamp, query)
}
*/
