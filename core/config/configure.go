package config

import (
	"fmt"

	"github.com/anoideaopen/foundation/proto"
)

// ConfigMapper defines a method for mapping Init arguments to a Config instance.
type ConfigMapper interface {
	// MapConfig maps the provided arguments to a proto.Config instance.
	MapConfig(args []string) (*proto.Config, error)
}

// ConfigMapperFunc is a function type that implements the ConfigMapper interface.
type ConfigMapperFunc func(args []string) (*proto.Config, error)

// MapConfig calls the underlying function to map the provided arguments to a proto.Config instance.
func (c ConfigMapperFunc) MapConfig(args []string) (*proto.Config, error) {
	return c(args)
}

// Configurator defines methods for validating, applying, and retrieving contract configuration.
type Configurator interface {
	// ValidateConfig validates the provided contract configuration data.
	ValidateConfig(config []byte) error

	// ApplyContractConfig applies the provided contract configuration.
	ApplyContractConfig(config *proto.ContractConfig) error

	// ContractConfig retrieves the current contract configuration.
	ContractConfig() *proto.ContractConfig
}

// TokenConfigurator defines methods for validating, applying, and retrieving token configuration.
type TokenConfigurator interface {
	// ValidateTokenConfig validates the provided token configuration data.
	ValidateTokenConfig(config []byte) error

	// ApplyTokenConfig applies the provided token configuration.
	ApplyTokenConfig(config *proto.TokenConfig) error

	// TokenConfig retrieves the current token configuration.
	TokenConfig() *proto.TokenConfig
}

// ExternalConfigurator defines methods for validating and applying external configuration.
type ExternalConfigurator interface {
	// ValidateExtConfig validates the provided external configuration data.
	ValidateExtConfig(cfgBytes []byte) error

	// ApplyExtConfig applies the provided external configuration to the chaincode.
	ApplyExtConfig(cfgBytes []byte) error
}

// Configure sets up the contract configuration for the given ContractConfigurable instance.
//
// This function attempts to perform the following steps:
// 1. If the ContractConfigurable instance implements the StubGetSetter interface, it sets the ChaincodeStub.
// 2. If the configuration bytes (cfgBytes) are nil, the function returns nil immediately.
// 3. It parses the configuration bytes into a ContractConfig instance.
// 4. If the ContractConfig instance has nil options, it initializes them.
// 5. It applies the parsed ContractConfig to the ContractConfigurable instance.
// 6. If the ContractConfigurable instance implements the TokenConfigurable interface, it parses and applies the TokenConfig.
// 7. If the ContractConfigurable instance implements the ExternalConfigurable interface, it applies the external configuration directly.
//
// Parameters:
// - cc: The ContractConfigurable instance to be configured.
// - stub: The ChaincodeStubInterface instance used for the contract.
// - cfgBytes: A byte slice containing the configuration data.
//
// Returns:
// - error: Returns an error if any step of the configuration process fails.
//
// Example:
//
//	err := config.Configure(myContract, myStub, configBytes)
//	if err != nil {
//	    log.Fatalf("Failed to configure contract: %v", err)
//	}
func Configure(contract Configurator, rawCfg []byte) error {
	cfg, err := FromBytes(rawCfg)
	if err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}

	if err = contract.ApplyContractConfig(cfg.GetContract()); err != nil {
		return fmt.Errorf("applying contract config: %w", err)
	}

	if tc, ok := contract.(TokenConfigurator); ok {
		if err = tc.ApplyTokenConfig(cfg.GetToken()); err != nil {
			return fmt.Errorf("applying token config: %w", err)
		}
	}

	if ec, ok := contract.(ExternalConfigurator); ok {
		if err = ec.ApplyExtConfig(rawCfg); err != nil {
			return fmt.Errorf("applying external config: %w", err)
		}
	}

	return nil
}

// Validate validates the contract configuration for the given ContractConfigurable instance.
//
// This function attempts to perform the following steps:
// 1. If the contract implements the Configurator interface, it validates the contract configuration.
// 2. If the contract implements the TokenConfigurator interface, it validates the token configuration.
// 3. If the contract implements the ExternalConfigurator interface, it validates the external configuration.
//
// Parameters:
// - contract: The Base instance to be validated.
// - rawCfg: A byte slice containing the configuration data.
//
// Returns:
// - error: Returns an error if any step of the validation process fails.
//
// Example:
//
//	err := contract.Validate(myContract, configBytes)
//	if err != nil {
//	    log.Fatalf("Failed to validate contract: %v", err)
//	}
func Validate(contract Configurator, rawCfg []byte) error {
	if err := contract.ValidateConfig(rawCfg); err != nil {
		return fmt.Errorf("validating base config: %w", err)
	}

	if tokenConfigurator, ok := contract.(TokenConfigurator); ok {
		if err := tokenConfigurator.ValidateTokenConfig(rawCfg); err != nil {
			return fmt.Errorf("validating token config: %w", err)
		}
	}

	if externalConfigurator, ok := contract.(ExternalConfigurator); ok {
		if err := externalConfigurator.ValidateExtConfig(rawCfg); err != nil {
			return fmt.Errorf("validating external config: %w", err)
		}
	}

	return nil
}
