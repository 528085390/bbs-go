#!/bin/bash
set -euo pipefail

# ============================================================
# bbs-go K8s 一键部署脚本 (Kustomize + k3s)
# ============================================================

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info()  { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"

cd "$PROJECT_DIR"

# ---------- Step 1: Build Docker images ----------
log_info "Building Docker images..."
declare -A IMAGES=(
  ["auth/Dockerfile"]=temp-auth
  ["user/Dockerfile"]=temp-user
  ["section/rpc/Dockerfile"]=temp-section
  ["post/rpc/Dockerfile"]=temp-post
  ["comment/rpc/Dockerfile"]=temp-comment
  ["interaction/rpc/Dockerfile"]=temp-interaction
  ["search/rpc/Dockerfile"]=temp-search
  ["file/rpc/Dockerfile"]=temp-file
  ["gateway/Dockerfile"]=temp-gateway
)

for df in "${!IMAGES[@]}"; do
  tag="${IMAGES[$df]}"
  log_info "  Building $tag..."
  docker build -f "$df" -t "$tag:latest" .
done

# ---------- Step 2: Import images into k3s ----------
log_info "Importing images into k3s containerd..."
for tag in "${IMAGES[@]}"; do
  log_info "  Importing $tag..."
  docker save "$tag:latest" | k3s ctr images import -
done

log_info "Cleaning up Docker build cache..."
docker image prune -f

# ---------- Step 3: Deploy with Kustomize ----------
log_info "Deploying to Kubernetes..."
kubectl apply -k "$SCRIPT_DIR/overlays/dev"

# ---------- Step 4: Wait for all pods ----------
log_info "Waiting for all pods to be ready..."
kubectl wait --for=condition=ready pod -l project=bbs-go -n temp --timeout=300s --all

# ---------- Step 5: Show status ----------
echo ""
log_info "=== Deployment Summary ==="
kubectl get pods -n temp -o wide
echo ""
kubectl get svc -n temp
echo ""
kubectl get ingress -n temp

echo ""
log_info "=== Done! ==="
log_info "Gateway should be accessible via the Ingress at your ECS IP."
log_info "Run 'kubectl get ingress -n temp' to check the status."
