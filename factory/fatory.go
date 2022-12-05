package factory

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
	muType factory.MuType
	muDriverConf any
}

func NewDistribMuElectDriverConf(key string, ttl time.Duration, muType factory.MuType, muDriverConf any) *DistribMuElectDriverConf {
	return &DistribMuElectDriverConf{
		key: key,
		ttl: ttl,
		muType: muType,
		muDriverConf: muDriverConf,
	}
}

func NewAutoElection(driver ElectDriver, driverConf any) (autoelect.AutoElection, error) {
	switch driver {
	case ElectDriverGitDistribMu:
		specConf := driverConf.(*DistribMuElectDriverConf)
		return gitdistribmu.New(specConf.key, specConf.ttl, specConf.muType, specConf.muDriverConf)
	}
	return nil, fmt.Errorf("invalid driver:%d", driver)
}