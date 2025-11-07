package analyzers

import (
	"context"
	"fmt"

	"github.com/dsreek/sreekdAgent/pkg/model"
)

type StorageCapacityConfig struct {
	WarnThreshold     float64
	CriticalThreshold float64
}

type StorageCapacityAnalyzer struct {
	cfg StorageCapacityConfig
}

func NewStorageCapacityAnalyzer(cfg StorageCapacityConfig) *StorageCapacityAnalyzer {
	if cfg.WarnThreshold == 0 {
		cfg.WarnThreshold = 0.75
	}
	if cfg.CriticalThreshold == 0 {
		cfg.CriticalThreshold = 0.90
	}
	return &StorageCapacityAnalyzer{cfg: cfg}
}

func (a *StorageCapacityAnalyzer) ID() string { return "storage-capacity" }

func (a *StorageCapacityAnalyzer) Diagnose(ctx context.Context, snap model.Snapshot) []model.Issue {
	var issues []model.Issue
	for _, volume := range snap.Volumes {
		util := volume.Utilization()
		switch {
		case util >= a.cfg.CriticalThreshold:
			rem := fmt.Sprintf("cloud-admin volumes resize %s --size %d", volume.ID, int(volume.SizeGB*1.2))
			issues = append(issues, a.buildIssue(volume, util, "high", rem))
		case util >= a.cfg.WarnThreshold:
			rem := fmt.Sprintf("cloud-admin volumes snapshot %s", volume.ID)
			issues = append(issues, a.buildIssue(volume, util, "medium", rem))
		}
	}
	return issues
}

func (a *StorageCapacityAnalyzer) buildIssue(volume model.Volume, util float64, severity, remediation string) model.Issue {
	message := fmt.Sprintf("volume %s at %.0f%% of %.1fGiB", volume.Name, util*100, volume.SizeGB)
	return model.Issue{
		ID:          a.ID(),
		Severity:    severity,
		Resource:    volume.ID,
		Message:     message,
		Remediation: remediation,
	}
}
