/*
 * Copyright (c) 2023. Liero AB
 *
 * Permission is hereby granted, free of charge, to any person obtaining
 * a copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software
 * is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
 * OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
 * IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
 * CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
 * TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package oids

// System
const (
	SysPrefix       = ".1.3.6.1.2.1.1"
	SysDescr        = ".1.3.6.1.2.1.1.1.0"
	SysObjectID     = ".1.3.6.1.2.1.1.2.0"
	SysUpTime       = ".1.3.6.1.2.1.1.3.0"
	SysContact      = ".1.3.6.1.2.1.1.4.0"
	SysName         = ".1.3.6.1.2.1.1.5.0"
	SysLocation     = ".1.3.6.1.2.1.1.6.0"
	SysServices     = ".1.3.6.1.2.1.1.7.0"
	SysORLastChange = ".1.3.6.1.2.1.1.8.0"
)

// Default ifTable Oids
const (
	// OIDS translations
	IfXEntryPrefix            = ".1.3.6.1.2.1.31.1.1.1"
	IfXEntry                  = ".1.3.6.1.2.1.31.1.1.1.1"
	IfEntryPrefix             = ".1.3.6.1.2.1.2.2.1"
	IfIndex            string = ".1.3.6.1.2.1.2.2.1.1"
	IfDescr                   = ".1.3.6.1.2.1.2.2.1.2"
	IfType                    = ".1.3.6.1.2.1.2.2.1.3"
	IfMtu                     = ".1.3.6.1.2.1.2.2.1.4"
	IfSpeed                   = ".1.3.6.1.2.1.2.2.1.5"
	IfPhysAddress             = ".1.3.6.1.2.1.2.2.1.6"
	IfAdminStatus             = ".1.3.6.1.2.1.2.2.1.7"
	IfOperStatus              = ".1.3.6.1.2.1.2.2.1.8"
	IfLastChange              = ".1.3.6.1.2.1.2.2.1.9"
	IfAlias                   = ".1.3.6.1.2.1.31.1.1.1.18"
	IfHighSpeed               = ".1.3.6.1.2.1.31.1.1.1.15"
	IfConnectorPresent        = ".1.3.6.1.2.1.31.1.1.1.17"
	IfEntPhysicalName         = ".1.3.6.1.2.1.47.1.1.1.1.7"

	// Input
	IfInOctets        = ".1.3.6.1.2.1.2.2.1.10"
	IfInUcastPkts     = ".1.3.6.1.2.1.2.2.1.11"
	IfInNUcastPkts    = ".1.3.6.1.2.1.2.2.1.12"
	IfInDiscards      = ".1.3.6.1.2.1.2.2.1.13"
	IfInErrors        = ".1.3.6.1.2.1.2.2.1.14"
	IfInMulticastPkts = ".1.3.6.1.2.1.31.1.1.1.2"
	IfInBroadcastPkts = ".1.3.6.1.2.1.31.1.1.1.3"

	// Output
	IfOutOctets        = ".1.3.6.1.2.1.2.2.1.16"
	IfOutUcastPkts     = ".1.3.6.1.2.1.2.2.1.17"
	IfOutNUcastPkts    = ".1.3.6.1.2.1.2.2.1.18"
	IfOutDiscards      = ".1.3.6.1.2.1.2.2.1.19"
	IfOutErrors        = ".1.3.6.1.2.1.2.2.1.20"
	IfOutMulticastPkts = ".1.3.6.1.2.1.31.1.1.1.4"
	IfOutBroadcastPkts = ".1.3.6.1.2.1.31.1.1.1.5"

	// Output HC
	IfHCOutOctets        = ".1.3.6.1.2.1.31.1.1.1.10"
	IfHCOutUcastPkts     = ".1.3.6.1.2.1.31.1.1.1.11"
	IfHCOutMulticastPkts = ".1.3.6.1.2.1.31.1.1.1.12"
	IfHCOutBroadcastPkts = ".1.3.6.1.2.1.31.1.1.1.13"

	// Input HC
	IfHCInOctets        = ".1.3.6.1.2.1.31.1.1.1.6"
	IfHCInUcastPkts     = ".1.3.6.1.2.1.31.1.1.1.7"
	IfHCInMulticastPkts = ".1.3.6.1.2.1.31.1.1.1.8"
	IfHCInBroadcastPkts = ".1.3.6.1.2.1.31.1.1.1.9"

	// ifXTable
	IfLinkUpDownTrapEnable = ".1.3.6.1.2.1.31.1.1.1.14"
)

// OIDS that needs to be formatted with index
const (
	IfDescrF       = ".1.3.6.1.2.1.2.2.1.2.%d"
	IfTypeF        = ".1.3.6.1.2.1.2.2.1.3.%d"
	IfMtuF         = ".1.3.6.1.2.1.2.2.1.4.%d"
	IfSpeedF       = ".1.3.6.1.2.1.2.2.1.5.%d"
	IfPhysAddressF = ".1.3.6.1.2.1.2.2.1.6.%d"
	IfAdminStatusF = ".1.3.6.1.2.1.2.2.1.7.%d"
	IfOperStatusF  = ".1.3.6.1.2.1.2.2.1.8.%d"
	IfLastChangeF  = ".1.3.6.1.2.1.2.2.1.9.%d"

	IfHighSpeedF        = ".1.3.6.1.2.1.31.1.1.1.15.%d"
	IfConnectorPresentF = ".1.3.6.1.2.1.31.1.1.1.17.%d"
	IfAliasF            = ".1.3.6.1.2.1.31.1.1.1.18.%d"

	// Input
	IfInOctetsF        = ".1.3.6.1.2.1.2.2.1.10.%d"
	IfInUcastPktsF     = ".1.3.6.1.2.1.2.2.1.11.%d"
	IfInNUcastPktsF    = ".1.3.6.1.2.1.2.2.1.12.%d"
	IfInDiscardsF      = ".1.3.6.1.2.1.2.2.1.13.%d"
	IfInErrorsF        = ".1.3.6.1.2.1.2.2.1.14.%d"
	IfInMulticastPktsF = ".1.3.6.1.2.1.31.1.1.1.2.%d"
	IfInBroadcastPktsF = ".1.3.6.1.2.1.31.1.1.1.3.%d"

	// Output
	IfOutOctetsF        = ".1.3.6.1.2.1.2.2.1.16.%d"
	IfOutUcastPktsF     = ".1.3.6.1.2.1.2.2.1.17.%d"
	IfOutNUcastPktsF    = ".1.3.6.1.2.1.2.2.1.18.%d"
	IfOutDiscardsF      = ".1.3.6.1.2.1.2.2.1.19.%d"
	IfOutErrorsF        = ".1.3.6.1.2.1.2.2.1.20.%d"
	IfOutMulticastPktsF = ".1.3.6.1.2.1.31.1.1.1.4.%d"
	IfOutBroadcastPktsF = ".1.3.6.1.2.1.31.1.1.1.5.%d"

	// Out
	IfHCOutOctetsF        = ".1.3.6.1.2.1.31.1.1.1.10.%d"
	IfHCOutUcastPktsF     = ".1.3.6.1.2.1.31.1.1.1.11.%d"
	IfHCOutMulticastPktsF = ".1.3.6.1.2.1.31.1.1.1.12.%d"
	IfHCOutBroadcastPktsF = ".1.3.6.1.2.1.31.1.1.1.13.%d"

	// In
	IfHCInOctetsF        = ".1.3.6.1.2.1.31.1.1.1.6.%d"
	IfHCInUcastPktsF     = ".1.3.6.1.2.1.31.1.1.1.7.%d"
	IfHCInMulticastPktsF = ".1.3.6.1.2.1.31.1.1.1.8.%d"
	IfHCInBroadcastPktsF = ".1.3.6.1.2.1.31.1.1.1.9.%d"

	// ifXTable
	IfLinkUpDownTrapEnableF = ".1.3.6.1.2.1.31.1.1.1.14.%d"
)

// Huawei OIDS
const (
	HuaPrefix                       = "1.3.6.1.4.1.2011"
	HuaIfEtherStatInCRCPkts  string = ".1.3.6.1.4.1.2011.5.25.41.1.6.1.1.12"
	HuaIfEtherStatInCRCPktsF        = ".1.3.6.1.4.1.2011.5.25.41.1.6.1.1.12.%d"

	//Pause
	HuaIfEtherStatInPausePkts  = ".1.3.6.1.4.1.2011.5.25.41.1.6.1.1.18"
	HuaIfEtherStatOutPausePkts = "	1.3.6.1.4.1.2011.5.25.41.1.6.1.1.22"

	// Pause with format
	HuaIfEtherStatInPausePktsF  = ".1.3.6.1.4.1.2011.5.25.41.1.6.1.1.18.%d"
	HuaIfEtherStatOutPausePktsF = ".1.3.6.1.4.1.2011.5.25.41.1.6.1.1.22.%d"

	// Resets
	HuaIfEthIfStatReset  = ".1.3.6.1.4.1.2011.5.25.41.1.6.1.1.23"
	HuaIfEthIfStatResetF = ".1.3.6.1.4.1.2011.5.25.41.1.6.1.1.23.%d"

	// VRP
	HuaIfVRPF                   = ".1.3.6.1.4.1.2011.5.25.31.1.1.3.1.%d"
	HuaIfVRPOpticalVendorSNF    = ".1.3.6.1.4.1.2011.5.25.31.1.1.3.1.4.%d"
	HuaIfVRPOpticalTemperatureF = ".1.3.6.1.4.1.2011.5.25.31.1.1.3.1.5.%d"
	HuaIfVRPOpticalVoltageF     = ".1.3.6.1.4.1.2011.5.25.31.1.1.3.1.6.%d"
	HuaIfVRPOpticalBiasF        = ".1.3.6.1.4.1.2011.5.25.31.1.1.3.1.7.%d"
	HuaIfVRPOpticalRxPowerF     = ".1.3.6.1.4.1.2011.5.25.31.1.1.3.1.8.%d"
	HuaIfVRPOpticalTxPowerF     = ".1.3.6.1.4.1.2011.5.25.31.1.1.3.1.9.%d"
	HuaIfVRPVendorPNF           = ".1.3.6.1.4.1.2011.5.25.31.1.1.3.1.25.%d"

	HuaIfVRPOpticalVendorSN    = ".1.3.6.1.4.1.2011.5.25.31.1.1.3.1.4"
	HuaIfVRPOpticalTemperature = ".1.3.6.1.4.1.2011.5.25.31.1.1.3.1.5"
	HuaIfVRPOpticalVoltage     = ".1.3.6.1.4.1.2011.5.25.31.1.1.3.1.6"
	HuaIfVRPOpticalBias        = ".1.3.6.1.4.1.2011.5.25.31.1.1.3.1.7"
	HuaIfVRPOpticalRxPower     = ".1.3.6.1.4.1.2011.5.25.31.1.1.3.1.8"
	HuaIfVRPOpticalTxPower     = ".1.3.6.1.4.1.2011.5.25.31.1.1.3.1.9"
	HuaIfVRPVendorPN           = ".1.3.6.1.4.1.2011.5.25.31.1.1.3.1.25"
)
