'use strict'

const path = require('path');
const configFile = path.join(__dirname,'..','orgs','org1-config.yaml');

const util = require('util');

const FabricClient = require('fabric-client');
const serviceConfig = (process.env.ORG == 'Org1')? require('../orgs/org1.json') : require('../orgs/org2.json');;

const log4js = require('log4js');
const logger = log4js.getLogger('maejor-blockchain');
logger.level = 'ALL';

module.exports.getClient = async ()=>{

    FabricClient.setLogger(logger);

    const client = FabricClient.loadFromConfig(configFile);
    await client.initCredentialStores();
    let adminUserObj = await client.getUserContext('admin', true);
    if (adminUserObj == null){
        adminUserObj = await client.setUserContext({username: 'admin', password: 'adminpw'});
    }
    return client;
}


module.exports.invokeTransaction = async (fcn, args) =>	{

    let request = {
        targets: serviceConfig.blockchain.targets,
        chaincodeId: serviceConfig.blockchain.chaincodeName,
        fcn: fcn,
        args: args,
        chainId: serviceConfig.blockchain.channelName
    }

	logger.debug('request --->: ' + request);

	let errorMessage = null
    let txIDString = null;
    try{
        const client = await this.getClient();
        const txId = client.newTransactionID();
        request.txId = txId;
        txIDString = txId.getTransactionID();
        const channel = client.getChannel(serviceConfig.blockchain.channelName);

        let results = await channel.sendTransactionProposal(request);

        var proposalResponses = results[0];
		var proposal = results[1];

		var all_good = true;
		for (var i in proposalResponses) {
			let one_good = false;
			if (proposalResponses && proposalResponses[i].response &&
				proposalResponses[i].response.status === 200) {
				one_good = true;
				logger.info('invoke chaincode proposal was good');
			} else {
				logger.error('invoke chaincode proposal was bad');
			}
			all_good = all_good & one_good;
        }
        
        if (all_good) {

            logger.info(util.format(
				'Successfully sent Proposal and received ProposalResponse: Status - %s, message - "%s", metadata - "%s", endorsement signature: %s',
				proposalResponses[0].response.status, proposalResponses[0].response.message,
				proposalResponses[0].response.payload, proposalResponses[0].endorsement.signature));

			// wait for the channel-based event hub to tell us
			// that the commit was good or bad on each peer in our organization
			var promises = [];
			let event_hubs = channel.getChannelEventHubsForOrg();
			event_hubs.forEach((eh) => {
				logger.debug('invokeEventPromise - setting up event');
				let invokeEventPromise = new Promise((resolve, reject) => {
					let event_timeout = setTimeout(() => {
						let message = 'REQUEST_TIMEOUT:' + eh.getPeerAddr();
						logger.error(message);
						eh.disconnect();
					}, 3000);
					eh.registerTxEvent(txIDString, (tx, code, block_num) => {
						logger.info('The chaincode invoke chaincode transaction has been committed on peer %s',eh.getPeerAddr());
						logger.info('Transaction %s has status of %s in blocl %s', tx, code, block_num);
						clearTimeout(event_timeout);

						if (code !== 'VALID') {
							let message = util.format('The invoke chaincode transaction was invalid, code:%s',code);
							logger.error(message);
							reject(new Error(message));
						} else {
							let message = 'The invoke chaincode transaction was valid.';
							logger.info(message);
							resolve(message);
						}
					}, (err) => {
						clearTimeout(event_timeout);
						logger.error(err);
						reject(err);
					},
						// the default for 'unregister' is true for transaction listeners
						// so no real need to set here, however for 'disconnect'
						// the default is false as most event hubs are long running
						// in this use case we are using it only once
						{unregister: true, disconnect: true}
					);
					eh.connect();
				});
				promises.push(invokeEventPromise);
			});

			const ordererRequest = {
				txId: txId,
				proposalResponses: proposalResponses,
				proposal: proposal
			};

			logger.info(util.format('------->>> send transactions : %O', ordererRequest));
			const sendPromise = channel.sendTransaction(ordererRequest);
			// put the send to the orderer last so that the events get registered and
			// are ready for the orderering and committing
			promises.push(sendPromise);
			let results = await Promise.all(promises);
			logger.info(util.format('------->>> R E S P O N S E : %j', results));
			let response = results.pop(); //  orderer results are last in the results
			if (response.status === 'SUCCESS') {
				logger.info('Successfully sent transaction to the orderer.');
			} else {
				errorMessage = util.format('Failed to order the transaction. Error code: %s',response.status);
				logger.debug(errorMessage);
			}

			// now see what each of the event hubs reported
			for(let i in results) {
				let event_hub_result = results[i];
				let event_hub = event_hubs[i];
				logger.debug('Event results for event hub :%s',event_hub.getPeerAddr());
				if(typeof event_hub_result === 'string') {
					logger.debug(event_hub_result);
				} else {
					if(!errorMessage) errorMessage = event_hub_result.toString();
					logger.debug(event_hub_result.toString());
				}
			}
		} else {
			errorMessage = util.format('Failed to send Proposal and receive all good ProposalResponse');
			logger.debug(errorMessage);
		}

    }catch(error){
        logger.error('Failed to invoke due to error: ' + error.stack ? error.stack : error);
		errorMessage = error.toString();
    }

    if (!errorMessage) {
		let message = util.format(
			'Successfully invoked the chaincode %s to the channel \'%s\' for transaction ID: %s',
			'aladdin', 'mychannel', txIDString);
		logger.info(message);

		return txIDString;
	} else {
		let message = util.format('Failed to invoke chaincode. cause:%s',errorMessage);
		logger.error(message);
		throw new Error(message);
	}

}

