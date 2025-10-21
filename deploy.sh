#!/bin/bash

set -e

echo "========================================="
echo "Deploying viddl.me"
echo "========================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
PROJECT_DIR="/var/www/viddl.me"
BACKEND_DIR="$PROJECT_DIR/backend"
FRONTEND_DIR="$PROJECT_DIR/frontend"
SERVICE_NAME="viddl.service"

# Check if running as root or with sudo
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}Please run with sudo${NC}"
    exit 1
fi

# Step 1: Pull latest code
echo -e "\n${YELLOW}[1/6] Pulling latest code from GitHub...${NC}"
cd "$PROJECT_DIR"
git pull origin master
echo -e "${GREEN}✓ Code updated${NC}"

# Step 2: Build backend
echo -e "\n${YELLOW}[2/6] Building Go backend...${NC}"
cd "$BACKEND_DIR"
go build -o viddl-server main.go
echo -e "${GREEN}✓ Backend built successfully${NC}"

# Step 3: Build frontend
echo -e "\n${YELLOW}[3/6] Building Vue.js frontend...${NC}"
cd "$FRONTEND_DIR"
npm install
npm run build
echo -e "${GREEN}✓ Frontend built successfully${NC}"

# Step 4: Set correct permissions
echo -e "\n${YELLOW}[4/6] Setting permissions...${NC}"
chown -R www-data:www-data "$PROJECT_DIR"
chmod +x "$BACKEND_DIR/viddl-server"
chmod 755 "$BACKEND_DIR/tmp" 2>/dev/null || mkdir -p "$BACKEND_DIR/tmp" && chmod 755 "$BACKEND_DIR/tmp"
echo -e "${GREEN}✓ Permissions set${NC}"

# Step 5: Restart backend service
echo -e "\n${YELLOW}[5/6] Restarting backend service...${NC}"
systemctl restart "$SERVICE_NAME"
sleep 2

# Check if service started successfully
if systemctl is-active --quiet "$SERVICE_NAME"; then
    echo -e "${GREEN}✓ Service restarted successfully${NC}"
else
    echo -e "${RED}✗ Service failed to start${NC}"
    echo "Check logs with: sudo journalctl -u $SERVICE_NAME -n 50"
    exit 1
fi

# Step 6: Reload nginx
echo -e "\n${YELLOW}[6/6] Reloading nginx...${NC}"
nginx -t && systemctl reload nginx
echo -e "${GREEN}✓ Nginx reloaded${NC}"

# Final status check
echo -e "\n${GREEN}=========================================${NC}"
echo -e "${GREEN}Deployment completed successfully!${NC}"
echo -e "${GREEN}=========================================${NC}"
echo ""
echo "Service status:"
systemctl status "$SERVICE_NAME" --no-pager -l | head -n 10
echo ""
echo "View logs: sudo journalctl -u $SERVICE_NAME -f"
echo "Website: https://viddl.me"
