# install
go get -u github.com/cheggaaa/pb/v3

go install github.com/mrco24/web-archive@latest
cp -r /root/go/bin//web-archive /usr/local/bin

# use
web-archive -l sub.txt -o url.txt

