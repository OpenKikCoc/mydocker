cp /work/test/sources.list /etc/opt/

apt update

apt-get install -y aufs-tools

dpkg -l | grep aufs-tools