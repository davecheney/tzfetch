tzfetch
=======

`tzfetch` streams and unpacks `.tar.gz` files from http servers. `tzfetch` assumes all the remote URLs are `.tar.gz` files, it is not a general purpose tool, merely an example. 

Installation
------------

    go get github.com/davecheney/tzfetch

Usage
-----

    tzfetch [-v] [-C dir] url1 url2 ...

Licence
-------

This code is derived from the [Juju project](https://launchpad.net/juju-core) and is released under the GNU Affero General Public License.
