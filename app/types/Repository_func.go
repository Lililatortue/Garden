package types

func (r *Repository) GetBranch(name string) *Branch {
	for _, b := range r.Branches {
		if b.Name == name {
			return b
		}
	}
	return nil
}
