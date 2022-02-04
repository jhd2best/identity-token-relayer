// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ownership_validator

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// OwnershipValidatorMetaData contains all meta data concerning the OwnershipValidator contract.
var OwnershipValidatorMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"addrMap\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"ethAddress\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"ethAddress\",\"type\":\"string\"},{\"internalType\":\"address[]\",\"name\":\"owners\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"tokenIds\",\"type\":\"uint256[]\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"ethAddress\",\"type\":\"string\"}],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"ethAddress\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"ethAddress\",\"type\":\"string\"},{\"internalType\":\"contractIToken721\",\"name\":\"oneAddress\",\"type\":\"address\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"ethAddress\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"baseURI_\",\"type\":\"string\"}],\"name\":\"setBaseURI\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"ethAddress\",\"type\":\"string\"}],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"ethAddress\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"ethAddress\",\"type\":\"string\"}],\"name\":\"unregister\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"ethAddress\",\"type\":\"string\"},{\"internalType\":\"address[]\",\"name\":\"owners\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"tokenIds\",\"type\":\"uint256[]\"}],\"name\":\"updateOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// OwnershipValidatorABI is the input ABI used to generate the binding from.
// Deprecated: Use OwnershipValidatorMetaData.ABI instead.
var OwnershipValidatorABI = OwnershipValidatorMetaData.ABI

// OwnershipValidator is an auto generated Go binding around an Ethereum contract.
type OwnershipValidator struct {
	OwnershipValidatorCaller     // Read-only binding to the contract
	OwnershipValidatorTransactor // Write-only binding to the contract
	OwnershipValidatorFilterer   // Log filterer for contract events
}

// OwnershipValidatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type OwnershipValidatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnershipValidatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OwnershipValidatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnershipValidatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OwnershipValidatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnershipValidatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OwnershipValidatorSession struct {
	Contract     *OwnershipValidator // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// OwnershipValidatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OwnershipValidatorCallerSession struct {
	Contract *OwnershipValidatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// OwnershipValidatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OwnershipValidatorTransactorSession struct {
	Contract     *OwnershipValidatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// OwnershipValidatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type OwnershipValidatorRaw struct {
	Contract *OwnershipValidator // Generic contract binding to access the raw methods on
}

// OwnershipValidatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OwnershipValidatorCallerRaw struct {
	Contract *OwnershipValidatorCaller // Generic read-only contract binding to access the raw methods on
}

// OwnershipValidatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OwnershipValidatorTransactorRaw struct {
	Contract *OwnershipValidatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwnershipValidator creates a new instance of OwnershipValidator, bound to a specific deployed contract.
func NewOwnershipValidator(address common.Address, backend bind.ContractBackend) (*OwnershipValidator, error) {
	contract, err := bindOwnershipValidator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OwnershipValidator{OwnershipValidatorCaller: OwnershipValidatorCaller{contract: contract}, OwnershipValidatorTransactor: OwnershipValidatorTransactor{contract: contract}, OwnershipValidatorFilterer: OwnershipValidatorFilterer{contract: contract}}, nil
}

// NewOwnershipValidatorCaller creates a new read-only instance of OwnershipValidator, bound to a specific deployed contract.
func NewOwnershipValidatorCaller(address common.Address, caller bind.ContractCaller) (*OwnershipValidatorCaller, error) {
	contract, err := bindOwnershipValidator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnershipValidatorCaller{contract: contract}, nil
}

// NewOwnershipValidatorTransactor creates a new write-only instance of OwnershipValidator, bound to a specific deployed contract.
func NewOwnershipValidatorTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnershipValidatorTransactor, error) {
	contract, err := bindOwnershipValidator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnershipValidatorTransactor{contract: contract}, nil
}

// NewOwnershipValidatorFilterer creates a new log filterer instance of OwnershipValidator, bound to a specific deployed contract.
func NewOwnershipValidatorFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnershipValidatorFilterer, error) {
	contract, err := bindOwnershipValidator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnershipValidatorFilterer{contract: contract}, nil
}

