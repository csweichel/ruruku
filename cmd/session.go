package cmd

import (
	"github.com/spf13/cobra"
    "github.com/32leaves/ruruku/pkg/prettyprint"
    log "github.com/sirupsen/logrus"
    "os"
)

type sessionFlags struct {
	server  string
	timeout uint32
    format  string
    template string
}

var sessionFlagValues sessionFlags

// sessionCmd represents the session command
var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Controls a test session/run",
	Args:  cobra.ExactArgs(1),
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(sessionCmd)

	sessionCmd.PersistentFlags().StringVarP(&sessionFlagValues.server, "server", "", "localhost:1234", "Host address of the ruruku API server")
	sessionCmd.PersistentFlags().Uint32VarP(&sessionFlagValues.timeout, "timeout", "", 10, "Request timeout in seconds")
    sessionCmd.PersistentFlags().StringVarP(&sessionFlagValues.format, "output", "o", "", "Output format. One of: string|json|template")
    sessionCmd.PersistentFlags().StringVar(&sessionFlagValues.template, "output-template", "", "Output format Go template. Use with -o template")
}

func (s *sessionFlags) GetOutputFormat(obj interface{}, template string) *prettyprint.Content {
    format := prettyprint.TemplateFormat
    if sessionCmd.PersistentFlags().Lookup("output").Changed {
        if s.format == "json" {
            format = prettyprint.JSONFormat
        } else if s.format == "string" {
            format = prettyprint.StringFormat
        } else if s.format == "template" {
            format = prettyprint.TemplateFormat
        } else {
            log.WithField("format", s.format).Warn("Unknown format, falling back to template")
        }
    }
    if sessionCmd.PersistentFlags().Lookup("output-template").Changed {
        template = s.template
    }

    return &prettyprint.Content{
        Obj: obj,
        Template: template,
        Format: format,
        Writer: os.Stdout,
    }
}