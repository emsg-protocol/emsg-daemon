#!/bin/bash
# deploy-sandipwalke.sh
# Deployment script for EMSG daemon on sandipwalke.com

set -e

echo "🚀 EMSG Daemon Deployment for sandipwalke.com"
echo "=============================================="

# Configuration
DOMAIN="sandipwalke.com"
EMSG_SUBDOMAIN="emsg.${DOMAIN}"
EMSG_USER="emsg"
INSTALL_DIR="/opt/emsg"
DATA_DIR="/var/lib/emsg"
SERVICE_NAME="emsg-daemon"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   print_error "This script must be run as root (use sudo)"
   exit 1
fi

echo "📋 Pre-deployment checklist:"
echo "1. DNS TXT record: _emsg.${DOMAIN} → https://${EMSG_SUBDOMAIN}:8765"
echo "2. A record: ${EMSG_SUBDOMAIN} → your server IP"
echo "3. EMSG daemon binary ready"
echo ""
read -p "Have you completed the DNS setup? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    print_warning "Please complete DNS setup first, then run this script again"
    exit 1
fi

# Step 1: Create user and directories
echo "📁 Setting up user and directories..."
if ! id "$EMSG_USER" &>/dev/null; then
    useradd -r -s /bin/false $EMSG_USER
    print_status "Created user: $EMSG_USER"
else
    print_status "User $EMSG_USER already exists"
fi

mkdir -p $INSTALL_DIR $DATA_DIR
chown $EMSG_USER:$EMSG_USER $DATA_DIR
print_status "Created directories"

# Step 2: Install daemon binary
echo "📦 Installing EMSG daemon..."
if [[ -f "./daemon" ]]; then
    cp ./daemon $INSTALL_DIR/emsg-daemon
    chown $EMSG_USER:$EMSG_USER $INSTALL_DIR/emsg-daemon
    chmod +x $INSTALL_DIR/emsg-daemon
    print_status "Installed daemon binary"
elif [[ -f "./emsg-daemon" ]]; then
    cp ./emsg-daemon $INSTALL_DIR/emsg-daemon
    chown $EMSG_USER:$EMSG_USER $INSTALL_DIR/emsg-daemon
    chmod +x $INSTALL_DIR/emsg-daemon
    print_status "Installed daemon binary"
else
    print_error "Daemon binary not found. Please build it first: go build ./cmd/daemon"
    exit 1
fi

# Step 3: Create systemd service
echo "⚙️  Creating systemd service..."
cat > /etc/systemd/system/${SERVICE_NAME}.service << EOF
[Unit]
Description=EMSG Daemon for ${DOMAIN}
After=network.target

[Service]
Type=simple
User=${EMSG_USER}
Group=${EMSG_USER}
WorkingDirectory=${INSTALL_DIR}
ExecStart=${INSTALL_DIR}/emsg-daemon
Environment=EMSG_DOMAIN=${DOMAIN}
Environment=EMSG_DATABASE_URL=${DATA_DIR}/emsg.db
Environment=EMSG_PORT=8765
Environment=EMSG_LOG_LEVEL=info
Environment=EMSG_MAX_CONNECTIONS=1000
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable $SERVICE_NAME
print_status "Created systemd service"

# Step 4: Configure firewall
echo "🔥 Configuring firewall..."
if command -v ufw &> /dev/null; then
    ufw allow 22/tcp    # SSH
    ufw allow 80/tcp    # HTTP
    ufw allow 443/tcp   # HTTPS
    ufw allow 8765/tcp  # EMSG
    print_status "Configured UFW firewall"
elif command -v firewall-cmd &> /dev/null; then
    firewall-cmd --permanent --add-port=22/tcp
    firewall-cmd --permanent --add-port=80/tcp
    firewall-cmd --permanent --add-port=443/tcp
    firewall-cmd --permanent --add-port=8765/tcp
    firewall-cmd --reload
    print_status "Configured firewalld"
else
    print_warning "No firewall detected. Please configure manually"
fi

# Step 5: Install Nginx (optional)
echo "🌐 Setting up Nginx reverse proxy..."
read -p "Install and configure Nginx? (Y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]] || [[ -z $REPLY ]]; then
    if command -v apt &> /dev/null; then
        apt update && apt install -y nginx
    elif command -v yum &> /dev/null; then
        yum install -y nginx
    else
        print_warning "Please install Nginx manually"
    fi
    
    # Create Nginx config
    cat > /etc/nginx/sites-available/${EMSG_SUBDOMAIN} << EOF
server {
    listen 80;
    server_name ${EMSG_SUBDOMAIN};
    
    location / {
        proxy_pass http://localhost:8765;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
EOF
    
    ln -sf /etc/nginx/sites-available/${EMSG_SUBDOMAIN} /etc/nginx/sites-enabled/
    nginx -t && systemctl reload nginx
    print_status "Configured Nginx"
fi

# Step 6: Start the service
echo "🚀 Starting EMSG daemon..."
systemctl start $SERVICE_NAME
sleep 2

if systemctl is-active --quiet $SERVICE_NAME; then
    print_status "EMSG daemon started successfully"
else
    print_error "Failed to start EMSG daemon"
    echo "Check logs: journalctl -u $SERVICE_NAME"
    exit 1
fi

# Step 7: Test the setup
echo "🧪 Testing the setup..."
sleep 3

# Test local connection
if curl -s http://localhost:8765/api/user?address=test > /dev/null; then
    print_status "Local API responding"
else
    print_warning "Local API not responding"
fi

# Test DNS
if dig +short TXT _emsg.${DOMAIN} | grep -q "emsg"; then
    print_status "DNS TXT record found"
else
    print_warning "DNS TXT record not found or not propagated yet"
fi

echo ""
echo "🎉 Deployment completed!"
echo "=============================================="
echo "📊 Status:"
echo "   Service: systemctl status $SERVICE_NAME"
echo "   Logs:    journalctl -u $SERVICE_NAME -f"
echo "   Config:  /etc/systemd/system/${SERVICE_NAME}.service"
echo ""
echo "🔧 Next steps:"
echo "1. Set up SSL certificate:"
echo "   sudo apt install certbot python3-certbot-nginx"
echo "   sudo certbot --nginx -d ${EMSG_SUBDOMAIN}"
echo ""
echo "2. Test your setup:"
echo "   curl http://${EMSG_SUBDOMAIN}/api/user?address=test"
echo "   dig TXT _emsg.${DOMAIN}"
echo ""
echo "3. Register your first user:"
echo "   curl -X POST http://${EMSG_SUBDOMAIN}/api/user \\"
echo "     -H 'Content-Type: application/json' \\"
echo "     -d '{\"address\":\"sandip#${DOMAIN}\",\"pubkey\":\"...\",\"first_name\":\"Sandip\"}'"
echo ""
echo "📧 Your EMSG addresses:"
echo "   sandip#${DOMAIN}"
echo "   admin#${DOMAIN}"
echo "   contact#${DOMAIN}"
echo ""
print_status "Welcome to the EMSG network! 🌐"
