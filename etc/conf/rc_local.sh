#!/bin/sh
# vi /etc/rc.local
# source /home/shy/contexts/usr/local/work/20230511-nginx-story/etc/conf/rc_local.sh

mkdir /tmp/nginx; cd $CTX_ROOT/usr/local/daemon/10000 && ./sbin/nginx -p $PWD
cd $CTX_ROOT/../docker && ./dockerd --host unix://$PWD/docker.sock --pidfile $PWD/docker.pid --exec-root=$PWD/exec --data-root=$PWD/data --registry-mirror "https://ccr.ccs.tencentyun.com" &
# cd $CTX_ROOT/usr/install/docker && ./dockerd --host unix://$PWD/docker.sock --pidfile $PWD/docker.pid --exec-root=$PWD/exec --data-root=$PWD/data --registry-mirror "https://ccr.ccs.tencentyun.com" &
