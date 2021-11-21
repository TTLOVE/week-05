package circuit

import (
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

/**
Bucket为整个移动窗口的环的桶
Total 为该Bucket的总请求数
Success_num 为该Bucket的成功请求数
Fail_num 为该Bucket的失败请求数
Timeout_num 为该Bucket的超时请求数
Reject_num 为该Bucket的拒绝请求数
*/
type Bucket struct {
	Total, Success_num, Fail_num, Timeout_num, Reject_num uint32
	Start_time                                            int64
}

/**
Circuit为整个移动窗口的环
Size 为该Circuit的总Bucket - 1
Head 为该Circuit的头部位置
Tail 为该Circuit的尾部位置
Total_mtime 为该Circuit统计的总时长
Bucket_mtime 为该Circuit每个Bucket统计的时长
Buckets 为该Circuit的所有Bucket信息
*/
type Circuit struct {
	Size, Head, Tail int32
	Total_mtime      int
	Bucket_mtime     int
	Buckets          []*Bucket
}

// 设置最新的的Bucket
func (c *Circuit) newBucket(ts int64) {
	if c.Buckets[c.Tail] == nil {
		c.Buckets[c.Tail] = &Bucket{
			Start_time:  ts,
			Total:       0,
			Success_num: 0,
			Fail_num:    0,
			Timeout_num: 0,
			Reject_num:  0,
		}
	} else {
		c.resetBucket(c.Tail, ts)
	}
}

// 重置Bucket信息
func (c *Circuit) resetBucket(n int32, ts int64) {
	c.Buckets[n].Start_time = ts
	c.Buckets[n].Total = 0
	c.Buckets[n].Success_num = 0
	c.Buckets[n].Fail_num = 0
	c.Buckets[n].Timeout_num = 0
	c.Buckets[n].Reject_num = 0
}

// 获取当前的Bucket
func (c *Circuit) getCurrentBucket() *Bucket {
	nowTime := time.Now().UnixMilli()

	b := c.Buckets[c.Tail]
	// 如果是生成第一个桶
	if b == nil && c.Tail == 0 {
		c.newBucket(nowTime)
	} else {
		// 判断时间是否在当前的桶的时间内
		endTime := b.Start_time + int64(c.Bucket_mtime)
		if endTime > nowTime {
			c.Buckets[c.Tail] = b
			return b
		}

		// 超出了当前桶的时间,生成新的桶
		tail := c.Tail + 1
		// 如果尾部大于bucket长度
		if tail > c.Size {
			tail %= c.Size + 1
		}

		if tail-c.Head >= c.Size || tail < c.Head {
			head := c.Head + 1
			if head > c.Size {
				head %= c.Size + 1
			}

			atomic.StoreInt32(&c.Head, head)
		}

		atomic.StoreInt32(&c.Tail, tail)

		c.newBucket(endTime + 1)
	}

	return c.Buckets[c.Tail]
}

const (
	Success = iota
	Fail
	Timeout
	Reject
)

// 根据类型获取当前Circuit对应类型的总数
func (c *Circuit) GetNumByType(t int) (uint32, error) {
	var sum uint32

	for i := c.Head; i < c.Tail; i++ {
		switch t {
		case Success:
			sum += c.Buckets[i].Success_num
		case Fail:
			sum += c.Buckets[i].Fail_num
		case Timeout:
			sum += c.Buckets[i].Timeout_num
		case Reject:
			sum += c.Buckets[i].Reject_num
		default:
			return 0, errors.New("type not match")
		}
	}

	return sum, nil
}

// 根据类型给当前Bucket的数量加1
func (c *Circuit) AddNumByType(t int) error {
	cb := c.getCurrentBucket()
	var column *uint32

	switch t {
	case Success:
		column = &cb.Success_num
	case Fail:
		column = &cb.Fail_num
	case Timeout:
		column = &cb.Timeout_num
	case Reject:
		column = &cb.Reject_num
	default:
		return errors.New("type not match")
	}

	atomic.AddUint32(column, 1)
	atomic.AddUint32(&cb.Total, 1)

	return nil
}
