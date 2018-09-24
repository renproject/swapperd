#!/bin/sh

# define color escape codes
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

# install unzip if command not found
if ! [ -x "$(command -v unzip)" ];then
  sudo apt-get install unzip
fi

# creating working directory
mkdir -p $HOME/.swapper
cd $HOME/.swapper

# get system information
ostype="$(uname -s)"
cputype="$(uname -m)"

# download swapper binary depending on the system and architecture
if [ "$ostype" = 'Linux' -a "$cputype" = 'x86_64' ]; then
  curl -s 'https://releases.republicprotocol.com/swapper/swapper_linux_amd64.zip' > swapper.zip
elif [ "$ostype" = 'Darwin' -a "$cputype" = 'x86_64' ]; then
  curl -s 'https://releases.republicprotocol.com/swapper/swapper_darwin_amd64.zip' > swapper.zip
else
  echo 'unsupported OS type or architecture'
  cd ..
  rm -rf .swapper
  exit 1
fi

unzip swapper.zip
chmod +x bin/swapper
chmod +x bin/installer

./bin/installer < /dev/tty

# make sure the swapper service is started when booted
if [ "$ostype" = 'Linux' -a "$cputype" = 'x86_64' ]; then
  chmod +x install.sh
  ./install.sh
  rm install.sh
elif [ "$ostype" = 'Darwin' -a "$cputype" = 'x86_64' ]; then
  echo "<?xml version=\"1.0\" encoding=\"UTF-8\"?>
<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">
<plist version=\"1.0\">
  <dict>
    <key>Label</key>
    <string>com.republicprotocol.swapper</string>
    <key>ProgramArguments</key>
    <array>
        <string>$HOME/.swapper/bin/swapper</string>
        <string>-loc</string>
        <string>$HOME/.swapper</string>
    </array>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/var/log/swapper.log</string>
    <key>StandardErrorPath</key>
    <string>/var/log/swapper.log</string>
  </dict>
</plist>" >> com.republicprotocol.plist
  sudo mv com.republicprotocol.plist /Library/LaunchDaemons/com.republicprotocol.plist
  sudo chown root /Library/LaunchDaemons/com.republicprotocol.plist
  sudo launchctl load -w /Library/LaunchDaemons/com.republicprotocol.plist
else
  echo 'unsupported OS type or architecture'
  cd ..
  rm -rf .swapper
  exit 1
fi

# clean up
rm swapper.zip
rm bin/installer

# make sure the binary is installed in the path
if ! [ -x "$(command -v swapper)" ]; then
  path=$SHELL
  shell=${path##*/}

  if [ "$shell" = 'zsh' ] ; then
    if [ -f "$HOME/.zprofile" ] ; then
      echo 'export PATH=$PATH:$HOME/.swapper/bin' >> $HOME/.zprofile
      swapper_home=$HOME/.zprofile
    elif [ -f "$HOME/.zshrc" ] ; then
      echo 'export PATH=$PATH:$HOME/.swapper/bin' >> $HOME/.zshrc
      swapper_home=$HOME/.zshrc
    elif [ -f "$HOME/.profile" ] ; then
      echo 'export PATH=$PATH:$HOME/.swapper/bin' >> $HOME/.profile
      swapper_home=$HOME/.profile
    fi
  elif  [ "$shell" = 'bash' ] ; then
    if [ -f "$HOME/.bash_profile" ] ; then
      echo 'export PATH=$PATH:$HOME/.swapper/bin' >> $HOME/.bash_profile
      swapper_home=$HOME/.bash_profile
    elif [ -f "$HOME/.bashrc" ] ; then
      echo 'export PATH=$PATH:$HOME/.swapper/bin' >> $HOME/.bashrc
      swapper_home=$HOME/.bashrc
    elif [ -f "$HOME/.profile" ] ; then
      echo 'export PATH=$PATH:$HOME/.swapper/bin' >> $HOME/.profile
      swapper_home=$HOME/.profile
    fi
  elif [ -f "$HOME/.profile" ] ; then
    echo 'export PATH=$PATH:$HOME/.swapper/bin' >> $HOME/.profile
  fi

  echo ''
  echo 'If you are using a custom shell, make sure you update your PATH.'
  echo "${GREEN}export PATH=\$PATH:\$HOME/.swapper/bin ${NC}"
fi

echo "RenEx Atomic Swapper is installed now. Great!"
echo ''
echo "To get started you need RenEx Atomic Swapper's bin directory ($HOME/.swapper/bin) in your PATH"
echo "environment variable. Next time you log in this will be done"
echo "automatically."
echo ''
echo "To configure your current shell run 'source ${swapper_home}'"