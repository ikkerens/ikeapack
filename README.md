# Go packed serializer

This is a packed struct serializer that is mostly meant for a private project but was released as it may be useful to someone else.

Originally this package was made as an extension to binary.Read and binary.Write, but I soon found those functions didn't match my use case as they offered no support for strings nor compression.

#### Features
* Caches types for faster calls to the same type
* Compression support

#### Format
* All primitives are stored in big endian format
* All slices are stored with a uint32 prefix indicating their length
* Strings are stored with a uint32 prefix indicating their length
* Compression blocks are stored using deflate (level 9) with a uint32 prefixing the size of the compressed data blob