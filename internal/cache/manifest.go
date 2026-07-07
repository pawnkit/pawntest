package cache

type Manifest struct {
	Inputs []string `json:"inputs"`
	Output string   `json:"output"`
}
