package set

type Empty struct{}

type String map[string]Empty

func NewString(items ...string) String {
	ss := String{}
	ss.Insert(items...)
	return ss
}

func (s String) Insert(items ...string) String {
	for _, item := range items {
		s[item] = Empty{}
	}
	return s
}

func (s String) Delete(items ...string) String {
	for _, item := range items {
		delete(s, item)
	}
	return s
}

func (s String) Has(item string) bool {
	_, ok := s[item]
	return ok
}

func (s String) HasAny(items ...string) bool {
	for _, item := range items {
		if s.Has(item) {
			return true
		}
	}
	return false
}

func (s String) Slice() []string {
	slice := make([]string, 0, len(s))
	for item := range s {
		slice = append(slice, item)
	}
	return slice
}
