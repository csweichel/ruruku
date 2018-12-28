package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/32leaves/ruruku/pkg/cli"
	"github.com/32leaves/ruruku/pkg/server"
)

var cfgFile string
var verbose bool

const (
	bash_completion_func = `
__ruruku_get_output_formats()
{
    local ruruku_output out
    ruruku_output="string json jsonpath template"
    COMPREPLY=( $( compgen -W "${ruruku_output[*]}" -- "$cur" ) )
}

__ruruku_get_session()
{
    local ruruku_output out
    if ruruku_output=$(ruruku session list -o jsonpath 2>/dev/null); then
        COMPREPLY=( $( compgen -W "${ruruku_output[*]}" -- "$cur" ) )
    fi
}

__ruruku_get_user()
{
    local ruruku_output out
    if ruruku_output=$(ruruku user list -o jsonpath 2>/dev/null); then
        COMPREPLY=( $( compgen -W "${ruruku_output[*]}" -- "$cur" ) )
    fi
}

__ruruku_custom_func() {
    case ${last_command} in
        ruruku_session_close | ruruku_session_join | ruruku_session_describe)
            __ruruku_get_session
            return
            ;;
        ruruku_user_chpwd | ruruku_user_delete | ruruku_user_grant)
            __ruruku_get_user
            return
            ;;
        *)
            ;;
    esac
}
`
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ruruku",
	Short: "A simple manual test coordinator",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.InfoLevel)

		if verbose {
			log.SetLevel(log.DebugLevel)
			log.Debug("Set log level to debug")
		}
	},
	BashCompletionFunction: bash_completion_func,
}

func GetRoot() *cobra.Command {
	return rootCmd
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ruruku.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Sets the log level to debug")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".ruruku" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".ruruku")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debug("Using config file:", viper.ConfigFileUsed())
	}
}

type Config struct {
	Server server.Config
	CLI    cli.Config
}

func GetConfigFromViper() (*Config, error) {
	viper.SetDefault("server.ui.enabled", true)
	viper.SetDefault("server.ui.port", 8080)
	viper.SetDefault("server.grpc.enabled", true)
	viper.SetDefault("server.grpc.port", 1234)
	viper.SetDefault("server.db.Filename", "ruruku.db")

	viper.SetDefault("cli.host", "localhost:1234")
	viper.SetDefault("cli.timeout", 10)

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
