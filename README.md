# Jumply – Just a Media Player

**Jumply** is a minimalistic, self-hosted video player built with Go and HTMX. It recursively scans video files from defined directories, displays them in a simple web UI, and allows users to stream videos.

---

## ✨ Features

- 🎥 Plays `.mp4`, `.avi`, and other video files from your local directory
- 🔍 Recursive folder scanning with regex-based filtering
- 📂 Clean, minimal UI with pagination
- 🚀 HTMX support for partial page updates
- 🎨 Theme support via drop-in CSS files
- 📦 Single binary per platform

---

## 🛠 Tech Stack

- **Backend:** Go + Echo framework
- **Frontend:** HTMX, Go templates
- **Styling:** Default CSS with optional

---

## 📦 Getting Started

### 1. Clone the Repo
```bash
git clone https://github.com/Serbroda/jumply.git
cd jumply
```

### 2. Build the Binary (Linux/macOS/Windows)
```bash
make all
```

This builds platform-specific binaries in `bin/`.

### 3. Run Jumply
```bash
./bin/jumply-<VERSION>-<PLATFORM>
```

Then open [http://localhost:8080](http://localhost:8080) in your browser.

### 4. Configuration
- Place your video files under `./testdata/` or modify the path in `main.go`

---
## 🚧 Roadmap

- [x] Basic playback and listing
- [x] Pagination and filtering
- [x] HTMX integration for fragment updates
- [ ] SQLite index + file watcher support

---

## 📝 License

MIT License. See `LICENSE` file.

---

## 💬 Credits

Created by [Serbroda](https://github.com/Serbroda) – for developers who just want to **jump in and play** 🎬
