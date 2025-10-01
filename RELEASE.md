# Release Process

This document describes the release process for `gac` (Google Admin Client).

## Overview

Releases are automated using [GoReleaser](https://goreleaser.com/) and GitHub Actions. When a version tag is pushed to GitHub, the release workflow automatically:

1. Builds binaries for multiple platforms (Linux, macOS, Windows × amd64/arm64)
2. Creates Docker images for amd64 and arm64 architectures
3. Generates checksums for all artifacts
4. Creates a GitHub release with changelog
5. Publishes Docker images to GitHub Container Registry (ghcr.io)

## Prerequisites

- Write access to the repository
- Git configured with your GitHub credentials
- Clean working directory (no uncommitted changes)
- All tests passing

## Version Numbering

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR** version (v1.0.0): Incompatible API changes
- **MINOR** version (v0.1.0): New functionality, backwards-compatible
- **PATCH** version (v0.0.1): Backwards-compatible bug fixes

### Pre-release versions

- **Alpha**: `v0.1.0-alpha.1` - Early testing, unstable
- **Beta**: `v0.1.0-beta.1` - Feature complete, testing
- **RC**: `v0.1.0-rc.1` - Release candidate

## Release Steps

### 1. Prepare the Release

```bash
# Ensure you're on the main branch with latest changes
git checkout main
git pull origin main

# Ensure all tests pass
make test

# Run all checks (fmt, vet, lint, security scan)
make check

# Test the release configuration locally
make release-test
```

### 2. Create and Push the Tag

```bash
# Create an annotated tag with semantic version
git tag -a v0.1.0 -m "Release v0.1.0"

# Push the tag to GitHub (this triggers the release workflow)
git push origin v0.1.0
```

**Alternative: Create tag with detailed message**

```bash
git tag -a v0.1.0 -m "Release v0.1.0

New Features:
- Feature 1
- Feature 2

Bug Fixes:
- Fix 1
- Fix 2
"

git push origin v0.1.0
```

### 3. Monitor the Release

1. Go to **Actions** tab in GitHub: https://github.com/acockrell/google-admin-client/actions
2. Watch the "Release" workflow execution
3. Verify all jobs complete successfully

### 4. Verify the Release

Once the workflow completes:

1. **Check GitHub Release**: https://github.com/acockrell/google-admin-client/releases
   - Verify release notes and changelog
   - Verify all platform binaries are attached
   - Verify checksums.txt is present

2. **Check Docker Images**: https://github.com/acockrell?tab=packages
   ```bash
   docker pull ghcr.io/acockrell/gac:v0.1.0
   docker run --rm ghcr.io/acockrell/gac:v0.1.0 --help
   ```

3. **Test a Binary**:
   ```bash
   # Download and test (example for macOS arm64)
   curl -Lo gac.tar.gz https://github.com/acockrell/google-admin-client/releases/download/v0.1.0/gac_Darwin_arm64.tar.gz
   tar xzf gac.tar.gz
   ./gac --help
   ```

## Local Testing

Before creating a release, test the build process locally:

```bash
# Build snapshot (no tag required)
make release-snapshot

# Test full release process without publishing
make release-test

# Check generated artifacts
ls -lh dist/
```

## Changelog Generation

GoReleaser automatically generates changelogs from commit messages using [Conventional Commits](https://www.conventionalcommits.org/).

### Commit Message Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Commit Types

Commits are grouped in the changelog by type:

- **feat**: New features (appears under "Features")
- **fix**: Bug fixes (appears under "Bug Fixes")
- **sec**: Security updates (appears under "Security Updates")
- **perf**: Performance improvements (appears under "Performance Improvements")
- **docs**: Documentation changes (excluded from changelog)
- **test**: Test changes (excluded from changelog)
- **chore**: Maintenance tasks (excluded from changelog)

### Examples

```bash
# Feature
git commit -m "feat: add user suspension command"

# Bug fix
git commit -m "fix: correct email validation regex"

# Security update
git commit -m "sec: update dependencies with CVE fixes"

# Breaking change
git commit -m "feat!: remove deprecated --cache flag

BREAKING CHANGE: The --cache flag has been removed. Use --cache-file instead."
```

## Automated Release Workflow

The release workflow (`.github/workflows/release.yml`) performs these steps:

1. **Checkout**: Fetch repository with full git history
2. **Setup Go**: Install Go 1.25
3. **Cache**: Cache Go modules for faster builds
4. **Docker Setup**: Configure QEMU and Buildx for multi-arch images
5. **Registry Login**: Authenticate to GitHub Container Registry
6. **GoReleaser**: Build, package, and publish release
7. **Upload Artifacts**: Save build artifacts

### Build Targets

| Platform | Architectures |
|----------|---------------|
| Linux    | amd64, arm64  |
| macOS    | amd64, arm64  |
| Windows  | amd64         |

### Docker Images

Images are published to GitHub Container Registry:

- `ghcr.io/acockrell/gac:v0.1.0` (versioned)
- `ghcr.io/acockrell/gac:latest` (latest release)
- `ghcr.io/acockrell/gac:v0.1.0-amd64` (platform-specific)
- `ghcr.io/acockrell/gac:v0.1.0-arm64` (platform-specific)

## Troubleshooting

### Release workflow fails

**Check the workflow logs:**
1. Go to Actions tab
2. Click on the failed workflow
3. Expand failed job to see error details

**Common issues:**

- **Tests failing**: Run `make test` locally to reproduce
- **Build failing**: Run `make release-snapshot` locally to debug
- **Docker build failing**: Ensure `Dockerfile.goreleaser` is valid
- **Permission denied**: Verify GITHUB_TOKEN has appropriate permissions

### Tag already exists

```bash
# Delete local tag
git tag -d v0.1.0

# Delete remote tag (if pushed)
git push origin :refs/tags/v0.1.0

# Create new tag
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

### GoReleaser configuration errors

```bash
# Validate configuration
~/go/bin/goreleaser check

# Test with environment variables
export GITHUB_REPOSITORY_OWNER=acockrell
~/go/bin/goreleaser check
```

### Docker images not appearing

1. Verify workflow completed successfully
2. Check GitHub Packages: https://github.com/acockrell?tab=packages
3. Ensure Docker steps didn't fail in workflow logs
4. Verify GITHUB_TOKEN permissions include `packages: write`

## Hotfix Releases

For urgent bug fixes on a released version:

```bash
# Create hotfix branch from tag
git checkout -b hotfix/v0.1.1 v0.1.0

# Make fixes and commit
git commit -m "fix: critical bug in user creation"

# Merge back to main
git checkout main
git merge --no-ff hotfix/v0.1.1

# Tag and release
git tag -a v0.1.1 -m "Hotfix v0.1.1: Fix critical bug"
git push origin main v0.1.1

# Delete hotfix branch
git branch -d hotfix/v0.1.1
```

## Rolling Back a Release

If a release has critical issues:

1. **Delete the GitHub release** (via GitHub UI)
2. **Delete the git tag**:
   ```bash
   git push --delete origin v0.1.0
   git tag -d v0.1.0
   ```
3. **Delete Docker images** (via GitHub Packages UI)
4. **Fix the issue and create a new patch release**

## Manual Release (Emergency)

If GitHub Actions is unavailable:

```bash
# Ensure goreleaser is installed
go install github.com/goreleaser/goreleaser/v2@latest

# Set required environment variables
export GITHUB_TOKEN="your-github-token"
export GITHUB_REPOSITORY_OWNER=acockrell

# Create release (requires git tag)
git tag -a v0.1.0 -m "Release v0.1.0"
goreleaser release --clean
```

## Post-Release Tasks

After a successful release:

1. **Announce the release** (if applicable):
   - Update project README if needed
   - Post to team chat/slack
   - Update documentation site

2. **Update TODO.md**: Mark release-related tasks as complete

3. **Monitor for issues**: Watch for bug reports on the new version

## References

- [GoReleaser Documentation](https://goreleaser.com/)
- [Semantic Versioning](https://semver.org/)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [GitHub Actions: Publishing packages](https://docs.github.com/en/actions/publishing-packages)
