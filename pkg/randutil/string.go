package randutil

import "unicode"

func countR16(rt *unicode.RangeTable) int {
	count := 0
	for _, r := range rt.R16 {
		count += int((r.Hi-r.Lo)/r.Stride) + 1
	}
	return count
}

func countR32(rt *unicode.RangeTable) int {
	count := 0
	for _, r := range rt.R32 {
		count += int((r.Hi-r.Lo)/r.Stride) + 1
	}
	return count
}

func selectR16(rt *unicode.RangeTable, i int) rune {
	count := 0
	for _, r := range rt.R16 {

		ri := i - count
		count += int((r.Hi-r.Lo)/r.Stride) + 1

		if i < count {
			return rune(r.Lo + uint16(ri)*r.Stride)
		}
	}

	return -1
}

func selectR32(rt *unicode.RangeTable, i int) rune {
	count := 0
	for _, r := range rt.R32 {
		ri := i - count
		count += int((r.Hi-r.Lo)/r.Stride) + 1

		if i < count {
			return rune(r.Lo + uint32(ri)*r.Stride)
		}
	}

	return -1
}

func (r *Rand) Unicode(rt *unicode.RangeTable) rune {
	c16 := countR16(rt)
	c32 := countR32(rt)
	i := r.rand.Intn(c16 + c32)
	if i < c16 {
		return selectR16(rt, i)
	}

	return selectR32(rt, i-c16)
}

func RandomUnicode(*unicode.RangeTable) {

}
