# llmvlt

A CLI secret manager built for AI/ML engineers.

`llmvlt` knows what `OPENAI_API_KEY` is. It knows what `ANTHROPIC_API_KEY` looks like. It warns you when a key format is wrong, tracks which keys were active during an experiment, and injects them safely into your scripts — without ever touching your shell history.

```
$ llmvlt init --preset full-llm-stack
✓ Vault initialized (.llmvlt.store)
✓ Preset: full-llm-stack — All major LLM provider credentials combined

6 secrets to fill in:

  OPENAI_API_KEY                      (required) Your OpenAI API key
  ANTHROPIC_API_KEY                   (required) Your Anthropic Claude API key
  HF_TOKEN                            (required) Your Hugging Face access token
  REPLICATE_API_TOKEN                 (required) Your Replicate API token
  WANDB_API_KEY                       (required) Your W&B API key
  LANGCHAIN_API_KEY                   (optional) Your LangSmith API key

$ llmvlt set OPENAI_API_KEY sk-...
✓ OPENAI_API_KEY set (version 1)

$ llmvlt run -- python train.py
```

---

## Why not just use `.env` files?

`.env` files get committed. API keys get leaked. It happens to everyone.

`llmvlt` keeps your secrets in an AES-256 encrypted vault on disk — never in plaintext, never in your shell history, never accidentally pushed to GitHub. And unlike generic secret managers, it understands the specific keys that AI/ML engineers use every day.

---

## Installation

### macOS

```bash
brew install yourname/tap/llmvlt
```

### Linux

```bash
curl -sSL https://get.llmvlt.dev | sh
```

### Windows

