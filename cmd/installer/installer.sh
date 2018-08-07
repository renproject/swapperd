#!/bin/sh

# define color escape codes
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

# Install unzip if command not found
if ! [ -x "$(command -v unzip)" ];then
  sudo apt-get install unzip
fi

# creating working directory
mkdir -p $HOME/.swapper
cd $HOME/.swapper

# get system information
ostype="$(uname -s)"
cputype="$(uname -m)"

# download darknode binary depending on the system and architecture
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

# clean up zip files
rm swapper.zip

# clean up installer files
rm bin/installer

# make sure the binary is installed in the path
if ! [ -x "$(command -v swapper)" ]; then
  path=$SHELL
  shell=${path##*/}

  if [ "$shell" = 'zsh' ] ; then
    if [ -f "$HOME/.zprofile" ] ; then
      echo 'export PATH=$PATH:$HOME/.swapper/bin' >> $HOME/.zprofile
    elif [ -f "$HOME/.zshrc" ] ; then
      echo 'export PATH=$PATH:$HOME/.swapper/bin' >> $HOME/.zshrc
    elif [ -f "$HOME/.profile" ] ; then
      echo 'export PATH=$PATH:$HOME/.swapper/bin' >> $HOME/.profile
    fi
  elif  [ "$shell" = 'bash' ] ; then
    if [ -f "$HOME/.bash_profile" ] ; then
      echo 'export PATH=$PATH:$HOME/.swapper/bin' >> $HOME/.bash_profile
    elif [ -f "$HOME/.bashrc" ] ; then
      echo 'export PATH=$PATH:$HOME/.swapper/bin' >> $HOME/.bashrc
    elif [ -f "$HOME/.profile" ] ; then
      echo 'export PATH=$PATH:$HOME/.swapper/bin' >> $HOME/.profile
    fi
  elif [ -f "$HOME/.profile" ] ; then
    echo 'export PATH=$PATH:$HOME/.swapper/bin' >> $HOME/.profile
  fi

  echo ''
  echo 'If you are using a custom shell, make sure you update your PATH.'
  echo "${GREEN}export PATH=\$PATH:\$HOME/.swapper/bin ${NC}"
fi

echo ''
echo "${GREEN}Done! Please update run the following command to start the RenEx Swapper.${NC}"
echo ''
echo "${GREEN}swapper${NC}"