package model

import "time"

type Snapshot struct {
	Instances []Instance
	Tunnels   []OVSTunnel
	Volumes   []Volume
}

type Instance struct {
	Name            string    `json:"name"`
	Status          string    `json:"status"`
	Type            string    `json:"instance_type"`
	Zone            string    `json:"zone"`
	CreatedAt       time.Time `json:"created_at"`
	LastStateChange time.Time `json:"last_state_change"`
	GPUUtil         float64   `json:"gpu_util"`
	LastActive      time.Time `json:"last_active"`
	ErrorReason     string    `json:"error_reason"`
}

type OVSTunnel struct {
	ID           string    `json:"id"`
	Source       string    `json:"source"`
	Destination  string    `json:"destination"`
	Status       string    `json:"status"`
	PacketLoss   float64   `json:"packet_loss"`
	LastObserved time.Time `json:"last_observed"`
}

type Volume struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	SizeGB           float64   `json:"size_gb"`
	UsedGB           float64   `json:"used_gb"`
	Status           string    `json:"status"`
	AttachedInstance string    `json:"attached_instance"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (v Volume) Utilization() float64 {
	if v.SizeGB == 0 {
		return 0
	}
	return v.UsedGB / v.SizeGB
}

type Issue struct {
	ID          string
	Severity    string
	Resource    string
	Message     string
	Remediation string
}
