#!/bin/sh

remove_plist()
{
    launchctl unload -w "$HOME/Library/LaunchAgents/ren.swapperd.plist"
    launchctl remove "$HOME/Library/LaunchAgents/ren.swapperd.plist"
    rm "$HOME/Library/LaunchAgents/ren.swapperd.plist"
}

remove_service()
{
    systemctl --user stop swapperd.service
    rm $HOME/.config/systemd/user/swapperd.service
}

timestamp() {
  date +"%Y-%m-%d_%H-%M-%S"
}

# get system information
ostype="$(uname -s)"
cputype="$(uname -m)"

if [ "$ostype" = 'Linux' -a "$cputype" = 'x86_64' ]; then
  remove_service
elif [ "$ostype" = 'Darwin' -a "$cputype" = 'x86_64' ]; then
  remove_plist
fi

mkdir -p $HOME/.swapperd_backup
mv $HOME/.swapperd/testnet.json $HOME/.swapperd_backup/testnet-$(timestamp).json
mv $HOME/.swapperd/mainnet.json $HOME/.swapperd_backup/mainnet-$(timestamp).json
rm -rf $HOME/.swapperd