package text

import (
	"bytes"
)

// Truncate resizes the string with the given length. It ellipses with '...' when the string's length exceeds
// the desired length or pads spaces to the right of the string when length is smaller than desired
func Truncate(s string, length int) string {
	slen := StringWidth(s)
	n := length
	if slen == n {
		return s
	}
	// Pads only when length of the string smaller than len needed
	s = PadRight(s, n, ' ')

	if slen > n {
		var buf bytes.Buffer
		w := 0
		for _, r := range s {
			buf.WriteRune(r)
			rw := RuneWidth(r)
			if w+rw >= n-3 {
				break
			}
			w += rw
		}
		buf.WriteString("...")
		s = buf.String()
	}
	return s
}
