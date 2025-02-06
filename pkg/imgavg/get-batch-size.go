package imgavg

import (
	"math"
)

// getBatchSize will break the given number into
// equal sized batches
func getBatchSize(len int) int {

	if len < MaxBatchSize {
		return len
	}

	var batchSize int
	for i := 2; i < len; i++ {
		try := int(math.Ceil(float64(len) / float64(i)))
		if try <= MaxBatchSize {
			batchSize = try
			break
		}
	}

	return batchSize
}
