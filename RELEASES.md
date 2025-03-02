# Release Process for Task Tracker

This document explains how to create and manage releases for the Task Tracker application on GitHub.

## Versioning

We follow [Semantic Versioning](https://semver.org/) (SemVer) for this project:

- **MAJOR version** (x.0.0): Incompatible API changes
- **MINOR version** (0.x.0): Add functionality in a backward-compatible manner
- **PATCH version** (0.0.x): Backward-compatible bug fixes

## Automated Release Process

We use GitHub Actions to automate the build and release process. The workflow is triggered whenever a new Git tag with the format `v*` is pushed to the repository.

### Creating a New Release

1. Update the version number in your code if applicable
2. Commit your changes to the `main` branch
3. Create and push a new tag:
   ```
   git tag v1.0.0
   git push origin v1.0.0
   ```
4. The GitHub Actions workflow will automatically:
   - Run tests
   - Build binaries for multiple platforms (Windows, macOS, Linux)
   - Create a GitHub release
   - Upload the binaries as release assets

5. Once the workflow completes, navigate to the "Releases" section of your GitHub repository to see the new release

### Using the Release Helper Script

For Windows users, we provide a PowerShell script to help with the release process:

```
.\scripts\release.ps1 -Version 1.0.0
```

This script will:
1. Update the VERSION file
2. Update the release date in CHANGELOG.md
3. Show you the git commands to run for tagging and pushing

### The Release Workflow

The workflow defined in `.github/workflows/release.yml`:

1. Builds the application for multiple platforms:
   - Windows (amd64)
   - macOS (amd64, arm64)
   - Linux (amd64, arm64)
2. Archives the binaries with the README.md file
3. Creates a GitHub release with automatic release notes

## Manual Release Process

If you prefer to create releases manually:

1. Build the binaries locally:
   ```
   # Windows
   GOOS=windows GOARCH=amd64 go build -o task-tracker-windows-amd64.exe ./cmd/task-tracker

   # macOS
   GOOS=darwin GOARCH=amd64 go build -o task-tracker-darwin-amd64 ./cmd/task-tracker
   GOOS=darwin GOARCH=arm64 go build -o task-tracker-darwin-arm64 ./cmd/task-tracker

   # Linux
   GOOS=linux GOARCH=amd64 go build -o task-tracker-linux-amd64 ./cmd/task-tracker
   GOOS=linux GOARCH=arm64 go build -o task-tracker-linux-arm64 ./cmd/task-tracker
   ```

2. Archive the binaries:
   ```
   # Windows
   zip -j task-tracker-windows-amd64.zip task-tracker-windows-amd64.exe README.md

   # macOS and Linux
   tar czf task-tracker-darwin-amd64.tar.gz task-tracker-darwin-amd64 README.md
   tar czf task-tracker-darwin-arm64.tar.gz task-tracker-darwin-arm64 README.md
   tar czf task-tracker-linux-amd64.tar.gz task-tracker-linux-amd64 README.md
   tar czf task-tracker-linux-arm64.tar.gz task-tracker-linux-arm64 README.md
   ```

3. Create a new release on GitHub:
   - Go to your repository on GitHub
   - Click on "Releases"
   - Click "Draft a new release"
   - Enter the tag version (e.g., `v1.0.0`)
   - Fill in the release title and description
   - Upload the binary archives
   - Click "Publish release"

## Release Checklist

Before creating a release, ensure:

- [ ] All tests pass
- [ ] Documentation is up-to-date
- [ ] CHANGELOG.md is updated (if you maintain one)
- [ ] Version numbers are updated in relevant files
- [ ] The code is stable and ready for release

## Hotfixes

For urgent fixes to a released version:

1. Create a branch from the release tag:
   ```
   git checkout -b hotfix/v1.0.1 v1.0.0
   ```

2. Make your fixes on this branch

3. Tag and release as described above

4. Ensure the fixes are also merged back to the `main` branch 