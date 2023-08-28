package fctl

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/spf13/pflag"
)

const (
	MembershipURIFlag string = "membership-uri"
	ConfigFlag        string = "config"
	ProfileFlag       string = "profile"
	OutputFlag        string = "output"
	DebugFlag         string = "debug"
	InsecureTlsFlag   string = "insecure-tls"
	TelemetryFlag     string = "telemetry"
	MetadataFlag      string = "metadata"
)

var (
	insecureTlsV   = &fValue[bool]{value: false}
	telemetryFlagV = &fValue[bool]{value: false}
	debugFlagV     = &fValue[bool]{value: false}
	profileFlagV   = &fValue[string]{value: ""}
	configFlagV    = &fValue[string]{value: fmt.Sprintf("%s/.formance/fctl.config", getHomeDir())}
	outputFlagV    = &fValue[string]{value: "plain"}
	GlobalFlags    = withGlobalFlags(flag.NewFlagSet("global", flag.ContinueOnError))
)

func GetBool(flags *flag.FlagSet, flagName string) bool {
	f := flags.Lookup(flagName)
	if f == nil {
		return false
	}

	fromEnv := strings.ToLower(os.Getenv(strcase.ToScreamingSnake(flagName)))
	if fromEnv != "" {
		return fromEnv == "true" || fromEnv == "1"
	}

	value := f.Value.String()
	if value == "" {
		return false
	}
	v, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}
	return v
}

func GetString(flagSet *flag.FlagSet, flagName string) string {
	f := flagSet.Lookup(flagName)
	if f == nil {
		return ""
	}

	envVar := os.Getenv(strcase.ToScreamingSnake(flagName))
	if envVar != "" {
		return envVar
	}

	return f.Value.String()
}

func GetStringSlice(flagSet *flag.FlagSet, flagName string) []string {

	f := flagSet.Lookup(flagName)

	if f == nil {
		return []string{}
	}

	envVar := os.Getenv(strcase.ToScreamingSnake(flagName))

	if len(envVar) > 0 {
		return strings.Split(envVar, " ")
	}

	value := f.Value.String()
	if len(value) > 0 {
		return strings.Split(value, " ")
	}

	return []string{}
}

func GetInt(flagSet *flag.FlagSet, flagName string) int {

	f := flagSet.Lookup(flagName)
	if f == nil {
		return 0
	}

	v := os.Getenv(strcase.ToScreamingSnake(flagName))
	if v != "" {
		v, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0
		}
		return int(v)
	}

	value := f.Value.String()
	if value == "" {
		return 0
	}

	i, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		return 0
	}
	return int(i)
}

func WithConfirmFlag(flagSet *flag.FlagSet) *bool {
	return flagSet.Bool("confirm", false, "Confirm the action")
}

func WithScopesFlags(flagSet *flag.FlagSet, scopes ...*flag.Flag) *flag.FlagSet {
	for _, f := range scopes {
		flagSet.Var(f.Value, f.Name, f.Usage)
	}

	return flagSet
}

func withGlobalFlags(flagSet *flag.FlagSet) *flag.FlagSet {
	flagSet.BoolVar(insecureTlsV.Get(), InsecureTlsFlag, false, "insecure TLS")
	flagSet.BoolVar(telemetryFlagV.Get(), TelemetryFlag, false, "enable telemetry")
	flagSet.BoolVar(debugFlagV.Get(), DebugFlag, false, "debug mode")
	flagSet.StringVar(profileFlagV.Get(), ProfileFlag, "", "config profile to use")
	flagSet.StringVar(configFlagV.Get(), ConfigFlag, fmt.Sprintf("%s/.formance/fctl.config", getHomeDir()), "config file to use")
	flagSet.StringVar(outputFlagV.Get(), outputFlag, "plain", "output format (plain, json)")
	return flagSet
}

func WithMetadataFlag(flag *flag.FlagSet) *flag.FlagSet {
	flag.String(MetadataFlag, "", "Metadata to use")
	return flag
}
func ConvertPFlagSetToFlagSet(pFlagSet *pflag.FlagSet) *flag.FlagSet {

	flagSet := flag.NewFlagSet("fctl", flag.ExitOnError)

	pFlagSet.VisitAll(func(f *pflag.Flag) {
		flagSet.Var(f.Value, f.Name, f.Usage)
	})

	return flagSet
}
