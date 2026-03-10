# Installing Snap Locally from dist Folder

## Prerequisites

Ensure snapd is installed on your system:

```bash
# Ubuntu/Debian
sudo apt install snapd

# Fedora
sudo dnf install snapd

# Arch Linux
sudo pacman -S snapd
```

## Build the Snap

Build the snap package using GoReleaser:

```bash
goreleaser release --snapshot --clean
```

This creates snap files in the `dist/` folder.

## Install Local Snap

Install the snap with `--dangerous` flag (required for unsigned local snaps):

```bash
sudo snap install dist/tomba_*.snap --dangerous
```

Or specify the exact file:

```bash
sudo snap install dist/tomba_v1.0.7-next_linux_amd64.snap --dangerous
```

## Verify Installation

```bash
tomba version
snap list tomba
```

## Uninstall

```bash
sudo snap remove tomba
```

## Notes

- The `--dangerous` flag is required because local snaps are not signed by the Snap Store
- For development, you can also use `--devmode` to disable confinement:
    ```bash
    sudo snap install dist/tomba_v1.0.7-next_linux_amd64.snap --dangerous --devmode
    ```
- Classic confinement requires `--classic` flag if the snap uses it
