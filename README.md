# week-05

Bucket为整个移动窗口的环的桶

Total 为该Bucket的总请求数

Success_num 为该Bucket的成功请求数

Fail_num 为该Bucket的失败请求数

Timeout_num 为该Bucket的超时请求数

Reject_num 为该Bucket的拒绝请求数

```
type Bucket struct {
	Total, Success_num, Fail_num, Timeout_num, Reject_num uint32
	Start_time                                            int64
}
```

Circuit为整个移动窗口的环

Size 为该Circuit的总Bucket - 1

Head 为该Circuit的头部位置

Tail 为该Circuit的尾部位置

Total_mtime 为该Circuit统计的总时长

Bucket_mtime 为该Circuit每个Bucket统计的时长

Buckets 为该Circuit的所有Bucket信息

```
type Circuit struct {
	Size, Head, Tail int32
	Total_mtime      int
	Bucket_mtime     int
	Buckets          []*Bucket
}
```
