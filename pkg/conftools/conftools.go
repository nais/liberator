package conftools

import (
	"fmt"
	"sort"
	"strings"

	"slices"

	"github.com/go-viper/mapstructure/v2"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Initialize configuration system.
// Configuration will be read from the file `/etc/<appName>.yaml` first, then overwritten by
// environment variables in the form `<APPNAME>_<OPTION>`.
func Initialize(appName string) {
	// Automatically read configuration options from environment variables.
	// i.e. --proxy.address will be configurable using PROXY_ADDRESS.
	viper.AutomaticEnv()
	viper.SetEnvPrefix(strings.ToUpper(appName))
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	// Read configuration file from working directory and/or /etc.
	// File formats supported include JSON, TOML, YAML, HCL, envfile and Java properties config files
	viper.SetConfigName(appName)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc")
}

// Read config from file, command-line flags and environment variables (in that order)
// into the structured object `cfg`.
//
// The structured object must annotate its variables with `json:"name"`.
// Nested objects are allowed, and names will be separated by underscore when reading from environment, e.g.
// APPNAME_LEVEL1_LEVEL2_VARIABLE.
func Load(cfg any) error {
	var err error

	err = viper.ReadInConfig()
	if err != nil {
		if err.(viper.ConfigFileNotFoundError) != err {
			return err
		}
	}

	flag.Parse()

	err = viper.BindPFlags(flag.CommandLine)
	if err != nil {
		return err
	}

	decoderHook := func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "json"
		dc.ErrorUnused = true
	}

	err = viper.Unmarshal(cfg, decoderHook)
	if err != nil {
		return err
	}

	return nil
}

// Return a human-readable printout of all configuration options, except secret stuff.
func Format(disallowedKeys []string) []string {
	var keys sort.StringSlice = viper.AllKeys()

	printed := make([]string, 0)

	keys.Sort()
	for _, key := range keys {
		if slices.Contains(disallowedKeys, key) {
			printed = append(printed, fmt.Sprintf("%s: ***REDACTED***", key))
		} else {
			printed = append(printed, fmt.Sprintf("%s: %v", key, viper.Get(key)))
		}
	}

	return printed
}
