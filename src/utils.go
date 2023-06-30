package main

import (
  "fmt"
  "github.com/ktr0731/go-fuzzyfinder"
  "os"
  "os/exec"
  "os/user"
  "path/filepath"
  "strings"
)

func die(err error) {
  fmt.Printf("error: %s\n", err)
  os.Exit(1)
}

func edit(editor, path string) {
  cmd := exec.Command(editor, path)

  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr

  err := cmd.Run()

  if err != nil {
    die(err)
  }
}

func expand(path string) string {
  usr, _ := user.Current()

  dir := usr.HomeDir

  if path == "~" {
    return dir
  } else if strings.HasPrefix(path, "~/") {
    return filepath.Join(dir, path[2:])
  } else {
    return path
  }
}

func isFile(path string) bool {
  info, err := os.Stat(path)

  if err != nil {
    return false
  }

  return info.Mode().IsRegular()
}

func search(matches []Entry) (Entry, error) {
  var paths []string

  for _, match := range matches {
    paths = append(paths, match.Path)
  }

  index, err := fuzzyfinder.Find(
    paths,
    func(i int) string {
      return paths[i]
    },
  )

  if err != nil {
    return Entry{}, err
  }

  return matches[index], nil
}
