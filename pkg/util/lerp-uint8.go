package util

func LerpUint8(a, b, t uint8) uint8 {
	return uint8(uint16(a) + (uint16(t) * (uint16(b) - uint16(a)) / 255))
}
