Go ile yazdığım CLI tool: Ölü kodları ve kullanılmayan bağımlılıkları saniyeler içinde buluyor 👻


Merhaba,

Her projede zamanla şişkinlik oluşur.

Hiç çağrılmayan fonksiyonlar.
Altı ay önce eklenen ama kullanılmayan importlar.
Üç refactor'dan sağ kurtulan bağımlılıklar.

ghost bunların hepsini saniyeler içinde buluyor.

Ne yapıyor?
→ Ölü fonksiyonları tespit ediyor
→ Kullanılmayan importları buluyor
→ Kullanılmayan değişkenleri raporluyor
→ Hayalet bağımlılıkları yakalıyor

Dil desteği:
→ Go
→ JavaScript / TypeScript
→ Python

CI/CD entegrasyonu da mevcut, GitHub Actions ile direkt kullanabilirsiniz.

GitHub: github.com/ossiqn/ghost
Demo:   ossiqn.com.tr

Görüşlerinizi bekliyorum.


<div align="center">

```
  ██████  ██░ ██  ▒█████    ██████ ▄▄▄█████▓
▒██    ▒ ▓██░ ██▒▒██▒  ██▒▒██    ▒ ▓  ██▒ ▓▒
░ ▓██▄   ▒██▀▀██░▒██░  ██▒░ ▓██▄   ▒ ▓██░ ▒░
  ▒   ██▒░▓█ ░██ ▒██   ██░  ▒   ██▒░ ▓██▓ ░
▒██████▒▒░▓█▒░██▓░ ████▓▒░▒██████▒▒  ▒██▒ ░
```

# ghost 👻

**Dead code & unused dependency hunter**

[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-blue?style=flat)](LICENSE)
[![Release](https://img.shields.io/github/v/release/ossiqn/ghost?style=flat)](https://github.com/ossiqn/ghost/releases)
[![Stars](https://img.shields.io/github/stars/ossiqn/ghost?style=flat)](https://github.com/ossiqn/ghost/stargazers)

![ghost demo](assets/ghost-demo.gif)

</div>

---

## What is ghost?

Every codebase has them. Functions that were written but never called. Imports that made sense six months ago. Dependencies that somehow survived three refactors.

**ghost finds them all. In seconds.**

```bash
$ ghost scan .

👻 ghost scan complete
  47 files scanned

  Dead Functions:
────────────────────────────────────────────────────────────
  src/utils/helper.go:42
  └─ func calculateOldPrice()  → never called

  src/api/routes.go:89
  └─ func legacyHandler()      → never called

  Unused Imports:
────────────────────────────────────────────────────────────
  src/models/user.go:12
  └─ import "fmt"              → never used

  Ghost Dependencies:
────────────────────────────────────────────────────────────
  go.mod
  └─ github.com/old/package    → not imported anywhere
  └─ github.com/unused/lib     → not imported anywhere

────────────────────────────────────────────────────────────
  👻 2 dead functions
  👻 1 unused imports
  👻 0 unused variables
  👻 2 ghost dependencies

  💾 Est. size saved: ~2.4MB

Run 'ghost clean' to remove all? (y/n):
```

---

## Install

**Homebrew**
```bash
brew install ossiqn/tap/ghost
```

**Go Install**
```bash
go install github.com/ossiqn/ghost@latest
```

**Binary** → [Releases](https://github.com/ossiqn/ghost/releases)

---

## Usage

```bash
# Scan current directory
ghost scan .

# Scan specific path
ghost scan ./myproject

# Scan and auto clean
ghost scan . --clean

# JSON output (for CI/CD)
ghost scan . --json

# Specific language only
ghost scan . --lang go
ghost scan . --lang js
ghost scan . --lang python

# Verbose mode
ghost scan . --verbose
```

---

## Language Support

| Language | Dead Functions | Unused Imports | Unused Variables | Dependencies |
|----------|:--------------:|:--------------:|:----------------:|:------------:|
| Go       | ✅ | ✅ | ✅ | ✅ go.mod |
| JavaScript / TypeScript | ✅ | ✅ | ✅ | ✅ package.json |
| Python   | ✅ | ✅ | ✅ | ✅ requirements.txt |

---

## CI/CD Integration

**GitHub Actions**
```yaml
name: Ghost Scan

on: [push, pull_request]

jobs:
  ghost:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install ghost
        run: go install github.com/ossiqn/ghost@latest

      - name: Scan for dead code
        run: ghost scan . --json > ghost-report.json

      - name: Upload report
        uses: actions/upload-artifact@v4
        with:
          name: ghost-report
          path: ghost-report.json
```

---

## Why ghost?

| Tool | Dead Functions | Multi-Language | Dependencies | CLI |
|------|:--------------:|:--------------:|:------------:|:---:|
| **ghost** | ✅ | ✅ | ✅ | ✅ |
| ESLint | ❌ | JS only | ❌ | ✅ |
| gopls | Partial | Go only | ❌ | ❌ |
| depcheck | ❌ | JS only | ✅ | ✅ |

---

## Roadmap

- [x] Go support
- [x] JavaScript / TypeScript support  
- [x] Python support
- [ ] Rust support
- [ ] Ruby support
- [ ] Auto-fix mode
- [ ] VS Code extension
- [ ] Web dashboard

---

## Contributing

```bash
git clone https://github.com/ossiqn/ghost
cd ghost
go mod tidy
go test ./...
```

PRs are welcome. Open an issue first for major changes.

---

<div align="center">

made with 🖤 by [ossiqn](https://github.com/ossiqn)

</div>
