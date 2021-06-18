package main

import (
	"math"
	"sync"
)

type tuple struct{ x, y int }

type ha_array [][][][]float64

/***
*	It's the ha_array that is a bit stupid
* 	Using the fact that we have sparse data
* 	and use a different datastructure and
* 	just padd it instead would be an option
*	I only care about the indicies.
 */
type histogram struct {
	n_cells, width, height, dim int
	data                        *ha_array
	mc                          *[][][]event
	idx                         *[][]int
	evc                         *[][]int
}

func init_datastructure(prm *params) histogram {
	var pp = *prm
	n_cells := (pp.width / pp.K) * (pp.height / pp.K)

	mc := make([][][]event, 2)
	mc[0] = make([][]event, n_cells)
	mc[1] = make([][]event, n_cells)
	evc := make([][]int, 2)
	evc[0] = make([]int, n_cells)
	evc[1] = make([]int, n_cells)

	var arr = construct4darr(n_cells, 2*pp.R+1, 2*pp.R+1)
	var cidx = cell_idx(pp.width, pp.height, pp.K)

	ds := histogram{
		n_cells: n_cells,
		width:   pp.width,
		height:  pp.height,
		data:    &arr,
		mc:      &mc,
		idx:     &cidx,
		evc:     &evc,
		dim:     2 * (n_cells) * (2*pp.R + 1) * (2*pp.R + 1),
	}
	return ds
}

func cell_idx(width, height, K int) (out [][]int) {
	var arr = make([][]int, width)
	var cell_width = width / K
	// var cell_height = height / K
	for i := range arr {
		arr[i] = make([]int, height)
		for j := range arr[i] {
			var p_row = i / K
			var p_col = j / K
			arr[i][j] = p_row*cell_width + p_col
		}
	}
	return arr
}

func normalize(ds *histogram) *[]float64 {
	var vector = make([]float64, ds.dim)
	var hst = (*ds).data
	var evc = (*ds).evc
	var count int

	// This I want to dissappear
	// Better datastructure to begin with
	// go routines bring overhead here.
	// Maybe maybe no, I always have the pointer
	// don't care when it finishes as long as the reference is intact
	for i := range *hst {
		for j := range (*hst)[i] {
			for z := range (*hst)[i][j] {
				for k := range (*hst)[i][j][z] {
					vector[count] = (*hst)[i][j][z][k] / (float64((*evc)[i][j]) + 1e-9)
					count++
				}
			}
		}
	}
	ds = nil
	return &vector
}

func compute_time_surface(e event, mce *[]event, prm *params) *[][]float64 {
	var R = (*prm).R
	var tau = (*prm).tau
	tau = float64(tau)
	var time_surface [][]float64 = make([][]float64, 2*R+1)
	// Tried a map but this is faster
	for idx, _ := range time_surface {
		time_surface[idx] = make([]float64, 2*R+1)
	}

	for _, e_i := range *mce {
		var delta_t = e.t - e_i.t
		var num = math.Exp(-delta_t / tau)

		// center
		var y_shift = e_i.y - (e.y - R)
		var x_shift = e_i.x - (e.x - R)
		time_surface[y_shift][x_shift] += num
	}
	return &time_surface
}

func multiply(h *[][]float64, ts *[][]float64) {
	for i := range *h {
		for j := range (*h)[i] {
			((*h)[i][j]) += ((*ts)[i][j])
		}
	}
}

func process(e event, prm *params, ds *histogram) {
	var idx int = (*((*ds).idx))[e.x][e.y]
	var mce *[]event = &((*(*ds).mc)[e.p][idx])

	(*((*ds).evc))[e.p][idx]++ // increment by one, used later to normalize the cell

	if len(*mce) == 0 {
		*mce = []event{e}
	} else {
		*mce = append(*mce, e)
	}
	// Ignore events that are too far back in time
	var bp = e.t - (*prm).time_window
	var i = 0

	for ((*mce)[i].t < bp) && i < (len(*mce)-1) {
		i++
	}
	*mce = (*mce)[i:]
	var time_surface = compute_time_surface(e, mce, prm)
	multiply(&((*(*ds).data)[e.p][idx]), time_surface)
}

func process_all(es []event, prm params, ch chan *[]float64, wg *sync.WaitGroup, wgg *sync.WaitGroup) {

	var ds = init_datastructure(&prm)
	var count = 0
	for _, e := range es {
		count++
		process(e, &prm, &ds)
	}

	ch <- normalize(&ds)

	wgg.Done()
	wg.Done()
}
