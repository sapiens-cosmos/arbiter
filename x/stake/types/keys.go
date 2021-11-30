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
)

var (
	// KeyEpoch defines key for storing epochs
	KeyEpoch = []byte{0x01}

	// KeyPrefixLock defines key for storing individual locks
	KeyPrefixLock = []byte{0x02}
)

func GetAddressLockStoreKey(address string) []byte {
	return append(KeyPrefixLock, []byte(address)...)
}
