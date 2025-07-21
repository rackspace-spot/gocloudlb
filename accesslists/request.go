package accesslists

import (
	"log"
	"strconv"

	"github.com/gophercloud/gophercloud"
)

type AccessType string

const (
	Allow AccessType = "ALLOW"
	Deny  AccessType = "DENY"
)

// Get returns data about access list for a specific load balancer by its ID.
func Get(client *gophercloud.ServiceClient, lbID uint64) (r GetResult) {
	url := client.ServiceURL("loadbalancers", strconv.FormatUint(lbID, 10), "accesslist")
	log.Printf("GET %s", url)
	_, r.Err = client.Get(url, &r.Body, nil)
	return
}

func Create(client *gophercloud.ServiceClient, lbID uint64, opts []CreateOpts) (r CreateResult) {
	url := client.ServiceURL("loadbalancers", strconv.FormatUint(lbID, 10), "accesslist")
	log.Printf("POST %s", url)

	body := struct {
		AccessList []CreateOpts `json:"accessList"`
	}{
		opts,
	}
	_, r.Err = client.Post(url, body, nil, nil)
	return
}

// Delete deletes all access list entries for a specific load balancer by its ID.
func DeleteAll(client *gophercloud.ServiceClient, lbID uint64) (r DeleteResult) {
	url := client.ServiceURL("loadbalancers", strconv.FormatUint(lbID, 10), "accesslist")
	log.Printf("DELETE %s", url)

	_, r.Err = client.Delete(url, nil)
	return
}

// Delete deletes the access list entry by its ID for a specific load balancer.
func BulkDelete(client *gophercloud.ServiceClient, lbID uint64, idAsUrlParameters string) (r DeleteResult) {
	url := client.ServiceURL("loadbalancers", strconv.FormatUint(lbID, 10), "accesslist", idAsUrlParameters)
	log.Printf("DELETE %s", url)

	_, r.Err = client.Delete(url, nil)
	return
}

type CreateOpts struct {
	Address string `json:"address"`
	// The type of the access list, e.g., "ALLOW" or "DENY"
	Type AccessType `json:"type"`
}

type NetworkItem struct {
	Address string `json:"address"`
	// The type of the access list, e.g., "ALLOW" or "DENY"
	Type AccessType `json:"type"`
	ID   uint64     `json:"id"`
}
