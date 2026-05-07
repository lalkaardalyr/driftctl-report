package remediation

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// WriteJSON serialises the Plan to the given writer as JSON.
func WriteJSON(w io.Writer, plan Plan) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(plan)
}

// WriteText writes a human-readable remediation plan to the given writer.
func WriteText(w io.Writer, plan Plan) error {
	if len(plan.Actions) == 0 {
		_, err := fmt.Fprintln(w, "No remediation actions required.")
		return err
	}

	_, err := fmt.Fprintf(w, "Remediation Plan (%d action(s))\n", len(plan.Actions))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, strings.Repeat("-", 60))
	if err != nil {
		return err
	}

	for i, a := range plan.Actions {
		_, err = fmt.Fprintf(w, "[%d] [%s] %s\n", i+1, strings.ToUpper(string(a.Severity)), a.Description)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "    $ %s\n\n", a.Command)
		if err != nil {
			return err
		}
	}
	return nil
}
