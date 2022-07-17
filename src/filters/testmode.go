package filters

import (
	"bytes"
	"io"
	"regexp"
)

var (
	testStyle      = []byte(":style(border: 2px dashed red !important)\n")
	commentMatch   = regexp.MustCompile(`^[!#]`)
	forbiddenRules = regexp.MustCompile(`:(style|remove)\(`)
)

type TestModeTransformer struct {
	out io.Writer
	buf []byte
}

func NewTestModeTransformer(out io.Writer) *TestModeTransformer {
	return &TestModeTransformer{
		out: out,
	}
}

func (t *TestModeTransformer) Write(p []byte) (int, error) {
	input := p

	for {
		if len(input) == 0 {
			break
		}
		i := bytes.Index(input, newLine)
		if i == -1 {
			t.buf = append(t.buf, input...)
			break
		}
		t.buf = append(t.buf, input[:i]...)
		if err := t.writeLine(); err != nil {
			return 0, err
		}
		input = input[i+1:]
	}

	return len(p), nil
}

func (t *TestModeTransformer) writeLine() error {
	injectStyle := true

	if len(t.buf) == 0 || commentMatch.Match(t.buf) {
		injectStyle = false
	} else if pos := forbiddenRules.FindIndex(t.buf); pos != nil {
		// Trim the conflicting directives
		t.buf = t.buf[0:pos[0]]
	}

	_, err := t.out.Write(t.buf)
	if err != nil {
		return err
	}

	if injectStyle {
		_, err = t.out.Write(testStyle)
	} else {
		_, err = t.out.Write(newLine)
	}
	t.buf = nil
	return err
}
