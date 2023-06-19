package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Assigment struct {
	ID           int64
	Shifts       []Shift
	Weekdays     [7]bool
	HashedShifts [7]int64
}

type Shift struct {
	Start string
	End   string
}

func New(id int64, shifts []Shift, weekdays [7]bool) Assigment {
	a := Assigment{
		ID:       id,
		Shifts:   shifts,
		Weekdays: weekdays,
	}
	var (
		hs int64
	)
	for _, shift := range a.Shifts {
		s := getHaftHourIndex(shift.Start)
		e := getHaftHourIndex(shift.End)
		var b int64 = 1 << (48 - s)
		fmt.Printf("s: %d - e: %d\n", s, e)
		fmt.Printf("bb: %064b\n", b)
		for idx := s; idx < e; idx++ {
			hs += b
			b >>= 1
		}
		fmt.Printf("hs: %064b\n", hs)
		for wIdx, weekday := range a.Weekdays {
			if !weekday {
				continue
			}
			a.HashedShifts[wIdx] = hs
		}

	}
	return a
}

func main() {
	as := []Assigment{
		New(1, []Shift{{"08:30", "12:00"}}, [7]bool{false, true, false, true, false, true, false}), // 8h30-12h | t2, t4, t6
		New(2, []Shift{{"12:00", "18:00"}}, [7]bool{false, true, false, true, false, true, false}), // 12h-18h | t2, t4, t6
		New(3, []Shift{{"08:00", "18:00"}}, [7]bool{true, false, true, false, true, false, true}),  // 8h-18h | cn, t3, t5
		New(4, []Shift{{"17:30", "22:00"}}, [7]bool{false, true, false, true, false, true, false}), // 17h30-22h | t2, t4, t6 => must duplicate
	}

	nwd := time.Now().Weekday()
	hsBf := as[0].HashedShifts[nwd]
	fmt.Printf("\n////////////////////START///////////////////\n")
	for idx := 1; idx < len(as); idx++ {
		fmt.Printf("checking assigment with id %d and shifts %+v\n", as[idx].ID, as[idx].Shifts)
		fmt.Printf("hsBf %064b\n", hsBf)
		fmt.Printf("this %064b\n", as[idx].HashedShifts[nwd])
		if as[idx].HashedShifts[nwd]&hsBf != 0 {
			panic("duplicate")
		}
		hsBf |= as[idx].HashedShifts[nwd]
	}
	return
}

func getHaftHourIndex(timeStr string) int64 {
	t := strings.Split(timeStr, ":")
	hour, _ := strconv.ParseInt(t[0], 10, 64)
	min, _ := strconv.ParseInt(t[1], 10, 64)
	fmt.Printf("shifthour: %d - min: %d - haftHour: %d\n", hour, min, hour*2+min/30)
	return hour*2 + min/30
}
