package service

// Result contains agent execution result
type Result struct {
	Owner     string `json:"owner"`
	Repo     string `json:"repo"`
	Issue    int    `json:"issue,omitempty"`
	PR       int    `json:"pr,omitempty"`
	Build    int    `json:"build,omitempty"`
	Status   string `json:"status"`
	Fix      string `json:"fix,omitempty"`
	Review   string `json:"review,omitempty"`
	Tests    string `json:"tests,omitempty"`
	Message  string `json:"message"`
}

// Message returns human-readable message
func (r *Result) Message() string {
	switch r.Status {
	case "ready":
		return fmt.Sprintf("PR #%d created and ready for review", r.PR)
	case "approved":
		return "Code reviewed and approved ✅"
	case "failed":
		return "Tests failed - please check"
	case "deployed":
		return fmt.Sprintf("Deployed to %s", r.Message)
	default:
		return r.Message
	}
}