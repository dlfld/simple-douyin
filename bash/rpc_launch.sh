workdir=$(pwd)
for serve in $(ls rpcServer)
do
`gnome-terminal -t "$serve" -e "bash -c 'cd $workdir/rpcServer/$serve; go run .'" `
done