## HVIF-go
![master](https://github.com/aoyako/hvif/actions/workflows/build-lint.yml/badge.svg)

*HVIF (Haiku Vector Image File) is a file format used in Haiku OS for storage-efficient applications icon storage*

This library provides low-level access to the HVIF contents. Since this library is focused on file interactions, the returned types are mostly interfaces. Even thought originaly parser is tightly coupled with renderer, it is not a scope of this library.

The functionality includes:
1. Loading vector data to the memory
2. Storing vector data back to the file (TBI)
3. Modification of style, pathes, and shapes information

### Examples:
#### Reading image file
```go
filename := "testdata/ime.hvif"
file, _ := os.Open(filename)
img, err := ReadImage(file)
```

### Contributing
HVIF-go is an open-source library. Any contributions, such as issues and pull requests, are welcomed.

### License
HVIF-go is licensed under MIT license. See LICENSE.md file.