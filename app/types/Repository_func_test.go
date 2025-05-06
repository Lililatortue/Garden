package types

import "testing"

func FuzzNewRepository_default(f *testing.F) {
	f.Fuzz(func(t *testing.T) {
		repo := NewRepository()

		if repo == nil {
			t.Errorf("expected repo to be non-nil")
		}
		if repo.Name != "" {
			t.Errorf("expected repo name to be empty, got %s", repo.Name)
		}
		if len(repo.Branches) != 0 {
			t.Errorf("expected repo branches to be empty, got %d", len(repo.Branches))
		}
		if repo.ID != -1 {

		}
	})
}

func FuzzUser_GetRepository(f *testing.F) {
	testSet := []struct {
		in string
	}{
		{"test"},
		{"test2"},
		{"test3"},
	}
	for _, ts := range testSet {
		f.Add(ts.in)
	}
	f.Fuzz(func(t *testing.T, in string) {

	})
}
