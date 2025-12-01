# Github Receipt 🧾

A receipt-style Github receipt designed to work perfectly as an embeddable iframe, inspired by [git-receipt.com](https://git-receipt.com).

## The Story Behind This Project

I discovered git-receipt.com and absolutely loved the concept of displaying Github stats in a creative, receipt-style format. However, when I tried to embed it on my personal website as an iframe, it didn't work as expected.

According to the official website, the original repository was at [gitreceipt](https://github.com/ankitkr0/gitreceipt), but it appears the creator has either deleted their profile or changed their username, making the original project inaccessible.

So I decided to build my own version from scratch that:

- Maintains the aesthetic I loved from the original
- Works perfectly as an embeddable iframe
- Includes smart caching to handle Github API limits
- Is actively maintained and easy to deploy

## Features

- **🖼️ Iframe Ready** - Built from the ground up to work perfectly when embedded
- **⚡ High Performance** - Smart caching with Ristretto for fast responses
- **📊 Real-time Metrics** - Cache statistics and performance monitoring
- **🔐 Github API Support** - Works with or without authentication tokens
- **🚀 Easy Deployment** - Simple setup with Tailscale Funnel or any cloud platform

## Quick Start

### Prerequisites

- Go 1.24 or higher
- Github token (optional, for higher API limits)

## Usage

To use the Github Receipt on your website, simply embed it using the following URL format:

```html
<iframe
  src="https://santiirepair.tail012146.ts.net/:username"
  width="400"
  height="600"
  frameborder="0"
  style="border: 1px solid #e1e4e8; border-radius: 6px;"
>
</iframe>
```

Replace :username with the Github username you want to display. For example:

```html
<iframe
  src="https://santiirepair.tail012146.ts.net/octocat"
  width="400"
  height="600"
  frameborder="0"
>
</iframe>
```

Or if you want to see what it looks like, let's look at it using the [Octocat](https://github.com/octocat) as an [example](https://santiirepair.tail012146.ts.net/octocat).
