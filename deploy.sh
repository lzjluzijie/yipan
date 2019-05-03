#!/usr/bin/env bash
export TZ='Asia/Shanghai'
mkdir ~/.ssh/
echo $sshkey > ~/.ssh/key
export GIT_SSH_COMMAND="ssh -i ~/.ssh/key"
sed -i 's/一/\n/g' ~/.ssh/key
chmod 400 ~/.ssh/key
eval $(ssh-agent -s)
ssh-add ~/.ssh/key
echo "github.com,192.30.255.112 ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTtw7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvSAHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIFImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrTJvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8VJiS5ap43JXiUFFAaQ==" > ~/.ssh/known_hosts

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
