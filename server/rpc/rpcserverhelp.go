// Copyright (c) 2015-2017 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package rpc

import (
	"errors"
	"sort"
	"strings"
	"sync"

	"github.com/daglabs/btcd/btcjson"
)

// helpDescsEnUS defines the English descriptions used for the help strings.
var helpDescsEnUS = map[string]string{
	// DebugLevelCmd help.
	"debugLevel--synopsis": "Dynamically changes the debug logging level.\n" +
		"The levelspec can either a debug level or of the form:\n" +
		"<subsystem>=<level>,<subsystem2>=<level2>,...\n" +
		"The valid debug levels are trace, debug, info, warn, error, and critical.\n" +
		"The valid subsystems are AMGR, ADXR, BCDB, BMGR, BTCD, CHAN, DISC, PEER, RPCS, SCRP, SRVR, and TXMP.\n" +
		"Finally the keyword 'show' will return a list of the available subsystems.",
	"debugLevel-levelspec":   "The debug level(s) to use or the keyword 'show'",
	"debugLevel--condition0": "levelspec!=show",
	"debugLevel--condition1": "levelspec=show",
	"debugLevel--result0":    "The string 'Done.'",
	"debugLevel--result1":    "The list of subsystems",

	// AddManualNodeCmd help.
	"addManualNode--synopsis": "Attempts to add or remove a persistent peer.",
	"addManualNode-addr":      "IP address and port of the peer to operate on",
	"addManualNode-onetry":    "When enabled, will try a single connection to a peer",

	// NodeCmd help.
	"node--synopsis":     "Attempts to add or remove a peer.",
	"node-subcmd":        "'disconnect' to remove all matching non-persistent peers, 'remove' to remove a persistent peer, or 'connect' to connect to a peer",
	"node-target":        "Either the IP address and port of the peer to operate on, or a valid peer ID.",
	"node-connectsubcmd": "'perm' to make the connected peer a permanent one, 'temp' to try a single connect to a peer",

	// TransactionInput help.
	"transactioninput-txid": "The hash of the input transaction",
	"transactioninput-vout": "The specific output of the input transaction to redeem",

	// CreateRawTransactionCmd help.
	"createRawTransaction--synopsis": "Returns a new transaction spending the provided inputs and sending to the provided addresses.\n" +
		"The transaction inputs are not signed in the created transaction.\n" +
		"The signrawtransaction RPC command provided by wallet must be used to sign the resulting transaction.",
	"createRawTransaction-inputs":         "The inputs to the transaction",
	"createRawTransaction-amounts":        "JSON object with the destination addresses as keys and amounts as values",
	"createRawTransaction-amounts--key":   "address",
	"createRawTransaction-amounts--value": "n.nnn",
	"createRawTransaction-amounts--desc":  "The destination address as the key and the amount in BTC as the value",
	"createRawTransaction-locktime":       "Locktime value; a non-zero value will also locktime-activate the inputs",
	"createRawTransaction--result0":       "Hex-encoded bytes of the serialized transaction",

	// ScriptSig help.
	"scriptsig-asm": "Disassembly of the script",
	"scriptsig-hex": "Hex-encoded bytes of the script",

	// PrevOut help.
	"prevout-addresses": "previous output addresses",
	"prevout-value":     "previous output value",

	// VinPrevOut help.
	"vinprevout-coinbase":  "The hex-encoded bytes of the signature script (coinbase txns only)",
	"vinprevout-txid":      "The hash of the origin transaction (non-coinbase txns only)",
	"vinprevout-vout":      "The index of the output being redeemed from the origin transaction (non-coinbase txns only)",
	"vinprevout-scriptSig": "The signature script used to redeem the origin transaction as a JSON object (non-coinbase txns only)",
	"vinprevout-prevOut":   "Data from the origin transaction output with index vout.",
	"vinprevout-sequence":  "The script sequence number",

	// Vin help.
	"vin-coinbase":  "The hex-encoded bytes of the signature script (coinbase txns only)",
	"vin-txid":      "The hash of the origin transaction (non-coinbase txns only)",
	"vin-vout":      "The index of the output being redeemed from the origin transaction (non-coinbase txns only)",
	"vin-scriptSig": "The signature script used to redeem the origin transaction as a JSON object (non-coinbase txns only)",
	"vin-sequence":  "The script sequence number",

	// ScriptPubKeyResult help.
	"scriptpubkeyresult-asm":       "Disassembly of the script",
	"scriptpubkeyresult-hex":       "Hex-encoded bytes of the script",
	"scriptpubkeyresult-type":      "The type of the script (e.g. 'pubkeyhash')",
	"scriptpubkeyresult-reqSigs":   "The number of required signatures",
	"scriptpubkeyresult-addresses": "The bitcoin addresses associated with this script",

	// Vout help.
	"vout-value":        "The amount in BTC",
	"vout-n":            "The index of this transaction output",
	"vout-scriptPubKey": "The public key script used to pay coins as a JSON object",

	// TxRawDecodeResult help.
	"txrawdecoderesult-txid":     "The hash of the transaction",
	"txrawdecoderesult-version":  "The transaction version",
	"txrawdecoderesult-locktime": "The transaction lock time",
	"txrawdecoderesult-vin":      "The transaction inputs as JSON objects",
	"txrawdecoderesult-vout":     "The transaction outputs as JSON objects",

	// DecodeRawTransactionCmd help.
	"decodeRawTransaction--synopsis": "Returns a JSON object representing the provided serialized, hex-encoded transaction.",
	"decodeRawTransaction-hextx":     "Serialized, hex-encoded transaction",

	// DecodeScriptResult help.
	"decodescriptresult-asm":       "Disassembly of the script",
	"decodescriptresult-type":      "The type of the script (e.g. 'pubkeyhash')",
	"decodescriptresult-reqSigs":   "The number of required signatures",
	"decodescriptresult-addresses": "The bitcoin addresses associated with this script",
	"decodescriptresult-p2sh":      "The script hash for use in pay-to-script-hash transactions (only present if the provided redeem script is not already a pay-to-script-hash script)",

	// DecodeScriptCmd help.
	"decodeScript--synopsis": "Returns a JSON object with information about the provided hex-encoded script.",
	"decodeScript-hexscript": "Hex-encoded script",

	// EstimateFeeCmd help.
	"estimateFee--synopsis": "Estimate the fee per kilobyte in satoshis " +
		"required for a transaction to be mined before a certain number of " +
		"blocks have been generated.",
	"estimateFee-numblocks": "The maximum number of blocks which can be " +
		"generated before the transaction is mined.",
	"estimateFee--result0": "Estimated fee per kilobyte in satoshis for a block to " +
		"be mined in the next NumBlocks blocks.",

	// GenerateCmd help
	"generate--synopsis": "Generates a set number of blocks (simnet or regtest only) and returns a JSON\n" +
		" array of their hashes.",
	"generate-numblocks": "Number of blocks to generate",
	"generate--result0":  "The hashes, in order, of blocks generated by the call",

	// GetAllManualNodesInfoCmd help.
	"getAllManualNodesInfo--synopsis":   "Returns information about manually added (persistent) peers.",
	"getAllManualNodesInfo-details":     "Specifies whether the returned data is a JSON object including DNS and connection information, or just a list of added peers",
	"getAllManualNodesInfo--condition0": "details=false",
	"getAllManualNodesInfo--condition1": "details=true",
	"getAllManualNodesInfo--result0":    "List of added peers",

	// GetManualNodeInfoResultAddr help.
	"getmanualnodeinforesultaddr-address":   "The ip address for this DNS entry",
	"getmanualnodeinforesultaddr-connected": "The connection 'direction' (inbound/outbound/false)",

	// GetManualNodeInfoResult help.
	"getmanualnodeinforesult-manualnode": "The ip address or domain of the manually added peer",
	"getmanualnodeinforesult-connected":  "Whether or not the peer is currently connected",
	"getmanualnodeinforesult-addresses":  "DNS lookup and connection information about the peer",

	// GetManualNodeInfoCmd help.
	"getManualNodeInfo--synopsis":   "Returns information about manually added (persistent) peers.",
	"getManualNodeInfo-details":     "Specifies whether the returned data is a JSON object including DNS and connection information, or just a list of added peers",
	"getManualNodeInfo-node":        "Only return information about this specific peer instead of all added peers",
	"getManualNodeInfo--condition0": "details=false",
	"getManualNodeInfo--condition1": "details=true",
	"getManualNodeInfo--result0":    "List of added peers",

	// GetBestBlockResult help.
	"getbestblockresult-hash":   "Hex-encoded bytes of the best block hash",
	"getbestblockresult-height": "Height of the best block",

	// GetBestBlockCmd help.
	"getBestBlock--synopsis": "Get block height and hash of best block in the main chain.",
	"getBestBlock--result0":  "Get block height and hash of best block in the main chain.",

	// GetBestBlockHashCmd help.
	"getBestBlockHash--synopsis": "Returns the hash of the of the best (most recent) block in the longest block chain.",
	"getBestBlockHash--result0":  "The hex-encoded block hash",

	// GetBlockCmd help.
	"getBlock--synopsis":   "Returns information about a block given its hash.",
	"getBlock-hash":        "The hash of the block",
	"getBlock-verbose":     "Specifies the block is returned as a JSON object instead of hex-encoded string",
	"getBlock-verbosetx":   "Specifies that each transaction is returned as a JSON object and only applies if the verbose flag is true (btcd extension)",
	"getBlock--condition0": "verbose=false",
	"getBlock--condition1": "verbose=true",
	"getBlock--result0":    "Hex-encoded bytes of the serialized block",

	// GetBlockChainInfoCmd help.
	"getBlockDagInfo--synopsis": "Returns information about the current blockDAG state and the status of any active soft-fork deployments.",

	// GetBlockDAGInfoResult help.
	"getblockdaginforesult-dag":                   "The name of the DAG the daemon is on (testnet, mainnet, etc)",
	"getblockdaginforesult-blocks":                "The number of blocks in the best known chain",
	"getblockdaginforesult-headers":               "The number of headers that we've gathered for in the best known chain",
	"getblockdaginforesult-tiphashes":             "The block hashes for the tips in the DAG",
	"getblockdaginforesult-difficulty":            "The current chain difficulty",
	"getblockdaginforesult-mediantime":            "The median time from the PoV of the best block in the chain",
	"getblockdaginforesult-verificationprogress":  "An estimate for how much of the best chain we've verified",
	"getblockdaginforesult-pruned":                "A bool that indicates if the node is pruned or not",
	"getblockdaginforesult-pruneheight":           "The lowest block retained in the current pruned chain",
	"getblockdaginforesult-dagwork":               "The total cumulative work in the DAG",
	"getblockdaginforesult-softforks":             "The status of the super-majority soft-forks",
	"getblockdaginforesult-bip9_softforks":        "JSON object describing active BIP0009 deployments",
	"getblockdaginforesult-bip9_softforks--key":   "bip9_softforks",
	"getblockdaginforesult-bip9_softforks--value": "An object describing a particular BIP009 deployment",
	"getblockdaginforesult-bip9_softforks--desc":  "The status of any defined BIP0009 soft-fork deployments",

	// SoftForkDescription help.
	"softforkdescription-reject":  "The current activation status of the softfork",
	"softforkdescription-version": "The block version that signals enforcement of this softfork",
	"softforkdescription-id":      "The string identifier for the soft fork",
	"-status":                     "A bool which indicates if the soft fork is active",

	// TxRawResult help.
	"txrawresult-hex":           "Hex-encoded transaction",
	"txrawresult-txid":          "The hash of the transaction",
	"txrawresult-version":       "The transaction version",
	"txrawresult-locktime":      "The transaction lock time",
	"txrawresult-vin":           "The transaction inputs as JSON objects",
	"txrawresult-vout":          "The transaction outputs as JSON objects",
	"txrawresult-blockhash":     "Hash of the block the transaction is part of",
	"txrawresult-confirmations": "Number of confirmations of the block",
	"txrawresult-time":          "Transaction time in seconds since 1 Jan 1970 GMT",
	"txrawresult-blocktime":     "Block time in seconds since the 1 Jan 1970 GMT",
	"txrawresult-size":          "The size of the transaction in bytes",
	"txrawresult-hash":          "The wtxid of the transaction",

	// SearchRawTransactionsResult help.
	"searchrawtransactionsresult-hex":           "Hex-encoded transaction",
	"searchrawtransactionsresult-txid":          "The hash of the transaction",
	"searchrawtransactionsresult-hash":          "The wxtid of the transaction",
	"searchrawtransactionsresult-version":       "The transaction version",
	"searchrawtransactionsresult-locktime":      "The transaction lock time",
	"searchrawtransactionsresult-vin":           "The transaction inputs as JSON objects",
	"searchrawtransactionsresult-vout":          "The transaction outputs as JSON objects",
	"searchrawtransactionsresult-blockhash":     "Hash of the block the transaction is part of",
	"searchrawtransactionsresult-confirmations": "Number of confirmations of the block",
	"searchrawtransactionsresult-time":          "Transaction time in seconds since 1 Jan 1970 GMT",
	"searchrawtransactionsresult-blocktime":     "Block time in seconds since the 1 Jan 1970 GMT",
	"searchrawtransactionsresult-size":          "The size of the transaction in bytes",

	// GetBlockVerboseResult help.
	"getblockverboseresult-hash":              "The hash of the block (same as provided)",
	"getblockverboseresult-confirmations":     "The number of confirmations",
	"getblockverboseresult-size":              "The size of the block",
	"getblockverboseresult-height":            "The height of the block in the block chain",
	"getblockverboseresult-version":           "The block version",
	"getblockverboseresult-versionHex":        "The block version in hexadecimal",
	"getblockverboseresult-merkleroot":        "Root hash of the merkle tree",
	"getblockverboseresult-tx":                "The transaction hashes (only when verbosetx=false)",
	"getblockverboseresult-rawtx":             "The transactions as JSON objects (only when verbosetx=true)",
	"getblockverboseresult-time":              "The block time in seconds since 1 Jan 1970 GMT",
	"getblockverboseresult-nonce":             "The block nonce",
	"getblockverboseresult-bits":              "The bits which represent the block difficulty",
	"getblockverboseresult-difficulty":        "The proof-of-work difficulty as a multiple of the minimum difficulty",
	"getblockverboseresult-parentblockhashes": "The hashes of the parent blocks",
	"getblockverboseresult-nextblockhashes":   "The hashes of the next blocks (only if there are any)",

	// GetBlockCountCmd help.
	"getBlockCount--synopsis": "Returns the number of blocks in the longest block chain.",
	"getBlockCount--result0":  "The current block count",

	// GetBlockHashCmd help.
	"getBlockHash--synopsis": "Returns hash of the block in best block chain at the given height.",
	"getBlockHash-index":     "The block height",
	"getBlockHash--result0":  "The block hash",

	// GetBlockHeaderCmd help.
	"getBlockHeader--synopsis":   "Returns information about a block header given its hash.",
	"getBlockHeader-hash":        "The hash of the block",
	"getBlockHeader-verbose":     "Specifies the block header is returned as a JSON object instead of hex-encoded string",
	"getBlockHeader--condition0": "verbose=false",
	"getBlockHeader--condition1": "verbose=true",
	"getBlockHeader--result0":    "The block header hash",

	// GetBlockHeaderVerboseResult help.
	"getblockheaderverboseresult-hash":              "The hash of the block (same as provided)",
	"getblockheaderverboseresult-confirmations":     "The number of confirmations",
	"getblockheaderverboseresult-height":            "The height of the block in the block chain",
	"getblockheaderverboseresult-version":           "The block version",
	"getblockheaderverboseresult-versionHex":        "The block version in hexadecimal",
	"getblockheaderverboseresult-merkleroot":        "Root hash of the merkle tree",
	"getblockheaderverboseresult-time":              "The block time in seconds since 1 Jan 1970 GMT",
	"getblockheaderverboseresult-nonce":             "The block nonce",
	"getblockheaderverboseresult-bits":              "The bits which represent the block difficulty",
	"getblockheaderverboseresult-difficulty":        "The proof-of-work difficulty as a multiple of the minimum difficulty",
	"getblockheaderverboseresult-parentblockhashes": "The hashes of the parent blocks",
	"getblockheaderverboseresult-nextblockhashes":   "The hashes of the next blocks (only if there are any)",

	// TemplateRequest help.
	"templaterequest-mode":         "This is 'template', 'proposal', or omitted",
	"templaterequest-capabilities": "List of capabilities",
	"templaterequest-longpollid":   "The long poll ID of a job to monitor for expiration; required and valid only for long poll requests ",
	"templaterequest-sigoplimit":   "Number of signature operations allowed in blocks (this parameter is ignored)",
	"templaterequest-sizelimit":    "Number of bytes allowed in blocks (this parameter is ignored)",
	"templaterequest-maxversion":   "Highest supported block version number (this parameter is ignored)",
	"templaterequest-target":       "The desired target for the block template (this parameter is ignored)",
	"templaterequest-data":         "Hex-encoded block data (only for mode=proposal)",
	"templaterequest-workid":       "The server provided workid if provided in block template (not applicable)",

	// GetBlockTemplateResultTx help.
	"getblocktemplateresulttx-data":    "Hex-encoded transaction data (byte-for-byte)",
	"getblocktemplateresulttx-hash":    "Hex-encoded transaction hash (little endian if treated as a 256-bit number)",
	"getblocktemplateresulttx-depends": "Other transactions before this one (by 1-based index in the 'transactions'  list) that must be present in the final block if this one is",
	"getblocktemplateresulttx-fee":     "Difference in value between transaction inputs and outputs (in Satoshi)",
	"getblocktemplateresulttx-sigops":  "Total number of signature operations as counted for purposes of block limits",

	// GetBlockTemplateResultAux help.
	"getblocktemplateresultaux-flags": "Hex-encoded byte-for-byte data to include in the coinbase signature script",

	// GetBlockTemplateResult help.
	"getblocktemplateresult-bits":              "Hex-encoded compressed difficulty",
	"getblocktemplateresult-curtime":           "Current time as seen by the server (recommended for block time); must fall within mintime/maxtime rules",
	"getblocktemplateresult-height":            "Height of the block to be solved",
	"getblocktemplateresult-parentblockhashes": "Hex-encoded big-endian hashes of the parent blocks",
	"getblocktemplateresult-sigoplimit":        "Number of sigops allowed in blocks ",
	"getblocktemplateresult-sizelimit":         "Number of bytes allowed in blocks",
	"getblocktemplateresult-transactions":      "Array of transactions as JSON objects",
	"getblocktemplateresult-version":           "The block version",
	"getblocktemplateresult-coinbaseaux":       "Data that should be included in the coinbase signature script",
	"getblocktemplateresult-coinbasetxn":       "Information about the coinbase transaction",
	"getblocktemplateresult-coinbasevalue":     "Total amount available for the coinbase in Satoshi",
	"getblocktemplateresult-workid":            "This value must be returned with result if provided (not provided)",
	"getblocktemplateresult-longpollid":        "Identifier for long poll request which allows monitoring for expiration",
	"getblocktemplateresult-longpolluri":       "An alternate URI to use for long poll requests if provided (not provided)",
	"getblocktemplateresult-submitold":         "Not applicable",
	"getblocktemplateresult-target":            "Hex-encoded big-endian number which valid results must be less than",
	"getblocktemplateresult-expires":           "Maximum number of seconds (starting from when the server sent the response) this work is valid for",
	"getblocktemplateresult-maxtime":           "Maximum allowed time",
	"getblocktemplateresult-mintime":           "Minimum allowed time",
	"getblocktemplateresult-mutable":           "List of mutations the server explicitly allows",
	"getblocktemplateresult-noncerange":        "Two concatenated hex-encoded big-endian 64-bit integers which represent the valid ranges of nonces the miner may scan",
	"getblocktemplateresult-capabilities":      "List of server capabilities including 'proposal' to indicate support for block proposals",
	"getblocktemplateresult-reject-reason":     "Reason the proposal was invalid as-is (only applies to proposal responses)",

	// GetBlockTemplateCmd help.
	"getBlockTemplate--synopsis": "Returns a JSON object with information necessary to construct a block to mine or accepts a proposal to validate.\n" +
		"See BIP0022 and BIP0023 for the full specification.",
	"getBlockTemplate-request":     "Request object which controls the mode and several parameters",
	"getBlockTemplate--condition0": "mode=template",
	"getBlockTemplate--condition1": "mode=proposal, rejected",
	"getBlockTemplate--condition2": "mode=proposal, accepted",
	"getBlockTemplate--result1":    "An error string which represents why the proposal was rejected or nothing if accepted",

	// GetCFilterCmd help.
	"getCFilter--synopsis":  "Returns a block's committed filter given its hash.",
	"getCFilter-filtertype": "The type of filter to return (0=regular, 1=extended)",
	"getCFilter-hash":       "The hash of the block",
	"getCFilter--result0":   "The block's committed filter",

	// GetCFilterHeaderCmd help.
	"getCFilterHeader--synopsis":  "Returns a block's compact filter header given its hash.",
	"getCFilterHeader-filtertype": "The type of filter header to return (0=regular, 1=extended)",
	"getCFilterHeader-hash":       "The hash of the block",
	"getCFilterHeader--result0":   "The block's gcs filter header",

	// GetConnectionCountCmd help.
	"getConnectionCount--synopsis": "Returns the number of active connections to other peers.",
	"getConnectionCount--result0":  "The number of connections",

	// GetCurrentNetCmd help.
	"getCurrentNet--synopsis": "Get bitcoin network the server is running on.",
	"getCurrentNet--result0":  "The network identifer",

	// GetDifficultyCmd help.
	"getDifficulty--synopsis": "Returns the proof-of-work difficulty as a multiple of the minimum difficulty.",
	"getDifficulty--result0":  "The difficulty",

	// GetGenerateCmd help.
	"getGenerate--synopsis": "Returns if the server is set to generate coins (mine) or not.",
	"getGenerate--result0":  "True if mining, false if not",

	// GetHashesPerSecCmd help.
	"getHashesPerSec--synopsis": "Returns a recent hashes per second performance measurement while generating coins (mining).",
	"getHashesPerSec--result0":  "The number of hashes per second",

	// InfoDAGResult help.
	"infodagresult-version":         "The version of the server",
	"infodagresult-protocolversion": "The latest supported protocol version",
	"infodagresult-blocks":          "The number of blocks processed",
	"infodagresult-timeoffset":      "The time offset",
	"infodagresult-connections":     "The number of connected peers",
	"infodagresult-proxy":           "The proxy used by the server",
	"infodagresult-difficulty":      "The current target difficulty",
	"infodagresult-testnet":         "Whether or not server is using testnet",
	"infodagresult-relayfee":        "The minimum relay fee for non-free transactions in BTC/KB",
	"infodagresult-errors":          "Any current errors",

	// InfoWalletResult help.
	"infowalletresult-version":         "The version of the server",
	"infowalletresult-protocolversion": "The latest supported protocol version",
	"infowalletresult-walletversion":   "The version of the wallet server",
	"infowalletresult-balance":         "The total bitcoin balance of the wallet",
	"infowalletresult-blocks":          "The number of blocks processed",
	"infowalletresult-timeoffset":      "The time offset",
	"infowalletresult-connections":     "The number of connected peers",
	"infowalletresult-proxy":           "The proxy used by the server",
	"infowalletresult-difficulty":      "The current target difficulty",
	"infowalletresult-testnet":         "Whether or not server is using testnet",
	"infowalletresult-keypoololdest":   "Seconds since 1 Jan 1970 GMT of the oldest pre-generated key in the key pool",
	"infowalletresult-keypoolsize":     "The number of new keys that are pre-generated",
	"infowalletresult-unlocked_until":  "The timestamp in seconds since 1 Jan 1970 GMT that the wallet is unlocked for transfers, or 0 if the wallet is locked",
	"infowalletresult-paytxfee":        "The transaction fee set in BTC/KB",
	"infowalletresult-relayfee":        "The minimum relay fee for non-free transactions in BTC/KB",
	"infowalletresult-errors":          "Any current errors",

	// GetHeadersCmd help.
	"getHeaders--synopsis":     "Returns block headers starting with the first known block hash from the request",
	"getHeaders-blocklocators": "JSON array of hex-encoded hashes of blocks.  Headers are returned starting from the first known hash in this list",
	"getHeaders-hashstop":      "Block hash to stop including block headers for; if not found, all headers to the latest known block are returned.",
	"getHeaders--result0":      "Serialized block headers of all located blocks, limited to some arbitrary maximum number of hashes (currently 2000, which matches the wire protocol headers message, but this is not guaranteed)",

	// GetInfoCmd help.
	"getInfo--synopsis": "Returns a JSON object containing various state info.",

	// GetMempoolInfoCmd help.
	"getMempoolInfo--synopsis": "Returns memory pool information",

	// GetMempoolInfoResult help.
	"getmempoolinforesult-bytes": "Size in bytes of the mempool",
	"getmempoolinforesult-size":  "Number of transactions in the mempool",

	// GetMiningInfoResult help.
	"getmininginforesult-blocks":           "Height of the latest best block",
	"getmininginforesult-currentblocksize": "Size of the latest best block",
	"getmininginforesult-currentblocktx":   "Number of transactions in the latest best block",
	"getmininginforesult-difficulty":       "Current target difficulty",
	"getmininginforesult-errors":           "Any current errors",
	"getmininginforesult-generate":         "Whether or not server is set to generate coins",
	"getmininginforesult-genproclimit":     "Number of processors to use for coin generation (-1 when disabled)",
	"getmininginforesult-hashespersec":     "Recent hashes per second performance measurement while generating coins",
	"getmininginforesult-networkhashps":    "Estimated network hashes per second for the most recent blocks",
	"getmininginforesult-pooledtx":         "Number of transactions in the memory pool",
	"getmininginforesult-testnet":          "Whether or not server is using testnet",

	// GetMiningInfoCmd help.
	"getMiningInfo--synopsis": "Returns a JSON object containing mining-related information.",

	// GetNetworkHashPSCmd help.
	"getNetworkHashPs--synopsis": "Returns the estimated network hashes per second for the block heights provided by the parameters.",
	"getNetworkHashPs-blocks":    "The number of blocks, or -1 for blocks since last difficulty change",
	"getNetworkHashPs-height":    "Perform estimate ending with this height or -1 for current best chain block height",
	"getNetworkHashPs--result0":  "Estimated hashes per second",

	// GetNetTotalsCmd help.
	"getNetTotals--synopsis": "Returns a JSON object containing network traffic statistics.",

	// GetNetTotalsResult help.
	"getnettotalsresult-totalbytesrecv": "Total bytes received",
	"getnettotalsresult-totalbytessent": "Total bytes sent",
	"getnettotalsresult-timemillis":     "Number of milliseconds since 1 Jan 1970 GMT",

	// GetPeerInfoResult help.
	"getpeerinforesult-id":             "A unique node ID",
	"getpeerinforesult-addr":           "The ip address and port of the peer",
	"getpeerinforesult-services":       "Services bitmask which represents the services supported by the peer",
	"getpeerinforesult-relaytxes":      "Peer has requested transactions be relayed to it",
	"getpeerinforesult-lastsend":       "Time the last message was received in seconds since 1 Jan 1970 GMT",
	"getpeerinforesult-lastrecv":       "Time the last message was sent in seconds since 1 Jan 1970 GMT",
	"getpeerinforesult-bytessent":      "Total bytes sent",
	"getpeerinforesult-bytesrecv":      "Total bytes received",
	"getpeerinforesult-conntime":       "Time the connection was made in seconds since 1 Jan 1970 GMT",
	"getpeerinforesult-timeoffset":     "The time offset of the peer",
	"getpeerinforesult-pingtime":       "Number of microseconds the last ping took",
	"getpeerinforesult-pingwait":       "Number of microseconds a queued ping has been waiting for a response",
	"getpeerinforesult-version":        "The protocol version of the peer",
	"getpeerinforesult-subver":         "The user agent of the peer",
	"getpeerinforesult-inbound":        "Whether or not the peer is an inbound connection",
	"getpeerinforesult-startingheight": "The latest block height the peer knew about when the connection was established",
	"getpeerinforesult-currentheight":  "The current height of the peer",
	"getpeerinforesult-banscore":       "The ban score",
	"getpeerinforesult-feefilter":      "The requested minimum fee a transaction must have to be announced to the peer",
	"getpeerinforesult-syncnode":       "Whether or not the peer is the sync peer",

	// GetPeerInfoCmd help.
	"getPeerInfo--synopsis": "Returns data about each connected network peer as an array of json objects.",

	// GetRawMempoolVerboseResult help.
	"getrawmempoolverboseresult-size":             "Transaction size in bytes",
	"getrawmempoolverboseresult-fee":              "Transaction fee in bitcoins",
	"getrawmempoolverboseresult-time":             "Local time transaction entered pool in seconds since 1 Jan 1970 GMT",
	"getrawmempoolverboseresult-height":           "Block height when transaction entered the pool",
	"getrawmempoolverboseresult-startingpriority": "Priority when transaction entered the pool",
	"getrawmempoolverboseresult-currentpriority":  "Current priority",
	"getrawmempoolverboseresult-depends":          "Unconfirmed transactions used as inputs for this transaction",

	// GetRawMempoolCmd help.
	"getRawMempool--synopsis":   "Returns information about all of the transactions currently in the memory pool.",
	"getRawMempool-verbose":     "Returns JSON object when true or an array of transaction hashes when false",
	"getRawMempool--condition0": "verbose=false",
	"getRawMempool--condition1": "verbose=true",
	"getRawMempool--result0":    "Array of transaction hashes",

	// GetRawTransactionCmd help.
	"getRawTransaction--synopsis":   "Returns information about a transaction given its hash.",
	"getRawTransaction-txid":        "The hash of the transaction",
	"getRawTransaction-verbose":     "Specifies the transaction is returned as a JSON object instead of a hex-encoded string",
	"getRawTransaction--condition0": "verbose=false",
	"getRawTransaction--condition1": "verbose=true",
	"getRawTransaction--result0":    "Hex-encoded bytes of the serialized transaction",

	// GetTxOutResult help.
	"gettxoutresult-bestblock":     "The block hash that contains the transaction output",
	"gettxoutresult-confirmations": "The number of confirmations",
	"gettxoutresult-value":         "The transaction amount in BTC",
	"gettxoutresult-scriptPubKey":  "The public key script used to pay coins as a JSON object",
	"gettxoutresult-version":       "The transaction version",
	"gettxoutresult-coinbase":      "Whether or not the transaction is a coinbase",

	// GetTxOutCmd help.
	"getTxOut--synopsis":      "Returns information about an unspent transaction output..",
	"getTxOut-txid":           "The hash of the transaction",
	"getTxOut-vout":           "The index of the output",
	"getTxOut-includemempool": "Include the mempool when true",

	// HelpCmd help.
	"help--synopsis":   "Returns a list of all commands or help for a specified command.",
	"help-command":     "The command to retrieve help for",
	"help--condition0": "no command provided",
	"help--condition1": "command specified",
	"help--result0":    "List of commands",
	"help--result1":    "Help for specified command",

	// PingCmd help.
	"ping--synopsis": "Queues a ping to be sent to each connected peer.\n" +
		"Ping times are provided by getPeerInfo via the pingtime and pingwait fields.",

	// RemoveManualNodeCmd help.
	"removeManualNode--synopsis": "Removes a peer from the manual nodes list",
	"removeManualNode-addr":      "IP address and port of the peer to remove",

	// SearchRawTransactionsCmd help.
	"searchRawTransactions--synopsis": "Returns raw data for transactions involving the passed address.\n" +
		"Returned transactions are pulled from both the database, and transactions currently in the mempool.\n" +
		"Transactions pulled from the mempool will have the 'confirmations' field set to 0.\n" +
		"Usage of this RPC requires the optional --addrindex flag to be activated, otherwise all responses will simply return with an error stating the address index has not yet been built.\n" +
		"Similarly, until the address index has caught up with the current best height, all requests will return an error response in order to avoid serving stale data.",
	"searchRawTransactions-address":     "The Bitcoin address to search for",
	"searchRawTransactions-verbose":     "Specifies the transaction is returned as a JSON object instead of hex-encoded string",
	"searchRawTransactions--condition0": "verbose=0",
	"searchRawTransactions--condition1": "verbose=1",
	"searchRawTransactions-skip":        "The number of leading transactions to leave out of the final response",
	"searchRawTransactions-count":       "The maximum number of transactions to return",
	"searchRawTransactions-vinextra":    "Specify that extra data from previous output will be returned in vin",
	"searchRawTransactions-reverse":     "Specifies that the transactions should be returned in reverse chronological order",
	"searchRawTransactions-filteraddrs": "Address list.  Only inputs or outputs with matching address will be returned",
	"searchRawTransactions--result0":    "Hex-encoded serialized transaction",

	// SendRawTransactionCmd help.
	"sendRawTransaction--synopsis":     "Submits the serialized, hex-encoded transaction to the local peer and relays it to the network.",
	"sendRawTransaction-hextx":         "Serialized, hex-encoded signed transaction",
	"sendRawTransaction-allowhighfees": "Whether or not to allow insanely high fees (btcd does not yet implement this parameter, so it has no effect)",
	"sendRawTransaction--result0":      "The hash of the transaction",

	// SetGenerateCmd help.
	"setGenerate--synopsis":    "Set the server to generate coins (mine) or not.",
	"setGenerate-generate":     "Use true to enable generation, false to disable it",
	"setGenerate-genproclimit": "The number of processors (cores) to limit generation to or -1 for default",

	// StopCmd help.
	"stop--synopsis": "Shutdown btcd.",
	"stop--result0":  "The string 'btcd stopping.'",

	// SubmitBlockOptions help.
	"submitblockoptions-workid": "This parameter is currently ignored",

	// SubmitBlockCmd help.
	"submitBlock--synopsis":   "Attempts to submit a new serialized, hex-encoded block to the network.",
	"submitBlock-hexblock":    "Serialized, hex-encoded block",
	"submitBlock-options":     "This parameter is currently ignored",
	"submitBlock--condition0": "Block successfully submitted",
	"submitBlock--condition1": "Block rejected",
	"submitBlock--result1":    "The reason the block was rejected",

	// ValidateAddressResult help.
	"validateaddressresult-isvalid": "Whether or not the address is valid",
	"validateaddressresult-address": "The bitcoin address (only when isvalid is true)",

	// ValidateAddressCmd help.
	"validateAddress--synopsis": "Verify an address is valid.",
	"validateAddress-address":   "Bitcoin address to validate",

	// VerifyChainCmd help.
	"verifyDag--synopsis": "Verifies the block DAG database.\n" +
		"The actual checks performed by the checklevel parameter are implementation specific.\n" +
		"For btcd this is:\n" +
		"checklevel=0 - Look up each block and ensure it can be loaded from the database.\n" +
		"checklevel=1 - Perform basic context-free sanity checks on each block.",
	"verifyDag-checklevel": "How thorough the block verification is",
	"verifyDag-checkdepth": "The number of blocks to check",
	"verifyDag--result0":   "Whether or not the DAG verified",

	// VerifyMessageCmd help.
	"verifyMessage--synopsis": "Verify a signed message.",
	"verifyMessage-address":   "The bitcoin address to use for the signature",
	"verifyMessage-signature": "The base-64 encoded signature provided by the signer",
	"verifyMessage-message":   "The signed message",
	"verifyMessage--result0":  "Whether or not the signature verified",

	// -------- Websocket-specific help --------

	// Session help.
	"session--synopsis":       "Return details regarding a websocket client's current connection session.",
	"sessionresult-sessionid": "The unique session ID for a client's websocket connection.",

	// NotifyBlocksCmd help.
	"notifyBlocks--synopsis": "Request notifications for whenever a block is connected or disconnected from the main (best) chain.",

	// StopNotifyBlocksCmd help.
	"stopNotifyBlocks--synopsis": "Cancel registered notifications for whenever a block is connected or disconnected from the main (best) chain.",

	// NotifyNewTransactionsCmd help.
	"notifyNewTransactions--synopsis": "Send either a txaccepted or a txacceptedverbose notification when a new transaction is accepted into the mempool.",
	"notifyNewTransactions-verbose":   "Specifies which type of notification to receive. If verbose is true, then the caller receives txacceptedverbose, otherwise the caller receives txaccepted",

	// StopNotifyNewTransactionsCmd help.
	"stopNotifyNewTransactions--synopsis": "Stop sending either a txaccepted or a txacceptedverbose notification when a new transaction is accepted into the mempool.",

	// NotifyReceivedCmd help.
	"notifyReceived--synopsis": "Send a recvtx notification when a transaction added to mempool or appears in a newly-attached block contains a txout pkScript sending to any of the passed addresses.\n" +
		"Matching outpoints are automatically registered for redeemingtx notifications.",
	"notifyReceived-addresses": "List of address to receive notifications about",

	// StopNotifyReceivedCmd help.
	"stopNotifyReceived--synopsis": "Cancel registered receive notifications for each passed address.",
	"stopNotifyReceived-addresses": "List of address to cancel receive notifications for",

	// OutPoint help.
	"outpoint-hash":  "The hex-encoded bytes of the outPoint hash",
	"outpoint-index": "The index of the outPoint",

	// NotifySpentCmd help.
	"notifySpent--synopsis": "Send a redeemingtx notification when a transaction spending an outPoint appears in mempool (if relayed to this btcd instance) and when such a transaction first appears in a newly-attached block.",
	"notifySpent-outpoints": "List of transaction outpoints to monitor.",

	// StopNotifySpentCmd help.
	"stopNotifySpent--synopsis": "Cancel registered spending notifications for each passed outPoint.",
	"stopNotifySpent-outpoints": "List of transaction outpoints to stop monitoring.",

	// LoadTxFilterCmd help.
	"loadTxFilter--synopsis": "Load, add to, or reload a websocket client's transaction filter for mempool transactions, new blocks and rescanBlocks.",
	"loadTxFilter-reload":    "Load a new filter instead of adding data to an existing one",
	"loadTxFilter-addresses": "Array of addresses to add to the transaction filter",
	"loadTxFilter-outpoints": "Array of outpoints to add to the transaction filter",

	// Rescan help.
	"rescan--synopsis": "Rescan block chain for transactions to addresses.\n" +
		"When the endblock parameter is omitted, the rescan continues through the best block in the main chain.\n" +
		"Rescan results are sent as recvtx and redeemingtx notifications.\n" +
		"This call returns once the rescan completes.",
	"rescan-beginblock": "Hash of the first block to begin rescanning",
	"rescan-addresses":  "List of addresses to include in the rescan",
	"rescan-outpoints":  "List of transaction outpoints to include in the rescan",
	"rescan-endblock":   "Hash of final block to rescan",

	// RescanBlocks help.
	"rescanBlocks--synopsis":   "Rescan blocks for transactions matching the loaded transaction filter.",
	"rescanBlocks-blockhashes": "List of hashes to rescan.  Each next block must be a child of the previous.",
	"rescanBlocks--result0":    "List of matching blocks.",

	// RescannedBlock help.
	"rescannedblock-hash":         "Hash of the matching block.",
	"rescannedblock-transactions": "List of matching transactions, serialized and hex-encoded.",

	// Uptime help.
	"uptime--synopsis": "Returns the total uptime of the server.",
	"uptime--result0":  "The number of seconds that the server has been running",

	// Version help.
	"version--synopsis":       "Returns the JSON-RPC API version (semver)",
	"version--result0--desc":  "Version objects keyed by the program or API name",
	"version--result0--key":   "Program or API name",
	"version--result0--value": "Object containing the semantic version",

	// VersionResult help.
	"versionresult-versionstring": "The JSON-RPC API version (semver)",
	"versionresult-major":         "The major component of the JSON-RPC API version",
	"versionresult-minor":         "The minor component of the JSON-RPC API version",
	"versionresult-patch":         "The patch component of the JSON-RPC API version",
	"versionresult-prerelease":    "Prerelease info about the current build",
	"versionresult-buildmetadata": "Metadata about the current build",
}