// bindOwnershipValidator binds a generic wrapper to an already deployed contract.
func bindOwnershipValidator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnershipValidatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnershipValidator *OwnershipValidatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OwnershipValidator.Contract.OwnershipValidatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnershipValidator *OwnershipValidatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnershipValidator.Contract.OwnershipValidatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnershipValidator *OwnershipValidatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnershipValidator.Contract.OwnershipValidatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnershipValidator *OwnershipValidatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OwnershipValidator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnershipValidator *OwnershipValidatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnershipValidator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnershipValidator *OwnershipValidatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnershipValidator.Contract.contract.Transact(opts, method, params...)
}

// AddrMap is a free data retrieval call binding the contract method 0xac87478a.
//
// Solidity: function addrMap(string ) view returns(address)
func (_OwnershipValidator *OwnershipValidatorCaller) AddrMap(opts *bind.CallOpts, arg0 string) (common.Address, error) {
	var out []interface{}
	err := _OwnershipValidator.contract.Call(opts, &out, "addrMap", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AddrMap is a free data retrieval call binding the contract method 0xac87478a.
//
// Solidity: function addrMap(string ) view returns(address)
func (_OwnershipValidator *OwnershipValidatorSession) AddrMap(arg0 string) (common.Address, error) {
	return _OwnershipValidator.Contract.AddrMap(&_OwnershipValidator.CallOpts, arg0)
}

// AddrMap is a free data retrieval call binding the contract method 0xac87478a.
//
// Solidity: function addrMap(string ) view returns(address)
func (_OwnershipValidator *OwnershipValidatorCallerSession) AddrMap(arg0 string) (common.Address, error) {
	return _OwnershipValidator.Contract.AddrMap(&_OwnershipValidator.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x04031852.
//
// Solidity: function balanceOf(string ethAddress, address owner) view returns(uint256)
func (_OwnershipValidator *OwnershipValidatorCaller) BalanceOf(opts *bind.CallOpts, ethAddress string, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _OwnershipValidator.contract.Call(opts, &out, "balanceOf", ethAddress, owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x04031852.
//
// Solidity: function balanceOf(string ethAddress, address owner) view returns(uint256)
func (_OwnershipValidator *OwnershipValidatorSession) BalanceOf(ethAddress string, owner common.Address) (*big.Int, error) {
	return _OwnershipValidator.Contract.BalanceOf(&_OwnershipValidator.CallOpts, ethAddress, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x04031852.
//
// Solidity: function balanceOf(string ethAddress, address owner) view returns(uint256)
func (_OwnershipValidator *OwnershipValidatorCallerSession) BalanceOf(ethAddress string, owner common.Address) (*big.Int, error) {
	return _OwnershipValidator.Contract.BalanceOf(&_OwnershipValidator.CallOpts, ethAddress, owner)
}

// Name is a free data retrieval call binding the contract method 0x5b43bc99.
//
// Solidity: function name(string ethAddress) view returns(string)
func (_OwnershipValidator *OwnershipValidatorCaller) Name(opts *bind.CallOpts, ethAddress string) (string, error) {
	var out []interface{}
	err := _OwnershipValidator.contract.Call(opts, &out, "name", ethAddress)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x5b43bc99.
//
// Solidity: function name(string ethAddress) view returns(string)
func (_OwnershipValidator *OwnershipValidatorSession) Name(ethAddress string) (string, error) {
	return _OwnershipValidator.Contract.Name(&_OwnershipValidator.CallOpts, ethAddress)
}

// Name is a free data retrieval call binding the contract method 0x5b43bc99.
//
// Solidity: function name(string ethAddress) view returns(string)
func (_OwnershipValidator *OwnershipValidatorCallerSession) Name(ethAddress string) (string, error) {
	return _OwnershipValidator.Contract.Name(&_OwnershipValidator.CallOpts, ethAddress)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OwnershipValidator *OwnershipValidatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OwnershipValidator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OwnershipValidator *OwnershipValidatorSession) Owner() (common.Address, error) {
	return _OwnershipValidator.Contract.Owner(&_OwnershipValidator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OwnershipValidator *OwnershipValidatorCallerSession) Owner() (common.Address, error) {
	return _OwnershipValidator.Contract.Owner(&_OwnershipValidator.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x915622c1.
//
// Solidity: function ownerOf(string ethAddress, uint256 tokenId) view returns(address)
func (_OwnershipValidator *OwnershipValidatorCaller) OwnerOf(opts *bind.CallOpts, ethAddress string, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _OwnershipValidator.contract.Call(opts, &out, "ownerOf", ethAddress, tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x915622c1.
//
// Solidity: function ownerOf(string ethAddress, uint256 tokenId) view returns(address)
func (_OwnershipValidator *OwnershipValidatorSession) OwnerOf(ethAddress string, tokenId *big.Int) (common.Address, error) {
	return _OwnershipValidator.Contract.OwnerOf(&_OwnershipValidator.CallOpts, ethAddress, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x915622c1.
//
// Solidity: function ownerOf(string ethAddress, uint256 tokenId) view returns(address)
func (_OwnershipValidator *OwnershipValidatorCallerSession) OwnerOf(ethAddress string, tokenId *big.Int) (common.Address, error) {
	return _OwnershipValidator.Contract.OwnerOf(&_OwnershipValidator.CallOpts, ethAddress, tokenId)
}

// Symbol is a free data retrieval call binding the contract method 0x41bb0559.
//
// Solidity: function symbol(string ethAddress) view returns(string)
func (_OwnershipValidator *OwnershipValidatorCaller) Symbol(opts *bind.CallOpts, ethAddress string) (string, error) {
	var out []interface{}
	err := _OwnershipValidator.contract.Call(opts, &out, "symbol", ethAddress)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x41bb0559.
//
// Solidity: function symbol(string ethAddress) view returns(string)
func (_OwnershipValidator *OwnershipValidatorSession) Symbol(ethAddress string) (string, error) {
	return _OwnershipValidator.Contract.Symbol(&_OwnershipValidator.CallOpts, ethAddress)
}

// Symbol is a free data retrieval call binding the contract method 0x41bb0559.
//
// Solidity: function symbol(string ethAddress) view returns(string)
func (_OwnershipValidator *OwnershipValidatorCallerSession) Symbol(ethAddress string) (string, error) {
	return _OwnershipValidator.Contract.Symbol(&_OwnershipValidator.CallOpts, ethAddress)
}

// TokenURI is a free data retrieval call binding the contract method 0x9b330cb9.
//
// Solidity: function tokenURI(string ethAddress, uint256 tokenId) view returns(string)
func (_OwnershipValidator *OwnershipValidatorCaller) TokenURI(opts *bind.CallOpts, ethAddress string, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _OwnershipValidator.contract.Call(opts, &out, "tokenURI", ethAddress, tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0x9b330cb9.
//
// Solidity: function tokenURI(string ethAddress, uint256 tokenId) view returns(string)
func (_OwnershipValidator *OwnershipValidatorSession) TokenURI(ethAddress string, tokenId *big.Int) (string, error) {
	return _OwnershipValidator.Contract.TokenURI(&_OwnershipValidator.CallOpts, ethAddress, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0x9b330cb9.
//
// Solidity: function tokenURI(string ethAddress, uint256 tokenId) view returns(string)
func (_OwnershipValidator *OwnershipValidatorCallerSession) TokenURI(ethAddress string, tokenId *big.Int) (string, error) {
	return _OwnershipValidator.Contract.TokenURI(&_OwnershipValidator.CallOpts, ethAddress, tokenId)
}

// Initialize is a paid mutator transaction binding the contract method 0x24f3ae78.
//
// Solidity: function initialize(string ethAddress, address[] owners, uint256[] tokenIds) returns()
func (_OwnershipValidator *OwnershipValidatorTransactor) Initialize(opts *bind.TransactOpts, ethAddress string, owners []common.Address, tokenIds []*big.Int) (*types.Transaction, error) {
	return _OwnershipValidator.contract.Transact(opts, "initialize", ethAddress, owners, tokenIds)
}

// Initialize is a paid mutator transaction binding the contract method 0x24f3ae78.
//
// Solidity: function initialize(string ethAddress, address[] owners, uint256[] tokenIds) returns()
func (_OwnershipValidator *OwnershipValidatorSession) Initialize(ethAddress string, owners []common.Address, tokenIds []*big.Int) (*types.Transaction, error) {
	return _OwnershipValidator.Contract.Initialize(&_OwnershipValidator.TransactOpts, ethAddress, owners, tokenIds)
}

// Initialize is a paid mutator transaction binding the contract method 0x24f3ae78.
//
// Solidity: function initialize(string ethAddress, address[] owners, uint256[] tokenIds) returns()
func (_OwnershipValidator *OwnershipValidatorTransactorSession) Initialize(ethAddress string, owners []common.Address, tokenIds []*big.Int) (*types.Transaction, error) {
	return _OwnershipValidator.Contract.Initialize(&_OwnershipValidator.TransactOpts, ethAddress, owners, tokenIds)
}

// Register is a paid mutator transaction binding the contract method 0x1e59c529.
//
// Solidity: function register(string ethAddress, address oneAddress) returns()
func (_OwnershipValidator *OwnershipValidatorTransactor) Register(opts *bind.TransactOpts, ethAddress string, oneAddress common.Address) (*types.Transaction, error) {
	return _OwnershipValidator.contract.Transact(opts, "register", ethAddress, oneAddress)
}

// Register is a paid mutator transaction binding the contract method 0x1e59c529.
//
// Solidity: function register(string ethAddress, address oneAddress) returns()
func (_OwnershipValidator *OwnershipValidatorSession) Register(ethAddress string, oneAddress common.Address) (*types.Transaction, error) {
	return _OwnershipValidator.Contract.Register(&_OwnershipValidator.TransactOpts, ethAddress, oneAddress)
}

// Register is a paid mutator transaction binding the contract method 0x1e59c529.
//
// Solidity: function register(string ethAddress, address oneAddress) returns()
func (_OwnershipValidator *OwnershipValidatorTransactorSession) Register(ethAddress string, oneAddress common.Address) (*types.Transaction, error) {
	return _OwnershipValidator.Contract.Register(&_OwnershipValidator.TransactOpts, ethAddress, oneAddress)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_OwnershipValidator *OwnershipValidatorTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnershipValidator.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_OwnershipValidator *OwnershipValidatorSession) RenounceOwnership() (*types.Transaction, error) {
	return _OwnershipValidator.Contract.RenounceOwnership(&_OwnershipValidator.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_OwnershipValidator *OwnershipValidatorTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _OwnershipValidator.Contract.RenounceOwnership(&_OwnershipValidator.TransactOpts)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x6790a9de.
//
// Solidity: function setBaseURI(string ethAddress, string baseURI_) returns()
func (_OwnershipValidator *OwnershipValidatorTransactor) SetBaseURI(opts *bind.TransactOpts, ethAddress string, baseURI_ string) (*types.Transaction, error) {
	return _OwnershipValidator.contract.Transact(opts, "setBaseURI", ethAddress, baseURI_)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x6790a9de.
//
// Solidity: function setBaseURI(string ethAddress, string baseURI_) returns()
func (_OwnershipValidator *OwnershipValidatorSession) SetBaseURI(ethAddress string, baseURI_ string) (*types.Transaction, error) {
	return _OwnershipValidator.Contract.SetBaseURI(&_OwnershipValidator.TransactOpts, ethAddress, baseURI_)
}

// SetBaseURI is a paid mutator transaction binding the contract method 0x6790a9de.
//
// Solidity: function setBaseURI(string ethAddress, string baseURI_) returns()
func (_OwnershipValidator *OwnershipValidatorTransactorSession) SetBaseURI(ethAddress string, baseURI_ string) (*types.Transaction, error) {
	return _OwnershipValidator.Contract.SetBaseURI(&_OwnershipValidator.TransactOpts, ethAddress, baseURI_)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_OwnershipValidator *OwnershipValidatorTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _OwnershipValidator.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_OwnershipValidator *OwnershipValidatorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _OwnershipValidator.Contract.TransferOwnership(&_OwnershipValidator.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_OwnershipValidator *OwnershipValidatorTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _OwnershipValidator.Contract.TransferOwnership(&_OwnershipValidator.TransactOpts, newOwner)
}

// Unregister is a paid mutator transaction binding the contract method 0x6598a1ae.
//
// Solidity: function unregister(string ethAddress) returns()
func (_OwnershipValidator *OwnershipValidatorTransactor) Unregister(opts *bind.TransactOpts, ethAddress string) (*types.Transaction, error) {
	return _OwnershipValidator.contract.Transact(opts, "unregister", ethAddress)
}

// Unregister is a paid mutator transaction binding the contract method 0x6598a1ae.
//
// Solidity: function unregister(string ethAddress) returns()
func (_OwnershipValidator *OwnershipValidatorSession) Unregister(ethAddress string) (*types.Transaction, error) {
	return _OwnershipValidator.Contract.Unregister(&_OwnershipValidator.TransactOpts, ethAddress)
}

// Unregister is a paid mutator transaction binding the contract method 0x6598a1ae.
//
// Solidity: function unregister(string ethAddress) returns()
func (_OwnershipValidator *OwnershipValidatorTransactorSession) Unregister(ethAddress string) (*types.Transaction, error) {
	return _OwnershipValidator.Contract.Unregister(&_OwnershipValidator.TransactOpts, ethAddress)
}

// UpdateOwnership is a paid mutator transaction binding the contract method 0x7ff35e0c.
//
// Solidity: function updateOwnership(string ethAddress, address[] owners, uint256[] tokenIds) returns()
func (_OwnershipValidator *OwnershipValidatorTransactor) UpdateOwnership(opts *bind.TransactOpts, ethAddress string, owners []common.Address, tokenIds []*big.Int) (*types.Transaction, error) {
	return _OwnershipValidator.contract.Transact(opts, "updateOwnership", ethAddress, owners, tokenIds)
}

// UpdateOwnership is a paid mutator transaction binding the contract method 0x7ff35e0c.
//
// Solidity: function updateOwnership(string ethAddress, address[] owners, uint256[] tokenIds) returns()
func (_OwnershipValidator *OwnershipValidatorSession) UpdateOwnership(ethAddress string, owners []common.Address, tokenIds []*big.Int) (*types.Transaction, error) {
	return _OwnershipValidator.Contract.UpdateOwnership(&_OwnershipValidator.TransactOpts, ethAddress, owners, tokenIds)
}

// UpdateOwnership is a paid mutator transaction binding the contract method 0x7ff35e0c.
//
// Solidity: function updateOwnership(string ethAddress, address[] owners, uint256[] tokenIds) returns()
func (_OwnershipValidator *OwnershipValidatorTransactorSession) UpdateOwnership(ethAddress string, owners []common.Address, tokenIds []*big.Int) (*types.Transaction, error) {
	return _OwnershipValidator.Contract.UpdateOwnership(&_OwnershipValidator.TransactOpts, ethAddress, owners, tokenIds)
}

// OwnershipValidatorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the OwnershipValidator contract.
type OwnershipValidatorOwnershipTransferredIterator struct {
	Event *OwnershipValidatorOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OwnershipValidatorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnershipValidatorOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OwnershipValidatorOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OwnershipValidatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnershipValidatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnershipValidatorOwnershipTransferred represents a OwnershipTransferred event raised by the OwnershipValidator contract.
type OwnershipValidatorOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_OwnershipValidator *OwnershipValidatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*OwnershipValidatorOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _OwnershipValidator.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OwnershipValidatorOwnershipTransferredIterator{contract: _OwnershipValidator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_OwnershipValidator *OwnershipValidatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OwnershipValidatorOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _OwnershipValidator.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnershipValidatorOwnershipTransferred)
				if err := _OwnershipValidator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_OwnershipValidator *OwnershipValidatorFilterer) ParseOwnershipTransferred(log types.Log) (*OwnershipValidatorOwnershipTransferred, error) {
	event := new(OwnershipValidatorOwnershipTransferred)
	if err := _OwnershipValidator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
