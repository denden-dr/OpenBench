package utils

import (
	"strconv"
)

// IntToString converts an integer to a string.
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// FormatCurrency formats an int64 amount as IDR currency string (e.g. 150000 -> "150.000")
func FormatCurrency(amount int64) string {
	in := strconv.FormatInt(amount, 10)
	out := make([]byte, len(in)+(len(in)-1)/3)
	if in[0] == '-' {
		in, out[0] = in[1:], '-'
	}
	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = '.'
		}
	}
}
