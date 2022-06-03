# MCPSync

MemCard Pro Sync attempts to be a small utility to sync MemCard Pro save data via its web API to a local directory for backup

# Installation

Requires a fairly recent installation of golang

`go install github.com/gordallott/mcpsync/cmd/mcpsync`

# Usage

## sync saves to local directory

`mcpsync --ip 192.168.1.165 --dest /mnt/my-nas/saves/psx/`

## show info about the running MemCard Pro

`mcpsync show --ip 192.168.1.165`

# Notes

The MemCard Pro needs to be available on the same network as this utility, It is worth binding your mcp to a specific static IP address in your router settings.

The currently selected card can't be synced, That is if you are playing a game, that game's save data won't be synced. This is a MCP limitation.

To save on writes, hashes of existing save files are compared to hashes of mcp save files, so saving the same data twice may not result in a sync as the data may not have changed.

# Testing

Tested with MCP ver 1.1.1 on Ubuntu linux
