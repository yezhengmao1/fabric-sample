'use strict';

const logger = require("@hyperledger/caliper-core").CaliperUtils.getLogger("Test");

const contractID  = "money_demo";
const contractVer = "1.0";

let bc, contx;
let initmoney = "100";
let counts = [];

module.exports.init = async (blockchain, context, args) => {
    bc = blockchain;
    contx = context;
};

module.exports.run = async() => {
    let count = "count_" + Math.random().toString(36).substr(7);
    counts.push(count);
    
    let txArgs = {
        chaincodeFunction: "open",
        chaincodeArguments: [count, initmoney]
    };

    return bc.invokeSmartContract(contx, contractID, contractVer, txArgs, 10000); 
};

module.exports.end = async() => { 
};

module.exports.initmoney = initmoney
module.exports.contractID = contractID
module.exports.contractVer = contractVer
module.exports.counts = counts;