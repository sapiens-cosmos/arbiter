package types

const (
	// ModuleName defines the module name
	ModuleName = "stake"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	ParamsKey = "params"
)

var (
	// KeyEpoch defines key for storing epochs
	KeyEpoch = []byte{0x01}

	// KeyPrefixLock defines key for storing individual locks
	KeyPrefixLock = []byte{0x02}

	// KeyPrefixLock defines key for storing total reserve of the protocol
	KeyStakeState = []byte{0x03}
)

func GetAddressLockStoreKey(address string) []byte {
	return append(KeyPrefixLock, []byte(address)...)
}
