package text

import (
	"bytes"
	"strings"
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

		hyperlinkMatches := hyperlinkOSCRegexp.FindAllStringIndex(s, -1)
		if hyperlinkMatches == nil {
			hyperlinkMatches = [][]int{{0, 0}}
		}

		firstLinkPos := hyperlinkMatches[0][0]
		if firstLinkPos > 0 {
			for _, r := range s[:firstLinkPos] {
				buf.WriteRune(r)
				w += RuneWidth(r)
				if w >= n-3 {
					break
				}
			}
		}

		if w < n-3 {
			for i, match := range hyperlinkMatches {
				startPos := match[0]
				endPos := match[1]

				isClosingSequence := s[startPos:endPos] == "\x1b]8;;\x1b\\"

				if w < n-3 || isClosingSequence {
					buf.WriteString(s[startPos:endPos]) // this doesn't count against our character total
				}

				if w >= n-3 {
					break
				}

				// determine the substring containing the next chunk of non-hyperlink OSC text
				nextStartPos := len(s)
				if i+1 < len(hyperlinkMatches) {
					nextStartPos = hyperlinkMatches[i+1][0]
				}

				for _, r := range s[endPos:nextStartPos] {
					if w >= n {
						break
					} else {
						if w >= n-3 {
							buf.WriteRune('.')
							w++
						} else {
							buf.WriteRune(r)
							w += RuneWidth(r)
						}
					}
				}
			}
		}

		if w < n {
			buf.WriteString(strings.Repeat(".", n-w))
		}

		s = buf.String()
	}
	return s
}
