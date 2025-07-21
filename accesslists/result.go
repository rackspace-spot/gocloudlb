package accesslists

import (
	"github.com/gophercloud/gophercloud"
)

type GetResult struct {
	gophercloud.Result
}

// Extract interprets a GetResult as a Node.
func (r GetResult) Extract() ([]NetworkItem, error) {
	var s struct {
		AccessList []NetworkItem `json:"accessList"`
	}
	err := r.ExtractInto(&s)
	return s.AccessList, err
}

type CreateResult struct {
	gophercloud.ErrResult
}

// DeleteResult method to determine if the call succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}
