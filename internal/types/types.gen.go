package types

import (
	"database/sql/driver"
	"encoding"
	"errors"
	"github.com/google/uuid"
)

var (
	ErrEntityIsNil     = errors.New("err entity is nil")
	ErrUnknownScanType = errors.New("err unable to scan, unknown underlying type")
)

func Parse[T encoding.TextUnmarshaler](s string) (T, error) {
	val := new(T)
	err := T.UnmarshalText(val, []byte(s))
	return *val, err
}

func MustParse[T encoding.TextUnmarshaler](s string) T {
	res, err := Parse[T](s)
	if err != nil {
		panic(err)
	}
	return res
}

var ChatIDNil = ChatID(uuid.Nil)

func NewChatID() ChatID {
	return ChatID(uuid.New())
}

func (cid *ChatID) String() string {
	if cid == nil {
		return ""
	}
	return uuid.UUID(*cid).String()
}

func (cid *ChatID) MarshalText() ([]byte, error) {
	if cid == nil {
		return nil, ErrEntityIsNil
	}
	return []byte(uuid.UUID(*cid).String()), nil
}

func (cid *ChatID) UnmarshalText(text []byte) error {
	if cid == nil {
		return ErrEntityIsNil
	}
	val, err := uuid.FromBytes(text)
	if err != nil {
		return err
	}
	*cid = ChatID(val)
	return nil
}

func (cid *ChatID) Value() (driver.Value, error) {
	if cid == nil {
		return nil, ErrEntityIsNil
	}
	return cid.String(), nil
}

func (cid *ChatID) Scan(src any) error {
	if cid == nil {
		return ErrEntityIsNil
	}
	if src == nil {
		*cid = ChatIDNil
	}
	switch src.(type) {
	case string:
		return cid.UnmarshalText([]byte(src.(string)))
	case *string:
		return cid.UnmarshalText([]byte(*src.(*string)))
	case []byte:
		return cid.UnmarshalText(src.([]byte))
	case *[]byte:
		return cid.UnmarshalText(*src.(*[]byte))
	default:
		return ErrUnknownScanType
	}
}

func (cid *ChatID) Validate() error {
	if cid == nil {
		return ErrEntityIsNil
	}
	_, err := uuid.FromBytes([]byte(cid.String()))
	return err
}

func (cid *ChatID) Matches(x any) bool {
	if cid == nil {
		return false
	}
	switch x.(type) {
	case ChatID:
		if *cid == x.(ChatID) {
			return true
		}
	case *ChatID:
		if x.(*ChatID) != nil && *cid == *x.(*ChatID) {
			return true
		}
	case uuid.UUID:
		if uuid.UUID(*cid) == x.(uuid.UUID) {
			return true
		}
	case *uuid.UUID:
		if x.(*uuid.UUID) != nil && uuid.UUID(*cid) == *x.(*uuid.UUID) {
			return true
		}
	default:
		val := ChatID{}
		if err := val.Scan(x); err != nil {
			return false
		}
		if val == *cid {
			return true
		}
	}
	return false
}

func (cid *ChatID) IsZero() bool {
	if cid == nil {
		return true
	}
	if cid.String() == "" {
		return true
	}
	if *cid == ChatIDNil {
		return true
	}
	return false
}
