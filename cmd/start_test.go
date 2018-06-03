package cmd

import "time"

func Exampleisfob() {
	isfob(time.Date(2018, 5, 6, 2, 0, 0, 0, time.Local))
	//Output: true
}
