package types

func NewRepository(opts ...func(*Repository)) *Repository {
	repo := DefaultRepository

	for _, opt := range opts {
		opt(&repo)
	}

	return &repo
}

func (r *Repository) GetBranch(name string) *Branch {
	for _, b := range r.Branches {
		if b.Name == name {
			return b
		}
	}
	return nil
}

func (r *Repository) AddBranch(branch ...*Branch) {
	r.Branches = append(r.Branches, branch...)
}
