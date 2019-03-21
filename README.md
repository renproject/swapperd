# Swapperd

[![CircleCI](https://circleci.com/gh/renproject/swapperd/tree/master.svg?style=svg)](https://circleci.com/gh/renproject/swapperd/tree/master)

[Documentation](https://republicprotocol.github.io/swapperd)

Swapperd is built and officially supported by the Republic Protocol team. It is a daemon that can be used to execute cross-chain atomic swaps between Bitcoin, Ethereum, and ERC20 tokens.

## Installation

### macOS and Ubuntu

`curl https://git.io/test-swapperd -sSLf | sh`

### Windows

Download the latest version from [releases](https://github.com/renproject/swapperd/releases)!


## Development

To build locally, run:

```bash
make build
```

In order to cross-compile for all platforms, you will need Docker and [xgo](https://github.com/karalabe/xgo/).

To build for all platforms, run:

```bash
make
```

You can build for a specific platform by running:

```bash
make [platform]
```

where [platform] is one of `darwin`, `linux`, or `windows`.

