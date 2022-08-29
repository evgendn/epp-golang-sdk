package eppoclient

import "net/http"

var __version__ = "1.0.0"

func InitClient(config Config) *EppoClient {
	config.validate()
	sdkParams := SDKParams{apiKey: config.ApiKey, sdkName: "go", sdkVersion: __version__}

	httpClient := newHttpClient(config.BaseUrl, &http.Client{Timeout: REQUEST_TIMEOUT_SECONDS}, sdkParams)
	configStore := newConfigurationStore(MAX_CACHE_ENTRIES)
	requestor := newExperimentConfigurationRequestor(*httpClient, *configStore)
	assignmentLogger := NewAssignmentLogger()

	client := newEppoClient(requestor, assignmentLogger)

	client.poller.Start()

	return client
}
