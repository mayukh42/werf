package lib

import (
	"strconv"
	"strings"
	"time"
)

/** units
 * y, m, d, H, M, S
 * but day is floor unit
 */
func getTimeUnit(unit string) time.Duration {
	// default storage duration is 10 days
	d := time.Duration(24*10) * time.Hour
	if len(unit) != 1 {
		return d
	}

	u := strings.ToLower(unit)

	switch u {
	case "y":
		d = time.Duration(24*365) * time.Hour
	case "m":
		d = time.Duration(24*30) * time.Hour
	case "d":
		d = time.Duration(24) * time.Hour
	default:
		// nothing
	}

	return d

}

/** delta = 7d, 1y, etc
 */
func DateAfter(delta string) time.Time {
	// unit, suffix
	u := delta[len(delta)-1:]
	d := getTimeUnit(u)

	// number, prefix
	n, e := strconv.Atoi(delta[:len(delta)-1])
	if e == nil {
		d *= time.Duration(n)
	}

	return time.Now().Add(d)
}

/** AtoI() w/o the annoying err in stdlib
 */
func AtoI(a string) int {
	i, e := strconv.Atoi(a)
	if e != nil {
		i = 0
	}
	return i
}
