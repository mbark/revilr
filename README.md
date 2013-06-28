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

Finally run

	go get github.com/mattn/go-sqlite3

Additional downloads needed
---------------------------

[Mustache](https://github.com/hoisie/mustache)

	go get github.com/hoisie/mustache