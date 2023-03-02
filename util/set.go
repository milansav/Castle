package util

type Set struct {
	innerData map[string]struct{}
}

func MakeSet(n []string) Set {
	s := Set{innerData: make(map[string]struct{})}
	for _, value := range n {
		s.innerData[value] = struct{}{}
	}

	return s
}

func (s Set) Add(n string) Set {
	s.innerData[n] = struct{}{}
	return s
}

func (s Set) Remove(n string) Set {
	delete(s.innerData, n)
	return s
}

func (s Set) Has(n string) bool {
	_, ok := s.innerData[n]

	return ok
}
