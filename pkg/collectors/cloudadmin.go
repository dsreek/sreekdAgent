package collectors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/dsreek/sreekdAgent/pkg/model"
)

type CloudAdmin struct {
	binary  string
	profile string
	project string
	region  string
}

func NewCloudAdmin(cfg Config) *CloudAdmin {
	bin := cfg.Binary
	if bin == "" {
		bin = "cloud-admin"
	}
	return &CloudAdmin{
		binary:  bin,
		profile: cfg.Profile,
		project: cfg.ProjectID,
		region:  cfg.Region,
	}
}

func (c *CloudAdmin) Instances(ctx context.Context) ([]model.Instance, error) {
	var out []model.Instance
	err := c.run(ctx, []string{"instances", "list"}, &out)
	return out, err
}

func (c *CloudAdmin) OVSTunnels(ctx context.Context) ([]model.OVSTunnel, error) {
	var out []model.OVSTunnel
	err := c.run(ctx, []string{"network", "ovs-tunnels", "list"}, &out)
	return out, err
}

func (c *CloudAdmin) Volumes(ctx context.Context) ([]model.Volume, error) {
	var out []model.Volume
	err := c.run(ctx, []string{"volumes", "list"}, &out)
	return out, err
}

func (c *CloudAdmin) run(ctx context.Context, args []string, target any) error {
	cmdArgs := append(args, "--output", "json")
	if c.project != "" {
		cmdArgs = append(cmdArgs, "--project-id", c.project)
	}
	if c.region != "" {
		cmdArgs = append(cmdArgs, "--region", c.region)
	}
	if c.profile != "" {
		cmdArgs = append(cmdArgs, "--profile", c.profile)
	}

	cmd := exec.CommandContext(ctx, c.binary, cmdArgs...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &bytes.Buffer{}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cloud-admin %v: %w", cmdArgs, err)
	}
	return json.Unmarshal(stdout.Bytes(), target)
}
