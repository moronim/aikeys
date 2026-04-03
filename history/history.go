package history

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Entry represents a single run log.
type Entry struct {
	Timestamp time.Time `json:"ts"`
	Command   string    `json:"cmd"`
	Keys      []string  `json:"keys"`
	Tag       string    `json:"tag,omitempty"`
}

func historyPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".llmvlt_history.jsonl"
	}
	dir := filepath.Join(home, ".llmvlt")
	os.MkdirAll(dir, 0700)
	return filepath.Join(dir, "history.jsonl")
}

// Log appends an entry to the history file.
func Log(cmdArgs []string, keys []string, tag string) {
	entry := Entry{
		Timestamp: time.Now(),
		Command:   strings.Join(cmdArgs, " "),
		Keys:      keys,
		Tag:       tag,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return // silently skip — history is best-effort
	}

	f, err := os.OpenFile(historyPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer f.Close()

	f.Write(data)
	f.WriteString("\n")
}

// Read returns the last N history entries, most recent first.
func Read(n int) ([]Entry, error) {
	path := historyPath()

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("could not open history: %w", err)
	}
	defer f.Close()

	var all []Entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var e Entry
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			continue // skip malformed lines
		}
		all = append(all, e)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("could not read history: %w", err)
	}

	// Return last N, most recent first
	if len(all) > n {
		all = all[len(all)-n:]
	}

	// Reverse
	for i, j := 0, len(all)-1; i < j; i, j = i+1, j-1 {
		all[i], all[j] = all[j], all[i]
	}

	return all, nil
}
