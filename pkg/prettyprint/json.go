package prettyprint

import "encoding/json"

const JSONFormat Format = "json"

func formatJSON(pp *Content) error {
	out, err := json.MarshalIndent(pp.Obj, "", "  ")
	if err != nil {
		return err
	}

	off := 0
	for n, err := pp.Writer.Write(out[off:]); off < len(out); n, err = pp.Writer.Write(out[off:]) {
		if err != nil {
			return err
		}
		off += n
	}

	return nil
}
