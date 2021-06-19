package main

import (
	"fmt"
	"io/ioutil"
	"path"
)

type event struct {
	x, y int8
	t    float32
	p    int8 // 0 or 1
}

func construct4darr(x, y, z int) ha_array {
	out := make(ha_array, 2)
	for pp := range out {
		out[pp] = make([][][]float32, x)
		for xx := range out[pp] {
			out[pp][xx] = make([][]float32, y)
			for yy := range out[pp][xx] {
				out[pp][xx][yy] = make([]float32, z)
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

func process_single(buffer []byte) []event {
	/**
	* Structure:
	* bit 39 - 32: Xaddress (in pixels)
	* bit 31 - 24: Yaddress (in pixels)
	* bit 23: Polarity (0 for OFF, 1 for ON)
	* bit 22 - 0: Timestamp (in microseconds)
	*
	* Taken from Gochard's website, adapted from python.
	* Made it less for-loopy
	**/

	x_address := ((1 << 8) - 1) << 32
	y_address := ((1 << 8) - 1) << 24
	p_address := (1 << 23)
	t_address := (1 << 23) - 1

	var es = make([]event, 0)

	const time_increment = (1 << 13)
	const n_bytes = 5

	// Select time unit:
	// 1e-3 gives ms
	// 1e-6 gives s
	const unit_conv = 1e-3

	var offset int
	var multiple int

	for offset < len(buffer) {
		var bits int
		chunk := buffer[offset:(offset + n_bytes)]
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
			offset += n_bytes
			continue
		}
		t += multiple * time_increment

		var e = event{
			x: int8(x),
			y: int8(y),
			p: int8(p),
			t: float32(t) * unit_conv,
		}

		es = append(es, e)
		offset += n_bytes
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
