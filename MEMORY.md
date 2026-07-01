# Project Memory

- Release workflow: GitHub Actions workflow lives in `.github/workflows/release.yml`.
- Release tags use the `v1.0526.<number>` pattern and publish version-suffixed artifacts like `went2-windows-amd64-v.1.0526.1.exe`.
- Main branch pushes automatically create the next `v1.0526.<number>` tag and publish release assets in the same workflow run; tag pushes also publish release artifacts.
- Installer scripts live at `install.sh` and `install.ps1`; they resolve the latest tag from the GitHub API by default and append the install directory to PATH persistently.
- Repository convention: keep changes minimal and aligned with the existing Go/Cobra scaffold style.
- `install.sh` must parse GitHub API JSON with whitespace-tolerant tag extraction; the releases API is not guaranteed to emit minified one-line JSON.
