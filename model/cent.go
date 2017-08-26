// copied from libra

package model

import (
	"fmt"
	"math"
	"strconv"
)

type Cent int64

func (c Cent) ParseString(s string) (Cent, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	c = Cent(i)
	return c, nil
}

func (c Cent) ParseFloat(f float64) Cent {
	c = Cent(math.Ceil(f * float64(100)))
	return c
}

func (c Cent) Int64() int64 {
	return int64(c)
}

func (c Cent) Currency() (s string) {
	if c == 0 {
		return "0"
	}

	f := float64(c)
	defer func() {
		if f < 0 {
			s = fmt.Sprintf("-%s", s)
		}
	}()

	integers, digits := math.Modf(math.Abs(f) / 100.0)

	str := strconv.FormatFloat(integers, 'f', 0, 64)
	bytes := []byte(str)
	length := len(bytes)
	arr := make([]byte, 0)
	for idx, b := range bytes {
		arr = append(arr, b)
		if (length-idx) > 3 && (length-idx-1)%3 == 0 {
			arr = append(arr, 44)
		}
	}

	if digits == 0.0 {
		s = string(arr)
		return
	}

	brr := []byte(strconv.FormatFloat(digits, 'g', 2, 64))
	s = fmt.Sprintf("%s%s", string(arr), string(brr[1:]))
	return
}

func (c Cent) CurrencyWithoutComma() (s string) {
	if c == 0 {
		return "0"
	}
	return fmt.Sprintf("%.2f", c.ToFloat())
}

func (c Cent) ToFloat() float64 {
	return float64(c) / 100.0
}

func (c Cent) Inverse() Cent {
	if c == 0 {
		return c
	}
	return -1 * c
}

func (c Cent) RoundToBasicUnit() Cent {
	i := c.Int64()
	i = (i + 50) / 100 * 100
	return Cent(i)
}

// Abs ...
func (c Cent) Abs() Cent {
	if c >= 0 {
		return c
	}
	return c.Inverse()
}
