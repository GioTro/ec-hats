package main

import (
	"fmt"
	"sync"
	"time"
)

type params struct {
	R, K, width, height       int
	tau, delta_t, time_window float64
}

func signal_global(wg *sync.WaitGroup, ch chan *[]float64) {
	defer close(ch)
	wg.Wait()
}

func batch_process(ev *[][]event, par *params, ch chan *[]float64, wg *sync.WaitGroup) {
	// const batch_size = 10
	//start := time.Now()
	//var wg sync.WaitGroup
	//var count int
	//wg.Add(len(*ev))
	// var i interface{} = idx
	for _, e := range *ev {
		go process_all(e, *par, ch, wg)
		//count++
	}
	//wg.Wait()
	//stop := float64(time.Since(start).Milliseconds())
	//fmt.Println("Average", stop/float64(count), "ms") // 1.3ms per example
}

// func init() {
// 	numcpu := runtime.NumCPU()
// 	fmt.Println(numcpu)
// 	runtime.GOMAXPROCS(numcpu)
// }

func main() {
	// For mnist width height = 35, (R, K = 7 is good), (tau ~ 1/2, delta_t ~.1, time_window depend on the unit)
	// This is a bit messy
	filename := "../../dataset/train/"

	var par = params{
		R:           7,
		K:           7,
		width:       35,
		height:      35,
		tau:         .5,
		delta_t:     .1,
		time_window: 1,
	}

	all_files := load_files(filename)
	data := load_data(all_files[0])

	ev := process_buffer(data)

	fmt.Println("I'm here")
	hst := make([][]float64, len(ev))

	var ch = make(chan *[]float64)

	var wg sync.WaitGroup
	wg.Add(len(hst))
	// wgg.Add(1)
	go signal_global(&wg, ch)
	//var count int
	go batch_process(&ev, &par, ch, &wg)

	count := 0
	start := time.Now()
	//time.Sleep(20 * time.Millisecond)
	for p := range ch {
		hst[count] = *p
		count++
		//var max float64
		// for _, n := range hst[count] {
		// 	if n > max {
		// 		max = n
		// 	}
		//}
		// if count%100 == 0 {
		// 	fmt.Println(len(hst) - count)
		// }
		//fmt.Println(max)
	}
	done := time.Since(start).Seconds()
	fmt.Println(done)
}
