Setup Go Build in EC2 instance

sudo wget https://dl.google.com/go/go1.13.4.linux-amd64.tar.gz
tar -xzf go1.13.4.linux-amd64.tar.gz


#environment variables to set
export GOROOT=/home/ec2-user/go
export PATH=$PATH:$GOROOT/bin 
export GOBIN=$GOROOT/bin
export GOPATH=~/UserAPI-firebase/
export PATH=$GOPATH/bin:$PATH

sudo yum install git