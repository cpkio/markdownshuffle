# Markdown Shuffle

This program collects all markdown files in current folder into one document, sorted by header titles.

It uses a collection of document trees one per opened file and merging of those trees.

## Usage

```
mdshf [*.extension]
```

## Installation

Add `%GOPATH%/bin` to `PATH`.

```
go get github.com/cpkio/markdownshuffle
cd cmd/mdshf
go install
```

