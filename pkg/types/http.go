package types

type PingResponse struct {
	MachineName string `json:"machineName"`
}

type ReloadResponse struct {
	ConfigReloaded bool `json:"configReloaded"`
}
