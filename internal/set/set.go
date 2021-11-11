package set

type (
	SetF map[float64]struct{}
	SetS map[string]struct{}
)

func (s *SetF) Add(f float64) {
	(*s)[f] = struct{}{}
}

func (s *SetS) Add(str string) {
	(*s)[str] = struct{}{}
}

func (s *SetF) Remove(f float64) bool {
	if s.Contains(f) {
		delete(*s, f)
		return true
	}
	return false
}

func (s *SetS) Remove(str string) bool {
	if s.Contains(str) {
		delete(*s, str)
		return true
	}
	return false
}

func (s *SetF) Contains(f float64) bool {
	_, ok := (*s)[f]
	return ok
}

func (s *SetS) Contains(str string) bool {
	_, ok := (*s)[str]
	return ok
}
