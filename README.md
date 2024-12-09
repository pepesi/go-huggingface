# **go-huggingface**, download, tokenize and convert models from HuggingFace. 

[![GoDev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/gomlx/go-huggingface?tab=doc)

## 📖 Overview

Simple APIs for downloading (`hub`), tokenizing (`tokenizers`) and (**future work**) model conversion (`models`) of 
[HuggingFace🤗](huggingface.co) models using [GoMLX](https://github.com/gomlx/gomlx).

🚧 **EXPERIMENTAL and IN DEVELOPMENT**: While the `hub` package has been stable. The `tokenizers` only supports
SentencePiece models (saved as proto), but has been working. 

## Examples

### Preamble: Imports And Variables

```go
import (
    "github.com/gomlx/go-huggingface/hub"
    "github.com/gomlx/go-huggingface/tokenizers"
)

var (
	// HuggingFace authentication token read from environment.
	// It can be created in https://huggingface.co
	// Some files may require it for downloading.
	hfAuthToken = os.Getenv("HF_TOKEN")

	// Model IDs we use for testing.
	hfModelIDs = []string{
		"google/gemma-2-2b-it",
		"sentence-transformers/all-MiniLM-L6-v2",
		"protectai/deberta-v3-base-zeroshot-v1-onnx",
		"KnightsAnalytics/distilbert-base-uncased-finetuned-sst-2-english",
		"KnightsAnalytics/distilbert-NER",
		"SamLowe/roberta-base-go_emotions-onnx",
	}
)
```

### List files for each model

```go
for _, modelID := range hfModelIDs {
	fmt.Printf("\n%s:\n", modelID)
	repo := hub.New(modelID).WithAuth(hfAuthToken)
	for fileName, err := range repo.IterFileNames() {
		if err != nil { panic(err) }
		fmt.Printf("\t%s\n", fileName)
	}
}
```

The result looks like this:

```
google/gemma-2-2b-it:
	.gitattributes
	README.md
	config.json
	generation_config.json
	model-00001-of-00002.safetensors
	model-00002-of-00002.safetensors
	model.safetensors.index.json
	special_tokens_map.json
	tokenizer.json
	tokenizer.model
	tokenizer_config.json
…
```


### List tokenizer classes for each model

```go
for _, modelID := range hfModelIDs {
	fmt.Printf("\n%s:\n", modelID)
	repo := hub.New(modelID).WithAuth(hfAuthToken)
	config, err := tokenizers.GetConfig(repo)
	if err != nil { panic(err) }
	fmt.Printf("\ttokenizer_class=%s\n", config.TokenizerClass)
}
```

Results:

```
google/gemma-2-2b-it:
	tokenizer_class=GemmaTokenizer

sentence-transformers/all-MiniLM-L6-v2:
	tokenizer_class=BertTokenizer

protectai/deberta-v3-base-zeroshot-v1-onnx:
	tokenizer_class=DebertaV2Tokenizer
…
```


### Tokenize for [`google/gemma-2-2b-it`](https://huggingface.co/google/gemma-2-2b-it) using Go-only "SentencePiece" tokenizer

* The output "Downloaded" message happens only the tokenizer file is not yet cached, so only the first time:

```go
repo := hub.New("google/gemma-2-2b-it").WithAuth(hfAuthToken)
tokenizer, err := tokenizers.New(repo)
if err != nil { panic(err) }

sentence := "The book is on the table."
tokens := tokenizer.Encode(sentence)
fmt.Printf("Sentence:\t%s\n", sentence)
fmt.Printf("Tokens:  \t%v\n", tokens)
```

```
Downloaded 1/1 files, 4.2 MB downloaded         
Sentence:	The book is on the table.
Tokens:  	[651 2870 603 611 573 3037 235265]
```

### Tokenize for a [Sentence Transformer](https://www.sbert.net/) derived model, using Rust's based [github.com/daulet/tokenizers](https://github.com/daulet/tokenizers) tokenizer

For most tokenizers in HuggingFace though, there is no Go-only version yet, and for now we use the 
[github.com/daulet/tokenizers](https://github.com/daulet/tokenizers), which is based on a fast tokenizer written in Rust.

It requires installation of the built Rust library though, 
see [github.com/daulet/tokenizers](https://github.com/daulet/tokenizers) on how to install it, 
they provide prebuilt binaries.

> **Note**: `daulet/tokenizers` also provides a simple downloader, so `go-huggingface` is not strictly necessary -- 
> if you don't want the extra dependency and only need the tokenizer, you don't need to use it. `go-huggingface` 
> helps by allowing also downloading other files (models, datasets), and a shared cache across different projects 
> and `huggingface-hub` (the python downloader library).

```go
import dtok "github.com/daulet/tokenizers"

%%
modelID := "KnightsAnalytics/all-MiniLM-L6-v2"
repo := hub.New(modelID).WithAuth(hfAuthToken)
localFile := must.M1(repo.DownloadFile("tokenizer.json"))
tokenizer := must.M1(dtok.FromFile(localFile))
defer tokenizer.Close()
tokens, _ := tokenizer.Encode(sentence, true)

fmt.Printf("Sentence:\t%s\n", sentence)
fmt.Printf("Tokens:  \t%v\n", tokens)
```

```
Sentence:	The book is on the table.
Tokens:  	[101 1996 2338 2003 2006 1996 2795 1012 102 0 0 0…]
```

### Download and execute ONNX model for [`sentence-transformers/all-MiniLM-L6-v2`](https://huggingface.co/sentence-transformers/all-MiniLM-L6-v2)

Only the first 3 lines are actually demoing `go-huggingface`.
The remainder lines uses [`github.com/gomlx/onnx-gomlx`](https://github.com/gomlx/onnx-gomlx)
to parse and convert the ONNX model to GoMLX, and then
[`github.com/gomlx/gomlx`](github.com/gomlx/gomlx) to execute the converted model
for a couple of sentences.

```go
// Get ONNX model.
repo := hub.New("sentence-transformers/all-MiniLM-L6-v2").WithAuth(hfAuthToken)
onnxFilePath, err := repo.DownloadFile("onnx/model.onnx")
if err != nil { panic(err) }
onnxModel, err := onnx.ReadFile(onnxFilePath)
if err != nil { panic(err) }

// Convert ONNX variables to GoMLX context (which stores variables):
ctx := context.New()
err = onnxModel.VariablesToContext(ctx)
if err != nil { panic(err) }

// Test input.
sentences := []string{
	"This is an example sentence",
	"Each sentence is converted"}
inputIDs := [][]int64{
	{101, 2023, 2003, 2019, 2742, 6251,  102},
	{ 101, 2169, 6251, 2003, 4991,  102,    0}}
tokenTypeIDs := [][]int64{
	{0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0}}
attentionMask := [][]int64{
	{1, 1, 1, 1, 1, 1, 1},
	{1, 1, 1, 1, 1, 1, 0}}

// Execute GoMLX graph with model.
embeddings := context.ExecOnce(
	backends.New(), ctx,
	func (ctx *context.Context, inputs []*graph.Node) *graph.Node {
		modelOutputs := onnxModel.CallGraph(ctx, inputs[0].Graph(), map[string]*graph.Node{
			"input_ids": inputs[0],
			"attention_mask": inputs[1],
			"token_type_ids": inputs[2]})
		return modelOutputs[0]
	}, 
	inputIDs, attentionMask, tokenTypeIDs)

fmt.Printf("Sentences: \t%q\n", sentences)
fmt.Printf("Embeddings:\t%s\n", embeddings)
```

```
Sentences: 	["This is an example sentence" "Each sentence is converted"]
Embeddings:	[2][7][384]float32{
 {{0.0366, -0.0162, 0.1682, ..., 0.0554, -0.1644, -0.2967},
  {0.7239, 0.6399, 0.1888, ..., 0.5946, 0.6206, 0.4897},
  {0.0064, 0.0203, 0.0448, ..., 0.3464, 1.3170, -0.1670},
  ...,
  {0.1479, -0.0643, 0.1457, ..., 0.8837, -0.3316, 0.2975},
  {0.5212, 0.6563, 0.5607, ..., -0.0399, 0.0412, -1.4036},
  {1.0824, 0.7140, 0.3986, ..., -0.2301, 0.3243, -1.0313}},
 {{0.2802, 0.1165, -0.0418, ..., 0.2711, -0.1685, -0.2961},
  {0.8729, 0.4545, -0.1091, ..., 0.1365, 0.4580, -0.2042},
  {0.4752, 0.5731, 0.6304, ..., 0.6526, 0.5612, -1.3268},
  ...,
  {0.6113, 0.7920, -0.4685, ..., 0.0854, 1.0592, -0.2983},
  {0.4115, 1.0946, 0.2385, ..., 0.8984, 0.3684, -0.7333},
  {0.1374, 0.5555, 0.2678, ..., 0.5426, 0.4665, -0.5284}}}
```

## Download Dataset Files

We are going to use the [HuggingFaceFW/fineweb](https://huggingface.co/datasets/HuggingFaceFW/fineweb) as an example, download one of its sample files (~2.5Gb of data) and parse the `.parquet` file.

### Structure of file
First we define the structure of each entry, with the tags for the Parquet parser:

```go
var (
    FineWebID = "HuggingFaceFW/fineweb"
    FineWebSampleFile = "sample/10BT/000_00000.parquet"
)

// FineWebEntry: inspection of fields in parque file done with tool in 
// github.com/xitongsys/parquet-go/tool/parquet-tools.
//
// The parquet annotations are described in: https://pkg.go.dev/github.com/parquet-go/parquet-go#SchemaOf
type FineWebEntry struct {
    Text string `parquet:"text,snappy"`
    ID string `parquet:"id,snappy"`
    Dump string `parquet:"dump,snappy"`
    URL string `parquet:"url,snappy"`
    Score float64 `parquet:"language_score"`
}

// TrimString returns s trimmed to at most maxLength runes. If trimmed it appends "…" at the end.
func TrimString(s string, maxLength int) string {
    if utf8.RuneCountInString(s) <= maxLength {
        return s
    }
    runes := []rune(s)
    return string(runes[:maxLength-1]) + "…"
}
```

Now we read the `parquet` files using the library [github.com/parquet-go/parquet-go](https://github.com/parquet-go/parquet-go).

```go
import (
    parquet "github.com/parquet-go/parquet-go"
)

func main() {
    // Download repo file.
    repo := hub.New(FineWebID).WithType(hub.RepoTypeDataset).WithAuth(hfAuthToken)
    localSampleFile := must.M1(repo.DownloadFile(FineWebSampleFile))
    
    // Parquet reading using parquet-go: it's somewhat cumbersome (to open the file it needs its size!?), but it works.
    schema := parquet.SchemaOf(&FineWebEntry{})
    fSize := must.M1(os.Stat(localSampleFile)).Size()
    fReader := must.M1(os.Open(localSampleFile))
    fParquet := must.M1(parquet.OpenFile(fReader, fSize))
    reader := parquet.NewGenericReader[FineWebEntry](fParquet, schema)
    defer reader.Close()
    
    // Print first 10 rows:
    rows := make([]FineWebEntry, 10)
    n := must.M1(reader.Read(rows))
    fmt.Printf("%d rows read\n", n)
    for ii, row := range rows {
        fmt.Printf("Row %0d:\tScore=%.3f Text=[%q], URL=[%s]\n", ii, row.Score, TrimString(row.Text, 50), TrimString(row.URL, 40))
    }
}
```

Results:

```
10 rows read
Row 0:	Score=0.823 Text=["|Viewing Single Post From: Spoilers for the Week …"], URL=[http://daytimeroyaltyonline.com/single/…]
Row 1:	Score=0.974 Text=["*sigh* Fundamentalist community, let me pass on s…"], URL=[http://endogenousretrovirus.blogspot.co…]
Row 2:	Score=0.873 Text=["A novel two-step immunotherapy approach has shown…"], URL=[http://news.cancerconnect.com/]
Row 3:	Score=0.932 Text=["Free the Cans! Working Together to Reduce Waste\nI…"], URL=[http://sharingsolution.com/2009/05/23/f…]
…
```

## [Demo Notebook](https://github.com/gomlx/go-huggingface/blob/main/go-huggingface.ipynb)

All examples were taken from the [demo notebook](https://github.com/gomlx/go-huggingface/blob/main/go-huggingface.ipynb).
It works it also as an easy playground to try out the functionality.

You can try it out using the [GoMLX docker that includes JupyterLab](https://hub.docker.com/r/janpfeifer/gomlx_jupyterlab).