// rpcResultTypes specifies the result types that each RPC command can return.
// This information is used to generate the help.  Each result type must be a
// pointer to the type (or nil to indicate no return value).
var rpcResultTypes = map[string][]interface{}{
	"addManualNode":         nil,
	"createRawTransaction":  {(*string)(nil)},
	"debugLevel":            {(*string)(nil), (*string)(nil)},
	"decodeRawTransaction":  {(*btcjson.TxRawDecodeResult)(nil)},
	"decodeScript":          {(*btcjson.DecodeScriptResult)(nil)},
	"estimateFee":           {(*float64)(nil)},
	"generate":              {(*[]string)(nil)},
	"getAllManualNodesInfo": {(*[]string)(nil), (*[]btcjson.GetManualNodeInfoResult)(nil)},
	"getBestBlock":          {(*btcjson.GetBestBlockResult)(nil)},
	"getBestBlockHash":      {(*string)(nil)},
	"getBlock":              {(*string)(nil), (*btcjson.GetBlockVerboseResult)(nil)},
	"getBlockCount":         {(*int64)(nil)},
	"getBlockHash":          {(*string)(nil)},
	"getBlockHeader":        {(*string)(nil), (*btcjson.GetBlockHeaderVerboseResult)(nil)},
	"getBlockTemplate":      {(*btcjson.GetBlockTemplateResult)(nil), (*string)(nil), nil},
	"getBlockDagInfo":       {(*btcjson.GetBlockDAGInfoResult)(nil)},
	"getCFilter":            {(*string)(nil)},
	"getCFilterHeader":      {(*string)(nil)},
	"getConnectionCount":    {(*int32)(nil)},
	"getCurrentNet":         {(*uint32)(nil)},
	"getDifficulty":         {(*float64)(nil)},
	"getGenerate":           {(*bool)(nil)},
	"getHashesPerSec":       {(*float64)(nil)},
	"getHeaders":            {(*[]string)(nil)},
	"getInfo":               {(*btcjson.InfoDAGResult)(nil)},
	"getManualNodeInfo":     {(*string)(nil), (*btcjson.GetManualNodeInfoResult)(nil)},
	"getMempoolInfo":        {(*btcjson.GetMempoolInfoResult)(nil)},
	"getMiningInfo":         {(*btcjson.GetMiningInfoResult)(nil)},
	"getNetTotals":          {(*btcjson.GetNetTotalsResult)(nil)},
	"getNetworkHashPs":      {(*int64)(nil)},
	"getPeerInfo":           {(*[]btcjson.GetPeerInfoResult)(nil)},
	"getRawMempool":         {(*[]string)(nil), (*btcjson.GetRawMempoolVerboseResult)(nil)},
	"getRawTransaction":     {(*string)(nil), (*btcjson.TxRawResult)(nil)},
	"getTxOut":              {(*btcjson.GetTxOutResult)(nil)},
	"node":                  nil,
	"help":                  {(*string)(nil), (*string)(nil)},
	"ping":                  nil,
	"removeManualNode":      nil,
	"searchRawTransactions": {(*string)(nil), (*[]btcjson.SearchRawTransactionsResult)(nil)},
	"sendRawTransaction":    {(*string)(nil)},
	"setGenerate":           nil,
	"stop":                  {(*string)(nil)},
	"submitBlock":           {nil, (*string)(nil)},
	"uptime":                {(*int64)(nil)},
	"validateAddress":       {(*btcjson.ValidateAddressResult)(nil)},
	"verifyDag":             {(*bool)(nil)},
	"verifyMessage":         {(*bool)(nil)},
	"version":               {(*map[string]btcjson.VersionResult)(nil)},

	// Websocket commands.
	"loadTxFilter":              nil,
	"session":                   {(*btcjson.SessionResult)(nil)},
	"notifyBlocks":              nil,
	"stopNotifyBlocks":          nil,
	"notifyNewTransactions":     nil,
	"stopNotifyNewTransactions": nil,
	"notifyReceived":            nil,
	"stopNotifyReceived":        nil,
	"notifySpent":               nil,
	"stopNotifySpent":           nil,
	"rescan":                    nil,
	"rescanBlocks":              {(*[]btcjson.RescannedBlock)(nil)},
}

