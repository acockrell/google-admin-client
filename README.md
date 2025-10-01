
# gac

**G**oogle **A**dmin **C**lient.

Perform some admin fuctions against Google Apps.

## Overview

A Go CLI tool for performing administrative functions against Google Workspace (Google Apps).

### Core Functionality
- **User Management:** Create, list, and update users with support for groups, departments, managers, phone numbers, addresses, employee types, and custom employee IDs
- **Group Management:** List and manage Google Workspace groups
- **Calendar Operations:** Create, list, and update calendars
- **Transfer Operations:** Handle resource transfers between users

### Technical Details
- Written in Go 1.17+
- CLI built with Cobra framework
- Uses Google Admin SDK APIs and OAuth2 for authentication
- Configuration managed via Viper
- Produces cross-platform binaries (Linux/macOS) via Jenkins CI/CD pipeline

### Architecture
- Entry point: `main.go` delegates to `cmd.Execute()`
- Command structure organized in `cmd/` directory with separate files for each command group (user, group, calendar, transfer)


## Usage

```console
  gac --help
  gac init
  gac init -h
  gac group list > groups.csv
  gac user create newuser@example.com
  gac user create -g group1 newuser@example.com
  gac user create -g group1 -g group2 newuser@example.com
  gac user create -h
  gac user list
  gac user list jdoe@example.com
  gac user update --address "Columbus, OH" jdoe@example.com
  gac user update --dept Engineering jdoe@example.com
  gac user update --group dev jdoe@example.com
  gac user update -g dev -g info jdoe@example.com
  gac user update --id $(uuidgen) jdoe@example.com
  gac user update --id $(uuidgen) --force jdoe@example.com
  gac user update --manager manager@example.com jdoe@example.com
  gac user update --phone mobile:555-555-5555 jdoe@example.com
  gac user update --phone 'mobile:555-555-5555; work:555-123-4567,555' jdoe@example.com
  gac user update --title "Sales Engineer" jdoe@example.com
  gac user update --type staff jdoe@example.com
```

