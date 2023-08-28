cd ..
workdir=$(pwd)
for serve in $(ls rpcServer)
do
#  osascript <<END
  open -a Terminal.app cd $serve & go run .
#  go run .
  osascript -e 'tell app "Terminal"
      do script "cd '$serve';go run ."
  end tell'
done