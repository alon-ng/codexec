package fs

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
