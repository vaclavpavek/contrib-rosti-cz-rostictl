package rostiapi

// SSHAccess holds access via SSH protocol
type SSHAccess struct {
	Hostname string `json:"hostname"`
	Port     uint   `json:"port"`
	Username string `json:"username"`
}

// App is structure keeping backend data about an application
type App struct {
	ID           uint        `json:"id,omitmepty"`
	Date         string      `json:"date,omitmepty"`
	Name         string      `json:"name"`
	Enabled      bool        `json:"enabled,omitmepty"`
	Image        string      `json:"image,omitmepty"`
	Domains      []string    `json:"domains,omitmepty"`
	Mode         string      `json:"mode,omitmepty"`
	Plan         uint        `json:"plan,omitmepty"`
	SSHAccess    []SSHAccess `json:"ssh_access,omitmepty"`
	SSHKeys      []string    `json:"ssh_keys,omitempty,omitmepty"`
	SMTPUsername string      `json:"smtp_username,omitmepty"`
	SMTPToken    string      `json:"smtp_token,omitmepty"`
}

// ErrorResponse represents message returned by the API in case of non-200 response
type ErrorResponse struct {
	Message string                 `json:"message"`
	Errors  map[string]interface{} `json:"errors"`
	// TODO: add errors: {"message":"validation error","errors":{"domains":["Toto pole nesmí být prázdné (null)."]}}
}

// Action tells API what to do with the application
type Action struct {
	Action string `json:"action"` // Can be start, stop, restart, rebuild
}

// Plan defines RAW parameters of the service
type Plan struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	RAM      uint   `json:"ram"`
	Disk     uint   `json:"disk"`
	Price    uint   `json:"price"`
	CPUQuote uint   `json:"cpu_quota"`
}

// Company groups people around one project
type Company struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// Runtime is environment where the application is running
type Runtime struct {
	ID      uint   `json:"id"`
	Image   string `json:"image"`
	Default bool   `json:"default"` // default runtime
	Show    bool   `json:"show"`    // shown in admin
}

// AppStatus contains status information about one application
type AppStatus struct {
	Errors     []string `json:"errors"`
	Info       []string `json:"info"`
	DNSStatus  bool     `json:"dns_status"`
	HTTPStatus bool     `json:"http_status"`
	Running    bool     `json:"running"`
	Storage    struct {
		Usage     float64 `json:"usage"`
		Limit     float64 `json:"limit"`
		OverLimit float64 `json:"over_limit"`
	} `json:"storage"`
	Memory struct {
		Usage float64 `json:"usage"`
		Limit float64 `json:"limit"`
	} `json:"memory"`
}
