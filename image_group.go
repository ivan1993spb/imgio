package imgio

import (
	"io"
)

type ImageGroup struct {
	images []*Image
	cursor int
}

func NewImageGroup(images ...*Image) *ImageGroup {
	i := make([]*Image, len(images))
	copy(i, images)
	return &ImageGroup{
		images: i,
		cursor: 0,
	}
}

// Read implements io.Reader interface
func (ig *ImageGroup) Read(p []byte) (n int, err error) {
	if ig.cursor >= len(ig.images) {
		return 0, io.EOF
	}

	n, err = ig.images[ig.cursor].Read(p)

	if n > 0 || err != io.EOF {
		if err == io.EOF && len(ig.images) > ig.cursor {
			// Don't return io.EOF yet. More images remain.
			err = nil
			ig.cursor++
		}
		return
	}

	return 0, io.EOF
}

// Write implements io.Writer interface
func (ig *ImageGroup) Write(p []byte) (n int, err error) {
	var written int

	for ig.cursor < len(ig.images) && len(p) > 0 {
		written, err = ig.images[ig.cursor].Write(p)
		n += written
		p = p[written:]

		if err != nil {
			if err == ErrOverflow {
				ig.cursor++
				err = nil
			} else {
				return
			}
		}
	}

	if len(p) > 0 {
		return n, ErrOverflow
	}

	return n, nil
}

func (ig *ImageGroup) Size() (size int64) {
	for _, image := range ig.images {
		size += image.Size()
	}
	return
}

func (ig *ImageGroup) Rewind() {
	ig.cursor = 0
	for _, image := range ig.images {
		image.gen.Rewind()
	}
}
