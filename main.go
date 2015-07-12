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
					Name:  "name",
					Usage: "[Mandatory] Droplet name",
				},
				cli.StringFlag{
					Name:  "size",
					Usage: "[Optional] Size of Droplet",
				},
				cli.StringFlag{
					Name:  "region",
					Usage: "[Optional] Droplet will be created at this region",
				},
				cli.StringFlag{
					Name:  "image",
					Usage: "[Optional] create droplet from this image",
				},
			},
			Action: func(c *cli.Context) {
				createDroplet(c)
			},
		},
		{
			Name:  "destroy",
			Usage: "destroy a droplet",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "id",
					Usage: "[Mandatory] droplet-id to delete",
				},
			},
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
			Name:  "sizes",
			Usage: "list a size information related to droplets",
			Action: func(c *cli.Context) {
				listSizes(c)
			},
		},
		/*		{
					Name:  "resize",
					Usage: "resize a droplet",
					Action: func(c *cli.Context) {
						resizeDroplet(c)
					},
				},
		*/
		{
			Name:  "halt",
			Usage: "shutdown droplet",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "id",
					Usage: "[Mandatory] droplet-id ",
				},
			},
			Action: func(c *cli.Context) {
				haltDroplet(c)
			},
		},
		{
			Name:  "restart",
			Usage: "restarts droplet",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "id",
					Usage: "[Mandatory] droplet-id ",
				},
			},
			Action: func(c *cli.Context) {
				restartDroplet(c)
			},
		},
		{
			Name:  "info",
			Usage: "details of droplet",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "id",
					Usage: "[Mandatory] droplet-id ",
				},
			},
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
	image := c.String("image")
	name := c.String("name")
	size := c.String("size")
	region := c.String("region")

	if name == "" {
		fmt.Println("Enter name of droplet")
		cli.ShowCommandHelp(c, "create")
	}

	d := getDomClient()

	if image == "" {
		image = d.config.Imageid
	}
	if size == "" {
		size = d.config.Size
	}
	if region == "" {
		region = d.config.Region
	}

	info, err := d.CreateDropletFromImage(name, region, size, image)
	if err != nil {
		fmt.Println("Error while creating droplet.")
	} else {
		fmt.Printf("Droplet will be ready within 60 Seconds Droplet ID: %d\n", info.ID)
	}

}

func destroyDroplet(c *cli.Context) {
	droplet := c.String("id")
	if droplet == "" {
		fmt.Printf("You need to specify an droplet-id\n\n")
		cli.ShowCommandHelp(c, "destroy")
		return
	}

	d := getDomClient()

	err := d.DestroyDroplet(droplet)
	if err != nil {
		fmt.Printf("Unable to delete the Droplet : %v\n", err)
	}
}

func listDroplets(c *cli.Context) {
	d := getDomClient()
	_, err := d.DropletList("")
	if err != nil {
		fmt.Println("Error while retriving the Droplet List")
	}
}

func listRegions(c *cli.Context) {
	d := getDomClient()
	err := d.ListRegions()
	if err != nil {
		fmt.Println("Unable to fetch User Region List")
	}
}

func listSizes(c *cli.Context) {
	d := getDomClient()
	err := d.ListSizes()
	if err != nil {
		fmt.Println("Unable to fetch User Size List")
	}
}

func resizeDroplet(c *cli.Context) {
	fmt.Println("NOT IMPLEMENTED !")
}
func haltDroplet(c *cli.Context) {
	droplet := c.String("id")
	if droplet == "" {
		fmt.Printf("You need to specify an droplet-id\n\n")
		cli.ShowCommandHelp(c, "halt")
		return
	}

	d := getDomClient()

	err := d.DropletShutdown(droplet)
	if err != nil {
		fmt.Printf("Unable to Shutdown the Droplet : %v\n", err)
	}
}
func restartDroplet(c *cli.Context) {
	//TODO: Force reboot.
	droplet := c.String("id")
	if droplet == "" {
		fmt.Printf("You need to specify an droplet-id\n\n")
		cli.ShowCommandHelp(c, "restart")
		return
	}

	d := getDomClient()

	err := d.DropletRestart(droplet)
	if err != nil {
		fmt.Printf("Unable to Reboot the Droplet : %v\n", err)
	}
}
func infoDroplet(c *cli.Context) {
	droplet := c.String("id")
	if droplet == "" {
		fmt.Printf("You need to specify an droplet-id\n\n")
		cli.ShowCommandHelp(c, "info")
		return
	}

	d := getDomClient()

	err := d.DropletInfo(droplet)
	if err != nil {
		fmt.Printf("Unable to fetch info \n")
	}
}
