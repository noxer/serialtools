// Package serialtools offers functionality to make working with serial
// interfaces easier. It provides a LFNormalizer type which converts the
// different line endings into \n's.
package serialtools

// import the package "io" for the "Reader" interface
import "io"

// LFNormalizer converts instances of \n\r, \r\n and \r into \n in the data
// read through it.
type LFNormalizer struct {
	// the underlying reader
	r io.Reader

	// the type of the last line feed, 0 indicates that the last byte was not
	// a line feed or carriage return.
	c byte
}

// NewLFNormalizer creates a new instance of the line feed normalizer.
func NewLFNormalizer(r io.Reader) *LFNormalizer {
	return &LFNormalizer{r: r}
}

func (l *LFNormalizer) Read(p []byte) (int, error) {

	n, err := l.r.Read(p)
	if n > 0 {
		// we've received bytes from l.r, normalize them
		n = l.normalize(p[:n])
	}

	if n == 0 && len(p) != 0 && err == nil {
		// we've eaten all the bytes, we should try to get more...
		//
		// from the io.Reader docs:
		// Implementations of Read are discouraged from returning a zero byte
		// count with a nil error, except when len(p) == 0.
		n, err = l.r.Read(p)
		if n > 0 {
			// we've received bytes from l.r, normalize them
			n = l.normalize(p[:n])
		}
	}

	return n, err

}

// normalize removes all instances of "\r", "\r\n" and "\n\r" by "\n" and
// returns the length of the final buffer.
func (l *LFNormalizer) normalize(p []byte) int {

	// count the number of bytes we are skipping
	skipped := 0

	// we may need to re-read a position in the buffer, thus we can't use range.
	for i := 0; i < len(p); i++ {
		// the current byte in the buffer
		b := p[i]

		if b == '\r' || b == '\n' {
			// the character is a carriage return or line feed

			if l.c != 0 && l.c != b {
				// the last character was a carriage return or line feed
				// and not the same as the current one, remove this one.
				copy(p[i:], p[i+1:])
				l.c = 0
				i--
				skipped++
			} else {
				// the last character was the same, this seems to be a new
				// line. Save the character and write a \n into the buffer.
				p[i] = '\n'
				l.c = b
			}

		} else {
			// reset the saved character (faster to just set it to 0 than to
			// check if it was not 0 and then set it)
			l.c = 0
		}
	}

	// calculate the new size of the buffer
	return len(p) - skipped

}
