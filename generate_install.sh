#!/bin/sh

RELEASES_URL="https://github.com/renproject/swapperd/releases"
BRANCH=$(git branch | grep \* | cut -d ' ' -f2)
if [ "$BRANCH" == "nightly" ]; then
  VERSION=nightly
else
  VERSION="v$(make version)"
fi

cat  << EOF
#!/bin/sh

timestamp() {
  date +"%Y-%m-%d_%H-%M-%S"
}

# install unzip if command not found
if ! [ -x "\$(command -v unzip)" ];then
  echo "please install unzip"
  exit 0
fi

# creating working directory
mkdir -p \$HOME/.swapperd
cd \$HOME/.swapperd

# get system information
ostype="\$(uname -s)"
cputype="\$(uname -m)"

echo "Latest version of swapperd is ${VERSION}"
# download swapperd binary depending on the system and architecture
if [ "\$ostype" = 'Linux' -a "\$cputype" = 'x86_64' ]; then
  curl -#L "${RELEASES_URL}/download/${VERSION}/swapper_linux_amd64.zip" > swapper.zip
elif [ "\$ostype" = 'Darwin' -a "\$cputype" = 'x86_64' ]; then
  curl -#L "${RELEASES_URL}/download/${VERSION}/swapper_darwin_amd64.zip" > swapper.zip
else
  echo 'unsupported OS type or architecture'
  cd ..
  rm -rf .swapperd
  exit 1
fi

unzip -o swapper.zip

# do not run the installer if mainnet keystore file exists
if ls "\$HOME"/.swapperd/mainnet.json 1> /dev/null 2>&2; then
  echo "Swapperd has already been installed, updating..."
  rm swapper.zip
  rm bin/installer
  echo "Swapperd has been updated. Great!"
  exit 0
fi

chmod +x bin/swapperd
chmod +x bin/installer

if [ "\$1" = '' ]; then
  ./bin/installer
else 
  ./bin/installer --mnemonic "\$1"
fi

mkdir -p \$HOME/.swapperd_backup
cp \$HOME/.swapperd/testnet.json \$HOME/.swapperd_backup/testnet-\$(timestamp).json
cp \$HOME/.swapperd/mainnet.json \$HOME/.swapperd_backup/mainnet-\$(timestamp).json

# clean up
rm swapper.zip
rm bin/installer

echo "Swapperd is installed now. Great!"
EOF
