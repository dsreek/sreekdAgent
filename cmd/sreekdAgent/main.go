package main

import (
	"context"
	"log"

	"github.com/dsreek/sreekdAgent/pkg/analyzers"
	"github.com/dsreek/sreekdAgent/pkg/collectors"
	"github.com/dsreek/sreekdAgent/pkg/config"
	"github.com/dsreek/sreekdAgent/pkg/llm"
	"github.com/dsreek/sreekdAgent/pkg/output"
)

func main() {
	ctx := context.Background()
	cfg := config.FromFlags()

	runner := collectors.New(collectors.Config{
		ProjectID: cfg.Collectors.ProjectID,
		Region:    cfg.Collectors.Region,
		Profile:   cfg.Collectors.Profile,
		Binary:    cfg.Collectors.Binary,
	})

	snap, err := runner.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

	issues := analyzers.Execute(ctx, snap, cfg.Analyzers)

	if cfg.LLM.Enabled {
		client := llm.New(cfg.LLM)
		issues = client.Enrich(issues)
	}

	output.Render(issues, cfg.Output)
}
