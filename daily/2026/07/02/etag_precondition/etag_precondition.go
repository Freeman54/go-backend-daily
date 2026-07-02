package etagprecondition

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

var (
	ErrMalformedETag       = errors.New("malformed etag")
	ErrPreconditionFailed  = errors.New("precondition failed")
	ErrMissingPrecondition = errors.New("missing if-match header")
)

type Document struct {
	Body    string
	Version int64
}

type Store struct {
	mu  sync.Mutex
	doc Document
}

func NewStore(initial string) *Store {
	return &Store{doc: Document{Body: initial, Version: 1}}
}

func (s *Store) Snapshot() Document {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.doc
}

func (s *Store) Update(ifMatch string, nextBody string) (Document, error) {
	if ifMatch == "" {
		return Document{}, ErrMissingPrecondition
	}

	version, err := parseETag(ifMatch)
	if err != nil {
		return Document{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.doc.Version != version {
		return Document{}, ErrPreconditionFailed
	}

	s.doc.Body = nextBody
	s.doc.Version++
	return s.doc, nil
}

func ETag(version int64) string {
	return fmt.Sprintf(`"v%d"`, version)
}

func parseETag(value string) (int64, error) {
	trimmed := strings.Trim(value, `"`)
	if !strings.HasPrefix(trimmed, "v") {
		return 0, ErrMalformedETag
	}

	version, err := strconv.ParseInt(strings.TrimPrefix(trimmed, "v"), 10, 64)
	if err != nil || version <= 0 {
		return 0, ErrMalformedETag
	}
	return version, nil
}
