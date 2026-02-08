#!/bin/bash

# Android Server Brain (ASB) Automated Installation Script
# This script automates the complete setup process for ASB

# --- Color Definitions ---
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# --- Functions ---
print_header() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}  Android Server Brain - Setup Wizard  ${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo
}

print_step() {
    echo -e "${CYAN}▶ Step $1: $2${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}! $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

check_command() {
    if ! command -v "$1" &> /dev/null; then
        print_error "$1 is not installed. Please install it first."
        exit 1
    fi
}

print_header
echo -e "${CYAN}🚀 Starting automated ASB installation...${NC}"
echo

print_step "1" "Checking prerequisites and requesting storage access"

# Check if running in Termux
if [ ! -d "/data/data/com.termux" ]; then
    print_error "This script must be run in Termux environment"
    exit 1
fi

# Check/Request Storage Access
if [ ! -d "$HOME/storage" ]; then
    print_warning "Requesting storage access... Please click 'Allow' on the popup."
    termux-setup-storage
    sleep 5 # Give user time to click
    if [ ! -d "$HOME/storage" ]; then
        print_error "Storage access not granted. Please run 'termux-setup-storage' manually and try again."
        exit 1
    fi
fi
print_success "Storage access confirmed"

print_step "2" "Making scripts executable"

# Make all .sh files executable
find . -name "*.sh" -exec chmod +x {} \;
print_success "All shell scripts made executable"

print_step "3" "Installing system dependencies"

# Update package lists
print_warning "Updating package lists..."
pkg update -q

# Install required packages
print_warning "Installing required packages: golang, git, termux-api..."
pkg install golang git termux-api -y -q

# Verify installations
check_command "go"
check_command "git"
print_success "Dependencies installed successfully"

print_step "4" "Configuring ASB settings"

# Check if config.json already exists
if [ -f "config.json" ]; then
    print_warning "config.json already exists. Skipping configuration."
    print_warning "Delete config.json if you want to reconfigure."
else
    echo -e "${YELLOW}Please provide the following information:${NC}"
    echo -e "${YELLOW}(Get bot token from @BotFather on Telegram)${NC}"
    echo -e "${YELLOW}(Get your ID from @userinfobot on Telegram)${NC}"
    echo
    
    read -p "Enter your Telegram Bot Token: " TELE_TOKEN
    while [ -z "$TELE_TOKEN" ]; do
        print_error "Bot token cannot be empty"
        read -p "Enter your Telegram Bot Token: " TELE_TOKEN
    done
    
    read -p "Enter your Telegram Admin ID: " ADMIN_ID
    while ! [[ "$ADMIN_ID" =~ ^-?[0-9]+$ ]]; do
        print_error "Admin ID must be a number"
        read -p "Enter your Telegram Admin ID: " ADMIN_ID
    done
    
    read -p "Enter storage directory [downloads/server]: " STORAGE_DIR
    STORAGE_DIR=${STORAGE_DIR:-downloads/server}
    
    # Create the config.json file
    cat <<EOF > config.json
{
  "telegram_token": "$TELE_TOKEN",
  "admin_id": $ADMIN_ID,
  "storage_dir": "$STORAGE_DIR"
}
EOF
    
    print_success "Configuration saved to config.json"
fi

print_step "5" "Setting up storage directories"

# Read storage directory from config
if [ -f "config.json" ]; then
    STORAGE_DIR=$(grep '"storage_dir"' config.json | cut -d '"' -f 4)
fi
STORAGE_DIR=${STORAGE_DIR:-downloads/server}
FULL_STORAGE_PATH="$HOME/$STORAGE_DIR"

# Create storage directory
print_warning "Creating directory: $FULL_STORAGE_PATH"
mkdir -p "$FULL_STORAGE_PATH"

# Create symlink
print_warning "Creating symlink ~/server -> $STORAGE_DIR"
if [ -L "$HOME/server" ] || [ -e "$HOME/server" ]; then
    rm -f "$HOME/server"
fi
ln -s "$FULL_STORAGE_PATH" "$HOME/server"
print_success "Storage setup completed"

print_step "6" "Building ASB application"

# Initialize Go project if needed
if [ ! -f "go.mod" ]; then
    print_warning "Initializing Go module..."
    go mod init android-server-brain
fi

# Download dependencies
print_warning "Downloading Go dependencies..."
go mod tidy -q

# Build the application
print_warning "Building ASB binary..."
go build -o asb main.go

if [ $? -eq 0 ] && [ -f "asb" ]; then
    print_success "ASB binary built successfully"
    chmod +x asb
else
    print_error "Failed to build ASB. Check for Go compilation errors."
    exit 1
fi

print_step "7" "Setting up auto-start service"

# Setup Termux:Boot auto-start
BOOT_SCRIPT_DIR="$HOME/.termux/boot"
BOOT_SCRIPT="$BOOT_SCRIPT_DIR/asb"

mkdir -p "$BOOT_SCRIPT_DIR"

cat <<'EOF' > "$BOOT_SCRIPT"
#!/data/data/com.termux/files/usr/bin/sh

cd "$HOME/android-server-brain" || exit 1

# Wait for network connectivity
until ping -c1 google.com &>/dev/null; do
    sleep 5
done

# Start ASB in background
nohup ./asb > asb.log 2>&1 &
EOF

chmod +x "$BOOT_SCRIPT"
print_success "Auto-start service configured"

# Final summary
print_header
echo -e "${GREEN}🎉 INSTALLATION COMPLETE!${NC}"
echo
echo -e "${CYAN}Summary:${NC}"
echo -e "  • Binary: $(pwd)/asb"
echo -e "  • Config: $(pwd)/config.json"
echo -e "  • Storage: ~/server -> $STORAGE_DIR"
echo -e "  • Logs: $(pwd)/asb.log"
echo -e "  • Auto-start: $BOOT_SCRIPT"
echo
echo -e "${YELLOW}Next steps:${NC}"
echo -e "  1. Make sure Termux:Boot app is installed from F-Droid"
echo -e "  2. Grant Termux:Boot permission in Android settings"
echo -e "  3. Restart your device to test auto-start"
echo -e "  4. Or run manually: ./asb"
echo
echo -e "${GREEN}Need help? Check README.md for detailed usage instructions${NC}"
echo
echo -e "${YELLOW}🧹 To uninstall ASB later:${NC}"
echo -e "${YELLOW}  Run: ./uninstall.sh${NC}"
echo -e "${YELLOW}  This will safely remove all ASB components${NC}"
