# hf-gguf-finder

Find GGUF model files on Hugging Face Hub from the command line.

## Why this tool?

LLM inference engines like [llama.cpp](https://github.com/ggerganov/llama.cpp), [Ollama](https://ollama.com), and [LM Studio](https://lmstudio.ai) use GGUF format for quantized models. Hugging Face hosts thousands of GGUF files across hundreds of repositories, but finding them usually means browsing the website manually or digging through API responses.

This CLI tool queries the Hugging Face API directly and gives you **download URLs** for every GGUF file matching your search — in plain text or structured JSON.

```
$ hf-gguf-finder -q deepseek-coder -limit 5
https://huggingface.co/deepseek-ai/deepseek-coder-6.7b-instruct-GGUF/resolve/main/deepseek-coder-6.7b-instruct.Q4_K_M.gguf
https://huggingface.co/deepseek-ai/deepseek-coder-6.7b-instruct-GGUF/resolve/main/deepseek-coder-6.7b-instruct.Q5_K_M.gguf
...
```

## Installation

### Option 1: Install with Go (recommended)

Requires [Go 1.21+](https://go.dev/dl/).

```bash
go install github.com/jeanmachuca/hf-gguf-finder@latest
```

The binary is placed in `$GOPATH/bin` (or `~/go/bin` by default). Make sure it's in your PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Then run from anywhere:

```bash
hf-gguf-finder -q llama
```

### Option 2: Clone and build

```bash
git clone git@github.com:jeanmachuca/hf-gguf-finder.git
cd hf-gguf-finder
go build -o hf-gguf-finder

# (Optional) move to PATH
mv hf-gguf-finder /usr/local/bin/
```

### Option 3: Run directly (no install)

```bash
go run github.com/jeanmachuca/hf-gguf-finder@latest -q llama
```

Or from a local clone:

```bash
cd hf-gguf-finder
go run main.go -q llama
```

## Tutorial

For a step-by-step walkthrough focused on finding small models for edge
devices (Raspberry Pi, phones, laptops), see:

➡️ **[TUTORIAL.md](TUTORIAL.md)**

## Quick run: gguf-run

`gguf-run` is a companion tool that ties everything together:
**search → download → run** with [llama.cpp](https://github.com/ggerganov/llama.cpp)
in one command.

### Install

```bash
# From source
go build -o gguf-run ./cmd/gguf-run/
sudo mv gguf-run /usr/local/bin/
```

```bash
# Or run without installing
go run ./cmd/gguf-run/ -q tinyllama
```

### Usage

```bash
# Auto-installs llama.cpp if missing, searches for best Q4_K_M file,
# downloads, and launches an interactive chat
gguf-run -q tinyllama

# Single-shot prompt
gguf-run -q qwen2.5-1.5b -p "What is the capital of France?"

# Use a local file or URL
gguf-run -m ~/models/my-model.q4_k_m.gguf
gguf-run -m https://huggingface.co/org/model/resolve/main/file.gguf

# Pass extra llama.cpp flags (use -- separator)
gguf-run -q phi -- --temp 0.8 --ctx-size 4096 -ngl 999

# Override cache directory
gguf-run -q smollm --cache-dir /ssd/models
```

Models are cached in `~/.cache/gguf/` so repeated runs skip the download.

## Usage

```
hf-gguf-finder -q <query> [-limit <n>] [-json]
```

| Flag | Default | Description |
|------|---------|-------------|
| `-q`  | (required) | Search query — searches model names and descriptions on Hugging Face |
| `-limit` | `20` | Maximum number of models to return |
| `-json` | `false` | Output results as structured JSON instead of plain text |

### Examples

**Basic search — plain text URLs:**
```bash
hf-gguf-finder -q llama
```

**Search with multi-word query and custom limit:**
```bash
hf-gguf-finder -q "deepseek coder" -limit 50
```

**JSON output:**
```bash
hf-gguf-finder -q phi -json
```

**JSON output with custom limit:**
```bash
hf-gguf-finder -q "mistral 7b" -limit 10 -json
```

### Output formats

#### Plain text (default)
Each GGUF file's download URL is printed on a separate line:

```
https://huggingface.co/org/model/resolve/main/file.Q4_K_M.gguf
https://huggingface.co/org/model/resolve/main/file.Q5_K_M.gguf
https://huggingface.co/org/model/resolve/main/file.Q8_0.gguf
```

#### JSON (`-json`)
Structured output grouped by model:

```json
{
  "query": "phi",
  "limit": 5,
  "results": [
    {
      "modelId": "microsoft/Phi-3-mini-4k-instruct-gguf",
      "modelName": "microsoft/Phi-3-mini-4k-instruct-gguf",
      "files": [
        "https://huggingface.co/microsoft/Phi-3-mini-4k-instruct-gguf/resolve/main/Phi-3-mini-4k-instruct.Q4_K_M.gguf",
        "https://huggingface.co/microsoft/Phi-3-mini-4k-instruct-gguf/resolve/main/Phi-3-mini-4k-instruct.Q8_0.gguf"
      ]
    }
  ]
}
```

## How it works

1. Sends a request to `https://huggingface.co/api/models` with `library=gguf`, your search query, and a sort by downloads.
2. Decodes the JSON response.
3. Filters sibling files to those ending in `.gguf`.
4. Constructs download URLs in the format `https://huggingface.co/{modelId}/resolve/main/{filename}`.
5. Prints URLs (plain text) or a structured JSON object.

## License

MIT
