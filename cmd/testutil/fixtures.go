package testutil

// Project fixtures
const (
	SampleProjects       = `["Project1", "Project2", "Project3"]`
	SingleProject        = `["OnlyOne"]`
	EmptyProjects        = `[]`
	SampleProjectConfig  = `{
		"cohortSelector": {
			"trustCenterAgent": {
				"server": {"baseUrl": "https://tc.example.com"},
				"domain": "test-domain",
				"patientIdentifierSystem": "http://example.com/pid",
				"policySystem": "http://example.com/policy",
				"policies": ["policy1", "policy2"]
			}
		},
		"dataSelector": {
			"everything": {
				"fhirServer": {"baseUrl": "https://fhir.example.com"},
				"resolve": {"patientIdentifierSystem": "http://example.com/pid"}
			}
		},
		"deidentificator": {
			"deidentifhir": {
				"trustCenterAgent": {
					"server": {"baseUrl": "https://tc.example.com"},
					"domains": {
						"pseudonym": "pseudonym-domain",
						"salt": "salt-domain",
						"dateShift": "dateshift-domain"
					}
				},
				"maxDateShift": "P14D",
				"deidentifhirConfig": "config.json",
				"scraperConfig": "scraper.json"
			}
		},
		"bundleSender": {
			"researchDomainAgent": {
				"server": {"baseUrl": "https://research.example.com"},
				"project": "TestProject"
			}
		}
	}`
)

// Process fixtures
const (
	SampleProcessRunning = `{
		"processId": "proc-123",
		"phase": "RUNNING",
		"createdAt": [2024, 1, 15, 10, 30, 0, 0],
		"finishedAt": null,
		"totalPatients": 100,
		"totalBundles": 500,
		"deidentifiedBundles": 300,
		"sentBundles": 200,
		"skippedBundles": 10
	}`

	SampleProcessCompleted = `{
		"processId": "proc-456",
		"phase": "COMPLETED",
		"createdAt": [2024, 1, 15, 10, 30, 0, 0],
		"finishedAt": [2024, 1, 15, 12, 0, 0, 0],
		"totalPatients": 50,
		"totalBundles": 250,
		"deidentifiedBundles": 250,
		"sentBundles": 250,
		"skippedBundles": 0
	}`

	SampleProcessList = `[
		{
			"processId": "proc-123",
			"phase": "RUNNING",
			"createdAt": [2024, 1, 15, 10, 30, 0, 0],
			"finishedAt": null,
			"totalPatients": 100,
			"totalBundles": 500,
			"deidentifiedBundles": 300,
			"sentBundles": 200,
			"skippedBundles": 10
		},
		{
			"processId": "proc-456",
			"phase": "COMPLETED",
			"createdAt": [2024, 1, 15, 10, 30, 0, 0],
			"finishedAt": [2024, 1, 15, 12, 0, 0, 0],
			"totalPatients": 50,
			"totalBundles": 250,
			"deidentifiedBundles": 250,
			"sentBundles": 250,
			"skippedBundles": 0
		}
	]`

	EmptyProcessList = `[]`
)

// Error response fixtures
const (
	Error404     = `{"error": "not found"}`
	Error500     = `{"error": "internal server error"}`
	Error400     = `{"error": "bad request"}`
	MalformedJSON = `{invalid json}`
	EmptyBody    = ``
)
