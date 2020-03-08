'use strict';

const G = require("./open.js");

const logger = require("@hyperledger/caliper-core").CaliperUtils.getLogger("Test");

let bc, contx;
let index = 0;

module.exports.init = async (blockchain, context, args) => {
    bc = blockchain;
    contx = context;
    index = 0;
};

module.exports.run = async() => {
    let count = G.counts[index]; 
    index++;
    
    let txArgs = {
        chaincodeFunction: "delete",
        chaincodeArguments: [count]
    };

    return bc.invokeSmartContract(contx, G.contractID, G.contractVer, txArgs, 10000); 
};

module.exports.end = async() => {
};