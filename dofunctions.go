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
	Token   string `json:"token"`
	Region  string `json:"region"`
	Size    string `json:"size"`
	Imageid string `json:"image"`
}

type domclient struct {
	tr     *oauth.Transport
	client *godo.Client
	config Configuration
}

// ConfigureDOM configures the client with token
func ConfigureDOM(path string) {
	var token string
	var region string
	var imageid string
	var size string

	fmt.Println("\nYou need to obtain Personal Access Token from DigitalOcean")
	fmt.Printf("This can be generated from https://cloud.digitalocean.com/settings/applications \n\n")
	fmt.Printf("Token: ")
	fmt.Scan(&token)
	fmt.Printf("Token entered: %s\n", token)

	fmt.Println("Please enter 0 for default setting")
	fmt.Printf("Enter your default region ID (0 for default to New York(nyc3)): ")
	fmt.Scan(&region)
	fmt.Printf("Enter your default image ID (0 for default to 11836690 (Ubuntu 14.04 x64)): ")
	fmt.Scan(&imageid)
	fmt.Printf("Enter your default size ID (0 for default to 512MB): ")
	fmt.Scan(&size)

	if region == "0" {
		region = "nyc3"
	}
	if imageid == "0" {
		imageid = "11836690"
	}
	if size == "0" {
		size = "512MB"
	}

	config := Configuration{token, region, size, imageid}
	configB, _ := json.MarshalIndent(config, " ", "  ")
	fmt.Printf("%s\n", configB)
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
	d.config = config
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

func (d *domclient) CreateDropletFromImage(name string, region string, size string, imageID string) (*godo.Droplet, error) {

	imageid, _ := strconv.Atoi(imageID)
	fmt.Printf("Param : %s %s %s %d\n", name, region, size, imageid)
	createRequest := &godo.DropletCreateRequest{
		Name:   name,
		Region: region,
		Size:   size,
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

func (d *domclient) ListSizes() error {
	sizeInfo, _, err := d.client.Sizes.List(nil)
	if err != nil {
		fmt.Println("Error in retriving Size info.")
		return err
	}
	fmt.Printf("%5s %5s %10s %10s %8s %20s\n", "Slug", "Availbale", "Disk", "Memory", "Vcpus", "Price (Monthy)")
	for _, s := range sizeInfo {
		fmt.Printf("%5s %5v %12d GB %8d MB %5d %14.2f USD\n", s.Slug, s.Available, s.Disk, s.Memory, s.Vcpus, s.PriceMonthly)

	}
	return nil
}
