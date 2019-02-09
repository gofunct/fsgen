package modules

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Init22a926dc76789f65a26a7529c2dbc5e1635b336b = "package init\n\nvar (\n\thelp string\n)"

// Init returns go-assets FileSystem
var Init = assets.NewFileSystem(map[string][]string{"/": []string{"init.go"}}, map[string]*assets.File{
	"/": &assets.File{
		Path:     "/",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1549660769, 1549660769081675588),
		Data:     nil,
	}, "/init.go": &assets.File{
		Path:     "/init.go",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1549660769, 1549660769075695355),
		Data:     []byte(_Init22a926dc76789f65a26a7529c2dbc5e1635b336b),
	}}, "")
