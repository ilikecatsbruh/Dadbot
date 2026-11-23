package main

import (
	"context"
	"sync"

	"github.com/ilikecatsbruh/Dadbot/internal/engine"
)

type App struct {
	ctx context.Context
	mu sync.RWMutex
	targets []engine.Target
	scanner *engine.Engine
}

func NewApp() *App {
	return &App{
		targets: engine.DefaultTargets(),
		scanner: engine.New(20),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetTargets() []engine.Target {
	a.mu.RLock()
	defer a.mu.RUnlock()
	out := make([]engine.Target, len(a.targets))
	copy(out, a.targets)
	return out
}

func (a *App) AddTarget(target engine.Target) error {
	if err := target.Validate(); err != nil {
		return err
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.targets = append(a.targets, target)
	return nil
}

func (a *App) StartScan() error {
	a.mu.RLock()
	targets := append([]engine.Target(nil), a.targets...)
	a.mu.RUnlock()
	return a.scanner.Start(a.ctx, targets)
}

func (a *App) StopScan() {
	a.scanner.Stop()
}

func (a *App) GetStatus() engine.Status {
	return a.scanner.Status()
}

func (a *App) GetLogs() []engine.LogEntry {
	return a.scanner.Logs()
}
