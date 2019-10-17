/*

 */

package identity

import (
	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/identity/internal"
)

func ProcessAllJobs(key keys.PrivateKey, jobs []Job) {

	internal.NewService()

	for i := range jobs {

		if !jobs[i].IsMyJobDone(key) {

			jobs[i].DoMyJob()
		}

		if jobs[i].IsSufficient() {

			jobs[i].DoFinalize()
		}
	}
}
