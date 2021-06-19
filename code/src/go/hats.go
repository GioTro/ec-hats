package main

import (
	"math"
	"sync"
)

type ha_array [][][][]float32

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

func normalize(ds *histogram) []float32 {
	var vector = make([]float32, ds.dim)
	var hst = (*ds).data
	var evc = (*ds).evc
	var count int

	// Make this dissappear, better dstructure?
	for i := range *hst {
		for j := range (*hst)[i] {
			for z := range (*hst)[i][j] {
				for k := range (*hst)[i][j][z] {
					vector[count] = (*hst)[i][j][z][k] / (float32((*evc)[i][j]) + 1e-9)
					count++
				}
			}
		}
	}
	return vector
}

func compute_time_surface(e event, mce *[]event, prm *params, hst *[][]float32) {
	/** This version computes the cum sum for the histogram instead of the time surface,
	The histogram gets normalized later. This order of operation gives a substantial speed up */
	var R = int8((*prm).R)
	var tau = float64((*prm).tau)

	for _, e_i := range *mce {
		var delta_t = float64(e.t - e_i.t)
		var num = math.Exp(-delta_t / tau)

		// center
		var y_shift = e_i.y - (e.y - R)
		var x_shift = e_i.x - (e.x - R)
		(*hst)[y_shift][x_shift] += float32(num)
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
	compute_time_surface(e, mce, prm, &((*(*ds).data)[e.p][idx]))
}

func process_all(es []event, prm params, ch chan []float32, wg *sync.WaitGroup) {
	defer wg.Done()
	var ds = init_datastructure(&prm)
	for _, e := range es {
		process(e, &prm, &ds)
	}
	var out = normalize(&ds)
	ch <- out
}
