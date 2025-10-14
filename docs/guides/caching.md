# Caching Guide

**gac** includes built-in caching to improve performance and reduce API calls when listing users and groups.

## Table of Contents

- [Overview](#overview)
- [How Caching Works](#how-caching-works)
- [Configuration](#configuration)
- [Using the Cache](#using-the-cache)
- [Cache Management](#cache-management)
- [Performance Benefits](#performance-benefits)
- [Best Practices](#best-practices)

## Overview

Caching stores the results of expensive API calls locally, allowing subsequent requests to return results instantly without hitting the Google API. This provides:

- **Faster response times** - Cached results return in < 50ms vs 500-2000ms for API calls
- **Reduced API quota usage** - Fewer calls to Google APIs
- **Better user experience** - Near-instant results for repeated queries
- **Offline capability** - Access previously fetched data without network connectivity

## How Caching Works

### Cache Storage

Cache entries are stored as JSON files in `~/.cache/gac/` by default:

```
~/.cache/gac/
├── users-example.com-default.json          # User list for example.com
├── users-example.com-a1b2c3d4.json         # User list with filters
├── groups-example.com-default.json         # Group list for example.com
├── group-members-team@example.com.json     # Members of specific group
└── ...
```

### Cache Keys

Cache keys are generated based on:
- **Resource type** - users, groups, or group-members
- **Domain** - Your Google Workspace domain
- **Filters** - Any query filters (e.g., disabled-only)

This ensures different queries are cached separately.

### Time-To-Live (TTL)

Each cache entry has a TTL (Time-To-Live) that determines how long the data is valid:
- Default: **15 minutes**
- Configurable via config file or command-line flag
- After TTL expires, next request fetches fresh data from API

### Cached Commands

The following commands support caching:

- `gac user list` - Caches user listings
- `gac group list` - Caches group listings
- `gac group list <group> --get-members` - Caches group member lists

Single-item lookups (e.g., `gac user list user@example.com`) are **not cached** to ensure real-time data.

## Configuration

### Config File

Add caching configuration to `~/.google-admin.yaml`:

```yaml
cache:
  enabled: true
  ttl: 15m
  directory: ~/.cache/gac
```

**Options:**
- `enabled` - Enable/disable caching (default: `true`)
- `ttl` - Cache TTL duration (default: `15m`)
  - Formats: `30s`, `5m`, `1h`, `2h30m`
- `directory` - Cache storage location (default: `~/.cache/gac`)

### Command-Line Flags

Override cache settings per command:

```bash
# Disable cache for a single command
gac user list --no-cache

# Use custom TTL
gac user list --cache-ttl 30m
gac group list --cache-ttl 1h
```

## Using the Cache

### Basic Usage

Simply use list commands normally - caching works automatically:

```bash
# First call - fetches from API, writes to cache
gac user list

# Second call within TTL - reads from cache (instant)
gac user list

# After TTL expires - fetches fresh data from API
```

### Forcing Fresh Data

To bypass the cache and get fresh data:

```bash
# Disable cache for this request
gac user list --no-cache

# Or clear cache first, then fetch
gac cache clear users
gac user list
```

### Debug Mode

See cache hits/misses with verbose logging:

```bash
gac -v user list
# Output shows:
# DBG Cache miss, fetching from API key=users-example.com-default.json
# or
# DBG Using cached user list key=users-example.com-default.json count=125
```

## Cache Management

### View Cache Status

Check cache statistics:

```bash
gac cache status
```

**Example output:**
```
Cache Status
============
Enabled:       true
Location:      /Users/you/.cache/gac
Total Size:    2.4 MB
Entry Count:   12
Oldest Entry:  3 hours ago
Newest Entry:  2 minutes ago
Default TTL:   15m0s
```

**JSON output:**
```bash
gac cache status --format json
```

### Clear Cache

Clear specific or all cache entries:

```bash
# Clear user cache
gac cache clear users

# Clear group cache
gac cache clear groups

# Clear all caches
gac cache clear all
gac cache clear --all
```

### Disable Caching

**Temporarily** (single command):
```bash
gac user list --no-cache
```

**Globally** (all commands):

Edit `~/.google-admin.yaml`:
```yaml
cache:
  enabled: false
```

Or set environment variable:
```bash
export GAC_CACHE_ENABLED=false
```

## Performance Benefits

### Benchmark Results

Typical performance improvements with caching enabled:

| Operation | Without Cache | With Cache | Improvement |
|-----------|---------------|------------|-------------|
| List 100 users | 1,200ms | 35ms | **34x faster** |
| List 50 groups | 2,500ms | 28ms | **89x faster** |
| List group members | 800ms | 22ms | **36x faster** |

### API Quota Savings

For a typical workflow listing users 10 times per day:
- **Without cache**: 10 API calls/day
- **With cache (15m TTL)**: ~1-2 API calls/day
- **Quota savings**: 80-90% reduction

## Best Practices

### When to Use Caching

✅ **Good use cases:**
- Frequent queries for the same data
- Scripting and automation workflows
- Interactive development and testing
- Generating reports from user/group data
- Auditing and compliance checks

❌ **When to disable caching:**
- Critical operations requiring real-time data
- Immediately after making changes (users, groups)
- Troubleshooting synchronization issues
- Security audits requiring current state

### Recommended TTL Values

Choose TTL based on your needs:

| Use Case | Recommended TTL | Rationale |
|----------|----------------|-----------|
| Active development | `5m` | Frequent changes, need fresh data |
| Daily operations | `15m` (default) | Balance of freshness and performance |
| Reporting/analytics | `1h` | Data changes less frequently |
| Read-only audits | `2h` or more | Minimal changes expected |

### Cache Invalidation

Clear cache after making changes:

```bash
# After creating/updating users
gac user create john@example.com
gac cache clear users

# After modifying groups
gac group add-member team@example.com user@example.com
gac cache clear groups

# Or use --no-cache for the read operation
gac user create john@example.com
gac user list --no-cache  # Bypasses cache
```

### Automation Scripts

For scripts that make changes and then read data:

```bash
#!/bin/bash

# Make changes
gac user create new-user@example.com --first-name "New" --last-name "User"

# Clear cache to ensure fresh data
gac cache clear users

# Verify the change
gac user list | grep new-user@example.com
```

Or use `--no-cache` flag:

```bash
#!/bin/bash

# Make changes
gac user create new-user@example.com --first-name "New" --last-name "User"

# Read with fresh data (bypass cache)
gac user list --no-cache | grep new-user@example.com
```

### Monitoring Cache Health

Periodically check cache status:

```bash
# View cache stats
gac cache status

# If cache size grows too large, clear old entries
gac cache clear all

# Or adjust TTL to expire entries faster
gac user list --cache-ttl 10m
```

### Security Considerations

Cache files contain sensitive data:

1. **File Permissions**: Cache files are created with `0600` (owner read/write only)
2. **Location**: Default cache directory is in user's home directory
3. **Cleanup**: Consider clearing cache on shared/public systems:
   ```bash
   gac cache clear all
   ```

## Troubleshooting

### Cache Not Working

If caching doesn't seem to work:

1. **Check if enabled:**
   ```bash
   gac cache status
   # Verify: Enabled: true
   ```

2. **Check for --no-cache flag:**
   ```bash
   # Remove any --no-cache flags
   gac user list  # instead of: gac user list --no-cache
   ```

3. **Verify cache directory exists:**
   ```bash
   ls -la ~/.cache/gac/
   ```

4. **Enable debug logging:**
   ```bash
   gac -v user list
   # Look for "Cache hit" or "Cache miss" messages
   ```

### Stale Data

If you're seeing outdated information:

1. **Clear the cache:**
   ```bash
   gac cache clear all
   ```

2. **Reduce TTL:**
   ```yaml
   # ~/.google-admin.yaml
   cache:
     ttl: 5m  # Shorter TTL
   ```

3. **Use --no-cache for critical operations:**
   ```bash
   gac user list --no-cache
   ```

### Cache Permission Errors

If you see permission errors:

```bash
# Fix cache directory permissions
chmod 700 ~/.cache/gac
chmod 600 ~/.cache/gac/*.json
```

## Examples

### Example 1: Development Workflow

```bash
# First run - fetches from API (slow)
time gac user list > /dev/null
# real    0m1.234s

# Subsequent runs - uses cache (fast)
time gac user list > /dev/null
# real    0m0.032s

# After making changes, clear cache
gac user update john@example.com --dept Engineering
gac cache clear users

# Next run fetches fresh data
gac user list
```

### Example 2: Daily Report Script

```bash
#!/bin/bash
# daily-report.sh

# Use 1-hour cache for relatively static data
GAC_FLAGS="--cache-ttl 1h"

echo "Generating daily user report..."
gac user list $GAC_FLAGS --format csv > users-$(date +%Y%m%d).csv

echo "Generating group membership report..."
gac group list $GAC_FLAGS --format csv > groups-$(date +%Y%m%d).csv

echo "Cache statistics:"
gac cache status
```

### Example 3: Interactive Exploration

```bash
# Exploring users - enable caching for speed
gac user list | grep Engineering
gac user list | grep Marketing
gac user list | grep Sales
# All subsequent calls use cache

# When done exploring, clear cache
gac cache clear users
```

## See Also

- [User Management Guide](user-management.md)
- [Group Management Guide](group-management.md)
- [Configuration](../../README.md#configuration)
