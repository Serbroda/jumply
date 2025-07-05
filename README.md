# Jumply – Just a Media Player

**Jumply** is a minimalistic, self-hosted video player built with Go and HTMX. It recursively scans video files from defined directories, displays them in a simple web UI, and allows users to stream videos.

## 💡 Motivation

The idea behind Jumply is simple: I wanted to automatically list and access newly added media files across multiple root directories through a clean and minimal web interface.

Whenever a new video file is added—whether it's downloaded or copied—it should instantly become available to watch in the browser without any extra steps, configurations, or overhead. No frills. Just drop it in, and hit play.

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
- **Styling:** Default CSS

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

Jumply uses a `.env` file to configure runtime behavior.
Create a `.env` file in the project root with the following content:

```env
SERVER_PORT=8080
ROOT_DIRS=./testdata/root1;./testdata/root2
DEFAULT_PAGE_SIZE=20
```

- `SERVER_PORT`: The port your server will run on
- `ROOT_DIRS`: Semicolon-separated list of root directories to scan recursively for video files
- `DEFAULT_PAGE_SIZE`: How many videos to show per page

If no `.env` file is present, Jumply will fall back to built-in defaults.

---

## 🚧 Roadmap

- [x] Basic playback and listing
- [x] Pagination and filtering
- [x] HTMX integration for fragment updates
- [ ] Custom CSS support

---

## 📝 License

MIT License. See `LICENSE` file.

---

## 💬 Credits

Created by [Serbroda](https://github.com/Serbroda) – for developers who just want to **jump in and play** 🎬
