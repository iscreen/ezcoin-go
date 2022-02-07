package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
	// "ezcoin.cc/ezcoin-go/server/global"
)

const TimeFormat = "2006-01-02 15:04:05"

type LocalTime time.Time

func (t *LocalTime) UnmarshalJSON(data []byte) error {
	// if string(data) == "null" || string(data) == `""` {
	// 	return nil
	// }
	// global.GVA_LOG.Debug(fmt.Sprintf("data %s", string(data)))
	// return json.Unmarshal(data, (*time.Time)(t))
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	now, err := time.ParseInLocation("2006-01-02", s, time.UTC)
	// global.GVA_LOG.Debug(fmt.Sprintf("data %v, %s, %v", err, s, now))
	if err != nil {
		return err
	}
	*t = LocalTime(now)
	return nil

	// global.GVA_LOG.Debug(fmt.Sprintf("2222222 UnmarshalJSON %v", data))
	// if len(data) == 2 {
	// 	*t = LocalTime(time.Time{})
	// 	return
	// }
	// global.GVA_LOG.Debug(fmt.Sprintf("2222222 data %s", string(data)))
	// now, err := time.Parse(`"`+TimeFormat+`"`, string(data))
	// global.GVA_LOG.Debug(fmt.Sprintf("2222222 err %v", err))
	// *t = LocalTime(now)
	// return
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	// global.GVA_LOG.Debug("2222222 MarshalJSON")
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TimeFormat)
	b = append(b, '"')
	return b, nil
}

func (t LocalTime) Value() (driver.Value, error) {
	if t.String() == "0001-01-01 00:00:00" {
		return nil, nil
	}
	return []byte(time.Time(t).Format(TimeFormat)), nil
}

func (t *LocalTime) Scan(v interface{}) error {
	tTime, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", v.(time.Time).String())
	*t = LocalTime(tTime)
	return nil
}

func (t LocalTime) String() string {
	return time.Time(t).Format(TimeFormat)
}
