package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/32leaves/ruruku/pkg/prettyprint"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type remoteCmdFlags struct {
	format          string
	template        string
	templateChanged bool

	outputFlagset *flag.FlagSet
}

var remoteCmdValues remoteCmdFlags

func remoteCmdValuesPreRun(cmd *cobra.Command, args []string) error {
	// HACK: call the root persistent prerun
	for c := cmd; c != nil; c = c.Parent() {
		if !c.HasParent() && c.PersistentPreRun != nil {
			c.PersistentPreRun(cmd, args)
		}
	}

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

	remoteCmdValues.outputFlagset = flags
	viper.BindPFlag("cli.host", flags.Lookup("host"))
	viper.BindPFlag("cli.token", flags.Lookup("token"))
	viper.BindPFlag("cli.timeout", flags.Lookup("timeout"))
	viper.BindPFlag("cli.tlscert", flags.Lookup("tls"))
	viper.BindEnv("cli.host", "RURUKU_HOST")
	viper.BindEnv("cli.token", "RURUKU_TOKEN")
	viper.BindEnv("cli.timeout", "RURUKU_TIMEOUT")
	viper.BindEnv("cli.tlscert", "RURUKU_TLSCERT")

	remoteCmdValues.templateChanged = flags.Lookup("output-template").Changed
	return nil
}

func registerRemoteCmdValueFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("host", "localhost:1234", "Host address of the ruruku API server")
	cmd.PersistentFlags().String("token", "", "User token for authenticating with ruruku")
	cmd.PersistentFlags().Uint32("timeout", 10, "Request timeout in seconds")
	cmd.PersistentFlags().String("tls", "", "Path to the server TLS certificate")

	cmd.PersistentFlags().StringVarP(&remoteCmdValues.format, "output", "o", "", "Output format. One of: string|json|jsonpath|template")
	cmd.PersistentFlags().Lookup("output").Annotations = map[string][]string{cobra.BashCompCustom: {"__ruruku_get_output_formats"}}
	cmd.PersistentFlags().StringVar(&remoteCmdValues.template, "output-template", "", "Output format Go template or jsonpath. Use with -o template or -o jsonpath")
}

func (s *remoteCmdFlags) GetOutputFormatWithDefault(obj interface{}, defaultFormat prettyprint.Format, template string, jsonpath string) *prettyprint.Content {
	format := defaultFormat
	if remoteCmdValues.outputFlagset.Lookup("output").Changed {
		if s.format == "json" {
			format = prettyprint.JSONFormat
		} else if s.format == "string" {
			format = prettyprint.StringFormat
		} else if s.format == "template" {
			format = prettyprint.TemplateFormat
		} else if s.format == "jsonpath" {
			format = prettyprint.JSONPathFormat
			template = jsonpath
		} else {
			log.WithField("format", s.format).Warn("Unknown format, falling back to template")
		}
	}
	if remoteCmdValues.outputFlagset.Lookup("output-template").Changed {
		template = s.template
	}

	return &prettyprint.Content{
		Obj:      obj,
		Template: template,
		Format:   format,
		Writer:   os.Stdout,
	}
}

func (s *remoteCmdFlags) GetOutputFormat(obj interface{}, template string, jsonpath string) *prettyprint.Content {
	return s.GetOutputFormatWithDefault(obj, prettyprint.TemplateFormat, template, jsonpath)
}

func (cfg *Config) Connect() (*grpc.ClientConn, error) {
	var opts grpc.DialOption
	if cfg.CLI.TLSCert == "" {
		opts = grpc.WithInsecure()
	} else {
		creds, err := credentials.NewClientTLSFromFile(cfg.CLI.TLSCert, "")
		if err != nil {
			return nil, fmt.Errorf("could not load tls cert: %s", err)
		}
		opts = grpc.WithTransportCredentials(creds)
	}
	return grpc.Dial(cfg.CLI.Host, opts)
}

func (cfg *Config) GetContext(withTimeout bool) (context.Context, context.CancelFunc) {
	ctx := context.Background()
	cancelFunc := func() {}
	if withTimeout {
		ctx, cancelFunc = context.WithTimeout(ctx, time.Duration(cfg.CLI.Timeout)*time.Second)
	}
	if cfg.CLI.Token != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", cfg.CLI.Token)
	}
	return ctx, cancelFunc
}
