workdir=$(pwd)
for serve in $(ls ./rpcServer)
do
  cd $workdir/rpcServer/$serve&&go run .
done