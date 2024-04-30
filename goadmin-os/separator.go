package goadmin_os

import "os"

const (
	PathSeparator     string = string(os.PathSeparator)     // OS-specific path separator '\\'
	PathListSeparator string = string(os.PathListSeparator) // OS-specific path list separator ';'
)
