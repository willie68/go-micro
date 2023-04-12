// Package for server internal models
package model

/*
ConfigDescription describres all metadata of a config
*/
type ConfigDescription struct {
	StoreID  string `json:"storeid"`
	TenantID string `json:"tenantID"`
	Size     int    `json:"size"`
}
