package main

import (
	//"fmt"
	"fmt"
	"io/ioutil"
	"path"
	//"log"
)

type event struct {
	x, y int
	t    float64
	p    int // 0 or 1
}

func construct4darr(x, y, z int) ha_array {
	out := make(ha_array, 2)
	for pp := range out {
		out[pp] = make([][][]float64, x)
		for xx := range out[pp] {
			out[pp][xx] = make([][]float64, y)
			for yy := range out[pp][xx] {
				out[pp][xx][yy] = make([]float64, z)
			}
		}
	}
	return out
}

func load_data(fnames []string) [][]byte {
	buffer := make([][]byte, len(fnames))
	for idx, name := range fnames {
		data, err := ioutil.ReadFile(name)
		if err != nil {
			panic(err)
		}
		buffer[idx] = data
	}
	return buffer
}

func load_files(filename string) map[int][]string {
	dirs, err := ioutil.ReadDir(filename)
	if err != nil {
		panic(err)
	}

	all_files := make(map[int][]string)
	count := 0

	for _, d := range dirs {
		name := path.Join(filename, d.Name())
		fmt.Println(name)

		f, err := ioutil.ReadDir(name)
		if err != nil {
			panic(err)
		}
		s := make([]string, len(f))
		for idx := range s {
			s[idx] = path.Join(name, f[idx].Name())
		}
		all_files[count] = s
		count++
	}
	return all_files
}

// bit 39 - 32: Xaddress (in pixels)
// bit 31 - 24: Yaddress (in pixels)
// bit 23: Polarity (0 for OFF, 1 for ON)
// bit 22 - 0: Timestamp (in microseconds)
// The videos below show the conversion process in action and some of the resulting recordings.

func process_single(buffer []byte) []event {
	// This is taken from Gochard website, translated from the python file.
	// Adapted, closer to linear, original was reading buffer mltpl times.
	x_address := ((1 << 8) - 1) << 32
	y_address := ((1 << 8) - 1) << 24
	p_address := (1 << 23)
	t_address := (1 << 23) - 1

	var es = make([]event, 0)

	var idx = 0
	var time_increment = (1 << 13)
	var multiple = 0

	for idx < len(buffer) {
		var bits int
		chunk := buffer[idx:(idx + 5)]
		for _, c := range chunk {
			bits |= int(c)
			bits <<= 8
		}
		bits >>= 8 // Yolo

		x := (bits & x_address) >> 32
		y := (bits & y_address) >> 24
		p := (bits & p_address) >> 23
		t := (bits & t_address) // ms

		if y == 240 {
			// skip overflow
			// corrupted data
			multiple++
			idx += 5
			continue
		}
		t += multiple * time_increment

		var e = event{
			x: x,
			y: y,
			p: p,
			t: float64(t) * 1e-3,
		}

		es = append(es, e)
		idx += 5
	}
	return es
}

func process_buffer(buffer [][]byte) (es_array [][]event) {
	es_array = make([][]event, len(buffer))
	for idx, bfr := range buffer {
		es_array[idx] = process_single(bfr)
	}
	return es_array
}
