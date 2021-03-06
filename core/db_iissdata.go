package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/icon-project/rewardcalculator/common"
	"github.com/icon-project/rewardcalculator/common/codec"
	"github.com/icon-project/rewardcalculator/common/db"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type IISSHeader struct {
	Version     uint64
	BlockHeight uint64
}

func (ih *IISSHeader) ID() []byte {
	return []byte("")
}

func (ih *IISSHeader) Bytes() ([]byte, error) {
	var bytes []byte
	if bs, err := codec.MarshalToBytes(ih); err != nil {
		return nil, err
	} else {
		bytes = bs
	}
	return bytes, nil
}

func (ih *IISSHeader) String() string {
	b, err := json.Marshal(ih)
	if err != nil {
		return "Can't covert Message to json"
	}
	return string(b)
}

func (ih *IISSHeader) SetBytes(bs []byte) error {
	_, err := codec.UnmarshalFromBytes(bs, ih)
	if err != nil {
		return err
	}
	return nil
}

func loadIISSHeader(iissDB db.Database) (*IISSHeader, error) {
	bucket, _ := iissDB.GetBucket(db.PrefixIISSHeader)
	data, err := bucket.Get([]byte(""))
	if err != nil {
		return nil, err
	}
	if data == nil {
		err = fmt.Errorf("There is no header data in IISS data\n")
		return nil, err
	}
	header := new(IISSHeader)
	err = header.SetBytes(data)
	if err != nil {
		return nil, err
	}

	return header, nil
}

type IISSGVData struct {
	IncentiveRep  uint64
	RewardRep     uint64
}

type IISSGovernanceVariable struct {
	BlockHeight   uint64
	IISSGVData
}

func (gv *IISSGovernanceVariable) ID() []byte {
	bs := make([]byte, 8)
	id := common.Uint64ToBytes(gv.BlockHeight)
	copy(bs[len(bs)-len(id):], id)
	return bs
}

func (gv *IISSGovernanceVariable) Bytes() ([]byte, error) {
	var bytes []byte
	if bs, err := codec.MarshalToBytes(&gv.IISSGVData); err != nil {
		return nil, err
	} else {
		bytes = bs
	}
	return bytes, nil
}

func (gv *IISSGovernanceVariable) String() string {
	b, err := json.Marshal(gv)
	if err != nil {
		return "Can't covert Message to json"
	}
	return string(b)
}

func (gv *IISSGovernanceVariable) SetBytes(bs []byte) error {
	_, err := codec.UnmarshalFromBytes(bs, &gv.IISSGVData)
	if err != nil {
		return err
	}
	return nil
}

func loadIISSGovernanceVariable(iissDB db.Database) ([]*IISSGovernanceVariable, error) {
	gvList := make([]*IISSGovernanceVariable, 0)
	iter, err := iissDB.GetIterator()
	if err != nil {
		return nil, err
	}
	prefix := util.BytesPrefix([]byte(db.PrefixIISSGV))
	iter.New(prefix.Start, prefix.Limit)
	for entries := 0; iter.Next(); entries++ {
		gv := new(IISSGovernanceVariable)
		err = gv.SetBytes(iter.Value())
		if err != nil {
			return nil, err
		}
		gv.BlockHeight = common.BytesToUint64(iter.Key()[len(db.PrefixIISSGV):])
		gvList = append(gvList, gv)
	}
	iter.Release()

	return gvList, nil
}

type IISSBlockProduceInfoData struct {
	Generator common.Address
	Validator []common.Address
}

type IISSBlockProduceInfo struct {
	BlockHeight uint64
	IISSBlockProduceInfoData
}

func (bp *IISSBlockProduceInfo) ID() []byte {
	bs := make([]byte, 8)
	id := common.Uint64ToBytes(bp.BlockHeight)
	copy(bs[len(bs)-len(id):], id)
	return bs
}

func (bp *IISSBlockProduceInfo) Bytes() ([]byte, error) {
	var bytes []byte
	if bs, err := codec.MarshalToBytes(&bp.IISSBlockProduceInfoData); err != nil {
		return nil, err
	} else {
		bytes = bs
	}
	return bytes, nil
}

func (bp *IISSBlockProduceInfo) String() string {
	b, err := json.Marshal(bp)
	if err != nil {
		return "Can't covert Message to json"
	}
	return string(b)
}

func (bp *IISSBlockProduceInfo) SetBytes(bs []byte) error {
	_, err := codec.UnmarshalFromBytes(bs, &bp.IISSBlockProduceInfoData)
	if err != nil {
		return err
	}
	return nil
}

func loadIISSBlockProduceInfo(iissDB db.Database) ([]*IISSBlockProduceInfo, error) {
	bpInfoList := make([]*IISSBlockProduceInfo, 0, numMainPRep)
	iter, err := iissDB.GetIterator()
	if err != nil {
		return nil, err
	}
	prefix := util.BytesPrefix([]byte(db.PrefixIISSBPInfo))
	iter.New(prefix.Start, prefix.Limit)
	for entries := 0; iter.Next(); entries++ {
		bpInfo := new(IISSBlockProduceInfo)
		err = bpInfo.SetBytes(iter.Value())
		if err != nil {
			return nil, err
		}
		bpInfo.BlockHeight = common.BytesToUint64(iter.Key()[len(db.PrefixIISSBPInfo):])
		bpInfoList = append(bpInfoList, bpInfo)
	}
	iter.Release()

	return bpInfoList, nil
}

