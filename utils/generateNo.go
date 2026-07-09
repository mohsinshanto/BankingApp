package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)
func GenerateAccountNo()(string,error){
	// 8 digit i.e 10000000-99999999
	max := big.NewInt(90000000)
	n, err := rand.Int(rand.Reader,max)
	if err != nil{
		return "",err
	}
	number := n.Int64()+ 10000000
	accountNo:=fmt.Sprintf("ACC%d",number)
	return accountNo,nil

}