module.exports.queryChaincode = async (fcn, args) => {

	// send query
	const request = {
		targets : serviceConfig.blockchain.targets, //queryByChaincode allows for multiple targets
		chaincodeId: serviceConfig.blockchain.chaincodeName,
		fcn: fcn,
		args: args
	};

	try {
		// first setup the client for this org
		const client = await this.getClient();
		logger.debug('Successfully got the fabric client for the organization "%s"', serviceConfig.blockchain.org);
		var channel = client.getChannel(serviceConfig.blockchain.channelName);
		if(!channel) {
			const message = util.format('Channel %s was not defined in the connection profile', serviceConfig.blockchain.channelName);
			logger.error(message);
			throw new Error(message);
		}


		let response_payloads = await channel.queryByChaincode(request);
		if (response_payloads) {
			for (let i = 0; i < response_payloads.length; i++) {
				logger.info(args[0]+' now has ' + response_payloads[i].toString('utf8') +
					' after the move');
			}
			return args[0]+' now has ' + response_payloads[0].toString('utf8') +
				' after the move';
		} else {
			logger.error('response_payloads is null');
			return 'response_payloads is null';
		}
	} catch(error) {
		logger.error('Failed to query due to error: ' + error.stack ? error.stack : error);
		return error.toString();
	}
}

module.exports.getBlockByNumber = async (blockNumber) => {

	try {
		// first setup the client for this org
		const client = await this.getClient();
		logger.debug('Successfully got the fabric client for the organization "%s"', serviceConfig.blockchain.org);
		var channel = client.getChannel(serviceConfig.blockchain.channelName);
		if(!channel) {
			let message = util.format('Channel %s was not defined in the connection profile', serviceConfig.blockchain.channelName);
			logger.error(message);
			throw new Error(message);
		}

		let response_payload = await channel.queryBlock(parseInt(blockNumber, serviceConfig.blockchain.peer));
		if (response_payload) {
			logger.debug(response_payload);
			return response_payload;
		} else {
			logger.error('response_payload is null');
			return 'response_payload is null';
		}
	} catch(error) {
		logger.error('Failed to query due to error: ' + error.stack ? error.stack : error);
		return error.toString();
	}

}

module.exports.getBlockByHash = async (hash) => {
	try {
		const client = await this.getClient();
		logger.debug('Successfully got the fabric client for the organization "%s"', serviceConfig.blockchain.org);
		var channel = client.getChannel(serviceConfig.blockchain.channelName);
		if(!channel) {
			let message = util.format('Channel %s was not defined in the connection profile', serviceConfig.blockchain.channelName);
			logger.error(message);
			throw new Error(message);
		}

		let response_payload = await channel.queryBlockByHash(Buffer.from(hash),  serviceConfig.blockchain.peer);
		if (response_payload) {
			logger.debug(response_payload);
			return response_payload;
		} else {
			logger.error('response_payload is null');
			return 'response_payload is null';
		}
	} catch(error) {
		logger.error('Failed to query due to error: ' + error.stack ? error.stack : error);
		return error.toString();
	}
};

module.exports.getChainInfo = async () => {
	try {
		// first setup the client for this org
		const client = await this.getClient();
		logger.debug('Successfully got the fabric client for the organization "%s"', serviceConfig.blockchain.org);
		var channel = client.getChannel(serviceConfig.blockchain.channelName);
		if(!channel) {
			let message = util.format('Channel %s was not defined in the connection profile', serviceConfig.blockchain.channelName);
			logger.error(message);
			throw new Error(message);
		}

		let response_payload = await channel.queryInfo(serviceConfig.blockchain.peer);
		if (response_payload) {
			logger.debug(response_payload);
			return response_payload;
		} else {
			logger.error('response_payload is null');
			return 'response_payload is null';
		}
	} catch(error) {
		logger.error('Failed to query due to error: ' + error.stack ? error.stack : error);
		return error.toString();
	}
};