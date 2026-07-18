package postcommithooks

type Hook func() error

type Runner struct {
	hooks []Hook
}

func (r *Runner) Add(hook Hook) {
	if hook == nil {
		return
	}
	r.hooks = append(r.hooks, hook)
}

func (r *Runner) Commit() []error {
	hooks := append([]Hook(nil), r.hooks...)
	r.hooks = nil

	errs := make([]error, 0)
	for _, hook := range hooks {
		if err := hook(); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (r *Runner) Rollback() {
	r.hooks = nil
}
