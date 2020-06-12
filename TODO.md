# Todo: Go Compression

## Folders

Instead of creating a new file create a folder that stores the key

```
|- file.txt
|- file.txt.cmp
```

Will become:

```
|- file.txt.cmp
|-- file.txt.cmp
|-- .keystore
```

This will allow it to be decoded in a seperate runtime. The keystore could also hold a Huffman tree

## Performance

Many of the functions need performance tweaks, presently `LoadFile` is the only finished function
