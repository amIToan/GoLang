package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)
const alphabet = "asdfghjklmnopqrstuvwxyz"
func init()  {
	rand.Seed(time.Now().UnixNano())
}
func RandomInt(min,max int64) int64 {
	return min + rand.Int63n(max-min + 1)
}
func RandomStr(n int) string {
	var sb strings.Builder
	k :=len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		fmt.Print(c)
		sb.WriteByte(c)
	}
	return sb.String();
}
func RandomOwner() string	 {
	return RandomStr(6)
}
func RandomMoney() int64 {
	return RandomInt(0,100000)
}
func RandomCurrency() string {
	currencies :=[]string{"USDT","EUR","CAD"}
	n :=len(currencies);
	return currencies[rand.Intn(n)]
}