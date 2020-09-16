#  go-axis RPM package generate and installation/uninstallation

## centos 7

from the cloned resource github/axis-cash/go-axis
goto build directory
run 
```
./rpmbuild.sh
```

it will generate rpm package in github/axis-cash/go-axis/build/package/RPMS

## install with rpm
```
rpm -ivh ${rpmfile} --nodeps
```


## check install and env
exec following command in console
```
gaxis
```
it it complains with missing libboost_system ... etc.
please ref [installation guide](https://github.com/axis-cash/go-axis/wiki/Building-Axis)

for centos
```
sudo yum --setopt=group_package_types=mandatory,default,optional group install "Development Tools"
sudo yum install boost boost-devel boost-system boost-filesystem boost-thread
```
for ubuntu
```
sudo apt-get install -y build-essential golang
sudo apt-get install libboost-all-dev
```

## uninstall
if you met following error:

  install of gaxis-1.0-1.x86_64 conflicts with file from package gaxis-1.0-1.x86_64
or  you want to upgrade gaxis
you need to uninstall previous gaxis package

```
sudo dpkg --purge ${packagename}
```
or 
```
sudo rpm -e ${packagename}
```
