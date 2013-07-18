revilr
======

Minimalist link sharing for Chrome. Today.

Setting up sqlite3 on Windows
-----------------------------

* [Add git to your Paths variable](http://blog.countableset.ch/2012/06/07/adding-git-to-windows-7-path/)
* [Download TDM-GCC and install it](http://tdm-gcc.tdragon.net/)
* Add the bin folder of your TDM-GCC installation
* [Download SQLite3](http://mislav.uniqpath.com/rails/install-sqlite3/)
* [Download HG, to be able to get all go package](http://tortoisehg.bitbucket.org/download/index.html)

Your Path should now contain the following, if not add as appropriate:

	C:\Go\go\bin;C:\Program Files (x86)\Git\bin;C:\Program Files (x86)\Git\cmd;C:\MinGW64\bin

Here TDM-GCC is installed at MinGW64.

Additional setup
---------------------------

In order to be able to build the project you need to set the GOPATH environment variable to the project home. The project home is at

	revilr/server

To build the project run

	go build revilr

Packages used
---------------------------

[sqlite3](https://github.com/mattn/go-sqlite3)

	go get github.com/mattn/go-sqlite3

[Mustache](https://github.com/hoisie/mustache)

	go get github.com/hoisie/mustache

[Gorilla](http://www.gorillatoolkit.org/pkg/sessions)

	go get github.com/gorilla/sessions

[bcrypt](http://godoc.org/code.google.com/p/go.crypto/bcrypt)

	go get code.google.com/p/go.crypto/bcrypt