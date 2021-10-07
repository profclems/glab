package text

import "bytes"

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
		buf := bytes.Buffer{}
		w := 0
		dotsWritten := 0

		hyperlinkMatches := hyperlinkOSCRegexp.FindAllStringIndex(s, -1)

		// append a faked match so that we always have something in
		// hyperlinkMatches[0] and can avoid nil checks
		hyperlinkMatches = append(hyperlinkMatches, []int{len(s), len(s)})

		for i, r := range s {
			startPos := hyperlinkMatches[0][0]
			endPos := hyperlinkMatches[0][1]

			if i >= startPos && i < endPos {
				// write runes inside hyperlink OSC sequences - this doesn't count
				// against our grapheme total
				buf.WriteRune(r)
			} else if w == 0 {
				// always write the first character
				buf.WriteRune(r)
				w += RuneWidth(r)
			} else {
				rw := RuneWidth(r)

				if w+rw <= n-3 {
					// if we have room before the ellipsis, go ahead and write it
					buf.WriteRune(r)
					w += rw
				} else if dotsWritten < 3 {
					buf.WriteRune('.')
					w++
					dotsWritten++
				} else {
					break
				}
			}

			if i == endPos-1 {
				// we're at the end of this hyperlink OSC sequence - get ready
				// to check for the next one
				hyperlinkMatches = hyperlinkMatches[1:]
			}
		}

		// if we currently have an "open" hyperlink OSC sequence, we need to write
		// a closing sequence
		nextOSC := s[hyperlinkMatches[0][0]:hyperlinkMatches[0][1]]
		isNextOSCClosing := nextOSC == "\x1b]8;;\x1b\\"

		if isNextOSCClosing {
			buf.WriteString(nextOSC)
		}

		s = buf.String()
	}
	return s
}
