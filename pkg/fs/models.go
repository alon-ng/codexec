package fs

type Entry struct {
	Name     string  `json:"name" binding:"required" validate:"required"`
	Content  string  `json:"content,omitempty"`
	Children []Entry `json:"children,omitempty"`
}
