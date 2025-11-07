package analyzers

import (
	"context"
	"strings"

	"github.com/dsreek/sreekdAgent/pkg/model"
)

type Analyzer interface {
	ID() string
	Diagnose(ctx context.Context, snap model.Snapshot) []model.Issue
}

type Config struct {
	Enabled    []string
	IdleGPU    IdleGPUConfig
	OVSTunnel  OVSTunnelConfig
	StorageCap StorageCapacityConfig
}

// SetEnabledCSV populates Enabled from a comma-separated flag.
func (c *Config) SetEnabledCSV(csv string) {
	if csv == "" {
		return
	}
	fields := strings.Split(csv, ",")
	for i := range fields {
		fields[i] = strings.TrimSpace(fields[i])
	}
	c.Enabled = fields
}

func Execute(ctx context.Context, snap model.Snapshot, cfg Config) []model.Issue {
	var issues []model.Issue
	enabled := make(map[string]struct{})
	for _, id := range cfg.Enabled {
		if id == "" {
			continue
		}
		enabled[id] = struct{}{}
	}

	for _, analyzer := range catalog(cfg) {
		if len(enabled) > 0 {
			if _, ok := enabled[analyzer.ID()]; !ok {
				continue
			}
		}
		issues = append(issues, analyzer.Diagnose(ctx, snap)...)
	}
	return issues
}

func catalog(cfg Config) []Analyzer {
	return []Analyzer{
		NewIdleGPUAnalyzer(cfg.IdleGPU),
		NewOVSTunnelAnalyzer(cfg.OVSTunnel),
		NewStorageCapacityAnalyzer(cfg.StorageCap),
	}
}
