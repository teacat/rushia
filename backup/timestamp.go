package rushia

import "strings"

// Timestamp 是一個資料庫的時間戳輔助函式。
type Timestamp struct {
	value interface{}
	query string
}

// IsDate 會確保欄位的時間戳是某個指定的年月日期。
func (t Timestamp) IsDate(date string) Timestamp {
	t.query = "DATE(FROM_UNIXTIME(%s)) = %s "
	t.value = date
	return t
}

// IsYear 會確保欄位的時間戳是某個指定的年份。
func (t Timestamp) IsYear(year int) Timestamp {
	t.query = "YEAR(FROM_UNIXTIME(%s)) = %s "
	t.value = year
	return t
}

// IsMonth 會確保欄位的時間戳是某個指定的月份，
// 可以傳入字串的 `January`、`Jan` 或正整數的 `1` 作為月份。
func (t Timestamp) IsMonth(month interface{}) Timestamp {
	t.query = "MONTH(FROM_UNIXTIME(%s)) = %s "
	switch v := month.(type) {
	case int:
		t.value = v
	case string:
		list := map[string]int{
			"january":   1,
			"february":  2,
			"march":     3,
			"april":     4,
			"may":       5,
			"june":      6,
			"july":      7,
			"august":    8,
			"september": 9,
			"october":   10,
			"november":  11,
			"december":  12,
			"jan":       1,
			"feb":       2,
			"mar":       3,
			"apr":       4,
			"jun":       6,
			"jul":       7,
			"aug":       8,
			"sep":       9,
			"oct":       10,
			"nov":       11,
			"dec":       12,
		}
		t.value = list[strings.ToLower(v)]
	}
	return t
}

// IsDay 會確保欄位的時間戳是某個指定的天數。
func (t Timestamp) IsDay(day int) Timestamp {
	t.query = "DAY(FROM_UNIXTIME(%s)) = %s "
	t.value = day
	return t
}

// IsWeekday 會確保欄位的時間戳是某個指定的星期。
// 可以傳入字串的 `Monday`、`Mon` 或正整數的 `1` 作為星期。
func (t Timestamp) IsWeekday(weekday interface{}) Timestamp {
	t.query = "WEEKDAY(FROM_UNIXTIME(%s)) = %s "
	switch v := weekday.(type) {
	case int:
		t.value = v
	case string:
		list := map[string]int{
			"monday":    1,
			"tuesday":   2,
			"wednesday": 3,
			"thursday":  4,
			"friday":    5,
			"saturday":  6,
			"sunday":    7,
			"mon":       1,
			"tue":       2,
			"wed":       3,
			"thu":       4,
			"fri":       5,
			"sat":       6,
			"sun":       7,
		}
		t.value = list[strings.ToLower(v)]
	}
	return t
}

// IsHour 會確保欄位的時間戳是某個指定的時數。
func (t Timestamp) IsHour(hour int) Timestamp {
	t.query = "HOUR(FROM_UNIXTIME(%s)) = %s "
	t.value = hour
	return t
}

// IsMinute 會確保欄位的時間戳是某個指定的分鐘。
func (t Timestamp) IsMinute(minute int) Timestamp {
	t.query = "MINUTE(FROM_UNIXTIME(%s)) = %s "
	t.value = minute
	return t
}

// IsSecond 會確保欄位的時間戳是某個指定的秒數。
func (t Timestamp) IsSecond(second int) Timestamp {
	t.query = "SECOND(FROM_UNIXTIME(%s)) = %s "
	t.value = second
	return t
}
