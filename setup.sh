#!/usr/bin/env bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WORK_DIR="$SCRIPT_DIR/.safe-install"
SRC_DIR="$WORK_DIR/src"
BIN="$WORK_DIR/safe-install"

# ================================================================
# SELF-EXTRACTION
# ================================================================
_extract_payload() {
    echo "📦 Extracting embedded payload..."

    mkdir -p "$WORK_DIR"

    local PAYLOAD_START
    PAYLOAD_START=$(grep -n '^# __PAYLOAD_START__$' "$0" | cut -d: -f1)

    if [[ -z "$PAYLOAD_START" ]]; then
        echo "❌ Payload marker not found"
        exit 1
    fi

    tail -n +"$((PAYLOAD_START + 1))" "$0" \
        | base64 -d \
        | tar xzf - -C "$WORK_DIR" || {
            echo "❌ Extraction failed"
            exit 1
        }

    echo "✅ Payload extracted"
}

# ================================================================
# BUILD
# ================================================================
_build_binary() {
    echo "⚙️ Building Go binary..."

    cd "$SRC_DIR"

    if ! command -v go >/dev/null; then
        echo "📥 Installing Go..."
        sudo apt update
        sudo apt install -y golang-go
    fi

    go build -o "$BIN"

    echo "✅ Build complete"
}

# ================================================================
# DEPENDENCIES
# ================================================================
_check_deps() {
    echo "🔍 Checking dependencies..."

    command -v podman >/dev/null || {
        echo "📥 Installing Podman..."
        sudo apt install -y podman
    }

    command -v fpm >/dev/null || {
        echo "📥 Installing fpm..."
        sudo apt install -y ruby ruby-dev build-essential
        sudo gem install fpm
    }
}

# ================================================================
# FIRST RUN SETUP
# ================================================================
if [[ ! -f "$BIN" ]]; then
    _extract_payload
    _check_deps
    _build_binary
fi

# ================================================================
# EXECUTE
# ================================================================
exec "$BIN" "$@"

exit 0

# __PAYLOAD_START__
