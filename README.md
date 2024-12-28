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

| Command    | Description                                | Status       |
| ---------- | ------------------------------------------ | ------------ |
| bsm help   | Show help                                  | finished     |
| bsm config | Create bsm config file                     | finished     |
| bsm status | Show current world, backups, version, etc. | not finished |
| bsm health | Check server & backup storage space        | not finished |

## Server

| Command                     | Description           | Status       |
| --------------------------- | --------------------- | ------------ |
| bsm server setup {version}  | Setup server          | finished     |
| bsm server update {version} | Update server version | not finished |

## Worlds

| Command                  | Description            | Status       |
| ------------------------ | ---------------------- | ------------ |
| bsm worlds list          | List all worlds        | finished     |
| bsm worlds switch {name} | Switch to world {name} | finished     |
| bsm worlds create {name} | Create world {name}    | finished     |
| bsm worlds delete {name} | Delete world {name}    | not finished |
| bsm worlds export {name} | Export world {name}    | not finished |
| bsm worlds import {file} | Import world {file}    | not finished |

## Backups

| Command                   | Description           | Status       |
| ------------------------- | --------------------- | ------------ |
| bsm backup list           | List all backups      | not finished |
| bsm backup create {name}  | Create backup {name}  | not finished |
| bsm backup restore {name} | Restore backup {name} | not finished |
