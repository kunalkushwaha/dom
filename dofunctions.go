package main

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"fmt"
	"github.com/digitalocean/godo"
	"io/ioutil"
)

type Configuration struct {
	Token string `json:"token"`
}

type domclient struct {
	tr     *oauth.Transport
	client *godo.Client
}

func SetupClient(filepath string) *domclient {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("Configuration file not found")
		return nil
	}
	config := Configuration{}
	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Printf("Error in config file, Please reconfigure dom: %v\n", err)
		return nil
	}
	d := &domclient{}
	d.tr = &oauth.Transport{Token: &oauth.Token{AccessToken: config.Token}}
	d.client = godo.NewClient(d.tr.Client())
	return d
}

func (d *domclient) DropletList(filter string) ([]godo.Droplet, error) {
	// create a list to hold our droplets
	list := []godo.Droplet{}

	// create options. initially, these will be blank
	opt := &godo.ListOptions{}
	for {
		droplets, resp, err := d.client.Droplets.List(opt)
		if err != nil {
			return nil, err
		}

		// append the current page's droplets to our list
		for _, d := range droplets {
			list = append(list, d)
		}

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}
	fmt.Println(list)
	return list, nil
}

func (d *domclient) ImageList(fliter string, user bool) ([]godo.Image, error) {
	list := []godo.Image{}
	var images []godo.Image
	var err error
	if user == true {
		fmt.Println("User ImageLists ")
		images, _, err = d.client.Images.ListUser(nil)
	} else {
		fmt.Println("Global ImageLists ")
		images, _, err = d.client.Images.List(nil)
	}
	if err != nil {
		return nil, err
	}

	for _, img := range images {
		list = append(list, img)
		fmt.Printf("ID : %6d, Name:%s: %-20s , type : %s\n", img.ID, img.Distribution, img.Name, img.Type)

	}

	return list, nil

}

func (d *domclient) RebuildByImageID(imageId int) error {

	action, _, err := d.client.DropletActions.RebuildByImageID(imageId, imageId)
	if err != nil {
		fmt.Printf("DropletActions.Restore returned error: %v\n", err)
		return err
	}
	fmt.Printf("DropletActions.Restore returned %+v \n", action)
	return nil
}

func (d *domclient) CreateDropletFromImage(imageID int) error {
	createRequest := &godo.DropletCreateRequest{
		Name:   "testdr1",
		Region: "nyc3",
		Size:   "512MB",
		Image: godo.DropletCreateImage{
			ID: imageID,
		},
	}

	root, _, err := d.client.Droplets.Create(createRequest)
	if err != nil {
		fmt.Printf("Droplets.Create returned error: %v", err)
		return err
	}
	if id := root.Droplet.ID; id != imageID {
		fmt.Printf("expected id '%d', received '%d'\n", imageID, id)
	}
	return nil
}
