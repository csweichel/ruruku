package cmd

import (
	"fmt"
	"github.com/32leaves/ruruku/pkg/prettyprint"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"
)

type remoteCmdFlags struct {
	server          string
	tlsCert         string
	timeout         uint32
	format          string
	outputChanged   bool
	template        string
	templateChanged bool
}

var remoteCmdValues remoteCmdFlags

func remoteCmdValuesPreRun(cmd *cobra.Command, args []string) error {
	var flags *flag.FlagSet
	for c := cmd; c.HasParent(); c = c.Parent() {
		if c.PersistentFlags().Lookup("output") != nil {
			flags = c.PersistentFlags()
			break
		}
	}

	if flags == nil {
		return fmt.Errorf("Did not find remote command flags. Did you call registerRemoteCmdValueFlags?")
	}

	remoteCmdValues.outputChanged = flags.Lookup("output").Changed
	remoteCmdValues.templateChanged = flags.Lookup("output-template").Changed
	return nil
}

func registerRemoteCmdValueFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&remoteCmdValues.server, "server", "", "localhost:1234", "Host address of the ruruku API server")
	cmd.PersistentFlags().Uint32VarP(&remoteCmdValues.timeout, "timeout", "", 10, "Request timeout in seconds")
	cmd.PersistentFlags().StringVarP(&remoteCmdValues.format, "output", "o", "", "Output format. One of: string|json|template")
	cmd.PersistentFlags().StringVar(&remoteCmdValues.template, "output-template", "", "Output format Go template. Use with -o template")
	cmd.PersistentFlags().StringVar(&remoteCmdValues.tlsCert, "tls", "", "Path to the server TLS certificate")
}

func (s *remoteCmdFlags) GetOutputFormat(obj interface{}, template string) *prettyprint.Content {
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
		Obj:      obj,
		Template: template,
		Format:   format,
		Writer:   os.Stdout,
	}
}

func (s *remoteCmdFlags) Connect() (*grpc.ClientConn, error) {
	var opts grpc.DialOption
	if s.tlsCert == "" {
		opts = grpc.WithInsecure()
	} else {
		creds, err := credentials.NewClientTLSFromFile(s.tlsCert, "")
		if err != nil {
			return nil, fmt.Errorf("could not load tls cert: %s", err)
		}
		opts = grpc.WithTransportCredentials(creds)
	}
	return grpc.Dial(remoteCmdValues.server, opts)
}
