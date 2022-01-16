package datetime

import (
	"fmt"
	"time"
)

func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

func GetStringRangeOfMonth(d time.Time) (string, string) {
	d0 := GetFirstDateOfMonth(d)
	d1 := GetLastDateOfMonth(d)
	return d0.Format("20060102"), d1.Format("20060102")
}

//GetWeekOfYear
/*
获取传入日期参数，所在周一至周日的范围。并确定是该年的第几周。
*/
func GetWeekOfYear(d time.Time) (time.Time, time.Time, uint8) {
	d0 := []int{-7, 0, -1, -2, -3, -4, -5}
	d1 := []int{0, 6, 5, 4, 3, 2, 1}
	weekStart := GetZeroTime(d.AddDate(0, 0, d0[d.Weekday()]))
	weekEnd := GetZeroTime(d.AddDate(0, 0, d1[d.Weekday()]))

	d2 := []int{1, 7, 6, 5, 4, 3, 2}
	yearStart := time.Date(d.Year(), 1, 1, 0, 0, 0, 0, d.Location())
	yearFirstWeekDays := d2[yearStart.Weekday()]
	var weekIndex uint8
	if d.YearDay() <= yearFirstWeekDays {
		weekIndex = 1
	} else {
		weekIndex = uint8((d.YearDay()-yearFirstWeekDays)/7 + 1)
	}

	return weekStart, weekEnd, weekIndex
}

func GetStringWeekOfYear(d time.Time) (string, string, string) {
	s, e, i := GetWeekOfYear(d)
	w := fmt.Sprintf("%d", i)
	if len(w) < 2 {
		w = "0" + w
	}
	return s.Format("20060102"), e.Format("20060102"), w
}

func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}
