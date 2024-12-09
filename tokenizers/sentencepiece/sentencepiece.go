// Package sentencepiece implements a tokenizers.Tokenizer based on SentencePiece tokenizer.
package sentencepiece

import (
	esentencepiece "github.com/eliben/go-sentencepiece"
	"github.com/gomlx/go-huggingface/hub"
	"github.com/gomlx/go-huggingface/tokenizers/api"
	"github.com/pkg/errors"
)

// New creates a SentencePiece tokenizer based on the "tokenizer.model" file, which must be a
// SentencePiece Model proto (see protos.Model).
//
// It implements a tokenizer.TokenizerConstructor function signature.
func New(config *api.Config, repo *hub.Repo) (api.Tokenizer, error) {
	if !repo.HasFile("tokenizer.model") {
		return nil, errors.Errorf("\"tokenizer.model\" file not found in repo")
	}
	tokenizerFile, err := repo.DownloadFile("tokenizer.model")
	if err != nil {
		return nil, errors.Wrapf(err, "can't download tokenizer.json file")
	}
	proc, err := esentencepiece.NewProcessorFromPath(tokenizerFile)
	if err != nil {
		return nil, errors.Wrapf(err, "can't create sentencepiece tokenizer")
	}
	return &Tokenizer{
		Processor: proc,
		Info:      proc.ModelInfo(),
	}, nil
}

// Tokenizer implements tokenizers.Tokenizer interface based on SentencePiece tokenizer by Google.
type Tokenizer struct {
	*esentencepiece.Processor
	Info *esentencepiece.ModelInfo
}

// Compile time assert that sentencepiece.Tokenizer implements tokenizers.Tokenizer interface.
var _ api.Tokenizer = &Tokenizer{}

// Encode returns the text encoded into a sequence of ids.
// It implements sampler.Vocabulary.
func (p *Tokenizer) Encode(text string) []int {
	tokens := p.Processor.Encode(text)
	return sliceMap(tokens, func(t esentencepiece.Token) int { return t.ID })
}

// Decode returns the text from a sequence of ids.
// It implements sampler.Vocabulary.
func (p *Tokenizer) Decode(ids []int) string {
	return p.Processor.Decode(ids)
}

// SpecialTokenID returns the token for the given symbol, or an error if not known.
func (p *Tokenizer) SpecialTokenID(token api.SpecialToken) (int, error) {
	switch token {
	case api.TokUnknown:
		return p.Info.UnknownID, nil
	case api.TokPad:
		return p.Info.PadID, nil
	case api.TokBeginningOfSentence:
		return p.Info.BeginningOfSentenceID, nil
	case api.TokEndOfSentence:
		return p.Info.EndOfSentenceID, nil
	}
	return 0, errors.Errorf("unknown special token: %s (%d)", token, token)
}

// sliceMap executes the given function sequentially for every element on in, and returns a mapped slice.
func sliceMap[In, Out any](in []In, fn func(e In) Out) (out []Out) {
	out = make([]Out, len(in))
	for ii, e := range in {
		out[ii] = fn(e)
	}
	return
}
