#!/bin/bash

# Android Server Brain (ASB) Uninstallation Script
# This script removes all ASB components and configuration

# --- Color Definitions ---
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# --- Functions ---
print_header() {
    echo -e "${BLUE}====================================${NC}"
    echo -e "${BLUE}  Android Server Brain - Removal   ${NC}"
    echo -e "${BLUE}====================================${NC}"
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

confirm_removal() {
    echo -e "${YELLOW}This will remove all ASB components:${NC}"
    echo -e "${YELLOW}  • Binary executable${NC}"
    echo -e "${YELLOW}  • Configuration files${NC}"
    echo -e "${YELLOW}  • Storage directories${NC}"
    echo -e "${YELLOW}  • Auto-start scripts${NC}"
    echo -e "${YELLOW}  • Log files${NC}"
    echo
    read -p "Are you sure you want to proceed? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${CYAN}Removal cancelled.${NC}"
        exit 0
    fi
}

print_header
echo -e "${RED}⚠️  WARNING: This will completely remove ASB from your system${NC}"
echo

# Confirmation
confirm_removal

# Read configuration to get storage directory
STORAGE_DIR="downloads/server"
if [ -f "config.json" ]; then
    # Use jq if available, otherwise fall back to grep
    if command -v jq &> /dev/null; then
        STORAGE_DIR=$(jq -r '.storage_dir // "downloads/server"' config.json)
    else
        STORAGE_DIR=$(grep '"storage_dir"' config.json | cut -d '"' -f 4)
    fi
fi
FULL_STORAGE_PATH="$HOME/$STORAGE_DIR"

print_step "1" "Stopping ASB service"
# Kill any running ASB processes
if pgrep -f "asb" > /dev/null; then
    pkill -f "asb"
    print_success "ASB processes stopped"
else
    print_warning "No ASB processes found"
fi

print_step "2" "Removing auto-start configuration"
# Remove Termux:Boot script
BOOT_SCRIPT="$HOME/.termux/boot/asb"
if [ -f "$BOOT_SCRIPT" ]; then
    rm "$BOOT_SCRIPT"
    print_success "Removed auto-start script: $BOOT_SCRIPT"
else
    print_warning "No auto-start script found"
fi

print_step "3" "Cleaning up files and directories"
# Remove binary
if [ -f "asb" ]; then
    rm "asb"
    print_success "Removed ASB binary"
else
    print_warning "No ASB binary found"
fi

# Remove log files
if [ -f "asb.log" ]; then
    rm "asb.log"
    print_success "Removed log files"
else
    print_warning "No log files found"
fi

# Remove configuration
if [ -f "config.json" ]; then
    rm "config.json"
    print_success "Removed configuration file"
else
    print_warning "No configuration file found"
fi

# Remove Go module files
if [ -f "go.mod" ]; then
    rm "go.mod"
    print_success "Removed go.mod"
fi

if [ -f "go.sum" ]; then
    rm "go.sum"
    print_success "Removed go.sum"
fi

print_step "4" "Cleaning up storage"
# Remove symlink
if [ -L "$HOME/server" ]; then
    rm "$HOME/server"
    print_success "Removed symlink: ~/server"
elif [ -e "$HOME/server" ]; then
    print_warning "~/server exists but is not a symlink"
fi

# Ask about storage directory removal
echo
echo -e "${YELLOW}Storage directory: $FULL_STORAGE_PATH${NC}"
read -p "Remove uploaded files and scripts? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if [ -d "$FULL_STORAGE_PATH" ]; then
        rm -rf "$FULL_STORAGE_PATH"
        print_success "Removed storage directory"
    else
        print_warning "Storage directory not found"
    fi
else
    print_warning "Storage directory preserved"
fi

print_step "5" "Final cleanup"
# Remove any temporary files
find . -name "*.tmp" -delete 2>/dev/null
find . -name "*~" -delete 2>/dev/null

print_header
echo -e "${GREEN}🎉 ASB Removal Complete!${NC}"
echo
echo -e "${CYAN}What was removed:${NC}"
echo -e "  • ASB binary and executables"
echo -e "  • Configuration files"
echo -e "  • Auto-start scripts"
echo -e "  • Log files"
echo -e "  • Symbolic links"
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "  • Uploaded files and scripts"
fi
echo
echo -e "${YELLOW}Manual cleanup (if needed):${NC}"
echo -e "  • Check ~/.termux/boot/ for any remaining scripts"
echo -e "  • Verify no ASB processes are still running: pgrep asb"
echo
echo -e "${GREEN}ASB has been successfully uninstalled!${NC}"