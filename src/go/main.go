package main

import (
	"fmt"
	"sync"
	"time"
)

type params struct {
	R, K, width, height       int
	tau, delta_t, time_window float32
}

func batch_process(ev *[][]event, par *params, ch chan []float32, wg *sync.WaitGroup) {
	defer close(ch)
	for _, e := range *ev {
		go process_all(e, *par, ch, wg)
	}
	wg.Wait()
}

func sum(arr []float32) float32 {
	var sum float32
	for _, i := range arr {
		sum += i
	}
	return sum
}

/**
* TODO:
* The ha_array is causing some trouble
* Would be nice to avoid 4d
* I should be able to exploit the sparse nature of the data
* And only track the changes.
*
* Check the pointers I think some [de?]referencing is redundant
* Artifact of me not knowing golang all that well.
*
* Construct the svm and do some prediction, that is the final test anyway.
**/

func main() {
	filename := "../../dataset/train/"

	/**
	* Reasonable values for nmnist
	* width height = 35 (leaves some padding in the histogram but depends on K and R)
	* (R, K = 7, tau ~ 1/2, delta_t ~.1, time_window depend on the unit)
	* (width height should be evenly divisible by K)
	**/
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
	data := load_data(all_files[0]) // only the zeros for testing

	ev := process_buffer(data)

	hst := make([][]float32, len(ev))

	var ch = make(chan []float32)

	var wg sync.WaitGroup
	wg.Add(len(hst))

	go batch_process(&ev, &par, ch, &wg)

	count := 0
	start := time.Now()

	for p := range ch {
		hst[count] = p
		count++
	}
	// Sanity checks
	for _, arr := range hst {
		fmt.Println(sum(arr))
	}
	done := time.Since(start).Seconds()
	fmt.Println(done)
	fmt.Println(count)
}
