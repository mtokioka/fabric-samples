git clone https://github.com/mtokioka/fabric-samples.git
cd fabric-samples
curl -sSL https://goo.gl/5ftp2f | bash
export PATH=<path to download location>/bin:$PATH

cd get_and_set
./startFabric.sh

npm install
node enrollAdmin.js
node registerUser.js

node get.js
node set.js
node get.js

# chaincode 編集
./updateFabric.sh
node get_all.js