Download the latest `.exe` from the [releases page](https://github.com/yourname/llmvlt/releases) and add it to your PATH.

Or, if you have Go installed, build it yourself (see [Building from source](#building-from-source)).

---

## Getting Started

**1. Initialize a vault in your project directory:**

```bash
llmvlt init --preset openai-stack
```

**2. Set your secrets:**

```bash
llmvlt set OPENAI_API_KEY sk-...
llmvlt set OPENAI_ORG_ID org-...
```

**3. Run your script with secrets injected:**

```bash
llmvlt run -- python train.py
```

Secrets are injected into the subprocess only — they never appear in your parent shell, in `ps aux` output, or in your shell history.

---

## Presets

Presets are the core feature that makes `llmvlt` different. Each preset knows which keys a provider uses, what they look like, and how often they should be rotated.

| Preset | Description |
|---|---|
| `openai-stack` | OpenAI API credentials |
| `anthropic-stack` | Anthropic Claude API credentials |
| `huggingface-stack` | Hugging Face Hub credentials |
| `wandb-stack` | Weights & Biases experiment tracking |
| `replicate-stack` | Replicate API credentials |
| `langchain-stack` | LangChain and LangSmith tracing |
| `full-llm-stack` | All of the above combined |

List all available presets:

```bash
llmvlt preset list
```

You can also define your own presets in `~/.llmvlt/presets/`:

```yaml
# ~/.llmvlt/presets/my-company.yaml
name: my-company-stack
description: "Internal API credentials"
secrets:
  - key: MY_INTERNAL_API_KEY
    description: "Internal service key"
    required: true
    validation:
      pattern: "^int-[a-zA-Z0-9]{32}$"
      hint: "Should start with 'int-'"
```

---

## Commands

### `llmvlt init`

Initialize a new vault in the current directory.

```bash
llmvlt init                          # empty vault
llmvlt init --preset openai-stack    # with a preset
```

### `llmvlt set`

Store a secret. Validates the format against the active preset and warns if something looks wrong — but always stores the value regardless.

```bash
llmvlt set OPENAI_API_KEY sk-...
llmvlt set HF_TOKEN hf_...
```

### `llmvlt get`

Retrieve a secret value.

```bash
llmvlt get OPENAI_API_KEY
```

### `llmvlt list`

List all secret names stored in the vault. Never prints values.

```bash
llmvlt list
```

### `llmvlt run`

Inject all secrets into a subprocess. The hero command.

```bash
llmvlt run -- python train.py
llmvlt run -- jupyter notebook
llmvlt run -- pytest tests/
```

Secrets are available as environment variables inside the subprocess. They are never exposed to the parent shell.

### `llmvlt inject`

Export secrets in different formats for use in other contexts.

```bash
# Inject into current shell session
eval $(llmvlt inject)

# Generate a Jupyter cell to paste into a notebook
llmvlt inject --format jupyter

# Write a .env file (use with caution)
llmvlt inject --format dotenv --out .env
```

### `llmvlt check`

Check for keys that haven't been rotated recently.

```bash
llmvlt check

# Example output:
# ⚠  OPENAI_API_KEY — last rotated 94 days ago (recommended: 90 days)
# ✓  ANTHROPIC_API_KEY — rotated 12 days ago
# ✓  HF_TOKEN — rotated 45 days ago
```

### `llmvlt use`

Switch between providers when benchmarking the same script across multiple LLMs.

```bash
llmvlt use anthropic    # activates Anthropic keys
llmvlt use openai       # switches to OpenAI keys
```

### `llmvlt history`

Show which key versions were active during past runs (requires `--tag` on `run`).

```bash
llmvlt run --tag "gpt4-baseline" -- python eval.py

llmvlt history
# 2026-03-12 10:00  gpt4-baseline     OPENAI_API_KEY@v3
# 2026-03-11 18:22  claude-comparison ANTHROPIC_API_KEY@v1
```

---

## Building from Source

You need [Go 1.22+](https://go.dev/dl/) installed.

### macOS / Linux

```bash
git clone https://github.com/yourname/llmvlt
cd llmvlt
go build -o llmvlt .
```

### Windows (PowerShell or Command Prompt)

```powershell
git clone https://github.com/yourname/llmvlt
cd llmvlt
go build -o llmvlt.exe .
.\llmvlt.exe --help
```

If Go is not installed on Windows, download the `.msi` installer from [go.dev/dl](https://go.dev/dl). It adds Go to your PATH automatically — no manual setup needed.

### Cross-compiling from any platform

Go has built-in cross-compilation. No extra tools required.

**From macOS or Linux:**

```bash
# Windows 64-bit
GOOS=windows GOARCH=amd64 go build -o dist/llmvlt-windows-amd64.exe .

# Windows ARM (Surface, newer laptops)
GOOS=windows GOARCH=arm64 go build -o dist/llmvlt-windows-arm64.exe .

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o dist/llmvlt-darwin-arm64 .

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o dist/llmvlt-darwin-amd64 .

# Linux 64-bit
GOOS=linux GOARCH=amd64 go build -o dist/llmvlt-linux-amd64 .
```

**From Windows (PowerShell):**

```powershell
# Linux 64-bit
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o dist/llmvlt-linux-amd64 .

# macOS Apple Silicon
$env:GOOS="darwin"; $env:GOARCH="arm64"; go build -o dist/llmvlt-darwin-arm64 .
```

Note: on PowerShell, env vars use `$env:VAR="value"` syntax instead of the Unix `VAR=value` prefix.

**Build all platforms at once:**

```bash
#!/bin/bash
# build-all.sh
VERSION=$(git describe --tags --always --dirty)

platforms=(
  "windows/amd64/.exe"
  "windows/arm64/.exe"
  "darwin/amd64/"
  "darwin/arm64/"
  "linux/amd64/"
  "linux/arm64/"
)

for platform in "${platforms[@]}"; do
  IFS='/' read -r os arch ext <<< "$platform"
  output="dist/llmvlt-${VERSION}-${os}-${arch}${ext}"
  echo "Building $output..."
  GOOS=$os GOARCH=$arch go build \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o "$output" .
done

echo "Done."
```

The `-ldflags="-s -w"` flag strips debug symbols and reduces binary size by ~30%.

---

## Security Model

- Secrets are encrypted with **AES-256** using a key derived from your master password via **Argon2id**
- The vault is a single encrypted file (`.llmvlt.store`) — never plaintext
- `llmvlt run` forks a subprocess: secrets are injected into the child process only, never the parent shell
- Secret values are never written to shell history
- File permissions on the store are set to `0600` (owner read/write only)

Built on top of [scrt](https://github.com/loderunner/scrt)'s proven encryption layer.

---

## Roadmap

- [ ] Cloud sync across machines
- [ ] Team vaults with access control
- [ ] Key rotation automation
- [ ] GitHub Actions integration
- [ ] More presets (Cohere, Mistral, Together AI, Groq)

---

## License

Apache 2.0