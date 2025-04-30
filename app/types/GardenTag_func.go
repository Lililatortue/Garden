package types

import "iter"

func (gt *GardenTag) IterateToParent() iter.Seq[*GardenTag] {
	return func(yield func(*GardenTag) bool) {
		curr := gt
		for curr != nil {
			if !yield(curr) {
				return
			}
			curr = curr.Parent
		}
	}
}

func (gt *GardenTag) IterateWhile(stop func(node *GardenTag) bool, action func(node *GardenTag)) {
	for tag := range gt.IterateToParent() {
		if stop(tag) {
			return
		}
		action(tag)
	}
}
