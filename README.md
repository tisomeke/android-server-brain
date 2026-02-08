# 🧠 Android Server Brain (ASB)

**Turn your Android smartphone into a powerful, Telegram-managed autonomous server.**

> [!TIP]
> 📚 **Documentation has been moved to [docs/](docs/) directory**
> 
> 📖 [English Documentation](docs/README.md) | [🇷🇺 Русская документация](docs/README.ru.md)

## Quick Links

- 📋 [Full Documentation](docs/README.md)
- 📄 [License Information](docs/LICENSE)
- 🛠 [Installation Guide](docs/README.md#-installation)

## Project Overview

ASB is a lightweight Go-based framework designed to run inside **Termux**. It transforms a spare Android device into a remote-controlled server node that can be managed entirely via a Telegram Bot.

### 🚀 Key Features

* **System Monitoring:** Real-time stats (Battery, CPU, Storage, Uptime) via `termux-api`.
* **Remote Shell:** Execute Bash commands directly from Telegram with timeout protection.
* **Smart Storage:** Upload scripts/files via Telegram; they are saved to `~/downloads/server` and automatically linked to `~/server` with `+x` permissions.
* **Mesh Networking:** Integrated Tailscale support for secure remote access without public IPs.
* **Admin Security:** Strict ID-based white-listing.
* **Battery Watchdog:** Intelligent battery monitoring with smart notifications.

---

*For complete documentation, please see [docs/README.md](docs/README.md)*