# Deployment Guide

## Requirements

### Controller Server
- Linux (any distro)
- Go 1.21+
- SQLite (default) or MySQL
- Port 8080 open to agents
- PHP 8.x + web server for dashboard

### Router/Agent Requirements
- Linux kernel 5.10+ (5.15+ recommended for XDP native)
- For XDP native: driver support (Intel i40e, mlx5, ixgbe, etc.)
- For XDP generic: any NIC, kernel 4.8+
- For TC eBPF: kernel 4.1+
- nftables: kernel 3.13+ with nft userspace tool
- Go 1.21+ to build agent (or use pre-built binary)
- Root/CAP_NET_ADMIN capability
- Mounted `/sys/fs/bpf`

## Step 1: Deploy Controller

```bash
# Install dependencies
apt-get install -y golang sqlite3

# Build controller
cd controller/
go build -o /usr/local/bin/ddos-controller .

# Create config
mkdir -p /etc/ddos /var/lib/ddos
cat > /etc/ddos/controller.yaml << EOF
listen: "0.0.0.0:8080"
database: "/var/lib/ddos/controller.db"
admin_token: "$(openssl rand -hex 32)"
agent_token_secret: "$(openssl rand -hex 32)"
EOF

# Install systemd service
cp ../scripts/ddos-controller.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable --now ddos-controller
```

## Step 2: Deploy Dashboard

```bash
# Copy dashboard files
cp -r ../dashboard/ /var/www/html/ddos/

# Configure controller URL
cat > /var/www/html/ddos/config.php << EOF
<?php
define('CONTROLLER_URL', 'http://YOUR_CONTROLLER_IP:8080');
define('ADMIN_PASSWORD', 'changeme123');
define('SESSION_SECRET', '$(openssl rand -hex 32)');
EOF

# Set permissions
chown -R www-data:www-data /var/www/html/ddos/
chmod 640 /var/www/html/ddos/config.php
```

### Nginx config example
```nginx
server {
    listen 80;
    server_name dashboard.your-isp.net;
    root /var/www/html/ddos;
    index index.php;
    location ~ \.php$ {
        fastcgi_pass unix:/run/php/php8.1-fpm.sock;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include fastcgi_params;
    }
}
```

## Step 3: Install Agent on Router

### Automated (recommended)
```bash
curl http://YOUR_CONTROLLER_IP:8080/install.sh | bash
```

### Manual
```bash
# On controller, generate router token
curl -X POST http://controller:8080/api/routers \
  -H 'Authorization: Bearer ADMIN_TOKEN' \
  -H 'Content-Type: application/json' \
  -d '{"name": "router-01", "ip": "192.0.2.1"}'
# Returns: {"id": 1, "token": "abc123..."}

# On router
cd agent/
go build -o /usr/local/bin/ddos-agent .

# Mount BPF fs if not mounted
mount bpffs /sys/fs/bpf -t bpf
echo 'bpffs /sys/fs/bpf bpf defaults 0 0' >> /etc/fstab

# Register and start
ddos-agent \
  --controller http://YOUR_CONTROLLER_IP:8080 \
  --token abc123... \
  --iface eth0  # optional: auto-detected if omitted

# Install as service
cp ../scripts/ddos-agent.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable --now ddos-agent
```

## Step 4: Enable XDP via Dashboard

1. Open dashboard in browser
2. Log in with admin credentials
3. Click on router → **"Enable XDP Protection"**
4. Dashboard shows: mode, interface, and latency

## Upgrading

```bash
# Controller
go build -o /usr/local/bin/ddos-controller ./controller/
systemctl restart ddos-controller

# Agent (on router)
go build -o /tmp/ddos-agent-new ./agent/
systemctl stop ddos-agent
cp /tmp/ddos-agent-new /usr/local/bin/ddos-agent
systemctl start ddos-agent
# BPF maps are pinned - traffic protected during restart
```
