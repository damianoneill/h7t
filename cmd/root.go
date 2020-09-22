package cmd

import (
	"crypto/tls"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/damianoneill/h7t/dsl"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type buildInfo struct {
	version string
	commit  string
	date    string
}

var bi buildInfo

var ci dsl.ConnectionInfo

var cfgFile string

var filePathSeperator = string(filepath.Separator)

// AppFs - defined for testing
var AppFs = afero.NewOsFs()

// TRUE used for boolean comparisons
const TRUE = "true"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "h7t",
	Short: "Healthbot Command Line Interface",
	Long: `A tool for interacting with Healthbot over the REST API.
	
The intent with this tool is to provide bulk or aggregate functions, that simplify interacting with Healthbot.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("verbose").Value.String() == TRUE {
			resty.SetDebug(true) // will show rest calls
		}
		ci = dsl.ConnectionInfo{
			Authority: viper.GetString("authority"),
			Username:  viper.GetString("username"),
			Password:  viper.GetString("password"),
		}
		// setup resty
		resty.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // nolint:gosec
		viper.Set("restclient.RedirectPolicy", "always")
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version, commit, date string) {
	bi = buildInfo{version, commit, date}
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.h7t.yaml)")

	rootCmd.PersistentFlags().StringP("authority", "a", "localhost:8080", "healthbot HTTPS Authority")
	rootCmd.PersistentFlags().StringP("username", "u", "admin", "healthbot Username")
	rootCmd.PersistentFlags().StringP("password", "p", "****", "healthbot Password")
	_ = viper.BindPFlag("authority", rootCmd.PersistentFlags().Lookup("authority"))
	_ = viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	_ = viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "cause "+rootCmd.Use+" to be more verbose")
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

		// Search config in home directory with name ".h7t" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".h7t")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
