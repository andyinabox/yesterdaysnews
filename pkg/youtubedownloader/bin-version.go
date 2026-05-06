package youtubedownloader

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"code.andydayton.com/andy/yesterdaysnews/pkg/util"
)

type binVersion time.Time

func ParseVersion(s string) (v binVersion, err error) {

	parts := strings.Split(s, ".")

	if len(parts) != 3 {
		err = fmt.Errorf("expected version %q to have 3 parts, found %d", s, len(parts))
		return
	}

	y, err := strconv.Atoi(parts[0])
	if err != nil {
		err = fmt.Errorf("problem parsing version year: %w", err)
		return
	}

	m, err := strconv.Atoi(parts[1])
	if err != nil {
		err = fmt.Errorf("problem parsing version yearmonth: %w", err)
		return
	}

	d, err := strconv.Atoi(parts[2])
	if err != nil {
		err = fmt.Errorf("problem parsing version yearday: %w", err)
		return
	}

	return binVersion(time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)), nil
}

func NewVersion(y, m, d int) binVersion {
	return binVersion(time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC))
}

func (v binVersion) String() string {
	y, m, d := time.Time(v).Date()
	return fmt.Sprintf("%04d.%02d.%02d", y, m, d)
}

func (v binVersion) IsEqualOrGreaterThan(v2 binVersion) bool {
	t1 := time.Time(v)
	t2 := time.Time(v2)
	r := util.CompareDates(t1, t2)
	return r == 0 || r == 1
}
