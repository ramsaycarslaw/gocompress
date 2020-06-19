package arithmetic

import "math/big"

func CumulativeFrequenc(freq map[byte]int64) map[byte]int64 {
	total := int64(0)
	cf := make(map[byte]int64)
	for i := 0; i < 256; i++ {
		b := byte(i)
		if v, ok := freq[b]; ok {
			cf[b] = total
			total += v
		}
	}
	return cf
}

func GetFrequency(chars []byte) (map[byte]int64, map[byte]int64) {
	freq := make(map[byte]int64)

	// count chars
	for _, ch := range chars {
		freq[ch]++
	}

	cf := CumulativeFrequenc(freq)

	return cf, freq
}

func ArithmeticCoding(src string, radix int64) (*big.Int, *big.Int, map[byte]int64) {
	chars := []byte(src)

	cf, freq := GetFrequency(chars)

	length := len(chars)

	lower := big.NewInt(0)

	product := big.NewInt(1)

	bigLen := big.NewInt(int64(length))

	for _, ch := range chars {
		count := big.NewInt(cf[ch])

		lower.Mul(lower, bigLen)
		lower.Add(lower, count.Mul(count, product))
		product.Mul(product, big.NewInt(freq[ch]))
	}

	upper := big.NewInt(1)
	upper.Set(lower)
	upper.Add(upper, product)

	One := big.NewInt(1)
	Zero := big.NewInt(0)
	Radix := big.NewInt(radix)

	temp := big.NewInt(0).Set(product)
	power := big.NewInt(0)

	for {
		temp.Div(temp, Radix)
		if temp.Cmp(Zero) == 0 {
			break
		}
		power.Add(power, One)
	}

	// diff is the encoded number
	diff := big.NewInt(0)
	diff.Sub(upper, One)
	diff.Div(diff, big.NewInt(0).Exp(Radix, power, nil))

	return diff, power, freq
}
