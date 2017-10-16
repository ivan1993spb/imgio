package imgio

import (
	"io"
)

type Storage interface {
	io.ReadWriteSeeker
}
