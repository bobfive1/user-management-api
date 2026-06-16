package chrono

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

const dateLayout = "2006-01-02"

type DateOnly struct {
	time.Time
}

func (cd *DateOnly) Scan(value interface{}) error {
	if value == nil {
		cd.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		cd.Time = v
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into DateOnly", value)
	}
}

func (do DateOnly) Value() (driver.Value, error) {
	if do.IsZero() {
		return nil, nil
	}
	return do.Time, nil
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(d.Format(dateLayout))
}

func (d *DateOnly) UnmarshalJSON(data []byte) error {
	var value *string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	if value == nil || *value == "" {
		d = nil
		return nil
	}

	parsed, err := time.Parse(dateLayout, *value)
	if err != nil {
		return fmt.Errorf("must use YYYY-MM-DD format")
	}

	d.Time = parsed
	return nil
}

func (d DateOnly) DBValue() any {
	if d.Time.IsZero() {
		return nil
	}
	return d.Time
}
