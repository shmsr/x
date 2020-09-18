package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
)

var r = flag.String(
	"range",
	"1,100",
	"enter range [l,r] where l <= r in the format \"l,r\"",
)

type errs struct {
	u string
	w string
}

func (e errs) Error() string {
	return fmt.Sprintf("Error in %s due to %s", e.u, e.w)
}

// streamOdd gets the next odd numbers and sends it ch, immdiately
func streamOdd(l, r int, ch chan<- int) {
	if l%2 == 0 {
		l++
	}
	for i := l; i <= r; i += 2 {
		if i%2 != 0 {
			ch <- i
		}
	}
	close(ch)
}

// streamEven gets the next even numbers and sends it ch, immdiately
func streamEven(l, r int, ch chan<- int) {
	if l%2 != 0 {
		l++
	}
	for i := l; i <= r; i += 2 {
		if i%2 == 0 {
			ch <- i
		}
	}
	close(ch)
}

func main() {
	flag.Parse()

	var rs []string
	// Split the string by comma(s)
	if rs = strings.Split(*r, ","); len(rs) != 2 {
		log.Fatalln(errs{"flag range", "invalid arguments"})
	}

	var err error
	var lR, rR int
	// Convert string to int
	if lR, err = strconv.Atoi(rs[0]); err != nil {
		log.Fatalln(errs{"lower bound", err.Error()})
	}
	// Convert string to int
	if rR, err = strconv.Atoi(rs[1]); err != nil {
		log.Fatalln(errs{"upper bound", err.Error()})
	}

	// Make channels
	var cho = make(chan int, 1)
	var che = make(chan int, 1)

	// Spawn goroutines
	go streamOdd(lR, rR, cho)
	go streamEven(lR, rR, che)

	// Determine the length of odd and even streams
	// using four cases
	switch {
	case lR%2 == 0 && rR%2 == 0: // [even,even]
		// cho < che
		for odd := range cho {
			fmt.Println(<-che)
			fmt.Println(odd)
		}
		fmt.Println(<-che)
	case lR%2 != 0 && rR%2 != 0: // [odd,odd]
		// che < cheo
		for even := range che {
			fmt.Println(<-cho)
			fmt.Println(even)
		}
		fmt.Println(<-cho)
	case lR%2 != 0 && rR%2 == 0: // [odd,even]
		// cho == che
		for even := range che {
			fmt.Println(<-cho)
			fmt.Println(even)
		}
	case lR%2 == 0 && rR%2 != 0: // [even,odd]
		// cho == che
		for odd := range cho {
			fmt.Println(<-che)
			fmt.Println(odd)
		}
	}
}
