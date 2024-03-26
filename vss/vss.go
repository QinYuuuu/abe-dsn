package vss

import "math/big"

type Share struct {
	index int
	fi    *big.Int
}

type Param struct{
	
}

func Share(param Param, s *big.Int, t int, n int) []Share {
	f := make([]*big.Int, t)
	for i:=0; i<t+1 ; i++{
		f[i] = 
	}
}
