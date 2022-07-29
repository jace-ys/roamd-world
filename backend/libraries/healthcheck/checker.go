package healthcheck

import "github.com/alexliesenfeld/health"

func NewChecker(checks ...health.Check) health.Checker {
	opts := []health.CheckerOption{
		health.WithDisabledAutostart(),
	}

	for _, check := range checks {
		opts = append(opts, health.WithCheck(check))
	}

	return health.NewChecker(opts...)
}
