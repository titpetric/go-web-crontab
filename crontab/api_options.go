package crontab

type APIDependencies struct {
	cron *Crontab
}

type APIOptions struct {
	APIDependencies
}

func (APIOptions) New(deps APIDependencies) (*APIOptions, error) {
	opts := &APIOptions{}
	opts.APIDependencies = deps
	return opts, nil
}
