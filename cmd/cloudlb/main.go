/*
Copyright 2021 Rackspace, Inc.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License.  You may obtain a copy of the
License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied.  See the License for the
specific language governing permissions and limitations under the License.
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/pagination"

	"github.com/os-pc/gocloudlb"
	"github.com/os-pc/gocloudlb/loadbalancers"
	"github.com/os-pc/gocloudlb/nodes"
	"github.com/os-pc/gocloudlb/virtualips"
)

func main() {
	// allow the user to supply a region
	var region string
	flag.StringVar(&region, "region", os.Getenv("OS_REGION_NAME"), "Cloud Region")

	// create load balancer
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createCmdType := createCmd.String("type", "PUBLIC", "PUBLIC or PRIVATE")

	// list load balancers
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	// show load balancer
	showCmd := flag.NewFlagSet("show", flag.ExitOnError)

	// update load balancer
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)

	// delete load balancer
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)

	flag.Parse()

	if flag.NArg() < 0 {
		log.Printf("hurr", flag.Args())
		log.Printf("thing", createCmd.Parsed())
		log.Printf("Usage: %s create|list|show|update|delete ...", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	switch flag.Arg(0) {
	case "create":
		createCmd.Parse(flag.Args()[1:])
	case "list":
		listCmd.Parse(flag.Args()[1:])
	case "show":
		showCmd.Parse(flag.Args()[1:])
	case "update":
		updateCmd.Parse(flag.Args()[1:])
	case "delete":
		deleteCmd.Parse(flag.Args()[1:])
	default:
		log.Print("bleh")
		log.Printf("Usage: %s create|list|show|update|delete ...", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	opts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		log.Fatal(err)
	}

	service, err := gocloudlb.NewLB(provider, gophercloud.EndpointOpts{
		Region: region,
	})
	if err != nil {
		log.Fatal(err)
	}

	if createCmd.Parsed() {
		args := createCmd.Args()
		if len(args) < 3 {
			fmt.Fprintf(createCmd.Output(), "Usage: %s create NAME PROTOCOL PORT\n", os.Args[0])
			createCmd.PrintDefaults()
			os.Exit(2)
		}

		port, err := strconv.ParseUint(args[2], 10, 16)
		if err != nil {
			log.Fatal("invalid port number: %s", err)
		}

		opts := loadbalancers.CreateOpts{
			Name:     args[0],
			Protocol: args[1],
			Port:     int32(port),
			VirtualIps: []virtualips.CreateOpts{
				virtualips.CreateOpts{
					Type: *createCmdType,
				},
			},
			Nodes: []nodes.CreateOpts{},
		}

		lb, err := loadbalancers.Create(service, opts).Extract()
		if err != nil {
			log.Fatal(err)
		}
		json, err := json.MarshalIndent(lb, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(json))
	} else if listCmd.Parsed() {
		opts := loadbalancers.ListOpts{}

		pager := loadbalancers.List(service, opts)

		listErr := pager.EachPage(func(page pagination.Page) (bool, error) {
			loadbalancerList, err := loadbalancers.ExtractLoadBalancers(page)

			if err != nil {
				return false, err
			}

			for _, loadbalancer := range loadbalancerList {
				json, err := json.MarshalIndent(loadbalancer, "", "  ")
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println(string(json))
			}
			return true, err
		})

		if listErr != nil {
			log.Fatal(listErr)
		}
	} else if showCmd.Parsed() {
		args := showCmd.Args()
		if len(args) < 1 {
			fmt.Fprintf(showCmd.Output(), "Usage: %s show ID\n", os.Args[0])
			showCmd.PrintDefaults()
			os.Exit(2)
		}

		lbID, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			log.Fatal("invalid load balancer id: %s", err)
		}

		lb, err := loadbalancers.Get(service, lbID).Extract()
		if err != nil {
			log.Fatal(err)
		}

		json, err := json.MarshalIndent(lb, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(json))
	} else if updateCmd.Parsed() {
		log.Print("update")
	} else if deleteCmd.Parsed() {
		args := deleteCmd.Args()
		if len(args) < 1 {
			fmt.Fprintf(deleteCmd.Output(), "Usage: %s delete ID\n", os.Args[0])
			deleteCmd.PrintDefaults()
			os.Exit(2)
		}

		lbID, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			log.Fatal("invalid load balancer id: %s", err)
		}

		deleteErr := loadbalancers.Delete(service, lbID).ExtractErr()
		if deleteErr != nil {
			log.Fatal(deleteErr)
		}
		log.Printf("Successfully deleted: %d\n", lbID)
	} else {
		flag.PrintDefaults()
		os.Exit(1)
	}
}
