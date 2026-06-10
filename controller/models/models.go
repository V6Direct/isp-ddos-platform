package models

import "time"

type Router struct {
	ID       int64      `json:"id"`
	Name     string     `json:"name"`
	IP       string     `json:"ip"`
	Token    string     `json:"token,omitempty"`
	Mode     string     `json:"mode"`
	Iface    string     `json:"iface"`
	LastSeen *time.Time `json:"last_seen"`
	CreatedAt time.Time `json:"created_at"`
	Online   bool       `json:"online"`
}

type Rules struct {
	ID             int64     `json:"id"`
	RouterID       *int64    `json:"router_id"`
	Allowlist      []string  `json:"allowlist"`
	Denylist       []string  `json:"denylist"`
	RateLimits     RateLimits `json:"rate_limits"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type RateLimits struct {
	ICMPPPS    int `json:"icmp_pps"`
	UDPPPS     int `json:"udp_pps"`
	SYNPPS     int `json:"syn_pps"`
	BGPPPS     int `json:"bgp_pps"`
	WgPPS      int `json:"wg_pps"`
	PerIPPPS   int `json:"per_ip_pps"`
}

type Stats struct {
	ID        int64     `json:"id"`
	RouterID  int64     `json:"router_id"`
	Timestamp time.Time `json:"timestamp"`
	XDPDrops  int64     `json:"xdp_drops"`
	TCDrops   int64     `json:"tc_drops"`
	NFTDrops  int64     `json:"nft_drops"`
	PPS       int64     `json:"pps"`
	BPS       int64     `json:"bps"`
}

type AgentHeartbeat struct {
	Mode     string `json:"mode"`
	Iface    string `json:"iface"`
	XDPDrops int64  `json:"xdp_drops"`
	TCDrops  int64  `json:"tc_drops"`
	NFTDrops int64  `json:"nft_drops"`
	PPS      int64  `json:"pps"`
	BPS      int64  `json:"bps"`
}

type XDPEnableResponse struct {
	Status  string `json:"status"`
	Mode    string `json:"mode"`
	Iface   string `json:"iface"`
	Latency string `json:"latency"`
	Error   string `json:"error,omitempty"`
}
