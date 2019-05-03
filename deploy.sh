#!/usr/bin/env bash
export TZ='Asia/Shanghai'
echo $sshkey > ~/key
export GIT_SSH_COMMAND="ssh -i ~/key"
sed -i 's/一/\n/g' ~/key
chmod 400 ~/key
eval $(ssh-agent -s)
ssh-add ~/key

git config --global user.name yipan-config
git clone git@github.com:lzjluzijie/yipan-config.git

mv yipan-config/config config
go run yipan.go
mv config yipan-config/config

cd yipan-config
git add .
git commit -m "Config updated by netlify: `date +"%Y%m%d-%H:%M:%S"` UTC+8"
git push origin master
cd ..

mkdir public
mv _redirects public
