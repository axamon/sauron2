package cmd

import (
	"fmt"
	"time"
)

func ExampleIsfobLunedimattina() {
	date := time.Date(2018, 6, 4, 6, 0, 0, 0, time.Local)
	ok := isfob(date, 18)
	fmt.Println(date.Weekday())
	fmt.Println(date.Hour())
	fmt.Println(ok)
	//Output:
	//Monday
	//6
	//true
}

func ExampleIsfobLunedi() {
	date := time.Date(2018, 6, 4, 8, 0, 0, 0, time.Local)
	ok := isfob(date, 18)
	fmt.Println(date.Weekday())
	fmt.Println(date.Hour())
	fmt.Println(ok)
	//Output:
	//Monday
	//8
	//false
}

func ExampleIsfobLunedifob() {
	date := time.Date(2018, 6, 4, 18, 0, 0, 0, time.Local)
	ok := isfob(date, 18)
	fmt.Println(date.Weekday())
	fmt.Println(date.Hour())
	fmt.Println(ok)
	//Output:
	//Monday
	//18
	//true
}

func ExampleIsfobDomenica() {
	date := time.Date(2018, 6, 3, 6, 0, 0, 0, time.Local)
	ok := isfob(date, 18)
	fmt.Println(date.Weekday())
	fmt.Println(date.Hour())
	fmt.Println(ok)
	//Output:
	//Sunday
	//6
	//true
}

func ExampleIsfobSabato() {
	date := time.Date(2018, 6, 2, 12, 0, 0, 0, time.Local)
	ok := isfob(date, 18)
	fmt.Println(date.Weekday())
	fmt.Println(ok)
	//Output:
	//Saturday
	//true
}
