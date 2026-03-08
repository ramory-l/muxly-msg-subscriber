#!/bin/bash
# sync-schemas.sh — manage the muxly-schemas git submodule
#
# Usage:
#   ./scripts/sync-schemas.sh pull      # pull latest schemas from remote into this repo
#   ./scripts/sync-schemas.sh push      # commit & push local schema changes, then update parent pointer
#   ./scripts/sync-schemas.sh status    # show submodule status

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
SUBMODULE_PATH="schemas"
SUBMODULE_DIR="$ROOT_DIR/$SUBMODULE_PATH"

# ── helpers ───────────────────────────────────────────────────────────────────

die() { echo "error: $*" >&2; exit 1; }

require_clean_submodule() {
    if ! git -C "$SUBMODULE_DIR" diff --quiet HEAD; then
        die "submodule has unstaged changes — stage or stash them first"
    fi
}

# ── commands ──────────────────────────────────────────────────────────────────

cmd_pull() {
    echo "==> Pulling latest schemas from remote..."
    git -C "$ROOT_DIR" submodule update --remote --merge "$SUBMODULE_PATH"
    echo "==> Updating parent repo submodule pointer..."
    git -C "$ROOT_DIR" add "$SUBMODULE_PATH"
    if git -C "$ROOT_DIR" diff --cached --quiet; then
        echo "    Parent pointer already up to date — nothing to commit."
    else
        git -C "$ROOT_DIR" commit -m "chore: update schemas submodule"
        echo "    Parent pointer updated. Push the parent repo when ready."
    fi
    echo "Done."
}

cmd_push() {
    local msg="${1:-}"
    if [ -z "$msg" ]; then
        die "provide a commit message: sync-schemas.sh push \"your message\""
    fi

    echo "==> Committing changes inside submodule..."
    git -C "$SUBMODULE_DIR" add -A
    if git -C "$SUBMODULE_DIR" diff --cached --quiet; then
        echo "    Nothing to commit in submodule."
    else
        git -C "$SUBMODULE_DIR" commit -m "$msg"
    fi

    echo "==> Pushing submodule to remote..."
    git -C "$SUBMODULE_DIR" push origin HEAD

    echo "==> Updating parent repo submodule pointer..."
    git -C "$ROOT_DIR" add "$SUBMODULE_PATH"
    if git -C "$ROOT_DIR" diff --cached --quiet; then
        echo "    Parent pointer already up to date — nothing to commit."
    else
        git -C "$ROOT_DIR" commit -m "chore: update schemas submodule"
        echo "    Parent pointer updated. Push the parent repo when ready."
    fi
    echo "Done."
}

cmd_status() {
    echo "==> Submodule git status:"
    git -C "$SUBMODULE_DIR" status
    echo ""
    echo "==> Parent repo sees submodule as:"
    git -C "$ROOT_DIR" submodule status "$SUBMODULE_PATH"
}

# ── dispatch ──────────────────────────────────────────────────────────────────

case "${1:-}" in
    pull)   cmd_pull ;;
    push)   cmd_push "${2:-}" ;;
    status) cmd_status ;;
    *)
        echo "Usage: $(basename "$0") <command> [args]"
        echo ""
        echo "Commands:"
        echo "  pull               Pull latest schemas from remote and update parent pointer"
        echo "  push \"message\"   Commit + push schema changes, update parent pointer"
        echo "  status             Show submodule and parent repo status"
        exit 1
        ;;
esac
