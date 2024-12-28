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

| Command    | Description                                |
| ---------- | ------------------------------------------ |
| bsm help   | Show help                                  |
| bsm config | Create bsm config file                     |
| bsm status | Show current world, backups, version, etc. |
| bsm health | Check server & backup storage space        |

## Server

| Command                     | Description           |
| --------------------------- | --------------------- |
| bsm server setup {version}  | Setup server          |
| bsm server update {version} | Update server version |

## Worlds

| Command                  | Description            |
| ------------------------ | ---------------------- |
| bsm worlds list          | List all worlds        |
| bsm worlds switch {name} | Switch to world {name} |
| bsm worlds delete {name} | Delete world {name}    |
| bsm worlds export {name} | Export world {name}    |
| bsm worlds import {file} | Import world {file}    |

## Backups

| Command                   | Description           |
| ------------------------- | --------------------- |
| bsm backup list           | List all backups      |
| bsm backup create {name}  | Create backup {name}  |
| bsm backup restore {name} | Restore backup {name} |
