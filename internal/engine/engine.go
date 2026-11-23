package engine

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

const DefaultWorkers = 20

type Target struct {
	Name      string   `json:"name"`
	Keywords  []string `json:"keywords"`
	Location  string   `json:"location"`
	RadiusKM  int      `json:"radius_km"`
	MaxPrice  float64  `json:"max_price"`
	Currency  string   `json:"currency"`
	Enabled   bool     `json:"enabled"`
}

func (t Target) Validate() error {
	if strings.TrimSpace(t.Name) == "" {
		return errors.New("target name is required")
	}
	if len(t.Keywords) == 0 {
		return errors.New("at least one keyword is required")
	}
	if t.MaxPrice < 0 {
		return errors.New("maximum price cannot be negative")
	}
	if t.Currency == "" {
		return errors.New("currency is required")
	}
	return nil
}

func DefaultTargets() []Target {
	return []Target{
		{Name: "Athens CDs", Keywords: []string{"CD συλλογή", "μουσικά CD", "CD albums"}, Location: "Athens", RadiusKM: 30, MaxPrice: 40, Currency: "EUR", Enabled: true},
		{Name: "Greek vinyl", Keywords: []string{"δίσκοι βινυλίου", "συλλογή δίσκων", "LP vinyl"}, Location: "Greece", RadiusKM: 0, MaxPrice: 80, Currency: "EUR", Enabled: true},
	}
}

type LogEntry struct {
	Time    time.Time `json:"time"`
	Level   string    `json:"level"`
	Message string    `json:"message"`
}

type Status struct {
	Running   bool      `json:"running"`
	Workers   int       `json:"workers"`
	Queued    int       `json:"queued"`
	Completed int       `json:"completed"`
	StartedAt time.Time `json:"started_at"`
}

type Engine struct {
	workers int
	mu sync.RWMutex
	cancel context.CancelFunc
	status Status
	logs []LogEntry
}

func New(workers int) *Engine {
	if workers < 1 {
		workers = DefaultWorkers
	}
	return &Engine{workers: workers}
}

func (e *Engine) Start(parent context.Context, targets []Target) error {
	e.mu.Lock()
	if e.status.Running {
		e.mu.Unlock()
		return errors.New("a scan is already running")
	}
	ctx, cancel := context.WithCancel(parent)
	e.cancel = cancel
	e.status = Status{Running: true, Workers: e.workers, Queued: len(targets), StartedAt: time.Now()}
	e.logs = append(e.logs, LogEntry{Time: time.Now(), Level: "INFO", Message: fmt.Sprintf("Scan started with %d workers and %d targets", e.workers, len(targets))})
	e.mu.Unlock()

	go e.run(ctx, targets)
	return nil
}

func (e *Engine) run(ctx context.Context, targets []Target) {
	defer func() {
		e.mu.Lock()
		e.status.Running = false
		e.cancel = nil
		e.logs = append(e.logs, LogEntry{Time: time.Now(), Level: "INFO", Message: "Scan stopped"})
		e.mu.Unlock()
	}()

	jobs := make(chan Target)
	var wg sync.WaitGroup
	for i := 0; i < e.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case target, ok := <-jobs:
					if !ok {
						return
					}
					if !target.Enabled {
						continue
					}
					// Source adapters will enqueue permitted searches here.
					e.addLog("DEBUG", "Prepared target: "+target.Name)
					e.mu.Lock()
					e.status.Completed++
					e.mu.Unlock()
				}
			}
		}()
	}

sendLoop:
	for _, target := range targets {
		select {
		case <-ctx.Done():
			break sendLoop
		case jobs <- target:
		}
	}
	close(jobs)
	wg.Wait()
}

func (e *Engine) Stop() {
	e.mu.RLock()
	cancel := e.cancel
	e.mu.RUnlock()
	if cancel != nil {
		cancel()
	}
}

func (e *Engine) Status() Status {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.status
}

func (e *Engine) Logs() []LogEntry {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make([]LogEntry, len(e.logs))
	copy(out, e.logs)
	return out
}

func (e *Engine) addLog(level, message string) {
	e.mu.Lock()
	e.logs = append(e.logs, LogEntry{Time: time.Now(), Level: level, Message: message})
	if len(e.logs) > 2000 {
		e.logs = append([]LogEntry(nil), e.logs[len(e.logs)-2000:]...)
	}
	e.mu.Unlock()
}
