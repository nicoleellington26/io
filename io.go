package io

// ReaderWritter allows something to be read or written
type ReaderWritter interface {
	Reader
	Writter
}

// Reader allows something to be read
type Reader interface {
	Read([]byte) error
}

// Writter allows something to be written
type Writter interface {
	Write([]byte) error
}
