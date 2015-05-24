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

func TestImageUserList(t *testing.T) {
	path, _ := homedir.Expand("~/.dom.json")
	d := SetupClient(path)
	if d == nil {
		t.Errorf("Cannot setup the client")
		return
	}
	_, err := d.ImageList(" ", true)
	if err != nil {
		t.Errorf("Unable to fetch User ImageList\n")
	}
}
func TestImageGlobalList(t *testing.T) {
	path, _ := homedir.Expand("~/.dom.json")
	d := SetupClient(path)
	if d == nil {
		t.Errorf("Cannot setup the client")
		return
	}
	_, err := d.ImageList(" ", false)
	if err != nil {
		t.Errorf("Unable to fetch User ImageList\n")
	}
}


func TestRestoreImage(t *testing.T) {
	path, _ := homedir.Expand("~/.dom.json")
	d := SetupClient(path)
	if d == nil {
		t.Errorf("Cannot setup the client")
		return
	}
	err := d.CreateDropletFromImage(11876597)
	if err != nil {
		t.Errorf("Unable to create Droplet from Image\n")
	}

}
