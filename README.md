# Fly.io Disk Performance Test

This project tests disk performance on Fly.io by measuring the time it takes to perform different file operations on a 1GB SQLite database file:

- Symbolic link creation
- Hard link creation
- File copy (`cp`)

This tests runs against the `performance-2x` machine type in the `iad` region, and doesn't use volumes.

### Results

```
2025-08-14T00:10:08Z app[48e42d6f7573e8] iad [info]Starting symlink
2025-08-14T00:10:08Z app[48e42d6f7573e8] iad [info]- 523.923µs for symlink
2025-08-14T00:10:08Z app[48e42d6f7573e8] iad [info]Starting hardlink
2025-08-14T00:12:44Z app[48e42d6f7573e8] iad [info]- 2m35.715685394s for hardlink
2025-08-14T00:13:00Z app[d8d4e25a195538] iad [info]Starting cp
2025-08-14T00:13:33Z app[48e42d6f7573e8] iad [info]- 48.785219465s for cp
```

## Setup

1. Install the Fly.io CLI: https://fly.io/docs/hands-on/install-flyctl/
2. Authenticate with Fly.io: `fly auth login`
3. Deploy the app: `make deploy`
4. Run the test: `fly machines start <machine-id>`
5. View logs: `make logs`
