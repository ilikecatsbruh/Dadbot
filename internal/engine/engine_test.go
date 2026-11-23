package engine

import "testing"

func TestTargetValidation(t *testing.T) {
	valid := Target{Name:"Athens CDs", Keywords:[]string{"CD συλλογή"}, Location:"Athens", MaxPrice:40, Currency:"EUR", Enabled:true}
	if err := valid.Validate(); err != nil { t.Fatalf("valid target rejected: %v", err) }
	invalid := Target{Keywords:[]string{"vinyl"}, Currency:"EUR"}
	if err := invalid.Validate(); err == nil { t.Fatal("expected missing name to fail validation") }
}

func TestDefaultTargets(t *testing.T) {
	targets := DefaultTargets()
	if len(targets) == 0 { t.Fatal("expected starter targets") }
	for _, target := range targets {
		if err := target.Validate(); err != nil { t.Errorf("invalid starter target %q: %v", target.Name, err) }
	}
}

func TestWorkerDefault(t *testing.T) {
	if got := New(0).workers; got != DefaultWorkers { t.Fatalf("got %d workers, want %d", got, DefaultWorkers) }
}
