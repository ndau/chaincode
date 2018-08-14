package main

// Predefined constants available to chasm programs.

func predefinedConstants() map[string]string {
	k := map[string]string{
		"Event_Default":                "0",
		"Event_Transfer":               "1",
		"Event_ChangeTransferKey":      "2",
		"Event_ReleaseFromEndowment":   "3",
		"Event_ChangeSettlementPeriod": "4",
		"Event_Delegate":               "5",
		"Event_ComputeEAI":             "6",
		"Event_Lock":                   "7",
		"Event_Notify":                 "8",
		"Event_SetRewardsTarget":       "9",
		"Event_GTValidatorChange":      "0xff",
		"Tx_Timestamp":                 "0",
		"Tx_Source":                    "1",
		"Tx_Destination":               "2",
		"Tx_Qty":                       "3",
		"Tx_Signatures":                "4",
		"Account_Balance":              "0",
		"Account_HasLock":              "1",
		"Account_Lock":                 "2",
		"Account_HasStake":             "3",
		"Account_Stake":                "4",
		"Account_WeightedAverageAge":   "5",
		"Account_Keys":                 "6",
		"Signature_Key":                "0",
	}

	return k
}
