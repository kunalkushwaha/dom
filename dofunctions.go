package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"code.google.com/p/goauth2/oauth"
	"github.com/digitalocean/godo"
)

// Configuration dom configutation struct
type Configuration struct {
	Token string `json:"token"`
}

type domclient struct {
	tr     *oauth.Transport
	client *godo.Client
}

// ConfigureDOM configures the client with token
func ConfigureDOM(path string) {
	var token string

	fmt.Println("\nYou need to obtain Personal Access Token from DigitalOcean")
	fmt.Printf("This can be generated from https://cloud.digitalocean.com/settings/applications \n\n")
	fmt.Printf("Token: ")
	fmt.Scan(&token)
	fmt.Printf("Token entered: %s\n", token)

	config := Configuration{token}
	configB, _ := json.MarshalIndent(config, " ", "  ")

	ioutil.WriteFile(path, configB, 0644)

}

// SetupClient for dom.
func SetupClient(filepath string) *domclient {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("Error: dom is not configured !! \nPlease run \"dom authorize\"")
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
	droplets, _, err := d.client.Droplets.List(opt)
	if err != nil {
		return nil, err
	}

	//Print the List
	fmt.Printf("%5s %16s %12s %23s %30s\n", "Droplet ID", "NAME", "Region", "IP", "Creation time")
	for _, d := range droplets {
		fmt.Printf("%5d %20s %10s %30s %30s\n", d.ID, d.Name, d.Region.Slug, d.Networks.V4[0].IPAddress, d.Created)
	}
	return list, nil
}

func (d *domclient) ImageList(fliter string, global bool) ([]godo.Image, error) {
	list := []godo.Image{}
	var images []godo.Image
	var err error
	if global == false {
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

func (d *domclient) RebuildByImageID(imageID int) error {

	action, _, err := d.client.DropletActions.RebuildByImageID(imageID, imageID)
	if err != nil {
		fmt.Printf("DropletActions.Restore returned error: %v\n", err)
		return err
	}
	fmt.Printf("DropletActions.Restore returned %+v \n", action)
	return nil
}

func (d *domclient) CreateDropletFromImage(imageID string) (*godo.Droplet, error) {

	imageid, _ := strconv.Atoi(imageID)

	createRequest := &godo.DropletCreateRequest{
		Name:   "testdr1",
		Region: "nyc3",
		Size:   "512MB",
		Image: godo.DropletCreateImage{
			ID: imageid,
		},
	}

	root, _, err := d.client.Droplets.Create(createRequest)
	if err != nil {
		fmt.Printf("Droplets.Create returned error: %v", err)
		return nil, err
	}

	return root, nil
}

func (d *domclient) DestroyDroplet(dropletID string) error {
	dropletid, _ := strconv.Atoi(dropletID)

	_, err := d.client.Droplets.Delete(dropletid)
	return err
}

func (d *domclient) DropletShutdown(dropletID string) error {
	dropletid, _ := strconv.Atoi(dropletID)

	_, _, err := d.client.DropletActions.Shutdown(dropletid)
	return err
}

func (d *domclient) DropletRestart(dropletID string) error {
	dropletid, _ := strconv.Atoi(dropletID)

	_, _, err := d.client.DropletActions.Reboot(dropletid)
	return err
}

func (d *domclient) ListRegions() error {
	regions, _, err := d.client.Regions.List(nil)
	if err != nil {
		fmt.Printf("Cannot retrive regions: %v\n", err)
		return err
	}

	fmt.Printf("%5s %20s %20s\n", "Slug", "Name", "Available")
	for _, region := range regions {
		fmt.Printf("%5s %20s %20v\n", region.Slug, region.Name, region.Available)
	}
	return nil
}

func (d *domclient) DropletInfo(dropletID string) error {
	dropletid, _ := strconv.Atoi(dropletID)

	info, _, err := d.client.Droplets.Get(dropletid)
	if err != nil {
		fmt.Printf("Unable to fetch droplet info : %v\n", err)
		return err
	}

	fmt.Printf("\n%12s : %s\n", "Name", info.Name)
	fmt.Printf("%12s : %s ( %s )\n", "Region", info.Region.Name, info.Region.Slug)
	fmt.Printf("%12s : %d GB\n", "Size", info.Disk)
	fmt.Printf("%12s : %d MB\n", "Memory", info.Memory)
	fmt.Printf("%12s : %d\n", "Vcpus", info.Vcpus)
	fmt.Printf("%12s : %s\n", "IP Address", info.Networks.V4[0].IPAddress)
	fmt.Printf("%12s : %s\n", "Status", info.Status)

	return nil
}
