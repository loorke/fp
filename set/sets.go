package set

import "maps"

type Set[tA comparable] map[tA]bool

func New[tA comparable](a ...tA) Set[tA] {
	m := make(Set[tA], len(a))
	for _, e := range a {
		m[e] = true
	}
	return m
}

func (sa Set[tA]) Union(sb Set[tA]) Set[tA] {
	m := make(Set[tA], len(sa))
	maps.Copy(m, sa)
	maps.Copy(m, sb)
	return m
}

func (sa Set[tA]) Intersection(sb Set[tA]) Set[tA] {
	if len(sa) < len(sb) {
		sa, sb = sb, sa
	}

	m := Set[tA]{}
	for e := range sa {
		if sb[e] {
			m[e] = true
		}
	}
	return m
}

func (sa Set[tA]) Diff(sb Set[tA]) Set[tA] {
	m := Set[tA]{}
	for e := range sa {
		if !sb[e] {
			m[e] = true
		}
	}
	return m
}

func (sa Set[tA]) SymmetricDiff(sb Set[tA]) Set[tA] {
	return sa.Union(sb).Diff(sa.Intersection(sb))
}

func (sa Set[tA]) List() []tA {
	a := make([]tA, 0, len(sa))
	for e := range sa {
		a = append(a, e)
	}
	return a
}

func (sa Set[tA]) Contains(e tA) bool {
	return sa[e]
}

func (sa Set[tA]) Add(e tA) {
	sa[e] = true
}
