package committor

import (
	"encoding/base64"
	"encoding/json"
)

func (c Committor) Commit(returnValue string, transaction string) (string, error) {
	preTransactionMap := map[string]string{"returnValue": returnValue, "transaction": transaction}
	blob, error := json.Marshal(preTransactionMap)
	if error == nil {
		transactionEncoded := base64.StdEncoding.EncodeToString(blob)
		return transactionEncoded, nil
	} else {
		return "", error
	}
}
