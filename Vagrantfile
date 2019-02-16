# -*- mode: ruby -*-
# vi: set ft=ruby :

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

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/bionic64"

  config.vm.provision "shell",
    inline: "apt install -y cryptsetup"

  config.vm.provision "shell", privileged: false,
    inline: $create_luks_dev
end
