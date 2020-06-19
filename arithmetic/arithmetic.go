// Copyright (C) Ramsay Carslaw

package arithmetic

import "sort"

/* Letter stores a rune and meta data ascociated with it
to be used in arithmetic compression */
type Letter struct {
	// Value is the rune literal
	Value rune
	// Freq is the number of times it appears in the source
	Freq int
	// Prob = Freq / len
	Prob float64
	// Range High of probabilities
	RangeH float64
	// Range Low of probablilities
	RangeL float64
}

// Letters is a slice of Letter used for sorting
type Letters []Letter

func (l Letters) Len() int           { return len(l) }
func (l Letters) Less(i, j int) bool { return l[i].Value < l[j].Value }
func (l Letters) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

/* GetChars loads the source string into a hash map with the
character abd it's frequency in the source string */
func GetChars(src string) map[rune]int {
	m := make(map[rune]int)
	for _, v := range src {
		if m[v] == 0 {
			m[v] = 1
		} else {
			m[v]++
		}
	}
	return m
}

/* GetLetters returns a slice of Letter (type Letters),
with frequency initialised and the probability
initialised and the range */
func GetLetters(src string) Letters {
	// get a map of [rune]int for freq
	m := GetChars(src)
	var l Letters
	for key, value := range m {
		l = append(l, Letter{
			Value: key,
			Freq:  value,
		})
	}

	// init probabilities
	for i := 0; i < len(l); i++ {
		l[i].Prob = float64(l[i].Freq) / float64(len(src))
	}

	l = GetRange(l)

	return l
}

/* GetRange performs a sort on the slice of Letters and then
returns the slice of letters with thier rangeH and rangeL filled in */
func GetRange(l Letters) Letters {
	sort.Sort(l)

	// initialise all ranges
	for i := 0; i < len(l); i++ {
		if i == 0 {
			l[i].RangeH = l[i].Prob
			l[i].RangeL = 0
			continue
		}
		l[i].RangeL = l[i-1].RangeH
		l[i].RangeH = l[i-1].RangeH + l[i].Prob
	}

	return l
}

/* EncodeChar takes a rune and a list of letters and returns the range of that
rune once encoded */
func EncodeChar(ch rune, l Letters, high, low float64) (float64, float64) {
	var code Letter
	for _, v := range l {
		if v.Value == ch {
			code = v
			break
		}
	}

	codeRange := code.RangeH - code.RangeL
	high = low + (codeRange * code.RangeH)
	low = low + (codeRange * code.RangeL)

	return high, low
}

/* Encode takes a source string and returns the float64 representation
of it after it has been arithmetically encoded */
func Encode(src string) float64 {
	l := GetLetters(src)
	var high float64 = 1
	var low float64 = 0
	for _, v := range src {
		high, low = EncodeChar(v, l, high, low)
	}
	return low
}
