package main

import (
	"code.google.com/p/goauth2/oauth"
	"github.com/digitalocean/godo"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Configuration struct {
  Token   string `json:"token"`
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
	if err!= nil {
		fmt.Printf("Error in config file, Please reconfigure dom: %v\n",err )
		return nil
	}
	d := &domclient{}
	d.tr = &oauth.Transport{Token: &oauth.Token{AccessToken: config.Token}}
	d.client = godo.NewClient(d.tr.Client())
	return d
}

func (d *domclient) DropletList() ([]godo.Droplet, error) {
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
