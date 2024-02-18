#!/bin/bash

if [[ `uname` == 'Darwin' ]]; then
    echo "Mac OS"
    export PATH=${PWD}/hyperledger-fabric-darwin-amd64-1.4.12/bin:$PATH
fi
if [[ `uname` == 'Linux' ]]; then
    echo "Linux"
    export PATH=${PWD}/hyperledger-fabric-linux-arm64-2.5.4/bin:$PATH
fi

echo "一、清理环境"
./stop.sh

echo "二、生成证书和秘钥（ MSP 材料），生成结果将保存在 crypto-config 文件夹中"
cryptogen generate --config=./crypto-config.yaml

echo "三、创建排序通道创世区块"
configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./config/genesis.block -channelID firstchannel

echo "四、生成通道配置事务'appchannel.tx'"
configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./config/appchannel.tx -channelID appchannel

echo "五、为 BENZ 定义锚节点"
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./config/BENZAnchor.tx -channelID appchannel -asOrg BENZ

echo "六、为 TESLA 定义锚节点"
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./config/TESLAAnchor.tx -channelID appchannel -asOrg TESLA

echo "区块链 : 启动"
docker-compose up -d
echo "正在等待节点的启动完成, 等待5秒"
sleep 5

BENZPeer0Cli="CORE_PEER_ADDRESS=peer0.benz.com:7051 CORE_PEER_LOCALMSPID=BENZMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/benz.com/users/Admin@benz.com/msp"
BENZPeer1Cli="CORE_PEER_ADDRESS=peer1.benz.com:7051 CORE_PEER_LOCALMSPID=BENZMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/benz.com/users/Admin@benz.com/msp"
TESLAPeer0Cli="CORE_PEER_ADDRESS=peer0.tesla.com:7051 CORE_PEER_LOCALMSPID=TESLAMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/tesla.com/users/Admin@tesla.com/msp"
TESLAPeer1Cli="CORE_PEER_ADDRESS=peer1.tesla.com:7051 CORE_PEER_LOCALMSPID=TESLAMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/tesla.com/users/Admin@tesla.com/msp"

echo "七、创建通道"
docker exec cli bash -c "$BENZPeer0Cli peer channel create -o orderer.carunion.com:7050 -c appchannel -f /etc/hyperledger/config/appchannel.tx"

echo "八、将所有节点加入通道"
docker exec cli bash -c "$BENZPeer0Cli peer channel join -b appchannel.block"
docker exec cli bash -c "$BENZPeer1Cli peer channel join -b appchannel.block"
docker exec cli bash -c "$TESLAPeer0Cli peer channel join -b appchannel.block"
docker exec cli bash -c "$TESLAPeer1Cli peer channel join -b appchannel.block"

echo "九、更新锚节点"
docker exec cli bash -c "$BENZPeer0Cli peer channel update -o orderer.carunion.com:7050 -c appchannel -f /etc/hyperledger/config/BENZAnchor.tx"
docker exec cli bash -c "$TESLAPeer0Cli peer channel update -o orderer.carunion.com:7050 -c appchannel -f /etc/hyperledger/config/TESLAAnchor.tx"



echo "十、部署链吗"
docker exec cli peer lifecycle chaincode package fabric-realty.tar.gz --path /opt/gopath/src/chaincode --lang golang --label fabric-realty_1.0.0

echo "  安装链码到 BENZ Peer0"
docker exec cli bash -c "$BENZPeer0Cli peer lifecycle chaincode install fabric-realty.tar.gz"

echo "  安装链码到 TESLA Peer0"
docker exec cli bash -c "$TESLAPeer0Cli peer lifecycle chaincode install fabric-realty.tar.gz"

PACKAGE_ID=$(docker exec cli bash -c "$BENZPeer0Cli peer lifecycle chaincode queryinstalled" | grep "fabric-realty_1.0.0:" | sed 's/Package ID: //; s/, Label.*//')

echo "  组织批准链码"
docker exec cli bash -c "$BENZPeer0Cli peer lifecycle chaincode approveformyorg --channelID appchannel --name fabric-realty --version 1.0.0 --package-id $PACKAGE_ID --sequence 1 --waitForEvent"
docker exec cli bash -c "$TESLAPeer0Cli peer lifecycle chaincode approveformyorg --channelID appchannel --name fabric-realty --version 1.0.0 --package-id $PACKAGE_ID --sequence 1 --waitForEvent"

echo "  查询准备状态"
docker exec cli bash -c "$BENZPeer0Cli peer lifecycle chaincode checkcommitreadiness --channelID appchannel --name fabric-realty --version 1.0.0 --sequence 1 --output json"

echo "  提交链码定义"
docker exec cli bash -c "$BENZPeer0Cli peer lifecycle chaincode commit -o orderer.carunion.com:7050 --channelID appchannel --name fabric-realty --version 1.0.0 --sequence 1 --waitForEvent --peerAddresses peer0.tesla.com:7051 --peerAddresses peer0.benz.com:7051"

# echo "十一、初始化账本"
# docker exec cli bash -c "$BENZPeer0Cli peer chaincode invoke -C appchannel -n fabric-realty -c '{\"Args\":[\"Init\"]}' --peerAddresses peer0.benz.com:7051 --peerAddresses peer0.tesla.com:7051"

# echo "十二、验证链码"
# docker exec cli bash -c "$BENZPeer0Cli peer chaincode invoke -C appchannel -n fabric-realty -c '{\"Args\":[\"hello\"]}'"
# docker exec cli bash -c "$TESLAPeer0Cli peer chaincode invoke -C appchannel -n fabric-realty -c '{\"Args\":[\"hello\"]}'"
