# 🧠 Android Server Brain (ASB)

**Turn your Android smartphone into a powerful, Telegram-managed autonomous server.**

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
git clone https://github.com/yourusername/android-server-brain.git
cd android-server-brain
chmod +x install.sh
./install.sh
```

The script will automatically:
- Check and request storage permissions
- Install required dependencies (golang, git, termux-api)
- Configure your bot settings interactively
- Set up storage directories and symlinks
- Build the ASB binary
- Configure auto-start service

**Option 2: Manual Installation**

1. **Install Termux packages:**
   ```bash
   pkg update && pkg upgrade
   pkg install golang git termux-api
   ```

2. **Clone and build ASB:**
   ```bash
   git clone https://github.com/yourusername/android-server-brain.git
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
After setting up Termux:Boot, ASB will start automatically on device boot.

**Background Service:**
```bash
nohup ./asb > asb.log 2>&1 &
```

### 📖 Usage

#### Telegram Commands:

* `/start` - Welcome message and basic info
* `/status` - View system health (battery, storage, uptime)
* `/battery` - Check detailed battery status
* `/watchdog` - View watchdog monitoring status
* `/exec <command>` - Execute shell commands remotely

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
tail -f asb.log
```

### 🤝 Contributing

Feel free to submit issues, feature requests, or pull requests. For major changes, please open an issue first to discuss the proposed changes.

### 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
