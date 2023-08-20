cd ..
workdir=$(pwd)
for serve in $(ls rpcServer)
do
  cd $serve
  go run .
done