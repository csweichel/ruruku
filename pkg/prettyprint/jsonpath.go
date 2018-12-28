package prettyprint

import (
	"os"

	"k8s.io/client-go/util/jsonpath"
)

const (
	JSONPathFormat Format = "jsonpath"
)

func formatJSONPath(pp *Content) error {
	p := jsonpath.New("expr")
	if err := p.Parse(pp.Template); err != nil {
		return err
	}
	return p.Execute(os.Stdout, pp.Obj)
}
