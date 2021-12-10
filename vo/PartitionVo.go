package vo

type PartitionVo struct {
	ID      uint   `json:"id"`
	Content string `json:"content"`
}

type AllPartitionVo struct {
	ID           uint          `json:"id"`
	Content      string        `json:"content"`
	Subpartition []PartitionVo `json:"subpartition"`
}
