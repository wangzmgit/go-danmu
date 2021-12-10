package dto

type PartitionDto struct {
	Content string
	Fid     uint //父分区ID
}

type DeletePartitionDto struct {
	ID uint
}
