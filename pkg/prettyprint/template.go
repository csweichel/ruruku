package prettyprint

import (
    "text/template"
    "text/tabwriter"
)

const TemplateFormat Format = "tpl"

func formatTemplate(pp *Content) error {
	tmpl, err := template.New("prettyprint").Parse(pp.Template)
	if err != nil {
		return err
	}

    w := tabwriter.NewWriter(pp.Writer, 8, 8, 8, ' ', 0)
    if err := tmpl.Execute(w, pp.Obj); err != nil {
        return err
    }
    if err := w.Flush(); err != nil {
        return err
    }
    return nil
}
