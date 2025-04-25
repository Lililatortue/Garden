package types

func (u *User) GetRepository(name string) *Repository {
	for _, r := range u.Repositories {
		if r.Name == name {
			return r
		}
	}
	return nil
}
