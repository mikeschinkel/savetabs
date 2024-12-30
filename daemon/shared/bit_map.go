package shared

type BitMap int

// MakeBitMap takes up to 32 bool params and returns a bit-mapped int
func MakeBitMap(bb ...bool) (result BitMap) {
	for i, b := range bb {
		if b {
			result |= 1 << i
		}
	}
	return result
}

func IsBitSet(bm BitMap, pos int) bool {
	return (bm & (1 << pos)) != 0
}
