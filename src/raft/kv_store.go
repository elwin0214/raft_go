package raft

type KVStore struct {
	server *Server
	data   map[string]string
}

func (kv *KVStore) apply(entry LogEntry) {
	cmd := entry.Cmd.(*Command)
	if cmd.Op == "set" {
		kv.data[cmd.Key] = cmd.Value
	}
}

func (kv *KVStore) reply(entry LogEntry) {
	cmd := entry.Cmd.(*Command)
	if cmd.Op == "get" {
		cmd.Value = kv.data[cmd.Key]
	}
}

func setKVStore(server *Server) *KVStore {
	kv := &KVStore{server, make(map[string]string, 1024)}
	kv.server.store = kv
	return kv
}
