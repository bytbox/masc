package main

func lim(s string, i int) string {
	if len(s) > i {
		return s[:i]
	}
	return s
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
