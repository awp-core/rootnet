#!/usr/bin/env bash
# AWP — Remote API deployment (cross-compile + scp + restart)
#
# Builds static binaries locally, uploads to remote server, and restarts services.
# Remote server only needs the binaries + .env — no Go toolchain required.
#
# Usage:
#   ./scripts/deploy-remote.sh                    # Full deploy: build + upload + restart (~40s)
#   ./scripts/deploy-remote.sh --restart-only     # Restart without rebuild (uses last uploaded binaries)
#   ./scripts/deploy-remote.sh --status           # Check remote service status
#   ./scripts/deploy-remote.sh --logs [service]   # Tail remote logs (api/indexer/keeper)
#   ./scripts/deploy-remote.sh --stop             # Stop remote services
#
# Remote server uses systemd services (awp-api, awp-indexer, awp-keeper).
# EnvironmentFile=/home/ubuntu/awp/.env must include KEEPER_SKIP_SETTLE=true
# Manual control:
#   sudo systemctl stop awp-api awp-indexer awp-keeper
#   sudo systemctl start awp-api awp-indexer awp-keeper
#   sudo systemctl restart awp-api awp-indexer awp-keeper
#   journalctl -u awp-api -f
#
# Config (env vars or defaults):
#   REMOTE_HOST     SSH host (default: ubuntu@api.awp.sh)
#   REMOTE_DIR      Remote install directory (default: /home/ubuntu/awp)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
API_DIR="$ROOT_DIR/api"
TMP_DIR="/tmp/awp-deploy-$$"

REMOTE_HOST="${REMOTE_HOST:-ubuntu@api.awp.sh}"
REMOTE_DIR="${REMOTE_DIR:-/home/ubuntu/awp}"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

info()  { echo -e "${GREEN}[INFO]${NC} $*"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $*"; }
err()   { echo -e "${RED}[ERROR]${NC} $*" >&2; exit 1; }
step()  { echo -e "${CYAN}[STEP]${NC} $*"; }

ssh_cmd() { ssh -o ConnectTimeout=10 -o StrictHostKeyChecking=accept-new "$REMOTE_HOST" "$@"; }

# ─── Build ───

do_build() {
    step "Building static binaries (CGO_ENABLED=0)..."
    mkdir -p "$TMP_DIR"
    cd "$API_DIR"
    local version
    version=$(git rev-parse --short HEAD 2>/dev/null || echo "dev")
    local ldflags="-X github.com/cortexia/rootnet/api/internal/server/handler.Version=$version"
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$ldflags" -o "$TMP_DIR/awp-api"     ./cmd/api     &
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$ldflags" -o "$TMP_DIR/awp-indexer" ./cmd/indexer  &
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$ldflags" -o "$TMP_DIR/awp-keeper"  ./cmd/keeper   &
    wait
    info "Built: awp-api, awp-indexer, awp-keeper"
    ls -lh "$TMP_DIR"/awp-*
}

# ─── Upload ───

do_upload() {
    step "Uploading binaries to $REMOTE_HOST:$REMOTE_DIR..."
    scp -q "$TMP_DIR/awp-api" "$TMP_DIR/awp-indexer" "$TMP_DIR/awp-keeper" "$REMOTE_HOST:/tmp/"
    info "Upload complete"
}

# ─── Stop remote ───

do_stop() {
    step "Stopping remote services..."
    ssh_cmd "sudo systemctl stop awp-api awp-indexer awp-keeper 2>/dev/null; echo 'All services stopped'"
}

# ─── Install + Start ───

do_install_and_start() {
    step "Installing binaries and restarting..."

    # Install binaries
    ssh_cmd "rm -f $REMOTE_DIR/awp-api $REMOTE_DIR/awp-indexer $REMOTE_DIR/awp-keeper && \
             cp /tmp/awp-api /tmp/awp-indexer /tmp/awp-keeper $REMOTE_DIR/ && \
             chmod +x $REMOTE_DIR/awp-api $REMOTE_DIR/awp-indexer $REMOTE_DIR/awp-keeper"

    # Install start.sh on remote (start-stop-daemon for clean double-fork)
    # Start services via systemd (handles daemonization, logging, restart)
    ssh_cmd "sudo systemctl start awp-api awp-indexer awp-keeper"
    sleep 2
    ssh_cmd "for svc in awp-api awp-indexer awp-keeper; do printf '%s=%s\n' \$svc \$(systemctl show -p MainPID --value \$svc); done"

    # Health check with retry
    step "Waiting for health check..."
    for i in $(seq 1 10); do
        sleep 1
        if ssh_cmd "curl -sf http://localhost:8080/api/health" >/dev/null 2>&1; then
            info "Health check passed"
            do_status
            return 0
        fi
    done
    warn "Health check failed after 10s — check logs"
    do_status
    return 1
}

# ─── Status ───

do_status() {
    echo ""
    echo "  Remote Service Status ($REMOTE_HOST):"
    echo "  ──────────────────────────────────────"
    ssh_cmd bash -s << 'STATUS_SCRIPT'
        for svc in awp-api awp-indexer awp-keeper; do
            status=$(systemctl is-active $svc 2>/dev/null || echo "unknown")
            pid=$(systemctl show -p MainPID --value $svc 2>/dev/null || echo "?")
            echo "  $svc:    $status (PID $pid)"
        done
        if curl -sf http://localhost:8080/api/health >/dev/null 2>&1; then
            echo '  health:  OK'
        else
            echo '  health:  unreachable'
        fi
STATUS_SCRIPT
    echo ""
}

# ─── Logs ───

do_logs() {
    local svc="${1:-api}"
    ssh_cmd "sudo journalctl -u awp-${svc} -f -n 50"
}

# ─── Cleanup ───

cleanup() { rm -rf "$TMP_DIR" 2>/dev/null || true; }
trap cleanup EXIT

# ─── Main ───

case "${1:-}" in
    --restart-only)
        do_stop
        do_install_and_start
        ;;
    --stop)
        do_stop
        ;;
    --status)
        do_status
        ;;
    --logs)
        do_logs "${2:-api}"
        ;;
    --help|-h)
        head -14 "$0" | grep '^#' | sed 's/^# \?//'
        ;;
    "")
        do_build
        do_upload
        do_stop
        do_install_and_start
        cleanup
        info "Deployment complete"
        ;;
    *)
        err "Unknown flag: $1. Use --help for usage."
        ;;
esac
