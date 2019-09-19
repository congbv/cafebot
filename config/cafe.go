package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"
)

type CafeConfig struct {
	FirstOrderTime   CafeTime      `json:"first_order_time"`
	LastOrderTime    CafeTime      `json:"last_order_time"`
	TimeSlotInterval time.Duration `json:"time_slot_interval"`

	Menu map[string][]string `json:"menu"`
}

// CafeTime exists only because we need to unmarshal string of type HH:MM into
// time.Time. It is possible only by having custom named type with its own
// json.Unmarshaler implementation
type CafeTime time.Time

func (c *CafeTime) UnmarshalJSON(data []byte) error {
	dlen := len(data)
	if dlen < 2 {
		return errors.New("invalid time format")
	}

	t, err := time.Parse("15:04", string(data[1:dlen-1]))
	if err != nil {
		return err
	}

	ct := CafeTime(t)
	*c = ct

	return nil
}
func loadCafeConfig(f string) (conf CafeConfig, err error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &conf)
	return
}
