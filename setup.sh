#!/bin/bash
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
  exit 0
fi

mkdir -p /tmp/owtfMonitor
cd /tmp/owtfMonitor

##
# Install golang to /usr/local and setup GOPATH to $HOME/go_workspace
##
if [[ $1 = "-g" ]]; then
  if ! hash go 2>/dev/null ; then
    echo "Installing golang 1.6.3"
    wget -c "https://storage.googleapis.com/golang/go1.6.3.linux-amd64.tar.gz"
    sudo tar -C /usr/local -xzf go1.6.3.linux-amd64.tar.gz
    sudo echo "export PATH=$PATH:/usr/local/go/bin" >> /etc/profile
    sudo echo "export GOPATH=~/go_workspace" >> /etc/profile
    sudo echo "export GOROOT=/usr/local/go" >> /etc/profile
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
  echo "Installing ssdeep"
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
echo "Fetching required dependencies..."
wget -c "https://raw.githubusercontent.com/owtf/health_monitor/master/dependencies"
cat dependencies | xargs go get -u
echo "Cloning OWTF Health-Monitor to $GOPATH/src/health_monitor"
if [ ! "$(ls -A $GOPATH/src/health_monitor)" ]; then
  git clone https://github.com/owtf/health_monitor.git $GOPATH/src/health_monitor
else
    echo "Health-Monitor is found, skipping cloning of the repository"
fi

cd $GOPATH/src/health_monitor
echo "Building health_monitor"
go build -i
echo "Setup complete"
