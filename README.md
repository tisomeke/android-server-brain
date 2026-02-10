# 🧠 Android Server Brain (ASB)

**Turn your Android smartphone into a powerful, Telegram-managed autonomous server.**

> [!WARNING]
> **Project Status: ARCHIVED - Pending Docker Implementation**
> 
> This project is currently archived due to environment-specific issues that require a different approach:
> - Go 1.25.7 toolchain problems on Android/arm64 (use Go 1.21.x instead)
> - DNS resolution issues in Termux environment
> - Binary execution problems on some Android devices
> - Termux:Boot detection sensitivity
> 
> **Recommended Solution:** Docker containerization is needed for stable cross-device compatibility. This would provide:
> - Consistent runtime environment across all Android devices
> - Isolated dependencies and networking
> - Easier deployment and updates
> - Better reliability and maintainability
> 
> The project will be revived once Docker implementation is complete.

> [!TIP]
> [🇷🇺 *Читать на русском*](README.ru.md)

### 📌 Project Overview

ASB is a lightweight Go-based framework designed to run inside **Termux**. It transforms a spare Android device into a remote-controlled server node that can be managed entirely via a Telegram Bot.

### 🚀 Key Features

* **System Monitoring:** Real-time stats (Battery, CPU, Storage, Uptime) via `termux-api`.
* **Remote Shell:** Execute Bash commands directly from Telegram with timeout protection.
* **Smart Storage:** Upload scripts/files via Telegram; they are saved to `~/downloads/server` and automatically linked to `~/server` with `+x` permissions.
* **Mesh Networking:** Integrated Tailscale support for secure remote access without public IPs.
* **Admin Security:** Strict ID-based white-listing.

### 🛠 Prerequisites & Dependencies

Before installing ASB, ensure you have the following:

