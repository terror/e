package main

import (
  "encoding/json"
  "fmt"
  "github.com/spf13/cobra"
  "io/ioutil"
  "log"
  "os"
  "os/exec"
  "os/user"
  "path/filepath"
  "strings"
  "time"
)

type Entry struct {
  Path       string    `json:"path"`
  Score      float64   `json:"score"`
  LastAccess time.Time `json:"last_access"`
}

func NewEntry(path string) Entry {
  return Entry{Path: path, Score: 1.0, LastAccess: time.Now()}
}

func (e *Entry) Merge(other Entry) Entry {
  return Entry{
    Path:       e.Path,
    Score:      e.Score + other.Score,
    LastAccess: time.Now(),
  }
}

type Index struct {
  path string
}

func NewIndex(path string) Index {
  return Index{path: path}
}

func (i *Index) Update(entry Entry) {
  entries := i.read()

  found := false

  for index, e := range entries {
    if e.Path == entry.Path {
      found = true
      entries[index] = e.Merge(entry)
      break
    }
  }

  if !found {
    entries = append(entries, entry)
  }

  i.write(entries)
}

func (i *Index) Search(name string) []Entry {
  entries := i.read()

  var matches []Entry

  for _, entry := range entries {
    if filepath.Base(entry.Path) == name {
      matches = append(matches, entry)
    }
  }

  return matches
}

func (i *Index) read() []Entry {
  data, err := ioutil.ReadFile(i.path)

  if os.IsNotExist(err) {
    return []Entry{}
  } else if err != nil {
    log.Fatal(err)
  }

  var entries []Entry

  err = json.Unmarshal(data, &entries)

  if err != nil {
    log.Fatal(err)
  }

  return entries
}

func (i *Index) write(entries []Entry) {
  data, err := json.Marshal(entries)

  if err != nil {
    log.Fatal(err)
  }

  file, err := os.Create(i.path)

  if err != nil {
    log.Fatal(err)
  }

  defer file.Close()

  _, err = file.Write(data)

  if err != nil {
    log.Fatal(err)
  }
}

func fuzzySearch(query string, matches []Entry) Entry {
  var bestMatch Entry

  bestScore := float64(-1)

  now := time.Now()

  for _, entry := range matches {
    frecency := calculateFrecency(entry, now)
    if frecency > bestScore {
      bestMatch = entry
      bestScore = frecency
    }
  }

  return bestMatch
}

func calculateFrecency(entry Entry, now time.Time) float64 {
  duration := now.Sub(entry.LastAccess)

  score := entry.Score

  if duration < time.Hour {
    return score * 4
  } else if duration < 24*time.Hour {
    return score * 2
  } else if duration < 7*24*time.Hour {
    return score / 2
  }

  return score / 4
}

var root = &cobra.Command{
  Use:   "e",
  Short: "Edit files quickly",
  Run:   run,
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

func openInEditor(editor, path string) {
  cmd := exec.Command(editor, path)

  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr

  err := cmd.Run()

  if err != nil {
    log.Fatal(err)
  }
}

func run(cmd *cobra.Command, args []string) {
  editor := os.Getenv("EDITOR")

  if editor == "" {
    editor = "vim"
  }

  index := NewIndex(expand("~/.e.db"))

  if len(args) == 0 {
    fmt.Println("error: No file specified")
    os.Exit(1)
  }

  fp, err := filepath.Abs(args[0])

  if err != nil {
    log.Fatal(err)
  }

  index.Update(NewEntry(fp))

  matches := index.Search(filepath.Base(fp))

  if len(matches) > 1 {
    openInEditor(editor, fuzzySearch(fp, matches).Path)
  } else {
    openInEditor(editor, matches[0].Path)
  }
}

func main() {
  if err := root.Execute(); err != nil {
    fmt.Println(fmt.Sprintf("error: %s", err))
    os.Exit(1)
  }
}
