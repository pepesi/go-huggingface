# go-huggingface 

## 📖 Overview

Simple APIs for downloading (`hub`), tokenizing (`tokenizers`) and (future work) model conversion (`models`) of 
[HuggingFace🤗](huggingface.co) models using [GoMLX](https://github.com/gomlx/gomlx).

🚧 **EXPERIMENTAL and IN DEVELOPMENT**: While the `hub` package has been stable, the `tokenizers` and the future `models` are still
under intense development.

## Examples

### Preamble: Imports And Variables

```go
import (
    "github.com/gomlx/go-huggingface/hub"
    "github.com/gomlx/go-huggingface/tokenizers"
)

var (
	// Model ids for testing.
	hfModelIDs = []string{
		"google/gemma-2-2b-it",
		"sentence-transformers/all-MiniLM-L6-v2",
		"protectai/deberta-v3-base-zeroshot-v1-onnx",
		"KnightsAnalytics/distilbert-base-uncased-finetuned-sst-2-english",
		"KnightsAnalytics/distilbert-NER",
		"SamLowe/roberta-base-go_emotions-onnx",
	}
	hfAuthToken = "..."  // Create your HuggingFace authentication token in huggingface.co, to allow download of models.
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

### Tokenize for `google/gemma-2-2b-it`

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


### Download and execute ONNX model for `sentence-transformers/all-MiniLM-L6-v2`

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
