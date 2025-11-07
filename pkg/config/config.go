package config

import (
	"flag"
	"strings"

	"github.com/dsreek/sreekdAgent/pkg/analyzers"
	"github.com/dsreek/sreekdAgent/pkg/collectors"
	"github.com/dsreek/sreekdAgent/pkg/llm"
	"github.com/dsreek/sreekdAgent/pkg/output"
)

type Config struct {
	Collectors collectors.Config
	Analyzers  analyzers.Config
	LLM        llm.Config
	Output     output.Config
	Auto       AutoConfig
}

type AutoConfig struct {
	Enabled bool
	Allow   []string
}

func Default() Config {
	return Config{
		Collectors: collectors.Config{},
		Analyzers:  analyzers.Config{},
		LLM:        llm.Config{},
		Output:     output.Config{Format: "table"},
		Auto:       AutoConfig{},
	}
}

func FromFlags() Config {
	cfg := Default()

	flag.StringVar(&cfg.Collectors.ProjectID, "project-id", "", "Crusoe project ID to scope scans")
	flag.StringVar(&cfg.Collectors.Region, "region", "", "Crusoe region filter (optional)")
	flag.StringVar(&cfg.Collectors.Profile, "profile", "", "cloud-admin profile/credential name")

	flag.StringVar(&cfg.Output.Format, "output", cfg.Output.Format, "Output format (table|json|sarif)")

	flag.BoolVar(&cfg.LLM.Enabled, "llm", false, "Enable LLM enrichment")
	flag.StringVar(&cfg.LLM.Provider, "llm-provider", "openai", "LLM provider id")
	flag.StringVar(&cfg.LLM.Model, "llm-model", "", "LLM model name")

	flag.BoolVar(&cfg.Auto.Enabled, "auto", false, "Enable auto-remediation")
	autoAllow := flag.String("auto-allow", "", "Comma-separated analyzer IDs allowed to auto-remediate")

	analyzersCSV := flag.String("analyzers", "", "Comma-separated analyzer IDs to run (default: all)")

	flag.Parse()

	cfg.Analyzers.SetEnabledCSV(*analyzersCSV)

	if *autoAllow != "" {
		cfg.Auto.Allow = splitAndTrim(*autoAllow)
	}

	return cfg
}

func splitAndTrim(csv string) []string {
	if csv == "" {
		return nil
	}
	parts := strings.Split(csv, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
