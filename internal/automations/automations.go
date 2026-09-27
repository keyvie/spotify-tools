package automations

import (
	"fmt"
	"time"

	"github.com/keyvie/spotify-tools/internal/config"
)

const (
	SorterType string = "sorter"
)

type Automation struct {
	Name       string
	Type       string
	Account    config.Account
	Enabled    bool
	Interval   int64
	LastRunAt  int64
	Options    map[string]any
}

func (a *Automation) Run(force bool) error {
	if a.LastRunAt != 0 && force != true {
		if time.Now().Unix() - a.LastRunAt < a.Interval {
			return nil
		}
	}

	if !a.Enabled && force != true {
		return nil
	}

	switch a.Type {
		case SorterType:
			if err := RunSorter(a); err != nil {	
				return err
			}
			break
		default:
			return fmt.Errorf("unknown automation type: %s", a.Type)
	}

	if err := config.Lock(); err != nil {
		return err
	}
	defer config.Unlock()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	cfgAutomation := cfg.Automations[a.Name]
	cfgAutomation.LastRunAt = time.Now().Unix()
	cfg.Automations[a.Name] = cfgAutomation
	if err := config.Save(cfg); err != nil {
		return err
	}

	return nil
}

func build(name string, a config.Automation) (*Automation, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	account := cfg.Accounts[a.AccountID]
	return &Automation{
		Name:     	name,
		Type:     	a.Type,
		Account:    account,
		Enabled:  	a.Enabled,
		Interval: 	a.Interval,
		LastRunAt: 	a.LastRunAt,
		Options:  	a.Options,
	}, nil
}

func All() ([]Automation, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	automations := make([]Automation, 0, len(cfg.Automations))
	for name, a := range cfg.Automations {
		built, err := build(name, a)
		if err != nil {
			return nil, fmt.Errorf("failed to build automation %s: %w", name, err)
		}
		automations = append(automations, *built)
	}
	return automations, nil
}

func Get(name string) (*Automation, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	a, exists := cfg.Automations[name]
	if !exists {
		return nil, fmt.Errorf("automation with name %s does not exist", name)
	}

	automation, err := build(name, a)
	if err != nil {
		return nil, fmt.Errorf("failed to build automation %s: %w", name, err)
	}

	return automation, nil
}