const (
	TXDataTypeDelegate  = 0
	TXDataTypePrepReg   = 1
	TXDataTypePrepUnReg = 2
)

type IISSTXData struct {
	Address     common.Address
	BlockHeight uint64
	DataType    uint64
	Data        *codec.TypedObj
}

type IISSTX struct {
	Index uint64
	IISSTXData
}

func (tx *IISSTX) ID() []byte {
	bs := make([]byte, 8)
	id := common.Uint64ToBytes(tx.Index)
	copy(bs[len(bs)-len(id):], id)
	return bs
}

func (tx *IISSTX) Bytes() ([]byte, error) {
	var bytes []byte
	if bs, err := codec.MarshalToBytes(&tx.IISSTXData); err != nil {
		return nil, err
	} else {
		bytes = bs
	}
	return bytes, nil
}

func (tx *IISSTX) String() string {
	b, err := json.Marshal(tx)
	if err != nil {
		return "Can't covert Message to json"
	}

	return fmt.Sprintf("%s\n\t Data: %+v", string(b), common.MustDecodeAny(tx.Data))
}

func (tx *IISSTX) SetBytes(bs []byte) error {
	_, err := codec.UnmarshalFromBytes(bs, &tx.IISSTXData)
	if err != nil {
		return err
	}
	return nil
}

func loadIISSTX(iissDB db.Database) ([]*IISSTX, error) {
	txList := make([]*IISSTX, 0)
	iter, _ := iissDB.GetIterator()
	prefix := util.BytesPrefix([]byte(db.PrefixIISSTX))
	iter.New(prefix.Start, prefix.Limit)
	for entries := 0; iter.Next(); entries++ {
		tx := new(IISSTX)
		err := tx.SetBytes(iter.Value())
		if err != nil {
			return nil, err
		}
		tx.Index = common.BytesToUint64(iter.Key()[len(db.PrefixIISSTX):])
		txList = append(txList, tx)
	}
	iter.Release()

	return txList, nil
}

func LoadIISSData(dbPath string, verbose bool) (*IISSHeader, []*IISSGovernanceVariable,
	[]*IISSBlockProduceInfo, []*PRep, []*IISSTX) {
	dbPath = filepath.Clean(dbPath)
	dbDir, dbName := filepath.Split(dbPath)

	iissDB := db.Open(dbDir, string(db.GoLevelDBBackend), dbName)
	defer iissDB.Close()

	// Header
	header, err := loadIISSHeader(iissDB)
	if err != nil {
		log.Printf("Failed to read header from IISS Data. err=%+v\n", err)
		return nil, nil, nil, nil, nil
	}

	// Governance Variable
	gvList, err := loadIISSGovernanceVariable(iissDB)
	if err != nil {
		log.Printf("Failed to read governance variable from IISS Data. err=%+v\n", err)
		return nil, nil, nil, nil, nil
	}

	// Block produce Info.
	bpInfoList, err := loadIISSBlockProduceInfo(iissDB)
	if err != nil {
		log.Printf("Failed to read Block Produce Info. from IISS Data. err=%+v\n", err)
		return nil, nil, nil, nil, nil
	}

	// Main/Sub P-Rep
	pRepList, err := LoadPRep(iissDB)
	if err != nil {
		log.Printf("Failed to read P-Rep list from IISS Data. err=%+v\n", err)
		return nil, nil, nil, nil, nil
	}

	// TX list
	txList, err := loadIISSTX(iissDB)
	if err != nil {
		log.Printf("Failed to read TX list from IISS Data. err=%+v\n", err)
		return nil, nil, nil, nil, nil
	}

	if verbose {
		log.Printf("Header: %s\n", header.String())

		log.Printf("Governance variable:\n")
		for i, gv := range gvList {
			log.Printf("\t%d: %s", i, gv.String())
		}

		log.Printf("Block Produce Info.:\n")
		for i, bpInfo := range bpInfoList {
			log.Printf("\t%d: %s\n", i, bpInfo.String())
		}

		log.Printf("Main/Sub P-Rep list:\n")
		for i, preps:= range pRepList {
			log.Printf("\t%d: %s\n", i, preps.String())
		}

		log.Printf("TX:\n")
		for i, tx := range txList {
			log.Printf("\t%d: %s\n", i, tx.String())
		}
	}

	return header, gvList, bpInfoList, pRepList, txList
}

func findIISSData(dir string) []os.FileInfo {
	iissData := make([]os.FileInfo, 0)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil
	}

	for _, f := range files {
		if f.IsDir() == true && strings.HasPrefix(f.Name(), "iiss_") == true {
			iissData = append(iissData, f)
		}
	}

	return iissData
}
