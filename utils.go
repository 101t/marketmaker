package main

import (
	"github.com/shopspring/decimal"
	"strconv"
	"unicode"
)

func StartPosOfTime(unixTime int64, granularity int64) int64 {
	return unixTime / (granularity * 60) * (granularity * 60)
}

func StringToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func AToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func F64ToA(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func I64ToA(i int64) string {
	return strconv.FormatInt(i, 10)
}

func Int32ToString(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

func IntToA(i int) string {
	return strconv.Itoa(i)
}

func DToF64(d decimal.Decimal) float64 {
	f, _ := d.Float64()
	return f
}

func SnakeCase(s string) string {
	in := []rune(s)
	isLower := func(idx int) bool {
		return idx >= 0 && idx < len(in) && unicode.IsLower(in[idx])
	}

	out := make([]rune, 0, len(in)+len(in)/2)
	for i, r := range in {
		if unicode.IsUpper(r) {
			r = unicode.ToLower(r)
			if i > 0 && in[i-1] != '_' && (isLower(i-1) || isLower(i+1)) {
				out = append(out, '_')
			}
		}
		out = append(out, r)
	}

	return string(out)
}
