#!/bin/sh

generate_plist()
{
password=$1
echo "<?xml version=\"1.0\" encoding=\"UTF-8\"?>
<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">
<plist version=\"1.0\">
  <dict>
    <key>Label</key>
    <string>exchange.ren.swapper</string>
    <key>ProgramArguments</key>
    <array>
        <string>$HOME/.swapper/bin/swapper</string>
        <string>-loc</string>
        <string>$HOME/.swapper</string>$password
    </array>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>$HOME/.swapper/swapper.log</string>
    <key>StandardErrorPath</key>
    <string>$HOME/.swapper/swapper.log</string>
  </dict>
</plist>" > "$HOME/Library/LaunchAgents/exchange.ren.swapper.plist"
}

generate_service()
{
password=$1
echo "[Unit]
Description=RenEx's Swapper Daemon
After=network.target

[Service]
ExecStart=$HOME/.swapper/bin/swapper --loc $HOME/.swapper $password
Restart=on-failure
StartLimitBurst=0

# Specifies which signal to use when killing a service. Defaults to SIGTERM.
# SIGHUP gives parity time to exit cleanly before SIGKILL (default 90s)
KillSignal=SIGHUP

[Install]
WantedBy=default.target" > swapper.service
sudo mv swapper.service /etc/systemd/system/swapper.service
}

# install unzip if command not found
if ! [ -x "$(command -v unzip)" ];then
  echo "please install unzip"
  exit 0
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

# assume the service/plist file exists
if ls "$HOME"/.swapper/BTC*.json 1> /dev/null 2>&1; then
  if ls "$HOME"/.swapper/ETH*.json 1> /dev/null 2>&1; then
    echo "RenEx Atomic Swapper has already been installed, updating..."
    if [ "$ostype" = 'Linux' -a "$cputype" = 'x86_64' ]; then
      sudo systemctl daemon-reload
      sudo systemctl restart swapper.service
    elif [ "$ostype" = 'Darwin' -a "$cputype" = 'x86_64' ]; then
      if [ "$(launchctl list | grep exchange.ren.swapper | wc -l)" -le 1 ]; then
        launchctl unload -w "$HOME/Library/LaunchAgents/exchange.ren.swapper.plist"
        sleep 5
      fi
      launchctl load -w "$HOME/Library/LaunchAgents/exchange.ren.swapper.plist"
    fi
    rm swapper.zip
    rm bin/installer
    echo "RenEx Atomic Swapper has been updated. Great!"
    exit 0
  fi
fi

# get passphrase from user
while :
do
  echo "Please enter a passphrase to encrypt your key files"
  read -s PASSWORD </dev/tty
  echo "Please re-enter the passphrase"
  read -s PASSWORDCONFIRM </dev/tty
  if [ "$PASSWORD" = "$PASSWORDCONFIRM" ]; then
    if [ "$PASSWORD" = "" ]; then
      echo "You are trying to use an empty passphrase, this means your keystores will not be encrypted are you sure (y/N): "
      while :
      do
        read choice </dev/tty
        choice=$(echo "$choice" | tr '[:upper:]' '[:lower:]')
        echo
        if [ "$choice" = "y" ] || [ "$choice" = "yes" ]; then
          confirm="yes"
          break
        elif [ "$choice" = "" ] || [ "$choice" = "n" ] || [ "$choice" = "no" ]; then
          confirm="no"
          break
        else
         echo "Please enter (y/N)"
        fi
      done
    else
      break
    fi

    if [ "$confirm" = "yes" ]; then
      break
    fi
  else
    echo "The two passwords you enter are different. Try again."
  fi
done

if [ "$PASSWORD" = "" ]; then
  PASSPHRASE=""
else
  PASSPHRASE="--passphrase $PASSWORD"
fi
./bin/installer $PASSPHRASE < /dev/tty

# make sure the swapper service is started when booted
if [ "$ostype" = 'Linux' -a "$cputype" = 'x86_64' ]; then
  generate_service "$PASSPHRASE"
  sudo systemctl daemon-reload
  sudo systemctl enable swapper.service
  sudo systemctl start swapper.service
elif [ "$ostype" = 'Darwin' -a "$cputype" = 'x86_64' ]; then
  if [ "$PASSWORD" = "" ]
  then
    PASSPHRASE=""
  else
    PASSPHRASE="
        <string>-passphrase</string>
        <string>$PASSWORD</string>"
  fi
  generate_plist "$PASSPHRASE"
  chmod +x "$HOME/Library/LaunchAgents/exchange.ren.swapper.plist"
  if [ "$(launchctl list | grep exchange.ren.swapper | wc -l)" -le 1 ]; then
    launchctl unload -w "$HOME/Library/LaunchAgents/exchange.ren.swapper.plist"
  fi
  launchctl load -w "$HOME/Library/LaunchAgents/exchange.ren.swapper.plist"
else
  echo 'unsupported OS type or architecture'
  cd ..
  rm -rf "$HOME/.swapper"
  exit 1
fi

# clean up
rm swapper.zip
rm bin/installer

echo "RenEx Atomic Swapper is installed now. Great!"
