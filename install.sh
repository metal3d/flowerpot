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
make user-install
cd -
rm -rf $TMP

echo "FlowerPot is now installed in $HOME/.local/bin and $HOME/.local/share/flowerpot"
echo "You may find FlowerPot in your desltop menu, or you can run it from a terminal with 'flowerpot'"
echo "Enjoy!"

