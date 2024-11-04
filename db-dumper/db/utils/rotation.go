package utils

import (
	"os"
	"path/filepath"
	"sort"
	// "time"
)

func RotateBackups(directory string, maxBackups int) error {
	files, err := filepath.Glob(filepath.Join(directory, "*.sql.gz"))
	if err != nil {
		return err
	}

	if len(files) <= maxBackups {
		return nil
	}

	sort.Slice(files, func(i, j int) bool {
		iInfo, _ := os.Stat(files[i])
		jInfo, _ := os.Stat(files[j])
		return iInfo.ModTime().After(jInfo.ModTime())
	})

	for _, file := range files[maxBackups:] {
		if err := os.Remove(file); err != nil {
			return err
		}
	}

	return nil
}
