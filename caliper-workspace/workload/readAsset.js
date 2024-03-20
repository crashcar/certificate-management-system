'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class ReadAssetWorkload extends WorkloadModuleBase {
    constructor() {
        super();
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments) {
        this.roundArguments = roundArguments;
    }

    async submitTransaction() {
        const assetID = `asset${this.workerIndex}_${this.roundArguments.args[0]}`;
        const certID = `cert${this.workerIndex}`;
        let args = {
            chaincodeFunction: this.roundArguments.func,
            chaincodeArguments: [assetID, certID]
        };

        if (this.roundArguments.func === 'uploadCertOrg') {
            args.chaincodeArguments = [
                `hashFile_${assetID}`,
                `hashPath_${certID}`,
                certID,
                `holderID_${certID}`,
                `holderName_${certID}`,
                'certType',
                'reviewer',
                '2020-01-01',
                '2025-01-01',
                'issuingAuthority',
                'phone',
                'email',
                'address'
            ];
        }

        await this.sutAdapter.sendRequests(args);
    }
}

function createWorkloadModule() {
    return new ReadAssetWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
