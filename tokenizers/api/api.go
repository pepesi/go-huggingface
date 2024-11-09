// Package api defines the Tokenizer API.
// It's just a hack to break the cyclic dependency, and allow the users to import `tokenizers` and get the
// default implementations.
package api

// Tokenizer interface allows one convert test to "tokens" (integer ids) and back.
//
// It also allows mapping of special tokens: tokens with a common semantic (like padding) but that
// may map to different ids (int) for different tokenizers.
type Tokenizer interface {
	Encode(text string) []int
	Decode([]int) string

	// SpecialTokenID returns ID for given special token if registered, or an error if not.
	SpecialTokenID(token SpecialToken) (int, error)
}

// SpecialToken is an enum of commonly used special tokens.
type SpecialToken int

const (
	TokBeginningOfSentence SpecialToken = iota
	TokEndOfSentence
	TokUnknown
	TokPad
	TokMask
	TokClassification
	TokSpecialTokensCount
)

//go:generate enumer -type=SpecialToken -trimprefix=Tok -transform=snake -values -text -json -yaml tokenizers.go
