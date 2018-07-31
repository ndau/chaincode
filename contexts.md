# Contexts for Chaincode VMs

## 	0x00: "TEST"

This is a context used only for testing; input, output, and error behavior are all specified by the individual test.

##  0x01: "NODE_PAYOUT"

Calculates node payout awards proportionally by stake, with a fixed percentage to the node operator.

This script is run with these values pushed on the stack in this order:

2: a list of stakers, consisting of structs, where each struct has:
    * address of the staked account
    * amount staked in napu
    * timestamp when the stake occurred
1: the total node payout amount in napu
0: the address of the staked account for which the node reward is being calculated (the node operator's address)

The return value is a list containing addresses and a payout for each address.

If the script exits with an error, the entire payout goes to the node operator's address, and the occurrence of the error is recorded and associated with this node.


##  0x02: "EAI_TIMING"
##  0x03: "NODE_QUALITY"
##  0x04: "MARKET_PRICE"

##  0x05: "ACCOUNT"

Validates signatures for an account prior to executing a transaction.

This script is run with these values pushed on the stack in this order:

1: A struct corresponding to the transaction under evaluation
0: A struct corresponding to the state of the source account for this transaction (the account that includes this script)

If the top value of the stack upon return is anything other than exactly equal to True (1), the transaction fails. (Note that "truthiness" does not apply -- anything not equal to 1 is considered failure).
