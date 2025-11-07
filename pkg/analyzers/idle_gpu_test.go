package analyzers

import (
	"context"
	"testing"
	"time"

	"github.com/dsreek/sreekdAgent/pkg/model"
)

func TestIdleGPUAnalyzer(t *testing.T) {
	now := time.Now()
	snap := model.Snapshot{
		Instances: []model.Instance{
			{Name: "gpu-1", Status: "RUNNING", GPUUtil: 0.1, LastActive: now.Add(-1 * time.Hour)},
			{Name: "gpu-2", Status: "RUNNING", GPUUtil: 0.8, LastActive: now.Add(-10 * time.Minute)},
		},
	}
	analyzer := NewIdleGPUAnalyzer(IdleGPUConfig{
		MinUtil: 0.2,
		MaxIdle: 30 * time.Minute,
	})
	issues := analyzer.Diagnose(context.Background(), snap)
	if len(issues) != 1 || issues[0].Resource != "gpu-1" {
		t.Fatalf("unexpected issues: %#v", issues)
	}
}
