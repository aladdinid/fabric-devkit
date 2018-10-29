'use strict';
const express = require('express');
const bodyParser = require('body-parser');
const http = require('http');
const app = express();
const cors = require('cors');

const blockchain = require('./blockchain');

const log4js = require('log4js');
const logger = log4js.getLogger('maejor-api');
const serviceConfig = (process.env.ORG === 'Org1')? require('./orgs/org1.json') : require('./orgs/org2.json');
logger.level = serviceConfig.loggingLevel;

const host = 8080;
const port = 8080;

//
//------------------ SET CONFIGURATONS -------------------
//

app.options('*', cors());
app.use(cors());
//support parsing of application/json type post data
app.use(bodyParser.json());
//support parsing of application/x-www-form-urlencoded post data
app.use(bodyParser.urlencoded({
	extended: false
}));

app.use((req, res, next) => {
	logger.debug(' ------>>>>>> new request for %s', req.originalUrl);
  return next();
});

//------------------------ START SERVER -----------------------
//
var server = http.createServer(app).listen(port, function() {});
logger.info('***************  http://%s:%s  ******************',host,port);
//server.timeout = 240000;

//
// ----- REST ENDPOINTS START HERE ----
//
// Invoke transaction on chaincode on target peers
app.post('/invoke', async (req, res) => {
	logger.debug('==================== INVOKE ON CHAINCODE ==================');

	var fcn = req.body.fcn;
	var args = req.body.args;

	logger.debug('fcn  : ' + fcn);
	logger.debug('args  : ' + args);

  const txIDString = await blockchain.invokeTransaction(fcn,args);
  logger.debug('Trx created : ' + txIDString);
  res.json({success: true, trxId: txIDString});
});

// Invoke transaction on chaincode on target peers
app.post('/query', async (req, res) => {
	logger.debug('==================== QUERY ON CHAINCODE ==================');

	var fcn = req.body.fcn;
	var args = req.body.args;

	logger.debug('fcn  : ' + fcn);
	logger.debug('args  : ' + args);

  const data = await blockchain.queryChaincode(fcn,args);
  logger.debug('data found : ' + data);
  res.json({success: true, data: data});
});

// Invoke transaction on chaincode on target peers
app.get('/blocks/:blockId', async (req, res) => {
	logger.debug('==================== QUERY ON CHAINCODE ==================');

  const blockId = req.params.blockId;

	logger.debug('blockId  : ' + blockId);

  const blockData = await blockchain.getBlockByNumber(blockId);
  logger.debug('Block found : ' + blockData);
  res.json({success: true, data: blockData});
});

// Invoke transaction on chaincode on target peers
app.get('/transactions/:trxHash', async (req, res) => {
	logger.debug('==================== QUERY ON CHAINCODE ==================');

  const trxHash = req.params.trxHash;

	logger.debug('trxHash  : ' + trxHash);

  const blockData = await blockchain.getBlockByHash(trxHash);
  logger.debug('Block found : ' + blockData);
  res.json({success: true, data: blockData});
});