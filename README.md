by ossiqn

text


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

Every codebase has them.

Functions that were written but never called.
Imports that made sense six months ago.
Dependencies that somehow survived three refactors.

**ghost finds them all. In seconds.**
$ ghost scan .

👻 ghost scan complete
47 files scanned

Dead Functions:
────────────────────────────────────────────────────────────
src/utils/helper.go:42
└─ func calculateOldPrice() → never called

src/api/routes.go:89
└─ func legacyHandler() → never called

Unused Imports:
────────────────────────────────────────────────────────────
src/models/user.go:12
└─ import "fmt" → never used

Ghost Dependencies:
────────────────────────────────────────────────────────────
go.mod
└─ github.com/old/package → not imported anywhere
└─ github.com/unused/lib → not imported anywhere

────────────────────────────────────────────────────────────
👻 2 dead functions
👻 1 unused imports
👻 0 unused variables
👻 2 ghost dependencies

💾 Est. size saved: ~2.4MB

Run ghost clean to remove all? (y/n):

text


---

## Install

**Go Install**
go install github.com/ossiqn/ghost@latest

text


**Binary**

Download from [Releases](https://github.com/ossiqn/ghost/releases)

---

## Usage
ghost scan . scan current directory
ghost scan ./myproject scan specific path
ghost scan . --clean scan and auto clean
ghost scan . --json output as JSON
ghost scan . --lang go scan only Go files
ghost scan . --lang js scan only JS/TS files
ghost scan . --lang python scan only Python files
ghost scan . --lang rust scan only Rust files
ghost scan . --lang ruby scan only Ruby files
ghost scan . --lang java scan only Java files
ghost scan . --verbose verbose output
ghost version print version

text


---

## Language Support

| Language           | Dead Functions | Unused Imports | Unused Variables | Dependencies     |
|--------------------|:--------------:|:--------------:|:----------------:|:----------------:|
| Go                 | ✅             | ✅             | ✅               | ✅ go.mod        |
| JavaScript         | ✅             | ✅             | ✅               | ✅ package.json  |
| TypeScript         | ✅             | ✅             | ✅               | ✅ package.json  |
| Python             | ✅             | ✅             | ✅               | ✅ requirements  |
| Rust               | ✅             | ✅             | ✅               | ✅ Cargo.toml    |
| Ruby               | ✅             | ✅             | 🔜               | 🔜 Gemfile       |
| Java               | ✅             | ✅             | 🔜               | 🔜 pom.xml       |

---

## CI/CD Integration

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

      - name: Scan
        run: ghost scan . --json > ghost-report.json

      - name: Upload report
        uses: actions/upload-artifact@v4
        with:
          name: ghost-report
          path: ghost-report.json
Why ghost?
Tool	Dead Functions	Multi-Language	Dependencies	CLI
ghost	✅	✅	✅	✅
ESLint	❌	JS only	❌	✅
gopls	Partial	Go only	❌	❌
depcheck	❌	JS only	✅	✅
vulture	✅	Python only	❌	✅
Roadmap
 Go support
 JavaScript / TypeScript support
 Python support
 Rust support
 Ruby support
 Java support
 C# support
 PHP support
 Auto-fix mode
 VS Code extension
 Web dashboard
Contributing
text

git clone https://github.com/ossiqn/ghost
cd ghost
go mod tidy
go test ./...
PRs are welcome.
Open an issue first for major changes.

<div align="center">
made with 🖤 by ossiqn

ossiqn.com.tr

</div> ```
