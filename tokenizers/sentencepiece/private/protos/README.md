# Proto Files

* [`sentencepiece_model.proto`](https://github.com/google/sentencepiece/blob/master/src/sentencepiece_model.proto) is
  downloaded from the C++ original source, in [https://github.com/google/sentencepiece/](https://github.com/google/sentencepiece),
  but it should match the one used by the [github.com/eliben/go-sentencepiece](https://github.com/eliben/go-sentencepiece)
  library.

Because of protoc unique file naming requirement (!?), described in email thread in https://groups.google.com/g/protobuf/c/UWWuoRWz1Uk,
we compile by first creating a unique prefix directory. See `gen_protos.sh` script.
