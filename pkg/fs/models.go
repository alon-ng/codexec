package fs

type File struct {
	Name    string `json:"name" binding:"required" validate:"required"`
	Ext     string `json:"ext" binding:"required" validate:"required"`
	Content string `json:"content" binding:"required" validate:"required"`
}

type Directory struct {
	Name        string      `json:"name" binding:"required" validate:"required"`
	Directories []Directory `json:"directories" binding:"required" validate:"required"`
	Files       []File      `json:"files" binding:"required" validate:"required"`
}
