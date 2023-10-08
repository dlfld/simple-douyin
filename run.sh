
cd ./rpcServer/interaction && nohup ./interaction > interaction_logfile.txt &
cd ../message  && nohup ./message > message_logfile.txt & 
cd ../relation && nohup ./relation > relation_logfile.txt &
cd ../user && nohup ./user > user_logfile.txt &
cd ../video && nohup ./video > user_logfile.txt &  
cd ../../ && nohup ./douyin > douyin_logfile.txt &