cd ./network
./stop.sh


cd ../application
./stop.sh
./build.sh
./start.sh


cd ../network
./start.sh