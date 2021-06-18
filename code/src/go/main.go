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
	wg.Wait()

	fmt.Println("Not closed waiting")
	defer close(ch)
	fmt.Println("Closed")
}

func batch_process(ev *[][]event, par *params, ch chan *[]float64, wgg *sync.WaitGroup) {
	idx := 0
	const batch_size = 5
	var wg sync.WaitGroup
	var count int

	// var i interface{} = idx
	start := time.Now()
	for _, e := range *ev {
		wg.Add(1)
		go process_all(e, *par, ch, &wg, wgg)
		wg.Wait()
		if idx > batch_size {
			wg.Wait()
			idx = 0
		}
		idx++
		count++
	}
	stop := float64(time.Since(start).Microseconds())
	fmt.Println("Average", stop/float64(count), "micro-seconds")
}

func main() {
	// For mnist width height = 35, (R, K = 7 is good), (tau ~ 1/2, delta_t ~.1, time_window depends on units.)
	filename := "../../dataset/train/"

	var par = params{
		R:           7,
		K:           7,
		width:       35,
		height:      35,
		tau:         .5,
		delta_t:     .1,
		time_window: .1,
	}

	all_files := load_files(filename)
	data := load_data(all_files[0])

	ev := process_buffer(data)

	fmt.Println("I'm here")
	hst := make([][]float64, len(ev))

	var ch = make(chan *[]float64)

	var wgg sync.WaitGroup
	wgg.Add(len(hst))
	// wgg.Add(1)
	go signal_global(&wgg, ch)
	//var count int
	go batch_process(&ev, &par, ch, &wgg)

	count := 0
	start := time.Now()
	//time.Sleep(20 * time.Millisecond)
	for p := range ch {
		hst[count] = *p
		var max float64
		for _, n := range hst[count] {
			if n > max {
				max = n
			}
		}
		count++
		if count%100 == 0 {
			fmt.Println(len(hst) - count)
		}
		//fmt.Println(max)
	}
	done := time.Since(start).Seconds()
	fmt.Println(done)
}
