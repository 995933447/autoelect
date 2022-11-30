package fatory

import (
	"fmt"
	"github.com/995933447/autoelect"
	"github.com/995933447/autoelect/impl/distribmu/gitdistribmu"
	"github.com/995933447/distribmu/factory"
	"time"
)

type ElectDriver int

const (
	ElectDriverNil = iota
	ElectDriverGitDistribMu
)

type DistribMuElectDriverConf struct {
	key string
	ttl time.Duration
	muConfType factory.MuType
	muDriverConf interface{}
}

func NewAuthElection(driver ElectDriver, driverConf interface{}) (autoelect.AutoElection, error) {
	switch driver {
	case ElectDriverGitDistribMu:
		specConf := driverConf.(*DistribMuElectDriverConf)
		return _9distribmu.New(specConf.key, specConf.ttl, specConf.muConfType, specConf.muDriverConf)
	}
	return nil, fmt.Errorf("invalid driver:%d", driver)
}