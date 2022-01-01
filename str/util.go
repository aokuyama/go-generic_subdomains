package str

func Or(strings ...string) string {
	for _, s := range strings {
		if len(s) > 0 {
			return s
		}
	}
	return ""
}
