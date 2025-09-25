# ccstatus-go 🎭

[![Go Version](https://img.shields.io/github/go-mod/go-version/mirage20/ccstatus-go)](go.mod)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/mirage20/ccstatus-go)](https://github.com/mirage20/ccstatus-go/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/mirage20/ccstatus-go)](https://goreportcard.com/report/github.com/mirage20/ccstatus-go)
[![Build Status](https://github.com/mirage20/ccstatus-go/workflows/CI/badge.svg)](https://github.com/mirage20/ccstatus-go/actions)

> *One of the world's status lines ever made. Definitely exists. Technically functional.*

A magnificently over-engineered status line generator for Claude Code that proves once and for all that displaying a handful of information requires dozens of Go files, multiple layers of abstraction, and a plugin architecture that would make enterprise Java developers weep with envy.

## What Is This Monstrosity?

Remember when you just wanted to know which Claude model you're using? Well, we took that simple desire and turned it into a provider-component architecture with self-registering plugins, parallel goroutines, TTL-based caching, and YAML configuration files. Because why use `printf` when you can have **ENTERPRISE GRADE STATUS LINES**.

## Features That Nobody Asked For

- 🏗️ **Provider-Component Architecture**: Because MVC wasn't abstract enough
- ⚡ **Parallel Data Fetching**: Multiple data sources run simultaneously! Feel the speed!
- 💾 **Multi-Tier Caching**: File-based and null caching strategies for when you really need to cache that model name
- 🔌 **Self-Registering Plugins**: Components that register themselves via `init()` because explicit is for chumps
- 🎨 **ANSI Color Support**: 7 whole colors to choose from!
- 📝 **Hundreds of Lines of Configuration**: More documentation than actual config
- 🎯 **Koanf Integration**: Because `json.Unmarshal` is too mainstream

## The Numbers Don't Lie

- **Dozens of Go files** for displaying a few pieces of information
- **Multiple providers** fetching data in parallel (it's basically distributed computing)
- **Numerous components** each with their own configuration
- **Infinite customization** via YAML (hundreds of lines of possibilities)
- **Actual performance improvements** through caching (take that, bash one-liners!)

## Installation

### Option 1: Download from GitHub Releases (For Those Who Value Their Time)

First, grab the appropriate binary from [releases](https://github.com/mirage20/ccstatus-go/releases/latest) (all 5MB of enterprise-grade architecture):

- **macOS (Intel)**: `ccstatus-darwin-amd64` - For those still rocking x86
- **macOS (Apple Silicon)**: `ccstatus-darwin-arm64` - Living in the future
- **Linux (x64)**: `ccstatus-linux-amd64` - The penguin's choice
- **Linux (ARM)**: `ccstatus-linux-arm64` - For your Raspberry Pi status needs
- **Windows (x64)**: `ccstatus-windows-amd64.exe` - The path of suffering
- **Windows (ARM)**: `ccstatus-windows-arm64.exe` - Suffering, but modern

#### macOS/Linux Installation (The Enlightened Path)

```bash
# Bestow upon it a proper name (shed the platform-specific suffix)
mv ccstatus-darwin-arm64 ccstatus  # or whatever platform you downloaded

# Make it executable (grant it the power of execution)
chmod +x ccstatus

# Install globally (because greatness should be system-wide)
sudo mv ccstatus /usr/local/bin/
```

Then tell Claude Code about your life-changing decision by editing `~/.claude/settings.json`:
   ```json
   {
     "statusLine": {
       "type": "command",
       "command": "ccstatus"
     }
   }
   ```

#### Windows Installation (The Path of Suffering)

1. Rename it to something civilized:
   ```cmd
   ren ccstatus-windows-amd64.exe ccstatus.exe
   ```

2. Put it somewhere sensible (like `C:\tools\ccstatus.exe`)

3. Update your Claude settings at `%USERPROFILE%\.claude\settings.json`:
   ```json
   {
     "statusLine": {
       "type": "command",
       "command": "C:\\tools\\ccstatus.exe"
     }
   }
   ```

### Option 2: Build from Source (For the Adventurous)

```bash
# Clone this monument to over-engineering
git clone https://github.com/mirage20/ccstatus-go
cd ccstatus-go

# Summon the Go compiler to witness our hubris
make build

# Bestow this binary upon your system
sudo cp build/ccstatus /usr/local/bin/

# Or prove your worth with lint AND build
make all  # Runs linting AND builds!
```

## Usage

Once configured in Claude Code settings, ccstatus runs automatically every time your conversation updates!

## Configuration

The included [`config.yaml`](config.yaml) file contains **HUNDREDS** of lines documenting every possible configuration option with painstaking detail. It's your definitive guide to customizing this architectural marvel. Start there. Seriously, we documented *everything*.

Place your actual configuration in one of these locations (because choice is important):
- `.claude/ccstatus.local.yaml` - For your secret local configs
- `.claude/ccstatus.yaml` - For configs you're willing to share
- `~/.claude/ccstatus.yaml` - For when you want consistency across projects

Example configuration that definitely isn't overkill:

```yaml
# Choose your fighters (components)
active: ["model", "context", "activeblockusage", "activeblocktime"]

# Customize that separator like your life depends on it
separator:
  symbol: " | "  # Revolutionary
  color: "gray"   # Cutting edge

# Cache configuration for maximum enterprise
cache:
  dir: "/tmp/ccstatus-cache-of-doom"

# Component-specific settings with more options than a luxury car
components:
  model:
    template: "{{.Icon}} {{.ShortName}}"
    icons:
      opus: "🎭"
      sonnet: "✨"
      haiku: "🍃"
    colors:
      opus: "red"
      sonnet: "yellow"
      haiku: "green"
```

Pro tip: The [`config.yaml`](config.yaml) file has **400 lines** of configuration documentation. That's roughly 395 more lines than you need, but we're nothing if not thorough.

## Architecture (Yes, We Have One)

```
┌─────────────────────────────────────────┐
│           Claude Session JSON           │
└─────────────────┬───────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────┐
│         StatusLine Orchestrator         │
│     (The Conductor of This Symphony)    │
└─────────────────┬───────────────────────┘
                  │
         ┌────────┴────────┐
         ▼                 ▼
┌──────────────┐  ┌──────────────────┐
│   Providers  │  │    Components    │
│  (The Wise)  │  │ (The Beautiful)  │
└──────────────┘  └──────────────────┘
         │                 │
         └────────┬────────┘
                  │
                  ▼
        ╔═══════════════════╗
        ║  YOUR STATUS LINE ║
        ╚═══════════════════╝
```

## Why Though?

Because sometimes you need to transform this:
```json
{"model": {"display_name": "Opus 4.1"}, "cost": {"total_lines_added": 67}}
```

Into this:
```
🎭 Opus | 22k | +67-12 | v1.0.89
```

And obviously, that requires:
- A registry pattern for components
- Parallel provider execution
- Template-based rendering
- Hierarchical configuration management
- Abstract caching layers
- And thoroughly documented configuration options

## Performance

We've benchmarked this extensively and can confidently say it runs. Every time.

## Contributing

Feel free to add more abstraction layers. We're particularly interested in:
- A GraphQL API for status queries
- Kubernetes operator for status line orchestration
- Machine learning to predict what status you want to see
- Blockchain integration for immutable status history
- WebAssembly port for browser-based status lines

## License

This project is licensed under the MIT License - because even over-engineered status lines deserve freedom. See the [LICENSE](LICENSE) file for the legally binding version that lawyers actually understand.

In simpler terms: Use it, break it, fork it, ship it. Just don't blame us when your status line becomes sentient.

## Acknowledgments

- To all the simple bash scripts that could have done this in 5 lines
- To the `echo` command, forever in our hearts
- To YAGNI, whom we've never met
- To the hundreds of lines of configuration comments that nobody will read

---

*"It's not over-engineering if it works"* - Ancient Proverb (citation needed)

*Built with 🎭 and an unreasonable amount of Go packages*
