package htmlbuffer

import "slices"

var voidTags = []string{
	"area", "base", "br", "col", "html",
	"embed", "hr", "img", "input", "link",
	"meta", "param", "source", "track", "wbr",
}

type tag struct {
	name   string
	isVoid bool
}

func newTag(name string) tag {
	isVoid := slices.Contains(voidTags, name)
	return tag{
		name:   name,
		isVoid: isVoid,
	}
}

func (t tag) isZero() bool {
	return len(t.name) == 0
}
