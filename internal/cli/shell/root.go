package cmd

import (
	"fmt"
	"os"
	"time"

	filename "github.com/keepeye/logrus-filename"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	address   string
	port      string
	debug     bool
	logFormat string
	logLevel  string
)

// ViperFlagsRoot defines a struct to hold the values of cobra CLI flags and use viper to populate them
type ViperFlagsRoot struct {
	CfgFile   string
	Address   string `mapstructure:"address"`
	Port      bool   `mapstructure:"port"`
	Debug     bool   `mapstructure:"debug"`
	LogFormat string `mapstructure:"log-format"`
	LogLevel  string `mapstructure:"log-level"`
}

// Declare the viper CLI flag values buffer
var vprFlgsRoot ViperFlagsRoot

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cineguard",
	Short: "Cineguard service to manage your movies",
	Long:  `This is the Cineguard server to manage your movies using a database`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logrus.Warn("No subcommand provided. Please use one of the available subcommands. Showing help message.")
		return cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	filenameHook := filename.NewHook()
	filenameHook.Field = "line"
	logrus.AddHook(filenameHook)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set and populate CLI flags buffer thanks to viper
func initConfig() {
	// Parse config file with viper
	ReadViperConfigE(viper.GetViper(), rootCmd)

	// Initialize and populate cobra CLI root flags values with viper
	InitViperSubCmdE(viper.GetViper(), rootCmd, &vprFlgsRoot)

	// Set logs format
	switch vprFlgsRoot.LogFormat {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{
			PrettyPrint:      false,
			DisableTimestamp: false,
			TimestampFormat:  time.RFC3339,
		})
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:      true,
			DisableTimestamp: false,
			TimestampFormat:  time.DateTime,
		})
	default:
		logrus.WithError(fmt.Errorf("logrus unknown output format")).Error("unknown log format")
	}
	logrus.Debugf("logrus output format is set to: %s", vprFlgsRoot.LogFormat)

	// Initialize logrus log level and log format for all cobra commands and subcommands.
	debugFlagIsUsed := rootCmd.Flags().Lookup("debug").Changed

	switch {
	case debugFlagIsUsed:
		// harcode that the --debug flags set logrus level to debug
		logrus.SetLevel(logrus.DebugLevel)
	default:
		// get the log level from viper which is bind to the cobra flag --log-level
		level, err := logrus.ParseLevel(vprFlgsRoot.LogLevel)
		if err != nil {
			logrus.WithError(err).Error("unknown log level")
		}
		logrus.SetLevel(level)
	}
	logrus.Debugf("logrus log-level is set to: %s", logrus.GetLevel())

}

func init() {
	// Ensure initConfig runs before anything else
	cobra.OnInitialize(initConfig)

	// Since this project uses Viper bind with Cobra flags, we generally do not need to use "Flags().*Var"
	// (like StringVar, BoolVar, Uint16Var, etc...) as we do not need to access the cobra flag values directly. This is
	// because we use Viper to retrieve the values of the flags.
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.yaml", "Configuration file for Cineguard's server")

	// logging level
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Set logrus.SetLevel to \"debug\". This is equivalent to using --log-level=debug. Flags --log-level and --debug flag are mutually exclusive.")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Set logrus.SetLevel. Possible values: trace, debug, info, warning, error, fatal and panic. Flags --log-level and --debug flag are mutually exclusive.")
	rootCmd.RegisterFlagCompletionFunc("log-level", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"trace", "debug", "info", "warning", "error", "fatal", "panic"}, cobra.ShellCompDirectiveNoFileComp
	})
	rootCmd.PersistentFlags().StringVar(&logFormat, "log-format", "text", "Logrus log output format. Possible values: text, json.")
	rootCmd.RegisterFlagCompletionFunc("log-format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"text", "json"}, cobra.ShellCompDirectiveNoFileComp
	})
	rootCmd.MarkFlagsMutuallyExclusive("log-level", "debug") // --log-level and --debug flag are mutually exclusive since debug is an alias for log-level=debug
}
