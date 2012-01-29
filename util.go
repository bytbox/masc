package main

func lim(s string, i int) string {
	if len(s) > i {
		return s[:i]
	}
	return s
}
