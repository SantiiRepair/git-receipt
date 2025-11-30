# GitHub Receipt 🧾

A receipt-style GitHub receipt generator designed to work perfectly as an embeddable iframe, inspired by [git-receipt.com](https://git-receipt.com).

## What This Solves

While I loved the design of git-receipt.com, I noticed it couldn't be easily embedded as an iframe on other websites. So I built this version that maintains the aesthetic while being fully iframe-compatible. The design isn't an exact copy—I did my best to recreate the visual style while making it work seamlessly in embedded contexts.

## Features

- **🖼️ Iframe Ready** - Built from the ground up to work perfectly when embedded
- **⚡ High Performance** - Smart caching with Ristretto for fast responses
- **📊 Real-time Metrics** - Cache statistics and performance monitoring
- **🔐 GitHub API Support** - Works with or without authentication tokens
- **🚀 Easy Deployment** - Simple setup with Tailscale Funnel or any cloud platform

## Quick Start

### Prerequisites

- Go 1.24 or higher
- GitHub token (optional, for higher API limits)
