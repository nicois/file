package file

import (
	"log"
	"path/filepath"
	"strings"
)

type (
	Void  struct{}
	Paths map[string]Void
)

var Member Void

func CreatePaths(paths ...string) Paths {
	var err error
	result := make(Paths)
	for _, path := range paths {
		if !strings.HasPrefix(path, "/") {
			path, err = filepath.Abs(path)
			if err != nil {
				log.Panicf("Trying to add %v which is not an absolute path", path)
			}
		}
		result[path] = Member
	}
	return result
}

func (p *Paths) Add(items ...string) *Paths {
	for _, i := range items {
		if !strings.HasPrefix(i, "/") {
			log.Panicf("Trying to add %v which is not an absolute path", i)
		}
		(*p)[i] = Member
	}
	return p
}

func (p *Paths) AddRelative(items ...string) *Paths {
	/*
	   By default paths should be absolute. This allows for
	   those times when a relative path is actually desired
	*/
	for _, i := range items {
		(*p)[i] = Member
	}
	return p
}

func (p *Paths) Discard(other Paths) *Paths {
	for o := range other {
		delete(*p, o)
	}
	return p
}

func (p *Paths) Union(other Paths) *Paths {
	for o := range other {
		(*p)[o] = Member
	}
	return p
}

func (p *Paths) Intersection(other Paths) *Paths {
	for mine := range *p {
		if _, exists := other[mine]; !exists {
			delete(*p, mine)
		}
	}
	return p
}

func (s *Paths) Copy() Paths {
	result := make(Paths)
	for o := range *s {
		result[o] = Member
	}
	return result
}
