package serialtools

import (
	"bytes"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
)

func TestLFNormalizer_MultiRead(t *testing.T) {

	l := NewLFNormalizer(strings.NewReader("this\n\ris\r\na\ntest\r!\r\n"))

	buf := make([]byte, 5)
	n, err := l.Read(buf)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if n != 5 {
		t.Errorf("unexpected result length: %d != 5", n)
	}
	if string(buf[:n]) != "this\n" {
		t.Errorf("unexpected result: %s != %s", strconv.Quote(string(buf[:n])), strconv.Quote("this\n"))
	}

	buf = make([]byte, 1)
	n, err = l.Read(buf)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if n != 1 {
		t.Errorf("unexpected result length: %d != 1", n)
	}
	if string(buf[:n]) != "i" {
		t.Errorf("unexpected result: %s != %s", strconv.Quote(string(buf[:n])), strconv.Quote("i"))
	}

	buf = make([]byte, 3)
	n, err = l.Read(buf)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if n != 2 {
		t.Errorf("unexpected result length: %d != 2", n)
	}
	if string(buf[:n]) != "s\n" {
		t.Errorf("unexpected result: %s != %s", strconv.Quote(string(buf[:n])), strconv.Quote("s\n"))
	}

	buf = make([]byte, 9)
	n, err = l.Read(buf)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if n != 9 {
		t.Errorf("unexpected result length: %d != 9", n)
	}
	if string(buf[:n]) != "a\ntest\n!\n" {
		t.Errorf("unexpected result: %s != %s", strconv.Quote(string(buf[:n])), strconv.Quote("a\ntest\n!\n"))
	}

	buf = make([]byte, 10)
	n, err = l.Read(buf)
	if err == nil {
		t.Error("expected error")
	}
	if n != 0 {
		t.Errorf("unexpected result length: %d != 0", n)
	}
	if string(buf[:n]) != "" {
		t.Errorf("unexpected result: %s != %s", strconv.Quote(string(buf[:n])), strconv.Quote(""))
	}

}

func TestLFNormalizer_Read(t *testing.T) {

	type test struct {
		name   string
		i      string
		o      string
		s      int
		c      byte
		hasErr bool
	}

	tests := []test{
		{
			name:   "normal, fitting string",
			i:      "this a str",
			o:      "this a str",
			s:      10,
			hasErr: false,
		},
		{
			name:   "small buffer",
			i:      "this a str",
			o:      "this a",
			s:      6,
			hasErr: false,
		},
		{
			name:   "big buffer",
			i:      "this a",
			o:      "this a",
			s:      10,
			hasErr: false,
		},
		{
			name:   "eof",
			i:      "",
			o:      "",
			s:      6,
			hasErr: true,
		},
		{
			name:   "single LF",
			i:      "\n",
			o:      "\n",
			s:      6,
			hasErr: false,
		},
		{
			name:   "multiple LF",
			i:      "\n\n\n",
			o:      "\n\n\n",
			s:      6,
			hasErr: false,
		},
		{
			name:   "single CR",
			i:      "\r",
			o:      "\n",
			s:      6,
			hasErr: false,
		},
		{
			name:   "multiple CR",
			i:      "\r\r\r",
			o:      "\n\n\n",
			s:      6,
			hasErr: false,
		},
		{
			name:   "multiple CRLF",
			i:      "\r\n\r\n\r\n",
			o:      "\n\n\n",
			s:      6,
			hasErr: false,
		},
		{
			name:   "multiple LFCR",
			i:      "\n\r\n\r\n\r",
			o:      "\n\n\n",
			s:      6,
			hasErr: false,
		},
		{
			name:   "previous CR, now LF, eof",
			i:      "\n",
			o:      "",
			s:      6,
			c:      '\r',
			hasErr: true,
		},
		{
			name:   "previous CR, now LF",
			i:      "\nx",
			o:      "x",
			s:      6,
			c:      '\r',
			hasErr: false,
		},
		{
			name:   "previous CR, now CR",
			i:      "\r",
			o:      "\n",
			s:      6,
			c:      '\r',
			hasErr: false,
		},
		{
			name:   "previous LF, now CR, eof",
			i:      "\r",
			o:      "",
			s:      6,
			c:      '\n',
			hasErr: true,
		},
		{
			name:   "previous LF, now CR",
			i:      "\rx",
			o:      "x",
			s:      6,
			c:      '\n',
			hasErr: false,
		},
		{
			name:   "previous LF, now LF",
			i:      "\n",
			o:      "\n",
			s:      6,
			c:      '\n',
			hasErr: false,
		},
		{
			name:   "short read",
			i:      "h\n\r\n\r",
			o:      "h\n",
			s:      3,
			hasErr: false,
		},
		{
			name:   "short read 2",
			i:      "h\r\n\r",
			o:      "h\n",
			s:      3,
			hasErr: false,
		},
		{
			name:   "0 read",
			i:      "string",
			o:      "",
			s:      0,
			hasErr: false,
		},
		{
			name:   "example #1",
			i:      "this\n\ris\ra\ntest\r\n!",
			o:      "this\nis\na\ntest\n!",
			s:      40,
			c:      '\n',
			hasErr: false,
		},
		{
			name:   "example #2",
			i:      "this\n\r\nis\r\ra\n\ntest\r\n\n\r!",
			o:      "this\n\nis\n\na\n\ntest\n\n!",
			s:      40,
			c:      '\n',
			hasErr: false,
		},
		{
			name:   "example #3",
			i:      "this\n \r \nis\r \ra\n \ntest\r \n \n \r!",
			o:      "this\n \n \nis\n \na\n \ntest\n \n \n \n!",
			s:      40,
			c:      '\n',
			hasErr: false,
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {

			r := strings.NewReader(tst.i)
			p := make([]byte, tst.s)
			l := &LFNormalizer{r: r, c: tst.c}

			n, err := l.Read(p)

			if n != len(tst.o) {
				t.Errorf("unexpected return length: %d != %d", n, len(tst.o))
			}
			if tst.hasErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tst.hasErr && err != nil {
				t.Errorf("unexpected error: %s", err)
			}
			if string(p[:n]) != tst.o {
				t.Errorf("Unexpected result: %s != %s", strconv.Quote(string(p[:n])), strconv.Quote(tst.o))
			}

		})
	}

}

func BenchmarkRead(b *testing.B) {

	src, err := ioutil.ReadFile("testdata/benchmark.txt")
	if err != nil {
		b.Fatalf("Unable to read the test data: %s", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s := bytes.NewReader(src)
		r := NewLFNormalizer(s)

		b.StartTimer()
		io.Copy(ioutil.Discard, r)
		b.StopTimer()
	}

}
