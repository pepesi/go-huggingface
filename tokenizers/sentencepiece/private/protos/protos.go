// Package protos have the Proto Buffer code for the sentencepiece_model.proto file,
// downloaded from https://github.com/google/sentencepiece/blob/master/src/sentencepiece_model.proto.
//
// The Model
package protos

//go:generate protoc --go_out=. --go_opt=paths=source_relative ./gomlx_go_huggingface_sentencepiece_model.proto
