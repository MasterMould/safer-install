#!/usr/bin/env bash

set -euo pipefail

APP_NAME="safe-install"
REPO_URL="https://github.com/MasterMould/safer-install.git"
BRANCH="dev"

BASE_DIR="$HOME/.safe-install"
SRC_DIR="$BASE_DIR/src"
BIN="$BASE_DIR/$APP_NAME"

# ================================================================
# Clone or update repo
# ================================================================
_sync_repo() {
    if [[ ! -d "$SRC_DIR/.git" ]]; then
        echo "📥 Cloning repository..."
        git clone --branch "$BRANCH" "$REPO_URL" "$SRC_DIR"
    else
        echo "🔄 Updating repository..."
        git -C "$SRC_DIR" fetch origin
        git -C "$SRC_DIR" reset --hard "origin/$BRANCH"
    fi
}

# ================================================================
# Install dependencies
# ================================================================
_install_deps() {
    echo "📥 Installing dependencies..."

    sudo apt update
    sudo apt install -y \
        golang-go \
        podman \
        ruby \
        ruby-dev \
        build-essential \
        git

    sudo gem install --no-document fpm
}

# ================================================================
# Check dependencies
# ================================================================
_check_deps() {
    echo "🔍 Checking dependencies..."

    local missing=0

    for cmd in go podman fpm git; do
        if ! command -v "$cmd" >/dev/null; then
            missing=1
        fi
    done

    if [[ "$missing" -eq 1 ]]; then
        _install_deps
    fi
}

# ================================================================
# Build binary
# ================================================================
_build() {
    echo "⚙️ Building Go binary..."

    cd "$SRC_DIR"

    go build -o "$BIN"

    echo "✅ Build complete at $BIN"
}

_install_symlink() {
    local TARGET="/usr/local/bin/safe-install"

    if [[ ! -f "$TARGET" ]]; then
        echo "🔗 Installing safe-install globally..."
        sudo ln -sf "$BIN" "$TARGET"
        echo "✅ Now you can run: safe-install"
    fi
}

# ================================================================
# CLI options
# ================================================================
case "${1:-}" in
    --reset)
        echo "🧹 Resetting..."
        rm -rf "$BASE_DIR"
        exit 0
        ;;
    --update)
        _sync_repo
        _build
        exit 0
        ;;
    --rebuild)
        echo "🔁 Rebuilding..."
        rm -f "$BIN"
        ;;
esac

# ================================================================
# Bootstrap
# ================================================================
mkdir -p "$BASE_DIR"

_check_deps
_sync_repo

if [[ ! -f "$BIN" ]]; then
    _build
fi

_install_symlink

# ================================================================
# Execute
# ================================================================
exec "$BIN" "$@"
