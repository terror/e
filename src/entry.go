package main

import (
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

func (e *Entry) Frecency(now time.Time) float64 {
  duration := now.Sub(e.LastAccess)

  score := e.Score

  if duration < time.Hour {
    return score * 4
  } else if duration < 24*time.Hour {
    return score * 2
  } else if duration < 7*24*time.Hour {
    return score / 2
  }

  return score / 4
}
