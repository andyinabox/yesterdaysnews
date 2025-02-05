package imgavg

type pixelSum [][3]uint16

func (p pixelSum) Clear() pixelSum {
	for i := 0; i < len(p); i++ {
		p[i][0] = 127
		p[i][1] = 127
		p[i][2] = 127
	}
	return p
}
