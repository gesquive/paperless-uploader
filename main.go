package main

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	buildVersion = "v0.1.0-dev"
	buildCommit  = ""
	buildDate    = ""
)

var cfgFile string

var showVersion bool
var debug bool

func main() {
	Execute()
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:    "paperless-uploader",
	Short:  "A paperless-ng document uploader",
	Long:   `Watches a directory for files and uploads them to paperless-ng`,
	PreRun: preRun,
	Run:    run,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	RootCmd.SetHelpTemplate(fmt.Sprintf("%s\nVersion:\n  github.com/gesquive/paperless-uploader %s\n",
		RootCmd.HelpTemplate(), buildVersion))
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"Path to a specific config file (default \"./config.yml\")")
	RootCmd.PersistentFlags().StringP("log-file", "l", "",
		"Path to log file (default \"/var/log/paperless-uploader.log\")")

	RootCmd.PersistentFlags().StringP("watch-path", "p", "",
		"Directory to watch for files.")
	RootCmd.PersistentFlags().StringP("watch-filter", "f", "",
		"The inclusive file filter regex for uploads.")
	RootCmd.PersistentFlags().DurationP("watch-interval", "i", time.Second,
		"The interval between polling for changes.")

	RootCmd.PersistentFlags().StringSliceP("upload-path", "x", []string{},
		"Path to the file(s) to upload, can be entered multiple times or comma delimited.")

	RootCmd.PersistentFlags().StringP("paperless-url", "u", "",
		"The base URL for your paperless instance")
	RootCmd.PersistentFlags().StringP("paperless-token", "t", "",
		"Authenticate the paperless server with this user token")

	RootCmd.PersistentFlags().BoolVar(&showVersion, "version", false,
		"Display the version info and exit")

	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false,
		"Include debug statements in log output")
	RootCmd.PersistentFlags().MarkHidden("debug")

	viper.SetEnvPrefix("paperless-uploader")
	viper.AutomaticEnv()
	viper.BindEnv("config")
	viper.BindEnv("log_file")
	viper.BindEnv("watch_path")
	viper.BindEnv("watch_filter")
	viper.BindEnv("watch_interval")
	viper.BindEnv("upload_path")
	viper.BindEnv("paperless_url")
	viper.BindEnv("paperless_token")

	viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("log_file", RootCmd.PersistentFlags().Lookup("log-file"))
	viper.BindPFlag("watch.path", RootCmd.PersistentFlags().Lookup("watch-path"))
	viper.BindPFlag("watch.filter", RootCmd.PersistentFlags().Lookup("watch-filter"))
	viper.BindPFlag("watch.interval", RootCmd.PersistentFlags().Lookup("watch-interval"))
	viper.BindPFlag("upload_path", RootCmd.PersistentFlags().Lookup("upload-path"))
	viper.BindPFlag("paperless.url", RootCmd.PersistentFlags().Lookup("paperless-url"))
	viper.BindPFlag("paperless.token", RootCmd.PersistentFlags().Lookup("paperless-token"))

	viper.SetDefault("log_file", "/var/log/paperless-uploader.log")
	viper.SetDefault("watch_path", "")
	viper.SetDefault("watch_interval", time.Second)

	dotReplacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(dotReplacer)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cfgFile := viper.GetString("config")
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")                           // name of config file (without extension)
		viper.AddConfigPath(".")                                // add current directory as first search path
		viper.AddConfigPath("$HOME/.config/paperless-uploader") // add home directory to search path
		viper.AddConfigPath("/etc/paperless-uploader")          // add etc to search path
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if !showVersion {
			if !strings.Contains(err.Error(), "Not Found") {
				fmt.Printf("Error opening config: %s\n", err)
			}
		}
	}
}

func preRun(cmd *cobra.Command, args []string) {
	if showVersion {
		fmt.Printf("github.com/gesquive/paperless-uploader\n")
		fmt.Printf(" Version:    %s\n", buildVersion)
		if len(buildCommit) > 6 {
			fmt.Printf(" Git Commit: %s\n", buildCommit[:7])
		}
		if buildDate != "" {
			fmt.Printf(" Build Date: %s\n", buildDate)
		}
		fmt.Printf(" Go Version: %s\n", runtime.Version())
		fmt.Printf(" OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}
}

func run(cmd *cobra.Command, args []string) {
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if isTerminal() {
		log.SetFormatter(&log.TextFormatter{
			DisableTimestamp: true,
		})
	} else {
		log.Infof("running paperless-uploader %s", buildVersion)
		if len(buildCommit) > 6 {
			log.Infof("build: commit=%s", buildCommit[:7])
		}
		if buildDate != "" {
			log.Infof("build: date=%s", buildDate)
		}
		log.Infof("build: info=%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	}

	logFilePath := viper.GetString("log_file")
	log.Debugf("config: log_file=%s", logFilePath)
	if strings.ToLower(logFilePath) == "stdout" || logFilePath == "-" || logFilePath == "" {
		log.SetOutput(os.Stdout)
	} else {
		logFilePath = getLogFilePath(logFilePath)
		logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening log file=%v", err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	log.Debugf("config: file=%s", viper.ConfigFileUsed())
	if viper.ConfigFileUsed() == "" {
		log.Fatal("No config file found.")
	}

	paperlessUrl := viper.GetString("paperless.url")
	if _, err := url.Parse(paperlessUrl); err != nil {
		log.Fatalf("paperless base url is not valid '%s'", paperlessUrl)
	}
	log.Debugf("config: paperless_url=%v", paperlessUrl)

	paperlessToken := viper.GetString("paperless.token")
	if len(paperlessToken) <= 0 {
		log.Fatalf("paperless token is missing")
	} else if len(paperlessToken) < 40 {
		log.Fatalf("paperless token is invalid")
	}
	log.Debugf("config: paperless_token=****************************************")
	uploader := NewUploader(paperlessUrl, paperlessToken)

	uploadPaths := viper.GetStringSlice("upload_path")
	for _, uploadPath := range uploadPaths {
		log.Debugf("config: upload_path=%v", uploadPaths)
		uploader.uploadAll(uploadPath, false)
		os.Exit(0)
	}

	watchDir := viper.GetString("watch.path")
	watchFilter := viper.GetString("watch.filter")
	watchInterval := viper.GetDuration("watch.interval")
	if len(watchDir) > 0 {
		log.Debugf("config: watch_path=%s", watchDir)
		log.Debugf("config: watch_interval=%s", watchInterval)
		log.Debugf("config: watch_filter=%s", watchFilter)
		watcher := NewWatcher(uploader)
		watcher.Watch(watchDir, watchInterval, watchFilter)
	}
}

func getLogFilePath(defaultPath string) (logPath string) {
	fi, err := os.Stat(defaultPath)
	if err == nil && fi.IsDir() {
		logPath = path.Join(defaultPath, "paperless-uploader.log")
	} else {
		logPath = defaultPath
	}
	return
}

func isTerminal() bool {
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		return true
	}
	return false
}
