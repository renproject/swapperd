#!/bin/sh

RELEASES_URL="https://github.com/renproject/swapperd/releases"
latest_version() {
  curl -sL -o /dev/null -w %{url_effective} "$RELEASES_URL/latest" | 
    rev | 
    cut -f1 -d'/'| 
    rev
}

timestamp() {
  date +"%Y-%m-%d_%H-%M-%S"
}

# install unzip if command not found
if ! [ -x "$(command -v unzip)" ];then
  echo "please install unzip"
  exit 0
fi

# creating working directory
mkdir -p $HOME/.swapperd
cd $HOME/.swapperd

# get system information
ostype="$(uname -s)"
cputype="$(uname -m)"

# check whether the service/plist file exists
if ls "$HOME"/.swapperd/*.json 1> /dev/null 2>&1; then
  echo "Swapperd has already been installed"
  exit 0
fi

VERSION=$(latest_version)
if [ "$VERSION" = '' ]; then 
  echo "Cannot get the latest version from Github"
  exit 1
fi

echo "Latest version of swapperd is $VERSION"
# download swapperd binary depending on the system and architecture
if [ "$ostype" = 'Linux' -a "$cputype" = 'x86_64' ]; then
  curl -Ls "https://github.com/renproject/swapperd/releases/download/$VERSION/swapper_linux_amd64.zip" > swapper.zip
elif [ "$ostype" = 'Darwin' -a "$cputype" = 'x86_64' ]; then
  curl -Ls "https://github.com/renproject/swapperd/releases/download/$VERSION/swapper_darwin_amd64.zip" > swapper.zip
else
  echo 'unsupported OS type or architecture'
  cd ..
  rm -rf .swapperd
  exit 1
fi

echo "{ \"version\": \"$VERSION\", \"frequency\": 3600 }" > config.json

unzip -o swapper.zip
chmod +x bin/swapperd
chmod +x bin/installer

if [ "$1" = '' ]; then
  ./bin/installer
else 
  ./bin/installer --mnemonic "$1"
fi

mkdir -p $HOME/.swapperd_backup
cp $HOME/.swapperd/testnet.json $HOME/.swapperd_backup/testnet-$(timestamp).json
cp $HOME/.swapperd/mainnet.json $HOME/.swapperd_backup/mainnet-$(timestamp).json

# clean up
rm swapper.zip
rm bin/installer

echo "Swapperd is installed now. Great!"