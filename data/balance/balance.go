/*

   ____             _          _
  / __ \           | |        | |
 | |  | |_ __   ___| | ___  __| | __ _  ___ _ __
 | |  | | '_ \ / _ \ |/ _ \/ _` |/ _` |/ _ \ '__|
 | |__| | | | |  __/ |  __/ (_| | (_| |  __/ |
  \____/|_| |_|\___|_|\___|\__,_|\__, |\___|_|
                                  __/ |
                                 |___/

	Copyright 2017 - 2019 OneLedger

*/

package balance

// Wrap the amount with owner information
type Balance struct {
	Amounts map[string]Coin `json:"amounts"`
}

/*
	Generators
*/
func NewBalance() *Balance {
	amounts := make(map[string]Coin, 0)
	result := &Balance{
		Amounts: amounts,
	}
	return result
}

/*
	methods
*/
func (b *Balance) GetCoin(currency Currency) Coin {
	if coin, ok := b.Amounts[currency.Name]; ok {
		return coin
	}

	return currency.NewCoinFromInt(0)
}

func (b *Balance) setCoin(coin Coin) *Balance {
	b.Amounts[coin.Currency.Name] = coin
	return b
}

// String method used in fmt and Dump
func (b Balance) String() string {
	buffer := ""
	for _, coin := range b.Amounts {
		if buffer != "" {
			buffer += ", "
		}
		buffer += coin.String()
	}
	return buffer
}
