package main

import "testing"
import "github.com/mitchellh/go-homedir"

func TestDropletList(t *testing.T) {
	path, _ := homedir.Expand("~/.dom.json")
	d := SetupClient(path)
	if d == nil {
		t.Errorf("Cannot setup the client")
		return
	}
	_, err := d.DropletList(" ")
	if err != nil {
		t.Errorf("Unable to fetch DropletList\n")
	}
}

func TestImageList(t *testing.T) {
	path, _ := homedir.Expand("~/.dom.json")
	d := SetupClient(path)
	if d == nil {
		t.Errorf("Cannot setup the client")
		return
	}
	_, err := d.ImageList(" ")
	if err != nil {
		t.Errorf("Unable to fetch DropletList\n")
	}
}
