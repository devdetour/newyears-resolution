package util

func Is2XXStatus(code int) bool {
	return 200 <= code && code <= 300
}
