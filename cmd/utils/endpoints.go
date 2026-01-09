package utils

// API v2 endpoint constants for static paths.
const (
	// EndpointProjects is the endpoint for listing all projects.
	// GET /api/v2/projects
	EndpointProjects = "/api/v2/projects"

	// EndpointProcessStatuses is the endpoint for listing all process statuses.
	// GET /api/v2/process/statuses
	EndpointProcessStatuses = "/api/v2/process/statuses"
)

// EndpointProjectConfig returns the endpoint for fetching a specific project configuration.
// GET /api/v2/projects/{projectName}
func EndpointProjectConfig(projectName string) string {
	return "/api/v2/projects/" + projectName
}

// EndpointProcessStatus returns the endpoint for fetching a specific process status.
// GET /api/v2/process/status/{processID}
func EndpointProcessStatus(processID string) string {
	return "/api/v2/process/status/" + processID
}

// EndpointTransferStart returns the endpoint for starting a transfer for a project.
// POST /api/v2/process/{projectName}/start
func EndpointTransferStart(projectName string) string {
	return "/api/v2/process/" + projectName + "/start"
}
