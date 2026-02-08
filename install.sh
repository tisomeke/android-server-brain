#!/bin/bash

# --- Color Definitions ---
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${CYAN}🚀 Starting Android Server Brain (ASB) Setup...${NC}"

# 1. Check/Request Storage Access
if [ ! -d "$HOME/storage" ]; then
    echo -e "${YELLOW}📂 Requesting storage access... Please click 'Allow' on the popup.${NC}"
    termux-setup-storage
    sleep 5 # Give user time to click
fi

# 2. Install Dependencies
echo -e "${CYAN}📦 Installing dependencies (golang, termux-api)...${NC}"
pkg update && pkg upgrade -y
pkg install golang termux-api -y

# 3. Interactive Configuration
echo -e "${GREEN}⚙️ Configuration Setup${NC}"

read -p "Enter your Telegram Bot Token: " TELE_TOKEN
read -p "Enter your Telegram Admin ID: " ADMIN_ID
read -p "Enter storage directory [default: downloads/server]: " STORAGE_DIR
STORAGE_DIR=${STORAGE_DIR:-downloads/server}

# Create the config.json file
cat <<EOF > config.json
{
  "telegram_token": "$TELE_TOKEN",
  "admin_id": $ADMIN_ID,
  "storage_dir": "$STORAGE_DIR"
}
EOF

echo -e "${GREEN}✅ config.json created successfully.${NC}"

# 4. Directory and Symlink Setup
FULL_STORAGE_PATH="$HOME/$STORAGE_DIR"
echo -e "${CYAN}📁 Creating directory: $FULL_STORAGE_PATH${NC}"
mkdir -p "$FULL_STORAGE_PATH"

echo -e "${CYAN}🔗 Creating symlink ~/server -> $STORAGE_DIR${NC}"
if [ -L "$HOME/server" ]; then
    rm "$HOME/server"
fi
ln -s "$FULL_STORAGE_PATH" "$HOME/server"

# 5. Initialize Go Project
echo -e "${CYAN}🏗️ Initializing Go project...${NC}"
if [ ! -f "go.mod" ]; then
    go mod init android-server-brain
fi
go mod tidy

echo -e "${GREEN}--- SETUP COMPLETE ---${NC}"
echo -e "To start your server, run: ${YELLOW}go run main.go${NC}"
