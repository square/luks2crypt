Luks2Crypt
==========

- https://golang.org/cmd/cgo/

- https://gitlab.com/cryptsetup/cryptsetup/blob/v1_6_6/lib/libcryptsetup.h

Luks2crypt is used to manage luks client devices and allow escrowing to a
[crypt-server](https://github.com/grahamgilbert/Crypt-Server). Currently, it
impliments some functionality similar to [Crypt2](https://github.com/grahamgilbert/crypt2).

`postimaging`:

- gathers system info (serial number, username, hostname)

- generates a random password

- test if the password passed in on the cli unlocks the disk

- caches the new password to `/etc/luks2crypt/crypt_recovery_key.json`

- uploads the new password to your local crypt-server

- changes the luks password passed in on the cli to the newly generated one

Dependencies
------------

Luks2crypt requires a pre-existing crypt-server to escrow keys. Crypt-server is
a Django web service for centrally storing recovery keys for full disk
encryption. See: https://github.com/grahamgilbert/Crypt-Server for more details.

Usage
-----

Setting the admin password and escrowing it post imaging:

    sudo luks2crypt postimaging \
      --luksdevice "<device_to_manage>" \
      --currentpassword "<password_to_replace>" \
      --cryptserver "<cryptserver.example.com>"

Development
-----------

- This repository uses go modules (https://github.com/golang/go/wiki/Modules).
You should be able to simply `go get` the repo and the dependencies will
auto install. You will need to be using go version 1.11 or higher.

- The cryptsetup libs are required to build. Cryptsetup C libraries are used
through cgo to manage the encrypted devices. On debian/ubuntu you can run:

      sudo apt install libcryptsetup-dev

- To prepare for a release by cleaning up the unused dependencies run:

      make deps

- Use the `Makefile` to test and build luks2crypt:

      make

- If you would like to use a mock crypt server to test client changes on is
  included in this project:

      make mockserver

- If you need a test enviornment, the provided `Vagrantfile` creates an ubuntu
  vm. The vagrantfile has a provision script that creates a luks disk image at
  `/home/vagrant/luks-dev-disk.img`. The image is then encrypted with the password
  "devpassword" and mounted at `/mnt`.

      make devup       # create the dev vm
      make devssh      # connect to the consule of the vm
      make devclean    # delete the vm

  This also includes a mock implimentation of crypt-server to log the form
  data to stdout. You can launch the dev environment as follows:

      make devup
      make devssh
      sudo cryptservermock  # start the mock crypt-server
      
      # in a new term window test the client
      make devssh
      sudo /vagrant/bin/luks2crypt postimaging \
        -l ./luks-dev-disk.img \
        -p devpassword \
        -s ubuntu-bionic:8443

  You should then see the form post data printed to stdout from
  `cryptservermock`.

License
-------

      Copyright 2018 Square Inc.

      This program is free software: you can redistribute it and/or modify
      it under the terms of the GNU General Public License as published by
      the Free Software Foundation, either version 3 of the License, or
      (at your option) any later version.

      This program is distributed in the hope that it will be useful,
      but WITHOUT ANY WARRANTY; without even the implied warranty of
      MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
      GNU General Public License for more details.

      You should have received a copy of the GNU General Public License
      along with this program.  If not, see <http://www.gnu.org/licenses/>.
