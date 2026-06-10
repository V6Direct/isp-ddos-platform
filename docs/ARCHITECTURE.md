# Architecture Deep Dive

## Data Flow

### Normal Operation
1. Traffic arrives at router NIC
2. XDP program runs **before** kernel network stack
3. Bogon/blocklisted IPs → XDP_DROP (zero-copy, line rate)
4. Rate-limited IPs → token bucket check → drop or pass
5. Allowlisted IPs (BGP peers, etc.) → XDP_PASS immediately
6. Clean traffic → XDP_PASS to kernel

### Fallback Chain
If XDP is unavailable or detaches:
1. Agent detects XDP failure via map ping
2. Agent loads TC clsact eBPF program on the same interface
3. If TC fails → generate nftables ruleset and apply
4. If nftables fails → userspace pcap filter (emergency)

### Rule Sync
1. Operator updates rules in PHP dashboard
2. Dashboard POSTs to controller REST API
3. Controller writes rules to SQLite/MySQL
4. Controller pushes rules to all relevant agents via HTTP push or agent polling
5. Agent receives new rules, atomically updates BPF maps
6. Agent regenerates nftables rules as fallback
7. Agent reports acknowledgement back to controller

## Component Responsibilities

### PHP Dashboard
- Human-facing UI only
- Never talks directly to agents
- All operations go through controller REST API
- Polls `/api/routers/{id}/stats` every 3 seconds for graphs

### Go Controller
- Single binary, SQLite or MySQL backend
- Maintains router registry + rule database
- Generates install tokens
- Pushes rule deltas to agents
- Aggregates stats for dashboard
- Validates all rule changes

### Go Agent (per router)
- Runs as systemd service
- Probes kernel capabilities on startup
- Compiles/loads BPF programs from embedded byte-code
- Maintains BPF maps (allowlist, denylist, rate limits)
- Polls controller for rule updates every 10s
- Reports stats every 5s
- Performs capability detection and tier switching

### XDP Program (C + libbpf)
- Compiled once, embedded in agent binary as byte array
- CO-RE (Compile Once Run Everywhere) compatible
- Per-CPU maps for lock-free operation
- Handles: bogons, allowlist, denylist, per-IP token bucket, SYN flood, WireGuard, BGP fast-path

## BPF Map Layout

```
/sys/fs/bpf/ddos/
├── allowlist          # LPM trie: allowed prefixes
├── denylist           # LPM trie: blocked prefixes  
├── rate_limits        # Hash map: per-IP token buckets
├── stats              # Per-CPU array: drop/pass counters
└── config             # Array map: global config flags
```

## Failure Modes

| Failure | Detection | Response | Time |
|---------|-----------|----------|------|
| XDP detach | BPF map ping fails | Load TC eBPF | <5s |
| TC failure | Prog load error | Apply nftables | <10s |
| nftables failure | nft apply exit code | Userspace filter | <15s |
| Controller unreachable | HTTP timeout | Use cached rules | Immediate |
| Agent crash | systemd restart | Auto-restart | <3s |
| BPF map corruption | Verifier error | Reload maps from controller | <30s |
