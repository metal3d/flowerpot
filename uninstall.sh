#!/bin/bash

# Get the latest release from GitHub API
LATEST=$(curl -s https://api.github.com/repos/metal3d/flowerpot/releases/latest | grep "browser_download_url" | awk -F '"' '{print $4}')

echo "Downloading latest release: $LATEST"

# Create a temporary directory and install the latest release
TMP=$(mktemp -d -t flowerpot-XXXXXX)
echo "Working in $TMP"
cd $TMP
curl -L -o flowerpot.tar.xz $LATEST
tar xf flowerpot.tar.xz
make user-uninstall
rm -rf $HOME/.config/autostart/flowerpot.desktop
cd -
rm -rf $TMP

echo "FlowerPot is now uninstalled"


