# Fly.io Disk Performance Test

This project tests disk performance on Fly.io by measuring the time it takes to perform different file operations:

- File copy (`cp`)
- Hard link creation
- Symbolic link creation

## Prerequisites

1. Install the Fly.io CLI: https://fly.io/docs/hands-on/install-flyctl/
2. Authenticate with Fly.io: `fly auth login`

## Quick Start

### Deploy to Fly.io

```bash
make deploy
```

### View Logs

```bash
make logs
```

### Follow Logs in Real-time

```bash
make logs-follow
```
