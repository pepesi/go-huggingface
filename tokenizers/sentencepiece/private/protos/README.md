# Proto Files

* [`gomlx_go_huggingface_sentencepiece_model.proto`](https://github.com/google/sentencepiece/blob/master/src/sentencepiece_model.proto) is
  downloaded from the C++ original source, in [https://github.com/google/sentencepiece/](https://github.com/google/sentencepiece),
  but it should match the one used by the [github.com/eliben/go-sentencepiece](https://github.com/eliben/go-sentencepiece)
  library.

Because of odd naming conflict (see question in https://groups.google.com/g/protobuf/c/UWWuoRWz1Uk)
the file name (not the path) has to be changed from `sentencepiece_model.proto` to something **globally** unique (!?).

