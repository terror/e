package main

import (
  "errors"
  "github.com/spf13/cobra"
  "os"
  "path/filepath"
  "sort"
  "time"
)

var root = &cobra.Command{
  Use:   "e",
  Short: "Edit files quickly",
  Run:   run,
}

func run(cmd *cobra.Command, args []string) {
  editor := os.Getenv("EDITOR")

  editorOverride, err := cmd.Flags().GetString("editor")

  if err != nil {
    die(err)
  }

  if editorOverride != "" {
    editor = editorOverride
  }

  if editor == "" {
    editor = "vim"
  }

  index := NewIndex(expand("~/.e.db"))

  if len(args) == 0 {
    die(errors.New("No filename specified"))
  }

  fp, err := filepath.Abs(args[0])

  if err != nil {
    die(err)
  }

  if err := index.Update(NewEntry(fp)); err != nil {
    die(err)
  }

  if state(fp) != Unknown {
    edit(editor, fp)
    return
  }

  matches, err := index.Search(filepath.Base(fp))

  if err != nil {
    die(err)
  }

  if len(matches) == 0 {
    edit(editor, fp)
    return
  }

  if len(matches) == 1 {
    edit(editor, matches[0].Path)
    return
  }

  sort.Slice(matches, func(i, j int) bool {
    time := time.Now()
    return matches[i].Frecency(time) < matches[j].Frecency(time)
  })

  interactive, err := cmd.Flags().GetBool("interactive")

  if err != nil {
    die(err)
  }

  if !interactive {
    edit(editor, matches[0].Path)
    return
  }

  selected, err := search(matches)

  if err != nil {
    die(err)
  }

  edit(editor, selected.Path)
}

func main() {
  root.PersistentFlags().
    Bool("interactive", false, "Search through matches interactively")

  root.PersistentFlags().
    String("editor", "", "Command to use for editing files")

  if err := root.Execute(); err != nil {
    die(err)
  }
}
