#!/usr/bin/env sh

set -eu

repo="${WENT_REPOSITORY:-went-project/went}"
version_input="${WENT_VERSION:-${1:-}}"
install_dir="${XDG_BIN_HOME:-$HOME/.local/bin}"
binary_name="went2"
temp_dir="$(mktemp -d 2>/dev/null || mktemp -d -t went-install)"

log() {
  printf '%s\n' "$*"
}

fail() {
  printf 'Error: %s\n' "$*" >&2
  exit 1
}

download_file() {
  url="$1"
  output="$2"

  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$url" -o "$output"
    return
  fi

  if command -v wget >/dev/null 2>&1; then
    wget -qO "$output" "$url"
    return
  fi

  fail "curl or wget is required"
}

fetch_latest_tag() {
  api_url="https://api.github.com/repos/$repo/releases/latest"
  response_file="$temp_dir/latest-release.json"

  download_file "$api_url" "$response_file"

  latest_tag="$(sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' "$response_file" | head -n 1)"
  [ -n "$latest_tag" ] || fail "unable to determine the latest release tag"
  printf '%s' "$latest_tag"
}

ensure_path_entry() {
  profile_file="$1"
  marker="# Added by WENT installer"

  mkdir -p "$(dirname "$profile_file")"
  touch "$profile_file"

  if ! grep -Fq "$marker" "$profile_file"; then
    {
      printf '\n%s\n' "$marker"
      printf 'export PATH="%s:$PATH"\n' "$install_dir"
    } >> "$profile_file"
  fi
}

case "$(uname -s)" in
  Darwin)
    os_name="darwin"
    ;;
  Linux)
    os_name="linux"
    ;;
  *)
    fail "unsupported operating system: $(uname -s)"
    ;;
esac

case "$(uname -m)" in
  x86_64|amd64)
    arch_name="amd64"
    ;;
  arm64|aarch64)
    arch_name="arm64"
    ;;
  *)
    fail "unsupported architecture: $(uname -m)"
    ;;
esac

if [ -n "$version_input" ]; then
  version_tag="$version_input"
else
  version_tag="$(fetch_latest_tag)"
fi

case "$version_tag" in
  v*)
    version_tag="${version_tag#v}"
    ;;
esac

asset_name="went2-${os_name}-${arch_name}-v.${version_tag}"
if [ "$os_name" = "windows" ]; then
  asset_name="${asset_name}.exe"
fi

download_url="https://github.com/$repo/releases/download/v${version_tag}/$asset_name"

trap 'rm -rf "$temp_dir"' EXIT INT TERM

mkdir -p "$install_dir"
temp_file="$temp_dir/$asset_name"

log "Downloading $asset_name"
download_file "$download_url" "$temp_file"

install_path="$install_dir/$binary_name"
mv "$temp_file" "$install_path"
chmod +x "$install_path"

ensure_path_entry "$HOME/.zshrc"
ensure_path_entry "$HOME/.bashrc"
ensure_path_entry "$HOME/.profile"

case ":$PATH:" in
  *":$install_dir:"*)
    ;;
  *)
    export PATH="$install_dir:$PATH"
    ;;
esac

log "Kurulum başarılı! Aramıza hoşgeldin! Hemen başlamak için terminalini yeniden açabilir veya aşağıdaki komutu çalıştırabilirsin:"
log "  $binary_name --help"
