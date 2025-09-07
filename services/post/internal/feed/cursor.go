package feed

import (
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"time"
)

type Cursor struct {
	CreatedAt time.Time
	ID        uint
}

func EncodeCursor(c Cursor) string {
	raw := c.CreatedAt.Format(time.RFC3339Nano) + "|" + strconv.FormatUint(uint64(c.ID), 10)
	return base64.StdEncoding.EncodeToString([]byte(raw))
}

func DecodeCursor(s string) (*Cursor, error) {
	if s == "" {
		return nil, nil
	}
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(string(b), "|")
	if len(parts) != 2 {
		return nil, errors.New("bad cursor")
	}
	t, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return nil, err
	}
	id64, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return nil, err
	}
	c := Cursor{CreatedAt: t, ID: uint(id64)}
	return &c, nil
}
