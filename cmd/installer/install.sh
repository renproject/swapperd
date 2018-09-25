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

unzip -o swapper.zip
chmod +x bin/swapper
chmod +x bin/installer

# get passphrase from user
while :
do
  echo "Please enter a passphrase to encrypt your key files"
  read -s PASSWORDFirst
  echo "Please re-enter the passphrase"
  read -s PASSWORDSecond
  if [ $PASSWORDFirst = $PASSWORDSecond ]
  then
    if [ "$PASSWORDFirst" = "" ]
    then
      echo "${RED}You are trying to use an empty passphrase, this means your keystores will not be encrypted are you sure (y/N): ${NC}"
      while :
      do
        read choice
        choice=$(echo "$choice" | tr '[:upper:]' '[:lower:]')
        echo
        if [ "$choice" = "" ] || [ "$choice" = "y" ] || [ "$choice" = "yes" ]
        then
          break
        elif [ "$choice" = "n" ] || [ "$choice" = "no" ]
        then
          exit 0
        else
         echo "Please enter (y/N)"
        fi
      done
    fi
    break
  else
    echo "${RED}The two passwords you enter are different. Try again ${NC}"
  fi
done
./bin/installer -passphrase $PASSWORD< /dev/tty

# make sure the swapper service is started when booted
if [ "$ostype" = 'Linux' -a "$cputype" = 'x86_64' ]; then
  sudo echo "[Unit]
Description=RenEx's Swapper Daemon
After=network.target

[Service]
ExecStart=$HOME/.swapper/bin/swapper --loc $HOME/.swapper
Restart=on-failure
StartLimitBurst=0

# Specifies which signal to use when killing a service. Defaults to SIGTERM.
# SIGHUP gives parity time to exit cleanly before SIGKILL (default 90s)
KillSignal=SIGHUP

[Install]
WantedBy=default.target" > swapper.service

  sudo mv swapper.service /etc/systemd/system/swapper.service
  sudo systemctl daemon-reload
  sudo systemctl enable swapper.service
  sudo systemctl start swapper.service

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
</plist>" > com.republicprotocol.plist
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
      echo '\nexport PATH=$PATH:$HOME/.swapper/bin' >> $HOME/.zprofile
      swapper_home=$HOME/.zprofile
    elif [ -f "$HOME/.zshrc" ] ; then
      echo '\nexport PATH=$PATH:$HOME/.swapper/bin' >> $HOME/.zshrc
      swapper_home=$HOME/.zshrc
    elif [ -f "$HOME/.profile" ] ; then
      echo '\nexport PATH=$PATH:$HOME/.swapper/bin' >> $HOME/.profile
      swapper_home=$HOME/.profile
    fi
  elif  [ "$shell" = 'bash' ] ; then
    if [ -f "$HOME/.bash_profile" ] ; then
      echo '\nexport PATH=$PATH:$HOME/.swapper/bin' >> $HOME/.bash_profile
      swapper_home=$HOME/.bash_profile
    elif [ -f "$HOME/.bashrc" ] ; then
      echo '\nexport PATH=$PATH:$HOME/.swapper/bin' >> $HOME/.bashrc
      swapper_home=$HOME/.bashrc
    elif [ -f "$HOME/.profile" ] ; then
      echo '\nexport PATH=$PATH:$HOME/.swapper/bin' >> $HOME/.profile
      swapper_home=$HOME/.profile
    fi
  elif [ -f "$HOME/.profile" ] ; then
    echo '\nexport PATH=$PATH:$HOME/.swapper/bin' >> $HOME/.profile
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