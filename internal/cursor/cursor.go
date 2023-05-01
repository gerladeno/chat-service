package cursor

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func Encode(data any) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("encoding data to json: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func Decode(in string, to any) error {
	bytes, err := base64.URLEncoding.DecodeString(in)
	if err != nil {
		return fmt.Errorf("base64 decoding: %v", err)
	}
	if err = json.Unmarshal(bytes, to); err != nil {
		return fmt.Errorf("decoding to json: %v", err)
	}
	return nil
}
