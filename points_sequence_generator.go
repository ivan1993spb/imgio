package imgio

import (
	"image"
	"sync/atomic"
)

type PointsSequenceGenerator interface {
	Current() image.Point
	Next()
	Rewind()
	Valid() bool
	Seek(offset uint64)
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
	cursor := atomic.LoadUint64(&spsg.cursor)
	p.X = int(cursor % uint64(spsg.rect.Size().X))
	p.Y = int((cursor - uint64(p.X)) / uint64(spsg.rect.Size().X))
	return spsg.rect.Min.Add(p)
}

func (spsg *SimplePointsSequenceGenerator) Next() {
	atomic.AddUint64(&spsg.cursor, 1)
}

func (spsg *SimplePointsSequenceGenerator) Rewind() {
	atomic.StoreUint64(&spsg.cursor, 0)
}

func (spsg *SimplePointsSequenceGenerator) Valid() bool {
	return spsg.cursor < uint64(spsg.rect.Size().X*spsg.rect.Size().Y)
}

func (spsg *SimplePointsSequenceGenerator) Seek(offset uint64) {
	atomic.StoreUint64(&spsg.cursor, offset)
}
