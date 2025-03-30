package seal

import (
	"testing"
	"time"
	"fmt"
)



func TestRandOP(t *testing.T){


	for i:=1; i<=256; i++ {

		st:=time.Now()
		RandOP(i)
		t:=time.Now()
		fmt.Println("TestRandOP",i,t.Sub(st))
	}

}

func TestVeryRandOP(t *testing.T){
	for i:=0; i<10; i++{
		st:=time.Now()
		VeryRandOP()
		t:=time.Now()
		fmt.Println("TestVeryRandOP",i,t.Sub(st))

	}

}
