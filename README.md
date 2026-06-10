# ISP DDoS Protection Platform

A production-ready, ISP-grade DDoS protection platform supporting XDP (fast path), TC eBPF (mid path), nftables (VM/fallback path), and userspace fallback (emergency mode).

## Architecture

```
+--------------------------+
| PHP Web Dashboard        |
| (Admin + Router UI)      |
+-----------+--------------+
            |
     REST API (JSON)
            |
+-----------v--------------+
| Central Controller       |
| (Go backend)             |
+-----------+--------------+
            |
    +-------+-------+
    |               |
+---v---+       +---v---+
| Agent |       | Agent |
+---+---+       +-------+
    |
+---v-------------------------------+
| XDP native → XDP generic →        |
| TC clsact → nftables → userspace  |
+-----------------------------------+
```

## Quick Start

### 1. Deploy the Controller

```bash
git clone https://github.com/V6Direct/isp-ddos-platform
cd isp-ddos-platform/controller
go build -o ddos-controller .
./ddos-controller --db /var/lib/ddos/controller.db --listen 0.0.0.0:8080
```

### 2. Deploy the PHP Dashboard

```bash
cp -r dashboard/ /var/www/html/ddos-dashboard
# Edit dashboard/config.php - set CONTROLLER_URL
# Configure your webserver (Apache/Nginx) to serve dashboard/
```

### 3. Install Agent on Router

```bash
curl https://YOUR_CONTROLLER_IP:8080/install.sh | bash
# OR manually:
cd agent && go build -o ddos-agent .
./ddos-agent --controller https://YOUR_CONTROLLER_IP:8080 --token YOUR_TOKEN --register
```

### 4. Enable XDP Protection

From the PHP dashboard, click **"Enable XDP Protection"** on any router — or via API:

```bash
curl -X POST https://controller/api/routers/1/enable-xdp \
  -H 'Authorization: Bearer ADMIN_TOKEN'
```

## Directory Structure

```
├── controller/          # Go central controller + REST API
├── agent/               # Go router agent (runs on each router)
├── xdp/                 # XDP C program + Makefile
├── dashboard/           # PHP web dashboard
├── nftables/            # nftables template generator
├── scripts/             # install.sh and helpers
└── docs/                # Architecture + deployment docs
```

## Protection Tiers

| Tier | Technology | Throughput | Use Case |
|------|-----------|-----------|----------|
| 1 | XDP native | 10-25 Mpps/core | Bare metal, high-end NIC |
| 2 | XDP generic | 2-5 Mpps | Any Linux kernel |
| 3 | TC clsact eBPF | 1-3 Mpps | VMs, containers |
| 4 | nftables | 500k pps | Light protection |
| 5 | Userspace | 100k pps | Emergency only |

## Security

- Agent authenticates to controller via pre-shared token
- All API endpoints require Bearer token auth
- BGP peers are always allowlisted (fail-safe)
- Pinned BPF maps survive agent restart
- Emergency mode enforces drop rules immediately

## License

MIT
