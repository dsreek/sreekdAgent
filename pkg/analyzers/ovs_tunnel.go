  package analyzers

  import (
      "context"
      "fmt"
      "strings"
      "time"

      "github.com/dsreek/sreekdAgent/pkg/model"
  )

  type OVSTunnelConfig struct {
      CriticalStates []string
      MaxPacketLoss  float64
      MaxDownTime    time.Duration
  }

  type OVSTunnelAnalyzer struct {
      cfg         OVSTunnelConfig
      stateFilter map[string]struct{}
  }

  func NewOVSTunnelAnalyzer(cfg OVSTunnelConfig) *OVSTunnelAnalyzer {
      if len(cfg.CriticalStates) == 0 {
          cfg.CriticalStates = []string{"DOWN", "DEGRADED"}
      }
      if cfg.MaxPacketLoss == 0 {
          cfg.MaxPacketLoss = 0.05
      }
      if cfg.MaxDownTime == 0 {
          cfg.MaxDownTime = 5 * time.Minute
      }
      stateFilter := make(map[string]struct{}, len(cfg.CriticalStates))
      for _, state := range cfg.CriticalStates {
          stateFilter[strings.ToUpper(state)] = struct{}{}
      }
      return &OVSTunnelAnalyzer{cfg: cfg, stateFilter: stateFilter}
  }

  func (a *OVSTunnelAnalyzer) ID() string { return "ovs-tunnel" }

  func (a *OVSTunnelAnalyzer) Diagnose(ctx context.Context, snap model.Snapshot) []model.Issue {
      var issues []model.Issue
      for _, tunnel := range snap.Tunnels {
          state := strings.ToUpper(tunnel.Status)
          if _, bad := a.stateFilter[state]; !bad && tunnel.PacketLoss <= a.cfg.MaxPacketLoss {
              continue
          }
          if _, bad := a.stateFilter[state]; bad && time.Since(tunnel.LastObserved) < a.cfg.MaxDownTime {
              continue
          }
          message := fmt.Sprintf("OVS tunnel %s→%s %s (loss %.2f%%)", tunnel.Source, tunnel.Destination, tunnel.Status, tunnel.PacketLoss*100)
          issues = append(issues, model.Issue{
              ID:          a.ID(),
              Severity:    "high",
              Resource:    tunnel.ID,
              Message:     message,
              Remediation: fmt.Sprintf("cloud-admin network ovs-tunnels restart %s", tunnel.ID),
          })
          continue
      }
      return issues
  }
