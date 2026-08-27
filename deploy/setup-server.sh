#!/usr/bin/env bash
# One-time server bootstrap for adamndegwa.com. Run as a sudo-capable user.
# After this, deploys happen automatically via the CD workflow (see README).
set -euo pipefail

IMAGE="ghcr.io/${1:?usage: setup-server.sh <github-owner>}/adamndegwa"

# 1. Packages (dnf or apt).
if command -v dnf >/dev/null; then
  sudo dnf install -y docker caddy
elif command -v apt-get >/dev/null; then
  sudo apt-get update && sudo apt-get install -y docker.io caddy
fi
sudo systemctl enable --now docker caddy

# 2. Firewall: only HTTP/HTTPS are public; the app stays on localhost:8080.
if command -v firewall-cmd >/dev/null; then
  sudo firewall-cmd --permanent --add-service=http --add-service=https
  sudo firewall-cmd --reload
elif command -v ufw >/dev/null; then
  sudo ufw allow 80/tcp
  sudo ufw allow 443/tcp
fi

# 3. Caddy config.
sudo cp "$(dirname "$0")/Caddyfile" /etc/caddy/Caddyfile
sudo systemctl reload caddy

# 4. Registry login — needs a GitHub PAT with read:packages, once.
#    echo "$GHCR_PAT" | docker login ghcr.io -u <github-owner> --password-stdin

# 5. Initial run.
docker pull "$IMAGE:latest"
docker rm -f site 2>/dev/null || true
docker run -d --restart unless-stopped --name site -p 127.0.0.1:8080:8080 "$IMAGE:latest"

echo "Done. Verify with: curl -sI localhost:8080 && curl -s https://adamndegwa.com/sitemap.xml | head -2"
