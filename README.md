# bedrock server manager

A tool to help you with the annoying tasks of managing your bedrock server.

This tool is meant to be a backup manager and updater for your bedrock server. It will help you keep your server up to date and keep your backups organized.

## Features

-   Setup server (DONE)
-   List, create and switch between active worlds (DONE)
-   Update server (NOT DONE)
-   Backup & Restore worlds (NOT DONE)

# CLI INTERFACE

## General

| Command | Description                         | Status       |
| ------- | ----------------------------------- | ------------ |
| help    | Show help                           | finished     |
| config  | Create config file                  | finished     |
| start   | Start server in the background      | not finished |
| stop    | Stop server                         | not finished |
| status  | Show status of bsm and server       | not finished |
| health  | Check server & backup storage space | not finished |

files to update:
behavior_packs
config
definitions
resource_packs
bedrock_server
bedrock_server_how_to.html
profanity_filter.wlist
release-notes.txt

when setting up:
keep track of server version

when updating:
update only when version is different

other stuff:
when switching worlds, if the server is running, ask if it should be stopped, switch and start again
when restoring backup and the to-be-restored world is already active, ask if it should be stopped, switch and start again
sync all world's server.properties with the config file for server-name

## Server

| Command                 | Description           | Status       |
| ----------------------- | --------------------- | ------------ |
| server setup {version}  | Setup server          | finished     |
| server update {version} | Update server version | not finished |

## Worlds

| Command             | Description            | Status       |
| ------------------- | ---------------------- | ------------ |
| world list          | List all worlds        | finished     |
| world switch {name} | Switch to world {name} | finished     |
| world create {name} | Create world {name}    | finished     |
| world delete {name} | Delete world {name}    | not finished |

## Backups

| Command               | Description           | Status       |
| --------------------- | --------------------- | ------------ |
| backup list           | List all backups      | finished     |
| backup create {name}  | Create backup {name}  | finished     |
| backup restore {name} | Restore backup {name} | not finished |
