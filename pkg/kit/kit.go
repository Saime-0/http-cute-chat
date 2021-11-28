package kit

func IntPtr(i int) *int { return &i }

func LeastOne(args ...bool) (discover bool) {
	for _, arg := range args {
		if arg {
			return true
		}
	}
	return
}
