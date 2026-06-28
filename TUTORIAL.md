# Tutorial: Finding Edge‑Ready LLMs with hf-gguf-finder

A hands-on walkthrough for finding small, quantized models that run on
Raspberry Pi, phones, laptops, or any low-capacity device.

---

## What you'll learn

- How GGUF quantization works and why it matters for edge devices
- Which search queries find the smallest models
- How to read file names to pick the right quant
- An end-to-end workflow: search → download → run with llama.cpp

## Prerequisites

- [hf-gguf-finder installed](README.md#installation)
- A terminal

---

## 1. Why small models?

Most people reach for 7B–70B parameter models, but those need 4–40 GB of
RAM even after quantization. For edge devices you want models with
**≤ 3B parameters**, quantized to 4 or 5 bits — they need only
0.5–2 GB and run at usable speeds on a CPU.

hf-gguf-finder lets you search specifically for these tiny models
without clicking through web pages.

---

## 2. Understanding GGUF file names

GGUF files encode their capabilities in the name. A file like:

```
Qwen2.5-1.5B-Instruct.Q4_K_M.gguf
```

| Part | Meaning |
|------|---------|
| `Qwen2.5-1.5B` | Model architecture + parameter count (1.5B) |
| `Instruct` | Fine-tuned for instruction following |
| `Q4_K_M` | 4-bit quantization, K-quant medium |

### Quantization levels at a glance

| Suffix | Bits | Size vs FP16 | Quality | Best for |
|--------|------|--------------|---------|----------|
| `Q2_K` | 2 | ~12 % | Low | Extreme memory limits |
| `Q3_K_M` | 3 | ~17 % | Fair | <512 MB devices |
| `Q4_K_M` | 4 | ~22 % | Good | **Edge sweet spot** |
| `Q5_K_M` | 5 | ~28 % | Very good | Slightly more RAM |
| `Q8_0` | 8 | ~44 % | Excellent | ~2 GB+ devices |
| `F16` | 16 | 100 % | Original | Desktops only |

**Rule of thumb:** start with `Q4_K_M` — it offers the best balance of
size and quality for edge devices.

---

## 3. The search query table

These queries consistently return models that fit on small devices:

| Query | What it finds | Typical size |
|-------|---------------|--------------|
| `-q "1.5b"` | Any 1.5B param model | 0.9–1.2 GB Q4 |
| `-q "tinyllama"` | TinyLlama 1.1B | 0.6–0.8 GB |
| `-q "phi"` | Phi-3 (3.8B), Phi-2 (2.7B) | 1.5–2.5 GB |
| `-q "smollm"` | SmolLM 135M–1.7B | 0.1–1.0 GB |
| `-q "qwen2.5-1.5b"` | Qwen2.5-1.5B series | 0.9–1.2 GB |
| `-q "gemma-2b"` | Gemma-2-2B / 2B | 1.2–1.5 GB |
| `-q "stablelm-2-1.6b"` | StableLM-2 1.6B | 0.9–1.1 GB |
| `-q "q4_k_m"` | Models with Q4_K_M files | varies |
| `-q "mobile"` | MobileLLM / edge-optimized | 0.3–1.5 GB |

---

## 4. Hands-on tutorial

### Step 1 — Find the smallest models

```bash
go run main.go -q "smollm" -json
```

Output shows SmolLM models at 135M, 360M, and 1.7B params. These are
the smallest viable LLMs — they run on a Raspberry Pi 4.

### Step 2 — Search for 1.5 B sweet spot

```bash
go run main.go -q "qwen2.5-1.5b" -limit 10 -json
```

Qwen2.5-1.5B is widely regarded as the best quality-per-parameter model
at this size. Look for the `Instruct` variants.

### Step 3 — Filter to Q4_K_M only (smallest usable)

```bash
go run main.go -q "q4_k_m qwen2.5" -json
```

This narrows results to models that already have a Q4_K_M quantization
file — no need to quantize yourself.

### Step 4 — Download and run

Take a URL from the output and download it:

```bash
# Download
wget https://huggingface.co/Qwen/Qwen2.5-1.5B-Instruct-GGUF/resolve/main/qwen2.5-1.5b-instruct-q4_k_m.gguf

# Run with llama.cpp
./llama-cli -m qwen2.5-1.5b-instruct-q4_k_m.gguf -p "Hello, who are you?" -n 100
```

---

## 5. Scripting with JSON output

The `-json` flag makes the tool pipeable:

```bash
# Get the first URL from every result
go run main.go -q tinyllama -json | jq -r '.results[].files[]'
```

```bash
# Count total GGUF files found
go run main.go -q phi -json | jq '.results | map(.files | length) | add'
```

```bash
# Save to file and process later
go run main.go -q "gemma-2b" -json > models.json
```

---

## 6. Real-world edge device recipes

### Raspberry Pi 4 (4 GB RAM)

```bash
go run main.go -q smollm -json
# Pick SmolLM-360M-Instruct Q4_K_M (~200 MB)
```

### Smartphone / Tablet

```bash
go run main.go -q tinyllama -json
# Pick TinyLlama-1.1B Q4_K_M (~700 MB)
```

### Laptop (8 GB RAM)

```bash
go run main.go -q qwen2.5-1.5b -json
# Pick Qwen2.5-1.5B-Instruct Q4_K_M (~1 GB)
```

### Chromebook / Low-end PC

```bash
go run main.go -q "stablelm-2-1.6b" -json
# Pick StableLM-2-1.6B Q4_K_M (~1 GB)
```

---

## 7. Pro tips

- **Add `-limit 50`** to see more options — many models have 10+ quant files
- **Combine terms** — `go run main.go -q "instruct q4_k_m 1.5b"` narrows results
- **Use `-json` with `jq`** to extract just the URLs and pipe to `wget` or `curl`
- **Check the "siblings"** in JSON output — a single model often has 5–15 quant variants
- **Prefer "Instruct" or "Chat" variants** for interactive use

---

## 8. One-command runner: gguf-run

The companion script `gguf-run` automates the full workflow:

```bash
# Search, download best Q4_K_M, and chat
./gguf-run -q tinyllama

# Single prompt, no interaction
./gguf-run -q qwen2.5-1.5b -p "Hello!"

# Extra llama.cpp options
./gguf-run -q phi -- --temp 0.8 -ngl 999
```

It auto-installs llama.cpp via Homebrew if missing and caches
downloaded models in `~/.cache/gguf/`.

## Next steps

- Build [llama.cpp](https://github.com/ggerganov/llama.cpp) for your device
- Try [Ollama](https://ollama.com) — it handles quantization and serving
- Experiment with different quantization levels on the same model
