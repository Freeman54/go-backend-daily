package optionalfieldpatch

import (
	"encoding/json"
	"errors"
)

type Field[T any] struct {
	set   bool
	null  bool
	value T
}

func (f *Field[T]) UnmarshalJSON(data []byte) error {
	f.set = true
	if string(data) == "null" {
		f.null = true
		var zero T
		f.value = zero
		return nil
	}
	f.null = false
	return json.Unmarshal(data, &f.value)
}

func (f Field[T]) Present() bool {
	return f.set
}

func (f Field[T]) Null() bool {
	return f.set && f.null
}

func (f Field[T]) Value() (T, error) {
	if !f.set || f.null {
		var zero T
		return zero, errors.New("field has no concrete value")
	}
	return f.value, nil
}

type Patch struct {
	Nickname Field[string] `json:"nickname"`
	Age      Field[int]    `json:"age"`
}

type User struct {
	Nickname string
	Age      int
}

func Apply(user User, patch Patch) User {
	if patch.Nickname.Present() {
		if patch.Nickname.Null() {
			user.Nickname = ""
		} else if value, err := patch.Nickname.Value(); err == nil {
			user.Nickname = value
		}
	}
	if patch.Age.Present() {
		if patch.Age.Null() {
			user.Age = 0
		} else if value, err := patch.Age.Value(); err == nil {
			user.Age = value
		}
	}
	return user
}
