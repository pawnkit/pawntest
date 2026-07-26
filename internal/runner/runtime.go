package runner

type RuntimeMetadata struct {
	SchemaVersion int              `json:"schemaVersion"`
	RuntimeTier   string           `json:"runtimeTier"`
	Engine        RuntimeComponent `json:"engine"`
	Target        string           `json:"target"`
	Capabilities  []string         `json:"capabilities"`
	Limitations   []string         `json:"limitations"`
	Plugin        *RuntimePlugin   `json:"plugin,omitempty"`
	Server        *RuntimeServer   `json:"server,omitempty"`
}

type RuntimeComponent struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type RuntimePlugin struct {
	Name           string `json:"name"`
	Architecture   string `json:"architecture"`
	WorkerProtocol int    `json:"workerProtocol"`
}

type RuntimeServer struct {
	Artifact        string `json:"artifact"`
	Version         string `json:"version"`
	Profile         string `json:"profile"`
	AdapterProtocol int    `json:"adapterProtocol"`
	SessionComplete bool   `json:"sessionCompleted"`
}
