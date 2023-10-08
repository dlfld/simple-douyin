
env GOOS=linux GOARCH=amd64 go build ./
cd ./rpcServer/interaction  && env GOOS=linux GOARCH=amd64 go build ./
cd ../message && env GOOS=linux GOARCH=amd64 go build ./
cd ../relation && env GOOS=linux GOARCH=amd64 go build ./
cd ../user && env GOOS=linux GOARCH=amd64 go build ./
cd ../video && env GOOS=linux GOARCH=amd64 go build  ./

