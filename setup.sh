#!/bin/sh
#
# ---------------------------------------------------------------------
# OWASP-OWTF Health-Monitor development environment setup script
# ---------------------------------------------------------------------
#

if  [[ $1 = "-h" ]]; then
  echo "OWASP-OWTF Health-Monitor development environment setup script"
  echo "Usage:"
  echo "setup.sh <Optional arguments>"
  echo ""
  echo "   -g     to install golang and its environment setup also"
  echo "   -h     help (this output)"
  echo "* Note: Run this script with root privalages for setup purposes"
  exit 0
fi

# Check whether the script is running as root
if [[ $EUID -ne 0 ]]; then
  echo "This script must be run as root" 1>&2
  exit 1
fi

mkdir -p /tmp/owtfMonitor
cd /tmp/owtfMonitor

##
# Install golang to /usr/local and setup GOPATH to $HOME/go_workspace
##
if [[ $1 = "-g" ]]; then
  if ! hash go 2>/dev/null ; then
    wget -c "https://storage.googleapis.com/golang/go1.6.3.linux-amd64.tar.gz"
    sudo tar -C /usr/local -xzf go1.6.3.linux-amd64.tar.gz
    echo "export PATH=$PATH:/usr/local/go/bin" >> /etc/profile
    echo "export GOPATH=~/go_workspace" >> /etc/profile
    source /etc/profile
    mkdir -p ~/go_workspace
  else
    echo "golang is already installed. Skipping the installation"
  fi
elif [ $# -gt 0 ]; then
  echo "[!] Option not found"
  echo "Try -h option for help"
  exit 1
fi

##
# Install ssdeep for fuzzy searching
##
export LD_LIBRARY_PATH=/usr/local/lib
if ! hash ssdeep 2>/dev/null ; then
  wget -c "http://downloads.sourceforge.net/project/ssdeep/ssdeep-2.13/ssdeep-2.13.tar.gz"
  tar zxvf ssdeep-2.13.tar.gz
  cd ssdeep-2.13
  ./configure
  make
  sudo make install
else
  echo "ssdeep is already installed. Skipping the installation"
fi

##
# Install required golang dependencies
##
wget -c "https://raw.githubusercontent.com/owtf/health_monitor/master/dependencies"
sudo -u ${USER} cat dependencies | xargs go get -u
if [ ! "$(ls -A ~/go_workspace/src/health_monitor)" ]; then
  sudo -u ${USER} git clone https://github.com/owtf/health_monitor.git ~/go_workspace/src/health_monitor
  cd ~/go_workspace/src/health_monitor
else
    echo "Health-Monitor is found, skipping cloning of the repository"
fi

sudo -u ${USER} go build -i
