package main

import (
	"fmt"
	"github.com/icon-project/rewardcalculator/common/db"
	"github.com/icon-project/rewardcalculator/rewardcalculator"
)

func (cli *CLI) header(version uint64, blockHeight uint64) {
	bucket, _ := cli.DB.GetBucket(db.PrefixIISSHeader)

	header := new(rewardcalculator.IISSHeader)
	header.Version = uint16(version)
	header.BlockHeight = blockHeight

	key := []byte("")
	value, _ := header.Bytes()
	bucket.Set(key, value)

	fmt.Printf("Set header %s\n", header.String())
}