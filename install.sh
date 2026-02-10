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

print_step "1" "Verifying system dependencies and environment"

# Check if running in Termux
if [ ! -d "/data/data/com.termux" ]; then
    print_error "This script must be run in Termux environment"
    exit 1
fi

# Check for Termux:Boot installation and first run
print_warning "Checking Termux:Boot status..."

# Allow skipping check for advanced users
if [ -n "$SKIP_BOOT_CHECK" ]; then
    print_warning "Skipping Termux:Boot verification (SKIP_BOOT_CHECK set)"
else
    # More flexible Termux:Boot detection
    TERMUX_BOOT_INSTALLED=false

    # Check multiple possible locations and conditions
    if [ -d "$HOME/.termux/boot" ]; then
        TERMUX_BOOT_INSTALLED=true
    elif [ -d "/data/data/com.termux/files/home/.termux/boot" ]; then
        TERMUX_BOOT_INSTALLED=true
    elif command -v termux-boot-setup &>/dev/null; then
        TERMUX_BOOT_INSTALLED=true
    fi

    # Additional check: see if Termux:Boot app has been launched
    if [ -f "/data/data/com.termux.boot/shared_prefs/com.termux.boot_preferences.xml" ]; then
        TERMUX_BOOT_INSTALLED=true
    fi

    if [ "$TERMUX_BOOT_INSTALLED" = false ]; then
        print_error "Termux:Boot not properly configured!"
        echo -e "${YELLOW}Please follow these steps:${NC}"
        echo -e "${YELLOW}1. Install Termux:Boot from F-Droid${NC}"
        echo -e "${YELLOW}2. Open the Termux:Boot app once${NC}"
        echo -e "${YELLOW}3. Grant any requested permissions${NC}"
        echo -e "${YELLOW}4. Run this installer again${NC}"
        echo
        echo -e "${YELLOW}Alternative: Skip Termux:Boot check (not recommended)${NC}"
        echo -e "${YELLOW}Run: SKIP_BOOT_CHECK=1 ./install.sh${NC}"
        exit 1
    fi
fi
print_success "Termux:Boot verified"

# Check for essential system commands
print_warning "Verifying system commands..."
check_command "pkg"
check_command "curl"
check_command "ping"
print_success "System commands verified"

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

# New storage structure: use Android's system Downloads folder
print_warning "Setting up storage in Android Downloads..."

# Ensure termux-storage-create has been run
if [ ! -d "$HOME/storage" ]; then
    print_warning "Requesting storage access... Please click 'Allow' on the popup."
    termux-setup-storage
    sleep 5
    if [ ! -d "$HOME/storage" ]; then
        print_error "Storage access not granted. Please run 'termux-setup-storage' manually and try again."
        exit 1
    fi
fi

# Create asb_files directory in system Downloads
SYSTEM_DOWNLOADS="/storage/emulated/0/Download"
ASB_FILES_DIR="$SYSTEM_DOWNLOADS/asb_files"

if [ -d "$SYSTEM_DOWNLOADS" ]; then
    print_warning "Creating asb_files in Downloads..."
    mkdir -p "$ASB_FILES_DIR"
    
    # Create symlink ~/asb_files -> /storage/emulated/0/Download/asb_files
    print_warning "Creating symlink ~/asb_files -> Downloads/asb_files"
    if [ -L "$HOME/asb_files" ] || [ -e "$HOME/asb_files" ]; then
        rm -f "$HOME/asb_files"
    fi
    ln -s "$ASB_FILES_DIR" "$HOME/asb_files"
    print_success "Storage setup completed"
    print_success "Files will be saved to: $ASB_FILES_DIR"
else
    print_error "Cannot access system Downloads folder"
    print_warning "Falling back to home directory storage"
    mkdir -p "$HOME/asb_files"
    print_success "Fallback storage created at ~/asb_files"
fi

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

print_step "7" "Setting up enhanced auto-start service"

# Setup Termux:Boot auto-start with enhanced reliability
BOOT_SCRIPT_DIR="$HOME/.termux/boot"
BOOT_SCRIPT="$BOOT_SCRIPT_DIR/asb"

mkdir -p "$BOOT_SCRIPT_DIR"

# Create enhanced boot script with better error handling
cat <<'EOF' > "$BOOT_SCRIPT"
#!/data/data/com.termux/files/usr/bin/sh

# Enhanced ASB Auto-start Script
LOG_FILE="$HOME/android-server-brain/asb-boot.log"
echo "[$(date)] Starting ASB auto-start sequence" >> "$LOG_FILE"

cd "$HOME/android-server-brain" || {
    echo "[$(date)] ERROR: Cannot cd to android-server-brain directory" >> "$LOG_FILE"
    exit 1
}

echo "[$(date)] Waiting for network connectivity..." >> "$LOG_FILE"
# Wait for network connectivity with timeout
COUNT=0
while ! ping -c1 8.8.8.8 &>/dev/null; do
    COUNT=$((COUNT + 1))
    if [ $COUNT -gt 60 ]; then
        echo "[$(date)] ERROR: Network timeout after 5 minutes" >> "$LOG_FILE"
        exit 1
    fi
    sleep 5
done
echo "[$(date)] Network connectivity established" >> "$LOG_FILE"

# Additional delay to ensure system is fully ready
sleep 10

echo "[$(date)] Starting ASB service..." >> "$LOG_FILE"
# Start ASB with proper process management
nohup ./asb > asb.log 2>&1 &
ASB_PID=$!
echo "[$(date)] ASB started with PID: $ASB_PID" >> "$LOG_FILE"
EOF

chmod +x "$BOOT_SCRIPT"
print_success "Enhanced auto-start service configured"

# Final summary
print_header
echo -e "${GREEN}🎉 INSTALLATION COMPLETE!${NC}"
echo
echo -e "${CYAN}Summary:${NC}"
echo -e "  • Binary: $(pwd)/asb"
echo -e "  • Config: $(pwd)/config.json"
echo -e "  • Storage: ~/asb_files -> /storage/emulated/0/Download/asb_files/"
echo -e "  • Main Log: $(pwd)/asb.log"
echo -e "  • Boot Log: $(pwd)/asb-boot.log"
echo -e "  • Auto-start: $BOOT_SCRIPT"
echo
echo -e "${YELLOW}Next steps:${NC}"
echo -e "  1. Grant Termux:Boot permission in Android settings"
echo -e "  2. Restart your device to test auto-start"
echo -e "  3. Check boot log: cat asb-boot.log"
echo -e "  4. Or run manually: ./asb"
echo
echo -e "${GREEN}💡 Pro tip: The enhanced auto-start includes network detection${NC}"
echo -e "${GREEN}   and detailed logging for troubleshooting${NC}"
echo
echo -e "${GREEN}Need help? Check README.md for detailed usage instructions${NC}"
echo
echo -e "${YELLOW}🧹 To uninstall ASB later:${NC}"
echo -e "${YELLOW}  Run: ./uninstall.sh${NC}"
echo -e "${YELLOW}  This will safely remove all ASB components${NC}"
