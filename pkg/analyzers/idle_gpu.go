package analyzers

import (
	"context"
	"fmt"
	"time"

	"github.com/dsreek/sreekdAgent/pkg/model"
)

type IdleGPUConfig struct {
	MinUtil float64
	MaxIdle time.Duration
}

type IdleGPUAnalyzer struct {
	cfg IdleGPUConfig
}

func NewIdleGPUAnalyzer(cfg IdleGPUConfig) *IdleGPUAnalyzer {
	if cfg.MinUtil == 0 {
		cfg.MinUtil = 0.25
	}
	if cfg.MaxIdle == 0 {
		cfg.MaxIdle = 30 * time.Minute
	}
	return &IdleGPUAnalyzer{cfg: cfg}
}

func (a *IdleGPUAnalyzer) ID() string { return "idle-gpu" }

func (a *IdleGPUAnalyzer) Diagnose(ctx context.Context, snap model.Snapshot) []model.Issue {
	var issues []model.Issue
	for _, inst := range snap.Instances {
		if inst.Status != "RUNNING" {
			continue
		}
		idleFor := time.Since(inst.LastActive)
		if inst.GPUUtil <= a.cfg.MinUtil && idleFor >= a.cfg.MaxIdle {
			issues = append(issues, model.Issue{
				ID:          a.ID(),
				Severity:    "medium",
				Resource:    inst.Name,
				Message:     fmt.Sprintf("GPU util %.0f%% idle for %s", inst.GPUUtil*100, idleFor.Truncate(time.Minute)),
				Remediation: fmt.Sprintf("cloud-admin instances stop %s", inst.Name),
			})
		}
	}
	return issues
}
