package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

var (
	applicationName    = "delim"
	applicationVersion = "dev"
	applicationURL     = "https://github.com/smford/delim"
)

func getVersion() string {
	if applicationVersion != "" && applicationVersion != "dev" {
		return applicationVersion
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return applicationVersion
}

// buildDelimiter builds the delimiter string efficiently using strings.Repeat.
func buildDelimiter(char string, width int) string {
	if width <= 0 || len(char) == 0 {
		return ""
	}
	return strings.Repeat(char, width)
}

// getTerminalWidth attempts to detect terminal width from standard streams,
// falling back to COLUMNS env var or default 80 columns.
func getTerminalWidth(explicitWidth int) int {
	if explicitWidth > 0 {
		return explicitWidth
	}

	for _, f := range []*os.File{os.Stdout, os.Stderr, os.Stdin} {
		if f != nil && term.IsTerminal(int(f.Fd())) {
			if width, _, err := term.GetSize(int(f.Fd())); err == nil && width > 0 {
				return width
			}
		}
	}

	if colsStr := os.Getenv("COLUMNS"); colsStr != "" {
		if cols, err := strconv.Atoi(colsStr); err == nil && cols > 0 {
			return cols
		}
	}

	return 80
}

func printHelp(out io.Writer, defaultConfigFile string) {
	fmt.Fprintf(out, "%s %s\n%s\n\n", applicationName, getVersion(), applicationURL)
	fmt.Fprintf(out, "Usage: %s [options]\n\n", applicationName)
	fmt.Fprintln(out, "Options:")
	fmt.Fprintln(out, "  -c, --char string        Default line character or pattern (default \"=\")")
	fmt.Fprintf(out, "      --config string      Configuration file: /path/to/file.yaml (default %q)\n", defaultConfigFile)
	fmt.Fprintln(out, "  -w, --width int          Line width (default: terminal width)")
	fmt.Fprintln(out, "  -n, --newline            Print trailing newline (default false)")
	fmt.Fprintln(out, "      --displayconfig      Display configuration")
	fmt.Fprintln(out, "  -v, --version            Display version")
	fmt.Fprintln(out, "  -h, --help               Display help")
}

func printConfig(v *viper.Viper, out io.Writer) {
	settings := v.AllSettings()
	keys := make([]string, 0, len(settings))
	for k := range settings {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(out, "CONFIG: %s : %v\n", k, settings[k])
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	homeDir, err := os.UserHomeDir()
	defaultConfigFile := ""
	if err == nil {
		defaultConfigFile = filepath.Join(homeDir, ".delim")
	}

	flags := pflag.NewFlagSet(applicationName, pflag.ContinueOnError)
	flags.SetOutput(stderr)

	flags.StringP("char", "c", "=", "Default line character")
	flags.String("config", defaultConfigFile, "Configuration file path")
	flags.IntP("width", "w", 0, "Line width (default: terminal width)")
	flags.BoolP("newline", "n", false, "Print trailing newline")
	flags.Bool("displayconfig", false, "Display configuration")
	flags.BoolP("version", "v", false, "Display version")
	flags.BoolP("help", "h", false, "Display help")

	if err := flags.Parse(args); err != nil {
		return err
	}

	if help, _ := flags.GetBool("help"); help {
		printHelp(stdout, defaultConfigFile)
		return nil
	}

	if ver, _ := flags.GetBool("version"); ver {
		fmt.Fprintf(stdout, "%s %s\n", applicationName, getVersion())
		return nil
	}

	v := viper.New()
	v.SetEnvPrefix("DELIM")
	v.AutomaticEnv()

	v.SetDefault("char", "=")
	v.SetDefault("config", defaultConfigFile)
	v.SetDefault("width", 0)
	v.SetDefault("newline", false)

	_ = v.BindPFlag("char", flags.Lookup("char"))
	_ = v.BindPFlag("config", flags.Lookup("config"))
	_ = v.BindPFlag("width", flags.Lookup("width"))
	_ = v.BindPFlag("newline", flags.Lookup("newline"))

	configFile := v.GetString("config")
	if configFile != "" {
		if _, err := os.Stat(configFile); err == nil {
			v.SetConfigFile(configFile)
			v.SetConfigType("yaml")
			if err := v.ReadInConfig(); err != nil {
				return fmt.Errorf("reading config file %q: %w", configFile, err)
			}
		} else if flags.Changed("config") {
			return fmt.Errorf("config file not found: %s", configFile)
		}
	}

	if disp, _ := flags.GetBool("displayconfig"); disp {
		printConfig(v, stdout)
		return nil
	}

	explicitWidth := v.GetInt("width")
	if explicitWidth < 0 {
		return fmt.Errorf("invalid width %d: width must be non-negative", explicitWidth)
	}

	width := getTerminalWidth(explicitWidth)
	char := v.GetString("char")
	line := buildDelimiter(char, width)

	if _, err := io.WriteString(stdout, line); err != nil {
		return err
	}
	if v.GetBool("newline") {
		if _, err := io.WriteString(stdout, "\n"); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		if errors.Is(err, pflag.ErrHelp) {
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
