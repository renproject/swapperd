## Installation Instructions

### Ubuntu

1. Install Bitcoin node

`sudo apt-add-repository ppa:bitcoin/bitcoin`

`sudo apt-get update`

`sudo apt-get install bitcoind`

2. Start the Bitcoin node

`bitcoind -daemon`

3. Stop the Bitcoin node

`bitcoin-cli stop`

4. Update the Bitcoin config to use the Bitcoin Testnet

```sh 
echo "testnet=1
blocksonly=1
rest=1
server=1
listen=0
rpcallowip=0.0.0.0/0 
rpcuser=<enter_your_username>
rpcpassword=<enter_your_password>"
>> ~/.bitcoin/bitcoin.conf
```

5. Start the bitcoin node

`bitcoind -daemon`
