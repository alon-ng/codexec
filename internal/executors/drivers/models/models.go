package models

type ExecuteResponse struct {
	Stdout   string  `json:"stdout"`
	Stderr   string  `json:"stderr"`
	ExitCode int     `json:"exit_code"`
	Time     float64 `json:"time"`
	Memory   int64   `json:"memory"`
	CPU      float64 `json:"cpu"`
}

type File struct {
	Name    string `json:"name"`
	Ext     string `json:"ext"`
	Content string `json:"content"`
}

type Directory struct {
	Name        string      `json:"name"`
	Directories []Directory `json:"directories"`
	Files       []File      `json:"files"`
}
