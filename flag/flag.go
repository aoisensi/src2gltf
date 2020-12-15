package flag

import "flag"

// Flags
var (
	Scale float64
)

func init() {
	flag.Float64Var(&Scale, "scale", 0.02, "Scale.")
	flag.Parse()
}

// Args return flag.Args()
func Args() []string {
	return flag.Args()
}
