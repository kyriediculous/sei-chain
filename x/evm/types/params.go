package types

import (
	"encoding/hex"
	"errors"
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sei-protocol/sei-chain/x/evm/artifacts/native"
	"gopkg.in/yaml.v2"
)

var (
	KeyBaseDenom                              = []byte("KeyBaseDenom")
	KeyPriorityNormalizer                     = []byte("KeyPriorityNormalizer")
	KeyBaseFeePerGas                          = []byte("KeyBaseFeePerGas")
	KeyMinFeePerGas                           = []byte("KeyMinFeePerGas")
	KeyChainConfig                            = []byte("KeyChainConfig")
	KeyChainID                                = []byte("KeyChainID")
	KeyWhitelistedCodeHashesBankSend          = []byte("KeyWhitelistedCodeHashesBankSend")
	KeyWhitelistedCwCodeHashesForDelegateCall = []byte("KeyWhitelistedCwCodeHashesForDelegateCall")
)

const (
	DefaultBaseDenom = "usei"
)

var DefaultPriorityNormalizer = sdk.NewDec(1)

// DefaultBaseFeePerGas determines how much usei per gas spent is
// burnt rather than go to validators (similar to base fee on
// Ethereum).
var DefaultBaseFeePerGas = sdk.NewDec(0)
var DefaultMinFeePerGas = sdk.NewDec(1000000000)
var DefaultChainID = sdk.NewInt(713715)
var DefaultWhitelistedCodeHashesBankSend = generateDefaultWhitelistedCodeHashesBankSend()

var DefaultWhitelistedCwCodeHashesForDelegateCall = generateDefaultWhitelistedCwCodeHashesForDelegateCall()

var _ paramtypes.ParamSet = (*Params)(nil)

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func DefaultParams() Params {
	return Params{
		BaseDenom:                              DefaultBaseDenom,
		PriorityNormalizer:                     DefaultPriorityNormalizer,
		BaseFeePerGas:                          DefaultBaseFeePerGas,
		MinimumFeePerGas:                       DefaultMinFeePerGas,
		ChainConfig:                            DefaultChainConfig(),
		ChainId:                                DefaultChainID,
		WhitelistedCodehashesBankSend:          DefaultWhitelistedCodeHashesBankSend,
		WhitelistedCwCodeHashesForDelegateCall: DefaultWhitelistedCwCodeHashesForDelegateCall,
	}
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyBaseDenom, &p.BaseDenom, validateBaseDenom),
		paramtypes.NewParamSetPair(KeyPriorityNormalizer, &p.PriorityNormalizer, validatePriorityNormalizer),
		paramtypes.NewParamSetPair(KeyBaseFeePerGas, &p.BaseFeePerGas, validateBaseFeePerGas),
		paramtypes.NewParamSetPair(KeyMinFeePerGas, &p.MinimumFeePerGas, validateMinFeePerGas),
		paramtypes.NewParamSetPair(KeyChainConfig, &p.ChainConfig, validateChainConfig),
		paramtypes.NewParamSetPair(KeyChainID, &p.ChainId, validateChainID),
		paramtypes.NewParamSetPair(KeyWhitelistedCodeHashesBankSend, &p.WhitelistedCodehashesBankSend, validateWhitelistedCodeHashesBankSend),
		paramtypes.NewParamSetPair(KeyWhitelistedCwCodeHashesForDelegateCall, &p.WhitelistedCwCodeHashesForDelegateCall, validateWhitelistedCwHashesForDelegateCall),
	}
}

func (p Params) Validate() error {
	if err := validateBaseDenom(p.BaseDenom); err != nil {
		return err
	}
	if err := validatePriorityNormalizer(p.PriorityNormalizer); err != nil {
		return err
	}
	if err := validateBaseFeePerGas(p.BaseFeePerGas); err != nil {
		return err
	}
	if err := validateMinFeePerGas(p.MinimumFeePerGas); err != nil {
		return err
	}
	if p.MinimumFeePerGas.LT(p.BaseFeePerGas) {
		return errors.New("minimum fee cannot be lower than base fee")
	}
	if err := validateChainID(p.ChainId); err != nil {
		return err
	}
	if err := validateChainConfig(p.ChainConfig); err != nil {
		return err
	}
	if err := validateWhitelistedCodeHashesBankSend(p.WhitelistedCodehashesBankSend); err != nil {
		return err
	}
	return validateWhitelistedCwHashesForDelegateCall(p.WhitelistedCwCodeHashesForDelegateCall)
}

func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateBaseDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == "" {
		return errors.New("empty base denom")
	}

	return nil
}

func validatePriorityNormalizer(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !v.IsPositive() {
		return fmt.Errorf("nonpositive priority normalizer: %d", v)
	}

	return nil
}

func validateBaseFeePerGas(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("negative base fee per gas: %d", v)
	}

	return nil
}

func validateMinFeePerGas(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("negative min fee per gas: %d", v)
	}

	return nil
}

func validateChainConfig(i interface{}) error {
	v, ok := i.(ChainConfig)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return v.Validate()
}

func validateChainID(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("negative min fee per gas: %d", v)
	}

	return nil
}

func validateWhitelistedCodeHashesBankSend(i interface{}) error {
	_, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateWhitelistedCwHashesForDelegateCall(i interface{}) error {
	_, ok := i.([][]byte)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func generateDefaultWhitelistedCodeHashesBankSend() (res []string) {
	h := crypto.Keccak256Hash(native.GetBin())
	res = append(res, h.Hex())
	return
}

func generateDefaultWhitelistedCwCodeHashesForDelegateCall() [][]byte {
	cw20, _ := hex.DecodeString("A25D78D7ACD2EE47CC39C224E162FE79B53E6BBE6ED2A56E8C0A86593EBE6102")
	cw721, _ := hex.DecodeString("94CDD9C3E85C26F7CEC43C23BFB4B3B2B2D71A0A8D85C58DF12FFEC0741FEBC8")
	return [][]byte{cw20, cw721}
}
