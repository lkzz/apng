package apng

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestAPNG(t *testing.T) {
	// filename := "world_cup_2014_42.png"
	filename := "apng_file.png"
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("os.Open(%s) error(%v)", filename, err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatalf("ioutil.ReadAll error(%v)", err)
	}
	if !Hit(data) {
		t.Fatalf("%s is an apng file in fact", filename)
	}
}

func TestPNG(t *testing.T) {
	filename := "png_file.png"
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("os.Open(%s) error(%v)", filename, err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatalf("ioutil.ReadAll error(%v)", err)
	}
	if Hit(data) {
		t.Fatalf("%s is an apng file in fact", filename)
	}
}
