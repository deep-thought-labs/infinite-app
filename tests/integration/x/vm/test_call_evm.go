package vm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/evm/contracts"
	testconstants "github.com/cosmos/evm/testutil/constants"
	utiltx "github.com/cosmos/evm/testutil/tx"
	"github.com/cosmos/evm/x/erc20/types"
	"github.com/cosmos/evm/x/vm/statedb"
	evmtypes "github.com/cosmos/evm/x/vm/types"
)

func (s *KeeperTestSuite) TestCallEVM() {
	wcosmosEVMContract := common.HexToAddress(testconstants.WEVMOSContractMainnet)
	testCases := []struct {
		name     string
		method   string
		stateDB  *statedb.StateDB
		expPass  bool
		expError string
	}{
		{
			"unknown method",
			"",
			nil,
			false,
			"",
		},
		{
			"pass",
			"balanceOf",
			nil,
			true,
			"",
		},
		{
			"fail with nil statedb",
			"balanceOf",
			nil,
			false,
			"stateDB cannot be nil",
		},
	}
	for _, tc := range testCases {
		s.SetupTest() // reset

		erc20 := contracts.ERC20MinterBurnerDecimalsContract.ABI
		account := utiltx.GenerateAddress()

		var stateDB *statedb.StateDB
		if tc.stateDB == nil && tc.name != "fail with nil statedb" {
			stateDB = statedb.New(s.Network.GetContext(), s.Network.App.GetEVMKeeper(), statedb.NewEmptyTxConfig())
		} else {
			stateDB = tc.stateDB
		}

		res, err := s.Network.App.GetEVMKeeper().CallEVM(s.Network.GetContext(), stateDB, erc20, types.ModuleAddress, wcosmosEVMContract, false, false, nil, tc.method, account)
		if tc.expPass {
			s.Require().IsTypef(&evmtypes.MsgEthereumTxResponse{}, res, tc.name)
			s.Require().NoError(err)
		} else {
			s.Require().Error(err)
			if tc.expError != "" {
				s.Require().Contains(err.Error(), tc.expError)
			}
		}
	}
}

