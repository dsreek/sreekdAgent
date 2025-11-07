package analyzers

  import (
      "context"
      "testing"
      "time"

      "github.com/dsreek/sreekdAgent/pkg/model"
  )

  func TestOVSTunnelAnalyzer(t *testing.T) {
      now := time.Now()
      snap := model.Snapshot{
          Tunnels: []model.OVSTunnel{
              {ID: "tun-1", Source: "node-a", Destination: "node-b", Status: "DOWN", PacketLoss: 0.00, LastObserved: now.Add(-10 * time.Minute)},
              {ID: "tun-2", Source: "node-c", Destination: "node-d", Status: "UP", PacketLoss: 0.10, LastObserved: now.Add(-2 * time.Minute)},
              {ID: "tun-3", Source: "node-e", Destination: "node-f", Status: "UP", PacketLoss: 0.01, LastObserved: now},
          },
      }

      analyzer := NewOVSTunnelAnalyzer(OVSTunnelConfig{
          CriticalStates: []string{"DOWN"},
          MaxPacketLoss:  0.05,
          MaxDownTime:    5 * time.Minute,
      })

      issues := analyzer.Diagnose(context.Background(), snap)
      if len(issues) != 2 {
          t.Fatalf("expected 2 issues (down + high loss), got %#v", issues)
      }
      if issues[0].Resource != "tun-1" {
          t.Fatalf("first issue should be tun-1, got %#v", issues[0])
      }
      if issues[1].Resource != "tun-2" {
          t.Fatalf("second issue should be tun-2, got %#v", issues[1])
      }
  }
