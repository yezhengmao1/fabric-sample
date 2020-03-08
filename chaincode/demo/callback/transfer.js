'use strict';

const G = require("./open.js");

const logger = require("@hyperledger/caliper-core").CaliperUtils.getLogger("Test");

let bc, contx;
let total;
let index = 0;

module.exports.init = async (blockchain, context, args) => {
    bc = blockchain;
    contx = context;
    total = G.counts.length;
    index = 0;
};

module.exports.run = async() => {
    let money = 1;
    let srccount = G.counts[index];
    let dstcount = G.counts[total - index - 1];
    
    if(index < total / 2) {
        money = 1;
    }else {
        money = 20;
    }
    
    index++;
    
    let txArgs = {
        chaincodeFunction: "invoke",
        chaincodeArguments: [srccount, dstcount, money.toString()]
    };

    return bc.invokeSmartContract(contx, G.contractID, G.contractVer, txArgs, 10000); 
};

module.exports.end = async() => {
};