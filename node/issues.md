# Programming Issues

	- lock in the dependencies

		- specify versions in Gopkgs.toml, do a 'dep ensure' and then check Gopks.lock 

	- golang.org/x/net

	 	- multiple registrations of /deb/requests

		- vendor copy of tendermint is in conflict with repo version (links are starting both inits, which collide)

			- remove repo copy (or sync with vendor copy?)
		
