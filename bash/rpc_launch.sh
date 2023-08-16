for serve in $(ls rpcServer)
do
echo $(go run rpcServer/$serve)
done