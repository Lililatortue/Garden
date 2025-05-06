package types

import "iter"

func NewGardenTag(opts ...func(*GardenTag)) *GardenTag {
	tag := DefaultGardenTag

	for _, opt := range opts {
		opt(&tag)
	}

	return &tag
}

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

func (gt *GardenTag) GetParents() List[GardenTag] {
	var parents List[GardenTag] = make([]*GardenTag, 0, 10)

	for tag := range gt.IterateToParent() {
		parents.Push(tag)
	}

	return parents
}
