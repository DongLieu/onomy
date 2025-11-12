package types

const (
	// ModuleName defines the module name.
	ModuleName = "gnark"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName

	// RouterKey is the message route for gnark.
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key.
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key.
	MemStoreKey = "mem_gnark"
)

// KV prefixes.
var (
	CircuitKeyPrefix = []byte{0x01}
)

// CircuitKey returns the store key for the provided circuit identifier.
func CircuitKey(id string) []byte {
	key := make([]byte, len(CircuitKeyPrefix)+len(id))
	copy(key, CircuitKeyPrefix)
	copy(key[len(CircuitKeyPrefix):], []byte(id))
	return key
}
