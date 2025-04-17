package utils

import (
	"math"
	"splitwise-app/models"
)

func MapTransformationFunc(m map[int]float64) []models.Balance {
	count := 0
	balanceList := make([]models.Balance, len(m))
	for k1, _ := range m {
		if m[k1] < 0 {
			amt := math.Abs(m[k1])
			for k2, _ := range m {
				if k2 == k1 { // same user
					continue
				} else if m[k2] <= 0 { // this user k2 already owes money so k1 does not owe k2 anything
					continue
				} else if m[k2] > amt { // this user k2 is owed more than amt
					balanceList[count].FromUserID = k1
					balanceList[count].ToUserID = k2
					balanceList[count].Amount = amt
					m[k2] -= amt // amt deducted from k2 owed money
					m[k1] = 0    // now k1 owns nothing
					break
				} else { // k2 is owed less than amt
					balanceList[count].FromUserID = k1
					balanceList[count].ToUserID = k2
					balanceList[count].Amount = m[k2]
					m[k1] += m[k2]
					m[k2] = 0
				}
			}
		}
		count++
	}

	return balanceList
}
