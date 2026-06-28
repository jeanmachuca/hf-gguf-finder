package main

/*
Hugging Face GGUF Model Finder

Searches the Hugging Face Hub for GGUF format models and prints
download URLs.  By default output is plain text (one URL per line).
Use the -json flag to get structured JSON output.

Flags:
  -q <query>      Search query (required).  Example: "llama", "deepseek coder"
  -limit <n>      Max models to return (default 20).
  -json           Output results as JSON instead of plain text.

Output formats:

  Plain text (default):
    https://huggingface.co/org/model/resolve/main/model.q4_K_M.gguf
    https://huggingface.co/org/model/resolve/main/model.q8_0.gguf
    ...

  JSON (-json flag):
    {
      "query": "llama",
      "limit": 20,
      "results": [
        {
          "modelId": "org/model",
          "modelName": "org/model",
          "files": [
            "https://huggingface.co/org/model/resolve/main/model.q4_K_M.gguf",
            "https://huggingface.co/org/model/resolve/main/model.q8_0.gguf"
          ]
        }
      ]
    }

Examples:
  go run main.go -q llama
  go run main.go -q "deepseek coder" -limit 50
  go run main.go -q phi -json
  go run main.go -q llama -limit 10 -json
*/

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Model struct {
	ID       string    `json:"id"`
	ModelID  string    `json:"modelId"`
	Siblings []Sibling `json:"siblings"`
}

type Sibling struct {
	RFilename string `json:"rfilename"`
}

type Result struct {
	ModelID   string   `json:"modelId"`
	ModelName string   `json:"modelName"`
	Files     []string `json:"files"`
}

type Output struct {
	Query   string   `json:"query"`
	Limit   int      `json:"limit"`
	Results []Result `json:"results"`
}

func main() {
	searchQuery := flag.String("q", "", "search query")
	limit := flag.Int("limit", 20, "number of results")
	outputJSON := flag.Bool("json", false, "output as JSON")
	flag.Parse()

	if *searchQuery == "" {
		fmt.Fprintf(os.Stderr, "usage: %s -q <search-query> [-limit 20]\n", os.Args[0])
		os.Exit(1)
	}

	apiURL := fmt.Sprintf(
		"https://huggingface.co/api/models?library=gguf&search=%s&sort=downloads&direction=-1&limit=%d&full=true",
		url.QueryEscape(*searchQuery),
		*limit,
	)

	resp, err := http.Get(apiURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("request failed: %s\n%s", resp.Status, string(body)))
	}

	var models []Model
	if err := json.NewDecoder(resp.Body).Decode(&models); err != nil {
		panic(err)
	}

	var results []Result

	for _, model := range models {
		modelID := model.ModelID
		if modelID == "" {
			modelID = model.ID
		}

		var files []string
		for _, file := range model.Siblings {
			if strings.HasSuffix(strings.ToLower(file.RFilename), ".gguf") {
				files = append(files, fmt.Sprintf("https://huggingface.co/%s/resolve/main/%s", modelID, file.RFilename))
			}
		}

		if len(files) > 0 {
			results = append(results, Result{
				ModelID:   modelID,
				ModelName: modelID,
				Files:     files,
			})
		}
	}

	if *outputJSON {
		out := Output{
			Query:   *searchQuery,
			Limit:   *limit,
			Results: results,
		}
		json.NewEncoder(os.Stdout).Encode(out)
	} else {
		for _, r := range results {
			for _, f := range r.Files {
				fmt.Println(f)
			}
		}
	}
}