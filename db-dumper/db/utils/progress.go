package utils

import (
	"fmt"
	"time"
)

type ProgressReporter struct {
	startTime  time.Time
	lastUpdate time.Time
	written    int64
}

func NewProgressReporter() *ProgressReporter {
	return &ProgressReporter{
		startTime:  time.Now(),
		lastUpdate: time.Now(),
	}
}

// Implement io.Writer interface
func (p *ProgressReporter) Write(b []byte) (int, error) {
	p.written += int64(len(b))
	if time.Since(p.lastUpdate) > time.Second {
		speed := float64(p.written) / time.Since(p.startTime).Seconds()
		fmt.Printf("\rBackup in progress: %.2f MB (%.2f MB/s)", 
			float64(p.written)/1024/1024, 
			speed/1024/1024)
		p.lastUpdate = time.Now()
	}
	return len(b), nil
}

func (p *ProgressReporter) Finish() {
	fmt.Printf("\nBackup completed! Total size: %.2f MB\n", float64(p.written)/1024/1024)
}
