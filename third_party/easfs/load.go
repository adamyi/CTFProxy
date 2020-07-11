package main

import "github.com/shirou/gopsutil/load"

import "fmt"

func loadAverage() string {
	loadstr := "unknown"
	l, err := load.Avg()
	if err == nil {
		loadstr = fmt.Sprintf("%.2f %.2f %.2f", l.Load1, l.Load5, l.Load15)
	}
	return loadstr
}
