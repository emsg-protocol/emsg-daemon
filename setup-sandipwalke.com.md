# EMSG Protocol Setup for sandipwalke.com

## Quick Setup Guide

### 1. DNS Configuration

Add this TXT record to your DNS provider (wherever you manage sandipwalke.com DNS):

```dns
Name: _emsg.sandipwalke.com
Type: TXT
Value: https://emsg.sandipwalke.com:8765
TTL: 3600
```

### 2. Server Deployment

#### Option A: Linux Server (Ubuntu/Debian)

```bash
# 1. Create EMSG user
sudo useradd -r -s /bin/false emsg
sudo mkdir -p /opt/emsg /var/lib/emsg
sudo chown emsg:emsg /var/lib/emsg

# 2. Upload and setup daemon
sudo cp emsg-daemon /opt/emsg/
sudo chown emsg:emsg /opt/emsg/emsg-daemon
sudo chmod +x /opt/emsg/emsg-daemon

# 3. Create systemd service
sudo tee /etc/systemd/system/emsg-daemon.service > /dev/null <<EOF
[Unit]
Description=EMSG Daemon for sandipwalke.com
After=network.target

[Service]
Type=simple
User=emsg
Group=emsg
WorkingDirectory=/opt/emsg
ExecStart=/opt/emsg/emsg-daemon
Environment=EMSG_DOMAIN=sandipwalke.com
Environment=EMSG_DATABASE_URL=/var/lib/emsg/emsg.db
Environment=EMSG_PORT=8765
Environment=EMSG_LOG_LEVEL=info
Environment=EMSG_MAX_CONNECTIONS=1000
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# 4. Start the service
sudo systemctl enable emsg-daemon
sudo systemctl start emsg-daemon
sudo systemctl status emsg-daemon
```

#### Option B: Docker Deployment

```bash
# 1. Create docker-compose.yml
cat > docker-compose.yml <<EOF
version: '3.8'
services:
  emsg-daemon:
    build: .
    ports:
      - "8765:8765"
    environment:
      - EMSG_DOMAIN=sandipwalke.com
      - EMSG_DATABASE_URL=/data/emsg.db
      - EMSG_PORT=8765
      - EMSG_LOG_LEVEL=info
    volumes:
      - ./data:/data
    restart: unless-stopped
EOF

# 2. Build and run
docker-compose up -d
```

### 3. Nginx Reverse Proxy (Optional but Recommended)

```nginx
# /etc/nginx/sites-available/emsg.sandipwalke.com
server {
    listen 80;
    server_name emsg.sandipwalke.com;
    
    # Redirect HTTP to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name emsg.sandipwalke.com;
    
    # SSL Configuration (use Let's Encrypt)
    ssl_certificate /etc/letsencrypt/live/emsg.sandipwalke.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/emsg.sandipwalke.com/privkey.pem;
    
    location / {
        proxy_pass http://localhost:8765;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 4. SSL Certificate (Let's Encrypt)

```bash
# Install certbot
sudo apt install certbot python3-certbot-nginx

# Get certificate
sudo certbot --nginx -d emsg.sandipwalke.com

# Auto-renewal
sudo crontab -e
# Add: 0 12 * * * /usr/bin/certbot renew --quiet
```

### 5. Firewall Configuration

```bash
# Allow necessary ports
sudo ufw allow 22    # SSH
sudo ufw allow 80    # HTTP
sudo ufw allow 443   # HTTPS
sudo ufw allow 8765  # EMSG (if not using reverse proxy)
sudo ufw enable
```

### 6. Testing Your Setup

```bash
# Test DNS resolution
dig TXT _emsg.sandipwalke.com

# Test EMSG daemon
curl https://emsg.sandipwalke.com/api/user?address=test

# Test address validation
curl -X POST https://emsg.sandipwalke.com/api/route/validate \
  -H "Content-Type: application/json" \
  -d '{"addresses":["test@sandipwalke.com"]}'
```

## EMSG Addresses for Your Domain

Once set up, users can have EMSG addresses like:
- `sandip#sandipwalke.com`
- `admin#sandipwalke.com`
- `contact#sandipwalke.com`

## Example User Registration

```bash
# Register a user (replace with real Ed25519 public key)
curl -X POST https://emsg.sandipwalke.com/api/user \
  -H "Content-Type: application/json" \
  -d '{
    "address": "sandip#sandipwalke.com",
    "pubkey": "your-base64-ed25519-public-key",
    "first_name": "Sandip",
    "last_name": "Walke",
    "display_picture": "https://sandipwalke.com/avatar.jpg"
  }'
```

## Monitoring and Maintenance

```bash
# Check daemon status
sudo systemctl status emsg-daemon

# View logs
sudo journalctl -u emsg-daemon -f

# Check database size
ls -lh /var/lib/emsg/emsg.db

# Backup database
sudo cp /var/lib/emsg/emsg.db /backup/emsg-$(date +%Y%m%d).db
```

## Security Checklist

- [ ] DNS TXT record configured
- [ ] SSL certificate installed
- [ ] Firewall configured
- [ ] Daemon running as non-root user
- [ ] Database file permissions secured
- [ ] Regular backups scheduled
- [ ] Monitoring set up

## Next Steps

1. Set up the DNS TXT record
2. Deploy the daemon to your server
3. Configure SSL/TLS
4. Test the setup
5. Register your first EMSG user
6. Start using decentralized messaging!

Your domain will then be part of the EMSG network, allowing secure, decentralized messaging with other EMSG-enabled domains.
