//1. 서버 세팅
'use strict';

var express = require('express');
const { Gateway, Wallets } = require('fabric-network');
const FabricCAServices = require('fabric-ca-client');

const { buildCACLient, registerAndErollUser, enrollAdmin } =require('./CAUtils.js');

var app = express();

var path = require('path');
var fs = require('fs');

//2. fabric 연결설정
app.use('/public', express.static(path.join(__dirname,'public')));
//3. 미들웨어 설정
app.use(express.urlencoded({ extended : false}));
app.use(express.json());

const mspOrg1 = 'Org1MSP'
const walletPath= path.join(__dirname, 'wallet');

const caClient = buildCCPOrg1();
const caClient = buildCACLient(FabricCAServices, ccp, 'ca.org1.example.com');

//4. /asset POST 라우팅
app.post('/asset', async(req, res) =>{
    var key = req.body.key;
    var value = req.body.value;

    console.log("/asset post start --", key, value);
    const gateway = new gateway();

    try {
        const wallet = await buildWallet(Wallets, walletPath);

        await gateway.connect(ccp, {
            wallet,
            identity: "appUser",
            discovery: { enabled: true, asLocalHost: true} //using asLocalhost as this gateway is using a fabric network deployed locally
        });
        const network = await gateway,getNetwork("mychannel");
        const contract = network.getContract("simpleasset");
        await contract.submitTransaction('Set', key, name);

    }catch (error) {
        var result = `{"result":"fail", "message":"tx has NOT submitted"}`;
        var obj = JSON.parse(result);
        console.log("/asset end -- failed", error);
        res.status(200).send(obj);
        return;
    }finally {
        gateway.disconnect();
    }

    var result = `{"result":"success", "message":"tx has submitted"}`;
    var obj = JSON.parse(result);
    console.log("/asset end -- success");
    res.status(200).send(obj);
});

//5. /asset GET 라우팅
app.get('/asset', async(req, res) =>{
    var key = req.query.key;
    console.log("/asset get start --", key);

    const gateway = new Gateway();

    try {
        const wallet = await buildWallet(Wallets, walletPath);
        //GW->connect ->CH -> CC-> submitTrnasaction
        await gateway.connect(ccp, {
            wallet,
            identity: "appUser",
            discovery: { enabled: true, asLocalHost: true }
        });
        const network = await gateway.getNetwork("mychannel");
        const contract = network.getContract("simpleasset");
        var result = await contract.evaluateTransaction('Get', key);
        //result가 byte array라고 생각하고 
        var result = `{"result": "success", "message": ${result}}`;
        console.log("/asset get end -- success", result);
        var obj = JSON.parse(result);
        res.status(200).send(obj);
        } catch (error) {
            var result = `{"result":"fail", "message": "Get has a error"}`;
            var obj = JSON.parse(result);
            console.log("/asset get end -- failed", error);
            res.status(200).send(obj);
            return;
        } finally {
            gateway.disconnect();
        }
});

//6. 서버 listen (서버 시작)
app.listen(3000, () => {
    console.log('Express server is started: 3000')
});