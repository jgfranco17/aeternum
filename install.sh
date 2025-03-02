#!/usr/bin/env bash

REPO_OWNER="jgfranco17"
PROJECT_NAME="aeternum"
DEFAULT_VERSION="latest"
INSTALL_PATH="${HOME}/.local/bin"

get_latest_version() {
  curl --silent "https://api.github.com/repos/${REPO_OWNER}/${PROJECT_NAME}/releases/${DEFAULT_VERSION}" | \
    grep '"tag_name":' | \
    sed -E 's/.*"([^"]+)".*/\1/'
}

download_binary() {
  local version=$1
  local os=$2
  local arch=$3

  url="https://github.com/${REPO_OWNER}/${PROJECT_NAME}/releases/download/${version}/aeternum-${version}-${os}-${arch}.tar.gz"
  echo "Downloading Aeternum from $url"

  curl -L "$url" -o aeternum.tar.gz || {
    echo "Error: Download failed. Please check the version and try again."
    exit 1
  }
}

install_binary() {
  sudo tar -xzf aeternum.tar.gz -C "${INSTALL_PATH}" aeternum || {
    echo "Error: Installation failed."
    exit 1
  }
  rm aeternum.tar.gz
}

# =============== MAIN SCRIPT ===============

version="${1:-$DEFAULT_VERSION}"

# Detect OS and architecture
case "$(uname -s)" in
  Linux*) os="linux" ;;
  Darwin*) os="darwin" ;;
  *) echo "Error: Unsupported OS"; exit 1 ;;
esac

arch="$(uname -m)"
case "$arch" in
  x86_64) arch="amd64" ;;
  aarch64) arch="arm64" ;;
  *) echo "Error: Unsupported architecture"; exit 1 ;;
esac

# Resolve latest version if needed
if [ "$version" = "latest" ]; then
  version=$(get_latest_version)
  if [ -z "$version" ]; then
    echo "Error: Unable to fetch the latest version."
    exit 1
  fi
fi

echo "Installing Aeternum version $version for $os/$arch"

# Download and install
download_binary "$version" "$os" "$arch"
install_binary

echo "Aeternum installation complete!"
echo "Installed at: ${INSTALL_PATH}"
echo "You can now run 'aeternum --version' to verify the installation."
