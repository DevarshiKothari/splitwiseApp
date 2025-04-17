package utils

import (
	"math"
	"strings"
)

type Balance struct {
	FromUserID int
	ToUserID   int
	Amount     float64
}

func StringJoin(elems []string, sep string) string {
	return strings.Join(elems, sep)
}

func SimplifyBalances(original map[int]float64) []Balance {
	balances := make(map[int]float64)
	for k, v := range original {
		balances[k] = v
	}
	balanceList := []Balance{}
	for k1, _ := range balances {
		if balances[k1] < 0 {
			amountOwedByk1 := math.Abs(balances[k1])
			for k2, _ := range balances {
				if k2 == k1 { // same user
					continue
				} else if balances[k2] <= 0 { // user k2 already owes money so k1 does not owe k2 anything
					continue
				} else if balances[k2] > amountOwedByk1 { // user k2 is owed more than amountOwedByk1
					balanceList = append(balanceList, Balance{FromUserID: k1, ToUserID: k2, Amount: amountOwedByk1})
					balances[k2] -= amountOwedByk1 // amountOwedByk1 deducted from k2 owed money
					balances[k1] = 0               // now k1 owns nothing
					break
				} else { // k2 is owed less than amt
					balanceList = append(balanceList, Balance{FromUserID: k1, ToUserID: k2, Amount: balances[k2]})
					balances[k1] += balances[k2]
					balances[k2] = 0
				}
			}
		}
	}

	if len(balanceList) == 0 {
		return []Balance{}
	}

	return balanceList
}