**Required Apps (from F-Droid):**
* **[Termux](https://f-droid.org/packages/com.termux/)** - Terminal emulator and Linux environment
* **[Termux:Boot](https://f-droid.org/packages/com.termux.boot/)** - Auto-start services on boot
* **[Termux:API](https://f-droid.org/packages/com.termux.api/)** - System APIs for battery, sensors, etc.

**Additional Requirements:**
* Android device (API 21+)
* Stable internet connection
* Telegram account
* Bot token from [@BotFather](https://t.me/BotFather)

### 📦 Installation

**Option 1: Automated Installation (Recommended)**

Simply run the installation script:
```bash
git clone https://github.com/tisomeke/android-server-brain
cd android-server-brain
chmod +x install.sh
./install.sh
```

The script will automatically:
- **Verify system dependencies** (pkg, curl, ping)
- **Check Termux:Boot status** (requires first-time setup)
  - Uses enhanced detection with multiple verification methods
  - Can be bypassed with `SKIP_BOOT_CHECK=1` for advanced users
- Check and request storage permissions
- Install required dependencies (golang, git, termux-api)
- Configure your bot settings interactively
- Set up storage directories and symlinks
- Build the ASB binary
- Configure **enhanced auto-start service** with network detection

**Option 2: Manual Installation**

1. **Install Termux packages:**
   ```bash
   pkg update && pkg upgrade
   pkg install golang git termux-api
   ```

2. **Clone and build ASB:**
   ```bash
   git clone https://github.com/tisomeke/android-server-brain.git
   cd android-server-brain
   go build -o asb .
   ```

3. **Configure the application:**
   ```bash
   cp config.json.example config.json
   # Edit config.json with your bot token and admin ID
   nano config.json
   ```

4. **Setup auto-start with Termux:Boot:**
   ```bash
   mkdir -p ~/.termux/boot
   echo '#!/data/data/com.termux/files/usr/bin/sh' > ~/.termux/boot/asb
   echo 'cd ~/android-server-brain && ./asb' >> ~/.termux/boot/asb
   chmod +x ~/.termux/boot/asb
   ```

### ▶️ Deployment

**Manual Start:**
```bash
./asb
```

**Auto-start:**
After setting up Termux:Boot, ASB will start automatically on device boot with enhanced reliability:
- Network connectivity detection
- Detailed boot logging (`asb-boot.log`)
- Process monitoring and timeout protection
- Automatic retry mechanisms

**Background Service:**
```bash
nohup ./asb > asb.log 2>&1 &
```

### 🗑 Uninstallation

To completely remove ASB from your system:

```bash
./uninstall.sh
```

The uninstall script will:
- Stop any running ASB processes
- Remove the binary and configuration files
- Clean up auto-start scripts
- Remove storage directories (optional)
- Provide confirmation before destructive actions

### 📖 Usage

#### Telegram Commands:

**Basic Commands:**
* `/start` - Welcome message and basic info
* `/status` - View system health (battery, storage, uptime)
* `/battery` - Check detailed battery status (charge %, temperature, charging status)
* `/watchdog` - View watchdog monitoring status and configuration

**System Management:**
* `/reboot` - Reboot the Android device (requires confirmation)
* `/restart <service>` - Restart system services
  - Usage: `/restart ssh` or `/restart nginx`
  - Use `/restart` without arguments to see available services
* `/update` - Check for and install ASB updates
  - Usage: `/update` to check for updates, `/update now` to install

**Remote Execution:**
* `/exec <command>` - Execute shell commands remotely
  - Example: `/exec ps aux` or `/exec df -h`
  - Commands run with Termux user privileges
  - Includes timeout protection

**File Management:**
* **Upload files** - Simply send any file to the bot
  - Files are automatically saved to `~/downloads/server`
  - Symlinked to `~/server` with executable permissions (`chmod +x`)
  - Access via `~/server/filename` or direct path
  - Can be executed directly after upload

#### File Management:

* **Upload files** - Send any file to the bot
* Files are saved to `~/downloads/server`
* Automatically symlinked to `~/server` with executable permissions
* Access via `~/server/filename` or direct path

#### Example Workflows:

```bash
# Check system status
/status

# Run system commands
/exec ps aux
/exec df -h
/exec top -n 1

# Upload and run scripts
# 1. Send script file to bot
# 2. Script automatically becomes executable
# 3. Run: /exec ~/server/myscript.sh

# Check for and install updates
/update
/update now
```

### 🎮 Advanced Usage Examples

#### 1. Minecraft Server Deployment

**Setup Process:**

1. **Upload Minecraft server JAR:**
   - Download desired Minecraft server version (e.g., `paper-1.20.4.jar`)
   - Send the JAR file to your ASB bot
   - File will be saved as `~/server/paper-1.20.4.jar`

2. **Initial Configuration:**
   ```bash
   # Accept Minecraft EULA
   /exec echo "eula=true" > ~/server/eula.txt
   
   # Create basic server properties
   /exec echo 'server-port=25565\ngamemode=survival\ndifficulty=normal' > ~/server/server.properties
   ```

3. **Start the Server:**
   ```bash
   # Allocate 2GB RAM to server
   /exec java -Xmx2G -Xms1G -jar ~/server/paper-1.20.4.jar nogui
   ```

4. **Server Management:**
   ```bash
   # Check server status
   /exec ps aux | grep java
   
   # View server logs
   /exec tail -f ~/server/logs/latest.log
   
   # Stop server gracefully
   /exec pkill -f "java.*paper"
   ```

5. **Automated Startup Script:**
   Create `~/server/start-mc.sh` and upload it:
   ```bash
   #!/data/data/com.termux/files/usr/bin/bash
   cd ~/server
   java -Xmx2G -Xms1G -jar paper-1.20.4.jar nogui
   ```
   Then run: `/exec ~/server/start-mc.sh`

#### 2. Python Bot Hosting

**Deployment Workflow:**

1. **Upload Bot Files:**
   - Send your Python bot script (e.g., `mybot.py`)
   - Send `requirements.txt` for dependencies
   - Send `config.json` for configuration

2. **Environment Setup:**
   ```bash
   # Install Python dependencies
   /exec pip install -r ~/server/requirements.txt
   
   # Set up virtual environment (optional)
   /exec python -m venv ~/server/venv
   /exec ~/server/venv/bin/pip install -r ~/server/requirements.txt
   ```

3. **Configuration:**
   ```bash
   # Set up bot configuration
   /exec cat ~/server/config.json
   
   # Test bot connectivity
   /exec python ~/server/mybot.py --test
   ```

4. **Running the Bot:**
   ```bash
   # Direct execution
   /exec python ~/server/mybot.py
   
   # Background execution with logging
   /exec nohup python ~/server/mybot.py > ~/server/bot.log 2>&1 &
   
   # Using virtual environment
   /exec nohup ~/server/venv/bin/python ~/server/mybot.py > ~/server/bot.log 2>&1 &
   ```

5. **Bot Lifecycle Management:**
   ```bash
   # Check if bot is running
   /exec ps aux | grep mybot.py
   
   # View bot logs
   /exec tail -f ~/server/bot.log
   
   # Restart bot
   /exec pkill -f mybot.py
   /exec nohup python ~/server/mybot.py > ~/server/bot.log 2>&1 &
   
   # Update bot code
   # 1. Send updated files to bot
   # 2. Restart bot process
   ```

6. **Automated Restart Setup:**
   Create a restart script `~/server/restart-bot.sh`:
   ```bash
   #!/data/data/com.termux/files/usr/bin/bash
   pkill -f mybot.py
   sleep 2
   nohup ~/server/venv/bin/python ~/server/mybot.py > ~/server/bot.log 2>&1 &
   echo "Bot restarted at $(date)"
   ```
   Make it executable: `/exec chmod +x ~/server/restart-bot.sh`
   Use it: `/exec ~/server/restart-bot.sh`

#### 3. System Administration Tasks

**Storage Management:**
```bash
# Check disk usage
/exec df -h

# Clean up old logs
/exec find ~/server/logs -name "*.log" -mtime +7 -delete

# Backup important files
/exec tar -czf ~/server-backup-$(date +%Y%m%d).tar.gz ~/server/
```

**Process Monitoring:**
```bash
# Monitor resource usage
/exec top -n 1

# Check specific processes
/exec pgrep -f "java\|python"

# Kill hanging processes
/exec pkill -f "process_name"
```

**Network Operations:**
```bash
# Check network connectivity
/exec ping -c 4 google.com

# View active connections
/exec netstat -tuln

# Port monitoring
/exec ss -tuln | grep :25565  # Minecraft port
```

### 🏗 Project Structure

```
android-server-brain/
├── config/
│   └── config.go          # Configuration loading and validation
├── internal/
│   ├── bot/
│   │   └── router.go      # Telegram bot command handlers
│   ├── storage/
│   │   └── files.go       # File upload and management
│   └── system/
│       ├── monitor.go     # System status monitoring
│       ├── shell.go       # Command execution
│       └── watchdog.go    # Battery monitoring service
├── main.go                # Application entry point
├── config.json           # Configuration file
├── install.sh            # Automated installation script
└── go.mod               # Go module dependencies
```

### ⚙️ Configuration

Create `config.json` with the following structure:

```json
{
  "telegram_token": "YOUR_BOT_TOKEN_HERE",
  "admin_id": 123456789,
  "storage_dir": "downloads/server"
}
```

* `telegram_token`: Get from [@BotFather](https://t.me/BotFather)
* `admin_id`: Your Telegram user ID (use [@userinfobot](https://t.me/userinfobot) to find it)
* `storage_dir`: Directory for uploaded files (relative to home)

### 🔒 Security Notes

* Only the configured AdminID can control the server
* All commands execute with Termux user privileges
* File uploads are sanitized and stored in isolated directory
* Network access depends on your Telegram security settings

### 🐛 Troubleshooting

**Common Issues:**

1. **Bot not responding:**
   - Check internet connection
   - Verify bot token in config.json
   - Ensure correct AdminID

2. **Commands failing:**
   - Check Termux API permissions
   - Verify required packages are installed
   - Review logs in `asb.log`

3. **Auto-start not working:**
   - Grant Termux:Boot permission in Android settings
   - Check `~/.termux/boot/asb` script permissions
   - Test manual execution of boot script

**Logs:**
```bash
# Main application log
tail -f asb.log

# Boot sequence log (for auto-start issues)
cat asb-boot.log
```

### 🤝 Contributing

Feel free to submit issues, feature requests, or pull requests. For major changes, please open an issue first to discuss the proposed changes.

### 📄 License

This project is licensed under the GNU General Public License v3.0 - see the LICENSE file for details.
