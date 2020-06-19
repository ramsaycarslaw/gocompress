// Copyright (C) 2020 Ramsay Carslaw

package arithmetic

import "math/big"

func ArithmeticDecode(num *big.Int, radix int64, pow *big.Int, freq map[byte]int64) string {
	power := big.NewInt(radix)

	enc := big.NewInt(0).Set(num)
	enc.Mul(enc, power.Exp(power, pow, nil))

	base := int64(0)

	for _, v := range freq {
		base += v
	}

	cf := CumulativeFrequenc(freq)

	dict := make(map[int64]byte)

	for k, v := range cf {
		dict[v] = k
	}

	lchar := -1

	for i := int64(0); i < base; i++ {
		if v, ok := dict[i]; ok {
			lchar = int(v)
		} else if lchar != -1 {
			dict[i] = byte(lchar)
		}
	}

	decoded := make([]byte, base)
	Base := big.NewInt(base)

	for i := base - 1; i >= 0; i-- {
		pow := big.NewInt(0)
		pow.Exp(Base, big.NewInt(i), nil)

		div := big.NewInt(0)
		div.Div(enc, pow)

		c := dict[div.Int64()]
		fv := freq[c]
		cv := cf[c]

		prod := big.NewInt(0).Mul(pow, big.NewInt(cv))
		diff := big.NewInt(0).Sub(enc, prod)
		enc.Div(diff, big.NewInt(fv))

		decoded[base-i-1] = c
	}

	return string(decoded)
}
