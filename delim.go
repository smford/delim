package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"golang.org/x/term"
)

const applicationName string = "delim"
const applicationVersion string = "v0.1"
const applicationUrl string = "https://github.com/smford/delim"

var dirname2 string

func init() {

	dirname, err1 := os.UserHomeDir()
	if err1 != nil {
		log.Fatal(err1)
	}
	fmt.Println(dirname)
	dirname2 = dirname

	flag.String("char", "=", "Default line character")
	flag.String("config", dirname+"/.delim", "Configuration file: /path/to/file.yaml, default = "+dirname+"/.delim")
	flag.Bool("displayconfig", false, "Display configuration")
	flag.Bool("help", false, "Display help")
	flag.Bool("version", false, "Display version")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	checkErr(err)

	if viper.GetBool("help") {
		displayHelp()
		os.Exit(0)
	}

	if viper.GetBool("version") {
		fmt.Println(applicationName + " " + applicationVersion)
		os.Exit(0)
	}

	configdir, configfile := filepath.Split(viper.GetString("config"))

	// set default configuration directory to current directory
	if configdir == "" {
		configdir = "."
	}

	viper.SetConfigType("yaml")
	viper.AddConfigPath(configdir)

	config := strings.TrimSuffix(configfile, ".yaml")
	config = strings.TrimSuffix(config, ".yml")

	viper.SetConfigName(config)

	err = viper.ReadInConfig()
	checkErr(err)

	if viper.GetBool("displayconfig") {
		displayConfig()
		os.Exit(0)
	}
}

func main() {

	if !term.IsTerminal(0) {
		fmt.Println("Error: not a terminal")
		os.Exit(1)
	}

	width, _, err := term.GetSize(0)

	if err != nil {
		fmt.Println("Error: calculating terminal width")
		os.Exit(2)
	}

	for n := 0; n < width; n++ {
		fmt.Printf(viper.GetString("char"))
	}
}

// checks errors
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// displays help information
func displayHelp() {
	message := `
      --char [x]            Default line character (default: = ) 
      --config [file]       Configuration file: /path/to/file.yaml (default: "` + dirname2 + `/.delim")
      --help                Display help
      --version             Display version`
	fmt.Println(applicationName + " " + applicationVersion + "\n" + applicationUrl)
	fmt.Println(message)
}

// display configuration
func displayConfig() {
	allmysettings := viper.AllSettings()
	var keys []string
	for k := range allmysettings {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Println("CONFIG:", k, ":", allmysettings[k])
	}
}
