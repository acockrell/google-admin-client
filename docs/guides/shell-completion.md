# Shell Completion Guide

## Overview

Shell completion (also known as tab completion) allows you to press the Tab key to automatically complete `gac` commands, subcommands, and flags. This guide covers how to set up shell completion for bash, zsh, and fish shells.

## Table of Contents

- [Bash Completion](#bash-completion)
- [Zsh Completion](#zsh-completion)
- [Fish Completion](#fish-completion)
- [Testing Completion](#testing-completion)
- [Troubleshooting](#troubleshooting)

## Bash Completion

### Prerequisites

**Linux:**
```bash
# Install bash-completion package
sudo apt-get install bash-completion  # Debian/Ubuntu
sudo yum install bash-completion      # RHEL/CentOS
```

**macOS:**
```bash
# Install bash-completion via Homebrew
brew install bash-completion@2

# Add to ~/.bash_profile or ~/.bashrc:
export BASH_COMPLETION_COMPAT_DIR="/usr/local/etc/bash_completion.d"
[[ -r "/usr/local/etc/profile.d/bash_completion.sh" ]] && . "/usr/local/etc/profile.d/bash_completion.sh"
```

### Installation

**System-wide installation (Linux):**
```bash
# Requires root/sudo access
sudo gac completion bash > /etc/bash_completion.d/gac
```

**System-wide installation (macOS with Homebrew):**
```bash
gac completion bash > $(brew --prefix)/etc/bash_completion.d/gac
```

**User-specific installation:**
```bash
# Generate completion script
gac completion bash > ~/.gac-completion.bash

# Add to ~/.bashrc
echo 'source ~/.gac-completion.bash' >> ~/.bashrc

# Reload shell configuration
source ~/.bashrc
```

### Verification

After installation, restart your shell or reload your configuration:

```bash
source ~/.bashrc
```

Test completion by typing:
```bash
gac <TAB><TAB>
# Should show: audit, calendar, completion, config, group, user, version, etc.

gac user <TAB><TAB>
# Should show: create, list, suspend, unsuspend, update
```

## Zsh Completion

### Prerequisites

Ensure completion system is enabled in your `~/.zshrc`:

```bash
# Add to ~/.zshrc if not already present
autoload -U compinit
compinit
```

### Installation

**Option 1: Standard zsh completion directory**

```bash
# Find your zsh completion directory
echo ${fpath[1]}
# Common locations: /usr/local/share/zsh/site-functions or ~/.zsh/functions

# Generate completion script
gac completion zsh > "${fpath[1]}/_gac"

# Reload completions
autoload -U compinit && compinit
```

**Option 2: oh-my-zsh**

```bash
# Generate completion for oh-my-zsh
mkdir -p ~/.oh-my-zsh/completions
gac completion zsh > ~/.oh-my-zsh/completions/_gac

# Restart zsh
exec zsh
```

**Option 3: Custom completion directory**

```bash
# Create custom completion directory
mkdir -p ~/.zsh/completions

# Add to ~/.zshrc BEFORE compinit:
fpath=(~/.zsh/completions $fpath)
autoload -U compinit
compinit

# Generate completion script
gac completion zsh > ~/.zsh/completions/_gac

# Reload shell
exec zsh
```

### Troubleshooting Zsh

If completions don't work:

1. **Check permissions:**
   ```bash
   # Zsh may ignore world-writable directories
   chmod 755 ~/.oh-my-zsh/completions
   chmod 644 ~/.oh-my-zsh/completions/_gac
   ```

2. **Rebuild completion cache:**
   ```bash
   rm -f ~/.zcompdump
   compinit
   ```

3. **Check fpath:**
   ```bash
   # Verify completion directory is in fpath
   echo $fpath
   ```

### Verification

```bash
# Restart zsh
exec zsh

# Test completion
gac <TAB>
# Should show commands with descriptions

gac user <TAB>
# Should show subcommands with descriptions
```

## Fish Completion

### Prerequisites

Fish shell version 3.0 or higher recommended.

### Installation

Fish automatically loads completions from `~/.config/fish/completions/`:

```bash
# Create completions directory if it doesn't exist
mkdir -p ~/.config/fish/completions

# Generate completion script
gac completion fish > ~/.config/fish/completions/gac.fish
```

Fish will automatically load the completions immediately - no need to restart the shell!

### Verification

```bash
# Test completion
gac <TAB>
# Should show available commands with descriptions

gac user <TAB>
# Should show user subcommands
```

## Testing Completion

After installing completion for your shell, test the following scenarios:

### Command Completion

```bash
# Type partial command and press TAB
gac ver<TAB>
# Should complete to: gac version

gac comp<TAB>
# Should complete to: gac completion
```

### Subcommand Completion

```bash
# Type command and press TAB to see subcommands
gac user <TAB>
# Shows: create, list, suspend, unsuspend, update

gac group-settings <TAB>
# Shows: list, update
```

### Flag Completion

```bash
# Type command and double-dash, then TAB
gac user create --<TAB>
# Shows available flags: --email, --first-name, --groups, --last-name, etc.

gac --<TAB>
# Shows global flags: --domain, --client-secret, --verbose, --yes, etc.
```

## Troubleshooting

### Bash: Completions Not Working

**Problem:** Tab completion doesn't work after installation

**Solutions:**

1. **Verify bash-completion is installed:**
   ```bash
   # Check if bash-completion is installed
   type _init_completion
   # Should show: _init_completion is a function
   ```

2. **Check bash version:**
   ```bash
   bash --version
   # Bash 4.0+ recommended
   ```

3. **Reload bash configuration:**
   ```bash
   source ~/.bashrc
   # or
   source ~/.bash_profile
   ```

4. **Verify completion script is loaded:**
   ```bash
   type _gac_completion
   # Should show function definition
   ```

### Zsh: Completions Not Showing

**Problem:** Completions don't appear after installation

**Solutions:**

1. **Rebuild completion cache:**
   ```bash
   rm -f ~/.zcompdump*
   exec zsh
   ```

2. **Check file permissions:**
   ```bash
   ls -l ~/.oh-my-zsh/completions/_gac
   # Should be readable (644 or similar)
   ```

3. **Verify fpath includes completion directory:**
   ```bash
   echo $fpath | grep completions
   # Should show the directory containing _gac
   ```

4. **Enable verbose completion debugging:**
   ```bash
   zstyle ':completion:*' verbose yes
   ```

### Fish: Completion Script Not Loading

**Problem:** Fish doesn't load the completion script

**Solutions:**

1. **Verify file location:**
   ```bash
   ls -l ~/.config/fish/completions/gac.fish
   # File should exist
   ```

2. **Check for syntax errors:**
   ```bash
   fish --check ~/.config/fish/completions/gac.fish
   # Should have no output if syntax is correct
   ```

3. **Manually reload completions:**
   ```bash
   fish_update_completions
   ```

### General Issues

**Problem:** Partial completions or wrong suggestions

**Solution:** Regenerate the completion script with the latest version of `gac`:

```bash
# Bash
gac completion bash > ~/.gac-completion.bash
source ~/.gac-completion.bash

# Zsh
gac completion zsh > ~/.oh-my-zsh/completions/_gac
exec zsh

# Fish
gac completion fish > ~/.config/fish/completions/gac.fish
```

## Uninstalling Completion

### Bash

```bash
# Remove completion script
rm ~/.gac-completion.bash
# or
sudo rm /etc/bash_completion.d/gac

# Remove source line from ~/.bashrc
# (manually edit the file)
```

### Zsh

```bash
# Remove completion script
rm ~/.oh-my-zsh/completions/_gac
# or
rm "${fpath[1]}/_gac"

# Rebuild completion cache
rm -f ~/.zcompdump
compinit
```

### Fish

```bash
# Remove completion script
rm ~/.config/fish/completions/gac.fish
```

## Advanced Configuration

### Custom Completion Behavior

The `gac` completion script is generated by Cobra and includes:

- **Command completion** - All commands and subcommands
- **Flag completion** - All global and command-specific flags
- **Flag value completion** - Some flags have value suggestions
- **Description support** - Most shells show command/flag descriptions

### Completion in Scripts

Completions only work in interactive shells. They are not available in shell scripts.

### Completion Performance

For very large environments, completion may be slow if it queries live data. The `gac` completion is static (generated from command definitions), so it should be fast.

## Examples

### Complete Workflow Example (Bash)

```bash
# 1. Generate completion script
gac completion bash > ~/.gac-completion.bash

# 2. Add to shell configuration
echo 'source ~/.gac-completion.bash' >> ~/.bashrc

# 3. Reload configuration
source ~/.bashrc

# 4. Test completion
gac <TAB><TAB>
# Output: audit calendar completion config group ou transfer user version ...

gac user cr<TAB>
# Completes to: gac user create

gac user create --<TAB><TAB>
# Output: --email --first-name --groups --last-name ...
```

### Complete Workflow Example (Zsh with oh-my-zsh)

```bash
# 1. Generate completion script
mkdir -p ~/.oh-my-zsh/completions
gac completion zsh > ~/.oh-my-zsh/completions/_gac

# 2. Reload zsh
exec zsh

# 3. Test completion
gac <TAB>
# Output: Shows commands with descriptions

gac user <TAB>
# Output: Shows user subcommands with descriptions
```

### Complete Workflow Example (Fish)

```bash
# 1. Generate completion script
mkdir -p ~/.config/fish/completions
gac completion fish > ~/.config/fish/completions/gac.fish

# 2. Test completion (no reload needed!)
gac <TAB>
# Output: Shows available commands

gac user <TAB>
# Output: Shows user subcommands
```

## Related Resources

- [Command Reference](../reference/commands.md) - Complete command list
- [Installation](../installation.md) - Installation instructions
- [Troubleshooting](../reference/troubleshooting.md) - Common issues
- [Cobra Documentation](https://github.com/spf13/cobra/blob/main/shell_completions.md) - Shell completion details

## Need Help?

- Check the [Troubleshooting Guide](../reference/troubleshooting.md)
- Report issues on [GitHub](https://github.com/acockrell/google-admin-client/issues)
- Review the [FAQ](../reference/troubleshooting.md#frequently-asked-questions)
