package yangkai

import (
	"math/rand"
	"time"
)

func RandInt(length int) int {
	if length > 9 || length < 1 {
		return 0
	}
	var number = 1
	var numberVerify int
	for i:=0; i<length ; i++ {
		number = number*10
	}
	number = number-1
	numberVerify = number/10
	rand.Seed(time.Now().UnixNano())
	randInt := rand.Intn(number)

	if randInt == 0 {
		randInt = 1
	}
	for i:=0; i<length ; i++ {
		if randInt > numberVerify {
			break
		}else{
			randInt = randInt*10
		}
	}
	return randInt
}

//验证0-9的字符集
func ASCIINumber(ByteDec byte) bool {
	if ByteDec >= 48 && ByteDec <= 57 {
		return true
	} else {
		return false
	}
}

//验证a-z A-Z的字符集
func ASCIILetter(ByteDec byte) bool {
	if (ByteDec >= 65 && ByteDec <= 90) || ( ByteDec >= 97 && ByteDec <= 122) {
		return true
	} else {
		return false
	}
}
