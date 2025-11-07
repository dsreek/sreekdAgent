package analyzers

import (
	"context"
	"testing"

	"github.com/dsreek/sreekdAgent/pkg/model"
)

func TestStorageCapacityAnalyzer(t *testing.T) {
	snap := model.Snapshot{
		Volumes: []model.Volume{
			{ID: "vol-1", Name: "data", SizeGB: 100, UsedGB: 95},
			{ID: "vol-2", Name: "logs", SizeGB: 100, UsedGB: 80},
			{ID: "vol-3", Name: "backup", SizeGB: 200, UsedGB: 50},
		},
	}

	analyzer := NewStorageCapacityAnalyzer(StorageCapacityConfig{
		WarnThreshold:     0.75,
		CriticalThreshold: 0.9,
	})

	issues := analyzer.Diagnose(context.Background(), snap)
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
	if issues[0].Resource != "vol-1" {
		t.Fatalf("expected vol-1 critical issue, got %#v", issues[0])
	}
}
