package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/mitchellh/go-homedir"
)

/*
Retrieve a list of your droplets
$ tugboat droplets

Fuzzy name matching
You can pass a unique fragment of a droplets name for interactions throughout tugboat.
$ tugboat restart admin

SSH into a droplet
You can configure an SSH username and key path in tugboat authorize, or by changing your ~/.tugboat.
This lets you ssh into a droplet by providing it's name, or a partial match.
$ tugboat ssh admin

Create a droplet
$ tugboat create pearkes-www-002 -s 64 -i 2676 -r 2 -k 11251

Info about a droplet
$ tugboat info admin

Destroy a droplet
$ tugboat destroy pearkes-www-002

Restart a droplet
$ tugboat restart admin

Shutdown a droplet
$ tugboat halt admin

Snapshot a droplet
$ tugboat snapshot test-admin-snaphot admin

Resize a droplet
$ tugboat resize admin -s 66

List Available Images
You can list images that you have created.
$ tugboat images

Optionally, list images provided by DigitalOcean as well.
$ tugboat images --global

List Available Sizes
$ tugboat sizes

List Available Regions
$ tugboat regions

List SSH Keys
$ tugboat keys

Wait for Droplet State
Sometimes you want to wait for a droplet to enter some state, for example "off".
$ tugboat wait admin --state off


*/
const domConfigPath string = "~/.dom.json"

func main() {

	dom := cli.NewApp()
	dom.Name = "digital-ocean-manager"
	dom.Usage = "Cli based digital-ocean client"
	dom.Commands = []cli.Command{
		{
			Name:  "authorize",
			Usage: "Authorize and configure dom client",
			Action: func(c *cli.Context) {
				path, _ := homedir.Expand(domConfigPath)
				ConfigureDOM(path)
			},
		},
		{
			Name:  "images",
			Usage: "List images",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "global",
					Usage: "List global images from digital-ocean",
				},
			},
			Action: func(c *cli.Context) {
				listImages(c)
			},
		},
		{
			Name:  "create",
			Usage: "creates a droplet",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "image",
					Usage: "create droplet from this image",
				},
			},
			Action: func(c *cli.Context) {
				createDroplet(c)
			},
		},
		{
			Name:  "destory",
			Usage: "destroy a droplet",
			Action: func(c *cli.Context) {
				destroyDroplet(c)
			},
		},
		{
			Name:  "droplets",
			Usage: "list  droplets",
			Action: func(c *cli.Context) {
				listDroplets(c)
			},
		},
		{
			Name:  "regions",
			Usage: "list a regions",
			Action: func(c *cli.Context) {
				listRegions(c)
			},
		},
		{
			Name:  "resize",
			Usage: "resize a droplet",
			Action: func(c *cli.Context) {
				resizeDroplet(c)
			},
		},
		{
			Name:  "halt",
			Usage: "shutdown droplet",
			Action: func(c *cli.Context) {
				haltDroplet(c)
			},
		},
		{
			Name:  "restart",
			Usage: "restart droplet",
			Action: func(c *cli.Context) {
				restartDroplet(c)
			},
		},
		{
			Name:  "info",
			Usage: "details of droplet",
			Action: func(c *cli.Context) {
				infoDroplet(c)
			},
		},
	}
	dom.Run(os.Args)
}

func getDomClient() *domclient {

	path, _ := homedir.Expand(domConfigPath)
	d := SetupClient(path)
	if d == nil {
		return nil
	}
	return d

}

func listImages(c *cli.Context) {
	global := c.Bool("global")

	d := getDomClient()
	_, err := d.ImageList(" ", global)
	if err != nil {
		fmt.Println("Unable to fetch User ImageList")
	}

}
func createDroplet(c *cli.Context) {
	fmt.Println("NOT IMPLEMENTED !")
}
func destroyDroplet(c *cli.Context) {
	fmt.Println("NOT IMPLEMENTED !")
}
func listDroplets(c *cli.Context) {
	fmt.Println("NOT IMPLEMENTED !")
}
func listRegions(c *cli.Context) {
	fmt.Println("NOT IMPLEMENTED !")
}
func resizeDroplet(c *cli.Context) {
	fmt.Println("NOT IMPLEMENTED !")
}
func haltDroplet(c *cli.Context) {
	fmt.Println("NOT IMPLEMENTED !")
}
func restartDroplet(c *cli.Context) {
	fmt.Println("NOT IMPLEMENTED !")
}
func infoDroplet(c *cli.Context) {
	fmt.Println("NOT IMPLEMENTED !")
}
