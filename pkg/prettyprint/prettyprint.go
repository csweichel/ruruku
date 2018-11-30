package prettyprint

import (
	"fmt"
	"io"
)

type Format string

const StringFormat Format = "string"

type formatterFunc func(*Content) error

var formatter = map[Format]formatterFunc{
	StringFormat:   formatString,
	TemplateFormat: formatTemplate,
	JSONFormat:     formatJSON,
}

func formatString(pp *Content) error {
	_, err := fmt.Fprintf(pp.Writer, "%s", pp.Obj)
	return err
}

type Content struct {
	Obj      interface{}
	Format   Format
	Writer   io.Writer
	Template string
}

func (pp *Content) Print() error {
	formatter, ok := formatter[pp.Format]
	if !ok {
		return fmt.Errorf("Unknown format: %s", pp.Format)
	}

	return formatter(pp)
}
