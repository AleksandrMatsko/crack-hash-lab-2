package cache

import (
	"distributed.systems.labs/shared/pkg/contracts"
	"distributed.systems.labs/worker/internal/config"
	"fmt"
	"sync"
)

func MakeKey(req contracts.TaskRequest) string {
	return fmt.Sprintf("%s_%v_%v_%s_%s_%v",
		req.RequestID, req.StartIndex, req.PartCount, req.Alphabet, req.ToCrack, req.MaxLength)
}

type RequestInfo struct {
	Req    contracts.TaskRequest
	Status config.RequestStatus
	Rsp    contracts.TaskResultRequest
}

type Cache struct {
	data map[string]RequestInfo
	mtx  sync.Mutex
}

func New() *Cache {
	return &Cache{
		data: make(map[string]RequestInfo),
		mtx:  sync.Mutex{},
	}
}

func (c *Cache) GetOrAdd(req contracts.TaskRequest) (RequestInfo, bool) {
	key := MakeKey(req)

	c.mtx.Lock()
	defer c.mtx.Unlock()

	val, ok := c.data[key]
	if ok {
		return val, ok
	} else {
		val = RequestInfo{
			Req:    req,
			Status: config.InProgress,
			Rsp:    contracts.TaskResultRequest{},
		}
		c.data[key] = val
	}
	return val, ok
}

func (c *Cache) SetDone(req contracts.TaskRequest, rsp contracts.TaskResultRequest) {
	key := MakeKey(req)

	c.mtx.Lock()
	defer c.mtx.Unlock()

	val, ok := c.data[key]
	if !ok {
		c.data[key] = RequestInfo{
			Req:    req,
			Status: config.Done,
			Rsp:    rsp,
		}
	} else {
		val.Rsp = rsp
		val.Status = config.Done
		c.data[key] = val
	}
}
