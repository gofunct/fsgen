{{ .app_name }}
{{ .project_short_description }}

Getting started
This project requires Go to be installed. On OS X with Homebrew you can just run brew install go.

Running it then should be as simple as:

$ make
$ ./bin/{{ .app_name }}
Testing
make test