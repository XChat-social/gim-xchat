package sequence

import (
	"fmt"
	"gim/pkg/db"
	"gim/pkg/gerrors"
)

const SeqKey = "seq:%s"

// GetNextSeq 获取指定key的下一个序列号
func GetNextSeq(key string) (int64, error) {
	num, err := db.RedisCli.Incr(fmt.Sprintf(SeqKey, key)).Result()
	if err != nil {
		return 0, gerrors.WrapError(err)
	}
	return num, nil
}
