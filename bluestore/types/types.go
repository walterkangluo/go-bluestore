package types

type ShardIdT struct {
	Id int8
}

func CreateShardIdT(id int8) *ShardIdT {
	return &ShardIdT{
		Id: id,
	}
}
