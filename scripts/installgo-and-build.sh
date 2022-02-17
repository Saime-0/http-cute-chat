export folder="./temp/go"
# install go
mkdir -p $folder \
&& wget https://go.dev/dl/go1.17.3.linux-amd64.tar.gz -O $folder/go.tar.gz \
&& tar -xzf $folder/go.tar.gz -C /usr/local \
&& echo "export PATH=$PATH:/usr/local/go/bin" >> /etc/profile \
&& source /etc/profile \
&& go version \
&& GOOS=linux go build -v ./server.go
