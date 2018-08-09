# RenEx Atomic Swapper

The RenEx Atomic Swapper is built, and officially supported by, the Republic Protocol team. It can used to execute atomic swaps between Ethereum and Bitcoin, and while it can be used independently of RenEx, it is designed for use with https://testnet.ren.exchange. Using this software, traders will be able to open Ethereum to Bitcoin orders on https://testnet.ren.exchange.
    
## Installation

### Mac/Ubuntu

#### Prerequisites

1. Bitcoin Node
2. Curl
3. Metamask

#### Steps

1. Run the following command

`curl https://releases.republicprotocol.com/swapper/install.sh -sSf | sh`

2. When prompted enter the Ethereum address on your MetaMask

3. When prompted enter the required Bitcoin Node information. You will need its IP address, the RPC port, and the RPC credentials (username / password).

### Windows

> Coming soon!

## Usage

The RenEx Atomic Swapper is designed for use with https://testnet.ren.exchange. 

IMPORTANT: The RenEx Atomic Swapper must be running at all times. If it is not running, it will not be able to execute atomic swaps. If you fail to execute an atomic swap for matching orders, your trading account being fined, resulting in the loss of funds.

To open an atomic swap on RenEx:

1. Select Ethereum / Bitcoin trading pair.

2. Click "Connect to atomic swapper".

3. Authorize your atomic swapping software to execute your atomic swaps.

4. Open the "Balances" tab and ensure that you have the necessary funds in the Ethereum and Bitcoin addresses.

5. Open the "Exchange" tab and open an order.