func (s *KeeperTestSuite) TestCallEVMWithData() {
	erc20 := contracts.ERC20MinterBurnerDecimalsContract.ABI
	wcosmosEVMContract := common.HexToAddress(testconstants.WEVMOSContractMainnet)
	testCases := []struct {
		name     string
		from     common.Address
		malleate func() []byte
		deploy   bool
<<<<<<< HEAD
<<<<<<< HEAD
		useNilDB bool
=======
		gasCap   *big.Int
>>>>>>> 0ca6371a (fix(vm): respect provided gasCap in CallEVMWithData (#992))
=======
		useNilDB bool
>>>>>>> 5b9f95e1 (chore: port v0.6.x onto krakatoa (#1053))
		expPass  bool
		expError string
	}{
		{
			name: "pass with unknown method",
			from: types.ModuleAddress,
			malleate: func() []byte {
				account := utiltx.GenerateAddress()
				data, _ := erc20.Pack("", account)
				return data
			},
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 5b9f95e1 (chore: port v0.6.x onto krakatoa (#1053))
			deploy:   false,
			useNilDB: false,
			expPass:  true,
			expError: "",
<<<<<<< HEAD
=======
			false,
			nil,
			true,
>>>>>>> 0ca6371a (fix(vm): respect provided gasCap in CallEVMWithData (#992))
=======
>>>>>>> 5b9f95e1 (chore: port v0.6.x onto krakatoa (#1053))
		},
		{
			name: "pass",
			from: types.ModuleAddress,
			malleate: func() []byte {
				account := utiltx.GenerateAddress()
				data, _ := erc20.Pack("balanceOf", account)
				return data
			},
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 5b9f95e1 (chore: port v0.6.x onto krakatoa (#1053))
			deploy:   false,
			useNilDB: false,
			expPass:  true,
			expError: "",
<<<<<<< HEAD
		},
		{
			name: "pass with empty data",
			from: types.ModuleAddress,
			malleate: func() []byte {
				return []byte{}
			},
			deploy:   false,
			useNilDB: false,
			expPass:  true,
			expError: "",
=======
			false,
			nil,
			true,
=======
>>>>>>> 5b9f95e1 (chore: port v0.6.x onto krakatoa (#1053))
		},
		{
			name: "pass with empty data",
			from: types.ModuleAddress,
			malleate: func() []byte {
				return []byte{}
			},
			deploy:   false,
			useNilDB: false,
			expPass:  true,
			expError: "",
		},
		{
			name: "fail empty sender",
			from: common.Address{},
			malleate: func() []byte {
				return []byte{}
			},
			deploy:   false,
			useNilDB: false,
			expPass:  false,
			expError: "",
		},
		{
			name: "fail with nil statedb",
			from: types.ModuleAddress,
			malleate: func() []byte {
				account := utiltx.GenerateAddress()
				data, _ := erc20.Pack("balanceOf", account)
				return data
			},
			deploy:   false,
			useNilDB: true,
			expPass:  false,
			expError: "stateDB cannot be nil",
		},
		{
<<<<<<< HEAD
			"pass with empty data",
			types.ModuleAddress,
			func() []byte {
				return []byte{}
			},
			false,
			nil,
			true,
>>>>>>> 0ca6371a (fix(vm): respect provided gasCap in CallEVMWithData (#992))
		},
		{
			name: "fail empty sender",
			from: common.Address{},
			malleate: func() []byte {
				return []byte{}
			},
<<<<<<< HEAD
			deploy:   false,
			useNilDB: false,
			expPass:  false,
			expError: "",
=======
			false,
			nil,
			false,
>>>>>>> 0ca6371a (fix(vm): respect provided gasCap in CallEVMWithData (#992))
		},
		{
			name: "fail with nil statedb",
			from: types.ModuleAddress,
			malleate: func() []byte {
				account := utiltx.GenerateAddress()
				data, _ := erc20.Pack("balanceOf", account)
				return data
			},
			deploy:   false,
			useNilDB: true,
			expPass:  false,
			expError: "stateDB cannot be nil",
		},
		{
=======
>>>>>>> 5b9f95e1 (chore: port v0.6.x onto krakatoa (#1053))
			name: "deploy",
			from: types.ModuleAddress,
			malleate: func() []byte {
				ctorArgs, _ := contracts.ERC20MinterBurnerDecimalsContract.ABI.Pack("", "test", "test", uint8(18))
				data := append(contracts.ERC20MinterBurnerDecimalsContract.Bin, ctorArgs...) //nolint:gocritic
				return data
			},
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 5b9f95e1 (chore: port v0.6.x onto krakatoa (#1053))
			deploy:   true,
			useNilDB: false,
			expPass:  true,
			expError: "",
<<<<<<< HEAD
=======
			true,
			nil,
			true,
>>>>>>> 0ca6371a (fix(vm): respect provided gasCap in CallEVMWithData (#992))
=======
>>>>>>> 5b9f95e1 (chore: port v0.6.x onto krakatoa (#1053))
		},
		{
			name: "fail deploy",
			from: types.ModuleAddress,
			malleate: func() []byte {
				params := s.Network.App.GetEVMKeeper().GetParams(s.Network.GetContext())
				params.AccessControl.Create = evmtypes.AccessControlType{
					AccessType: evmtypes.AccessTypeRestricted,
				}
				_ = s.Network.App.GetEVMKeeper().SetParams(s.Network.GetContext(), params)
				ctorArgs, _ := contracts.ERC20MinterBurnerDecimalsContract.ABI.Pack("", "test", "test", uint8(18))
				data := append(contracts.ERC20MinterBurnerDecimalsContract.Bin, ctorArgs...) //nolint:gocritic
				return data
			},
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 5b9f95e1 (chore: port v0.6.x onto krakatoa (#1053))
			deploy:   true,
			useNilDB: false,
			expPass:  false,
			expError: "",
		},
		{
			name: "fail deploy with nil statedb",
			from: types.ModuleAddress,
			malleate: func() []byte {
				ctorArgs, _ := contracts.ERC20MinterBurnerDecimalsContract.ABI.Pack("", "test", "test", uint8(18))
				data := append(contracts.ERC20MinterBurnerDecimalsContract.Bin, ctorArgs...) //nolint:gocritic
				return data
			},
			deploy:   true,
			useNilDB: true,
			expPass:  false,
			expError: "stateDB cannot be nil",
<<<<<<< HEAD
=======
			true,
			nil,
			false,
>>>>>>> 0ca6371a (fix(vm): respect provided gasCap in CallEVMWithData (#992))
=======
>>>>>>> 5b9f95e1 (chore: port v0.6.x onto krakatoa (#1053))
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			s.SetupTest() // reset

			data := tc.malleate()
			var res *evmtypes.MsgEthereumTxResponse
			var err error

			var stateDB *statedb.StateDB
			if !tc.useNilDB {
				stateDB = statedb.New(s.Network.GetContext(), s.Network.App.GetEVMKeeper(), statedb.NewEmptyTxConfig())
			}

			if tc.deploy {
<<<<<<< HEAD
<<<<<<< HEAD
				res, err = s.Network.App.GetEVMKeeper().CallEVMWithData(s.Network.GetContext(), stateDB, tc.from, nil, data, true, false, nil)
			} else {
				res, err = s.Network.App.GetEVMKeeper().CallEVMWithData(s.Network.GetContext(), stateDB, tc.from, &wcosmosEVMContract, data, false, false, nil)
=======
				res, err = s.Network.App.GetEVMKeeper().CallEVMWithData(s.Network.GetContext(), tc.from, nil, data, true, tc.gasCap)
			} else {
				res, err = s.Network.App.GetEVMKeeper().CallEVMWithData(s.Network.GetContext(), tc.from, &wcosmosEVMContract, data, false, tc.gasCap)
>>>>>>> 0ca6371a (fix(vm): respect provided gasCap in CallEVMWithData (#992))
=======
				res, err = s.Network.App.GetEVMKeeper().CallEVMWithData(s.Network.GetContext(), stateDB, tc.from, nil, data, true, false, nil)
			} else {
				res, err = s.Network.App.GetEVMKeeper().CallEVMWithData(s.Network.GetContext(), stateDB, tc.from, &wcosmosEVMContract, data, false, false, nil)
>>>>>>> 5b9f95e1 (chore: port v0.6.x onto krakatoa (#1053))
			}

			if tc.expPass {
				s.Require().IsTypef(&evmtypes.MsgEthereumTxResponse{}, res, tc.name)
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
				if tc.expError != "" {
					s.Require().Contains(err.Error(), tc.expError)
				}
			}
		})
	}
}
