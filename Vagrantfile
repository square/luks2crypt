# -*- mode: ruby -*-
# vi: set ft=ruby :

$golang_install = <<-SCRIPT
set -ux
GOLANGVER=1.12
GOLANGTAR=https://dl.google.com/go/go${GOLANGVER}.linux-amd64.tar.gz

pushd /tmp
rm -rf /usr/local/go /tmp/go*.tar.gz
wget --no-verbose "${GOLANGTAR}"
tar -C /usr/local -xzf go${GOLANGVER}.linux-amd64.tar.gz
echo 'export PATH=${PATH}:/usr/local/go/bin:${HOME}/go/bin' > /etc/profile.d/go-path.sh
popd
SCRIPT

$create_luks_dev = <<-SCRIPT
set -ux
pushd /home/vagrant
sudo umount /dev/mapper/luks-dev-disk
sudo cryptsetup close luks-dev-disk
rm -f luks-dev-disk.img
fallocate -l 1G luks-dev-disk.img
parted luks-dev-disk.img mklabel msdos --script
echo "devpassword" | cryptsetup --batch-mode luksFormat luks-dev-disk.img
echo "devpassword" | sudo cryptsetup open luks-dev-disk.img luks-dev-disk
sudo mkfs.ext4 /dev/mapper/luks-dev-disk
sudo mount /dev/mapper/luks-dev-disk /mnt

echo "A dev luks disk has been created at /home/vagrant/luks-dev-disk.img"
echo "The device has a password of \"devpassword\" and is mounted at /mnt"
popd
SCRIPT

$cryptservermock_install = <<-SCRIPT
pushd /vagrant
go build -v tools/cryptservermock/cryptservermock.go
cp cryptservermock /usr/local/bin/
popd
SCRIPT

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/bionic64"
  config.vm.provider "virtualbox" do |v|
    v.customize ["setextradata", :id, "VBoxInternal/Devices/pcbios/0/Config/DmiSystemSerial", "string:1234foobar"]
  end

  config.vm.provision "shell",
    inline: "apt install -y cryptsetup ssl-cert"

  config.vm.provision "shell", privileged: true,
    inline: $golang_install

  config.vm.provision "shell", privileged: false,
    inline: $create_luks_dev

  config.vm.provision "shell", privileged: true,
    inline: $cryptservermock_install
end
