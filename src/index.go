package main

import (
  "encoding/json"
  "os"
  "path/filepath"
  "sync"
)

type Index struct {
  path string
}

func NewIndex(path string) Index {
  return Index{path: path}
}

func (i *Index) Search(name string) ([]Entry, error) {
  entries, err := i.read()

  if err != nil {
    return nil, err
  }

  var wg sync.WaitGroup

  matches := make([]Entry, 0, len(entries))

  matchesMutex := &sync.Mutex{}

  for _, entry := range entries {
    wg.Add(1)

    go func(e Entry) {
      defer wg.Done()
      if filepath.Base(e.Path) == name && state(e.Path) != Unknown {
        matchesMutex.Lock()
        matches = append(matches, e)
        matchesMutex.Unlock()
      }
    }(entry)
  }

  wg.Wait()

  return matches, nil
}

func (i *Index) Update(entry Entry) error {
  entries, err := i.read()

  if err != nil {
    return err
  }

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

  if err := i.write(entries); err != nil {
    return err
  }

  return nil
}

func (i *Index) read() ([]Entry, error) {
  data, err := os.ReadFile(i.path)

  if os.IsNotExist(err) {
    return []Entry{}, nil
  }

  if err != nil {
    return nil, err
  }

  var entries []Entry

  if err := json.Unmarshal(data, &entries); err != nil {
    return nil, err
  }

  return entries, nil
}

func (i *Index) write(entries []Entry) error {
  data, err := json.Marshal(entries)

  if err != nil {
    return err
  }

  file, err := os.Create(i.path)

  if err != nil {
    return err
  }

  defer file.Close()

  if _, err := file.Write(data); err != nil {
    return err
  }

  return nil
}
