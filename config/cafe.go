package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"
)

type CafeConfig struct {
	OpenTime  CafeTime `json:"open_time"`
	CloseTime CafeTime `json:"close_time"`

	TimeSlotInterval time.Duration `json:"time_slot_interval"`
	TimeSlotNumber   int           `json:"time_slot_number"`
	TimeSlotNumInRow int           `json:"time_slot_num_in_row"`

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
