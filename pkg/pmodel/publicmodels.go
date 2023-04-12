package pmodel

//ConfigDescription describe the config of a tenant
type ConfigDescription struct {
	StoreID  string `json:"storeid"`
	TenantID string `json:"tenantID"`
	Size     int    `json:"size"`
}
