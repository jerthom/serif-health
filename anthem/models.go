package anthem

type ReportingStructure struct {
	Plans        []Plan         `json:"reporting_plans"`
	NetworkFiles []FileLocation `json:"in_network_files"`
}

type Plan struct {
	Name       string `json:"plan_name"`
	IdType     string `json:"plan_id_type"`
	Id         string `json:"plan_id"`
	MarketType string `json:"plan_market_type"`
}

type FileLocation struct {
	Description string `json:"description"`
	Location    string `json:"location"`
}
