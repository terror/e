## e 🏃

![CI](https://github.com/terror/e/actions/workflows/ci.yml/badge.svg)

**e** is a tool let's you find and edit recently used files quickly.

It works by maintaining an index of
[file access information](https://github.com/terror/e/blob/745b598e8ccbb5af654f695812750018252736c3/src/main.go#L20)
and computing best matches for a given query.

Matches are ranked based on
[frecency](https://en.wikipedia.org/wiki/Frecency?useskin=vector), meaning that
we use information about the files total access count _in addition to_ how long
ago it was last accessed.

For more information about exactly how this works, check out
[the algorithm](https://github.com/terror/e/blob/745b598e8ccbb5af654f695812750018252736c3/src/main.go#L38).

### Installation

You can download one of the pre-built binaries from the [releases](https://github.com/terror/e/releases)
page or you can build from source:

```bash
git clone https://github.com/terror/e.git && cd e
go build -o e ./src
./e # Move this into your $PATH
```

### Usage

Interacting with `e` is as simple as passing it a file to access:

```
Edit files quickly

Usage:
  e [flags]

Flags:
      --editor string   Command to use for editing files
  -h, --help            help for e
      --interactive     Search through matches interactively
```

Direct matches will get accessed immediately, otherwise you have a few options:

- Specify no flags for a frecency-based match
- Add the `--interactive` flag to fuzzy search through matches

### Prior Art

`e` was inspired by the well-known command-line tool
[`zoxide`](https://github.com/ajeetdsouza/zoxide), a directory-based access by
frecency utility.
