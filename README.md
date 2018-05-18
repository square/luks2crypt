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

- Install golang dep (https://github.com/golang/dep/blob/master/README.md)

- The cryptsetup libs are required to build. Cryptsetup C libraries are used
through cgo to manage the encrypted devices. On debian/ubuntu you can run:

      sudo apt install libcryptsetup-dev

- Install deps required for project with `dep`:

      dep ensure


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
