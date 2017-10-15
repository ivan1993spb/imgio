package imgio

import "image"

type PointsSequenceGenerator interface {
	Current() image.Point
	Next()
	Rewind()
	Valid() bool
}

type SimplePointsSequenceGenerator struct {
	rect   image.Rectangle
	cursor uint64
}

func NewSimplePointsSequenceGenerator(rect image.Rectangle) *SimplePointsSequenceGenerator {
	return &SimplePointsSequenceGenerator{
		rect:   rect,
		cursor: 0,
	}
}

func (spsg *SimplePointsSequenceGenerator) Current() image.Point {
	p := image.Point{}
	p.X = int(spsg.cursor % uint64(spsg.rect.Size().Y))
	p.Y = int((spsg.cursor - uint64(p.X)) / uint64(spsg.rect.Size().Y))
	return spsg.rect.Min.Add(p)
}

func (spsg *SimplePointsSequenceGenerator) Next() {
	spsg.cursor++
}

func (spsg *SimplePointsSequenceGenerator) Rewind() {
	spsg.cursor = 0
}

func (spsg *SimplePointsSequenceGenerator) Valid() bool {
	return spsg.cursor < uint64(spsg.rect.Size().X*spsg.rect.Size().Y)
}
