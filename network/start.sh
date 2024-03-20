#!/bin/bash

if [[ `uname` == 'Darwin' ]]; then
    echo "Mac OS"
    export PATH=${PWD}/hyperledger-fabric-darwin-arm64-2.5.4/bin:$PATH
fi
if [[ `uname` == 'Linux' ]]; then
    if [[ `uname -m` == 'x86_64' ]]; then
        echo "Linux x86_64"
        export PATH=${PWD}/hyperledger-fabric-linux-amd64-2.5.4/bin:$PATH
    else
        echo "Linux arm64"
        export PATH=${PWD}/hyperledger-fabric-linux-arm64-2.5.4/bin:$PATH
    fi
fi


echo "一、清理环境"
./stop.sh

echo "二、生成证书和秘钥（ MSP 材料），生成结果将保存在 crypto-config 文件夹中"
cryptogen generate --config=./crypto-config.yaml

echo "三、创建排序通道创世区块"
configtxgen -profile ThreeOrgsOrdererGenesis -outputBlock ./config/genesis.block -channelID firstchannel

echo "四、生成通道配置事务'appchannel.tx'"
configtxgen -profile ThreeOrgsChannel -outputCreateChannelTx ./config/appchannel.tx -channelID appchannel

echo "五、为 CERTMAN 定义锚节点"
configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./config/CERTMANAnchor.tx -channelID appchannel -asOrg CERTMAN

echo "六、为 CET 定义锚节点"
configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./config/CETAnchor.tx -channelID appchannel -asOrg CET

echo "六、为 NCRE 定义锚节点"
configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./config/NCREAnchor.tx -channelID appchannel -asOrg NCRE

echo "区块链 : 启动"
docker-compose up -d
echo "正在等待节点的启动完成, 等待5秒"
sleep 5

CERTMANPeer0Cli="CORE_PEER_ADDRESS=peer0.certman.com:7051 CORE_PEER_LOCALMSPID=CERTMANMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/certman.com/users/Admin@certman.com/msp"
CERTMANPeer1Cli="CORE_PEER_ADDRESS=peer1.certman.com:7051 CORE_PEER_LOCALMSPID=CERTMANMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/certman.com/users/Admin@certman.com/msp"
CETPeer0Cli="CORE_PEER_ADDRESS=peer0.cet.com:7051 CORE_PEER_LOCALMSPID=CETMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/cet.com/users/Admin@cet.com/msp"
CETPeer1Cli="CORE_PEER_ADDRESS=peer1.cet.com:7051 CORE_PEER_LOCALMSPID=CETMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/cet.com/users/Admin@cet.com/msp"
NCREPeer0Cli="CORE_PEER_ADDRESS=peer0.ncre.com:7051 CORE_PEER_LOCALMSPID=NCREMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/ncre.com/users/Admin@ncre.com/msp"
NCREPeer1Cli="CORE_PEER_ADDRESS=peer1.ncre.com:7051 CORE_PEER_LOCALMSPID=NCREMSP CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/peer/ncre.com/users/Admin@ncre.com/msp"


echo "七、创建通道"
docker exec cli bash -c "$CERTMANPeer0Cli peer channel create -o orderer.certmanorder.com:7050 -c appchannel -f /etc/hyperledger/config/appchannel.tx"

echo "八、将所有节点加入通道"
docker exec cli bash -c "$CERTMANPeer0Cli peer channel join -b appchannel.block"
docker exec cli bash -c "$CERTMANPeer1Cli peer channel join -b appchannel.block"
docker exec cli bash -c "$CETPeer0Cli peer channel join -b appchannel.block"
docker exec cli bash -c "$CETPeer1Cli peer channel join -b appchannel.block"
docker exec cli bash -c "$NCREPeer0Cli peer channel join -b appchannel.block"
docker exec cli bash -c "$NCREPeer1Cli peer channel join -b appchannel.block"

echo "九、更新锚节点"
docker exec cli bash -c "$CERTMANPeer0Cli peer channel update -o orderer.certmanorder.com:7050 -c appchannel -f /etc/hyperledger/config/CERTMANAnchor.tx"
docker exec cli bash -c "$CETPeer0Cli peer channel update -o orderer.certmanorder.com:7050 -c appchannel -f /etc/hyperledger/config/CETAnchor.tx"
docker exec cli bash -c "$NCREPeer0Cli peer channel update -o orderer.certmanorder.com:7050 -c appchannel -f /etc/hyperledger/config/NCREAnchor.tx"



echo "十、部署链吗"
docker exec cli peer lifecycle chaincode package fabric-realty.tar.gz --path /opt/gopath/src/chaincode --lang golang --label fabric-realty_1.0.0

echo "  安装链码到 CERTMAN Peer0"
docker exec cli bash -c "$CERTMANPeer0Cli peer lifecycle chaincode install fabric-realty.tar.gz"

echo "  安装链码到 CET Peer0"
docker exec cli bash -c "$CETPeer0Cli peer lifecycle chaincode install fabric-realty.tar.gz"

echo "  安装链码到 NCRE Peer0"
docker exec cli bash -c "$NCREPeer0Cli peer lifecycle chaincode install fabric-realty.tar.gz"

PACKAGE_ID=$(docker exec cli bash -c "$CERTMANPeer0Cli peer lifecycle chaincode queryinstalled" | grep "fabric-realty_1.0.0:" | sed 's/Package ID: //; s/, Label.*//')

echo "  组织批准链码"
docker exec cli bash -c "$CERTMANPeer0Cli peer lifecycle chaincode approveformyorg --channelID appchannel --name fabric-realty --version 1.0.0 --package-id $PACKAGE_ID --sequence 1 --waitForEvent"
docker exec cli bash -c "$CETPeer0Cli peer lifecycle chaincode approveformyorg --channelID appchannel --name fabric-realty --version 1.0.0 --package-id $PACKAGE_ID --sequence 1 --waitForEvent"
docker exec cli bash -c "$NCREPeer0Cli peer lifecycle chaincode approveformyorg --channelID appchannel --name fabric-realty --version 1.0.0 --package-id $PACKAGE_ID --sequence 1 --waitForEvent"

echo "  查询准备状态"
docker exec cli bash -c "$CERTMANPeer0Cli peer lifecycle chaincode checkcommitreadiness --channelID appchannel --name fabric-realty --version 1.0.0 --sequence 1 --output json"

echo "  提交链码定义"
docker exec cli bash -c "$CERTMANPeer0Cli peer lifecycle chaincode commit -o orderer.certmanorder.com:7050 --channelID appchannel --name fabric-realty --version 1.0.0 --sequence 1 --waitForEvent --peerAddresses peer0.cet.com:7051 --peerAddresses peer0.certman.com:7051 --peerAddresses peer0.ncre.com:7051"

echo "十一、初始化账本"
docker exec cli bash -c "$CERTMANPeer0Cli peer chaincode invoke -C appchannel -n fabric-realty -c '{\"Args\":[\"Init\"]}' --peerAddresses peer0.cet.com:7051 --peerAddresses peer0.certman.com:7051 --peerAddresses peer0.ncre.com:7051"

echo "十二、验证链码"
docker exec cli bash -c "$CERTMANPeer0Cli peer chaincode invoke -C appchannel -n fabric-realty -c '{\"Args\":[\"hello\"]}'"
docker exec cli bash -c "$CETPeer0Cli peer chaincode invoke -C appchannel -n fabric-realty -c '{\"Args\":[\"hello\"]}'"
docker exec cli bash -c "$NCREPeer0Cli peer chaincode invoke -C appchannel -n fabric-realty -c '{\"Args\":[\"hello\"]}'"
