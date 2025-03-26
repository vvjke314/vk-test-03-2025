package entities

type VaultItem struct {
	Key   string `msgpack:"key"`
	Value string `msgpack:"value"`
}