// helpCacher provides a concurrent safe type that provides help and usage for
// the RPC server commands and caches the results for future calls.
type helpCacher struct {
	sync.Mutex
	usage      string
	methodHelp map[string]string
}

// rpcMethodHelp returns an RPC help string for the provided method.
//
// This function is safe for concurrent access.
func (c *helpCacher) rpcMethodHelp(method string) (string, error) {
	c.Lock()
	defer c.Unlock()

	// Return the cached method help if it exists.
	if help, exists := c.methodHelp[method]; exists {
		return help, nil
	}

	// Look up the result types for the method.
	resultTypes, ok := rpcResultTypes[method]
	if !ok {
		return "", errors.New("no result types specified for method " +
			method)
	}

	// Generate, cache, and return the help.
	help, err := btcjson.GenerateHelp(method, helpDescsEnUS, resultTypes...)
	if err != nil {
		return "", err
	}
	c.methodHelp[method] = help
	return help, nil
}

// rpcUsage returns one-line usage for all support RPC commands.
//
// This function is safe for concurrent access.
func (c *helpCacher) rpcUsage(includeWebsockets bool) (string, error) {
	c.Lock()
	defer c.Unlock()

	// Return the cached usage if it is available.
	if c.usage != "" {
		return c.usage, nil
	}

	// Generate a list of one-line usage for every command.
	usageTexts := make([]string, 0, len(rpcHandlers))
	for k := range rpcHandlers {
		usage, err := btcjson.MethodUsageText(k)
		if err != nil {
			return "", err
		}
		usageTexts = append(usageTexts, usage)
	}

	// Include websockets commands if requested.
	if includeWebsockets {
		for k := range wsHandlers {
			usage, err := btcjson.MethodUsageText(k)
			if err != nil {
				return "", err
			}
			usageTexts = append(usageTexts, usage)
		}
	}

	sort.Sort(sort.StringSlice(usageTexts))
	c.usage = strings.Join(usageTexts, "\n")
	return c.usage, nil
}

// newHelpCacher returns a new instance of a help cacher which provides help and
// usage for the RPC server commands and caches the results for future calls.
func newHelpCacher() *helpCacher {
	return &helpCacher{
		methodHelp: make(map[string]string),
	}
}