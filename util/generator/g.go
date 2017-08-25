package generator

import (
	"fmt"
	"sync"
	"time"
)

//  		------------------------------------------------
//elements  |timestamp| datacenter(dc) | node | sequence   | total 64bit
//variable  | ts      | dc             | node |  seq       |

func init() {
	var err error
	Default, err = NewGenerator()
	if err != nil {
		panic(err)
	}
}

var (
	nanosInMilli           = time.Millisecond.Nanoseconds()
	defaultDCBits   uint64 = 2
	defaultNodeBits uint64 = 6
	defaultSeqBits  uint64 = 10

	Default *Generator
)
var (
	MaxExcludeTS  uint64 = 22 //set the sequence+node+data center <=22. all the rest spaces are for timestamp
	DefaultDCID   uint64 = 1
	DefaultNodeID uint64 = 1
	t, _                 = time.Parse("2006-01-02 15:04:05 -0700 UTC", "2017-08-03 18:48:03 +0800 UTC")
	defaultEpoch         = uint64(t.UnixNano() / nanosInMilli)
)

type Generator struct {
	*sync.Mutex
	epoch    uint64
	dc       uint64
	dcBits   uint64
	node     uint64
	nodeBits uint64
	seq      uint64
	seqBits  uint64
	lastTS   uint64
}
type Options struct {
	DataCenterID   uint64
	DataCenterBits uint64
	NodeID         uint64
	NodeBits       uint64
	SequenceBits   uint64
	Epoch          uint64
}

func mask(bits uint64) uint64 {
	return 1<<bits - 1
}

func (opts *Options) build() (*Generator, error) {
	if opts.Epoch <= 0 {
		opts.Epoch = defaultEpoch
	}
	// use default,if not specified.
	if opts.DataCenterBits <= 0 {
		opts.DataCenterBits = defaultDCBits

	}
	if opts.NodeBits <= 0 {
		opts.NodeBits = defaultNodeBits
	}
	if opts.SequenceBits <= 0 {
		opts.SequenceBits = defaultSeqBits
	}

	//guarantee Legitimacy
	maxNodeID := mask(opts.NodeBits)
	if opts.NodeID > maxNodeID {
		return nil, fmt.Errorf("node bit:%d, node id should be <=%d", opts.NodeBits, maxNodeID)
	}
	maxDataCenterID := mask(opts.DataCenterBits)
	if opts.DataCenterID > maxDataCenterID {
		return nil, fmt.Errorf("data center bit:%d,  data center id  should be <=%d", opts.DataCenterBits, maxNodeID)
	}
	currentExcludeTS := opts.DataCenterBits + opts.NodeBits + opts.SequenceBits
	if currentExcludeTS > MaxExcludeTS {
		return nil, fmt.Errorf("opts.DataCenterBits + opts.NodeBits+opts.SequenceBits=%d, should be <= %d", currentExcludeTS, MaxExcludeTS)
	}

	//use default dc & node if not specified
	if opts.DataCenterID < 0 {
		opts.DataCenterID = DefaultDCID
	}
	if opts.NodeID < 0 {
		opts.NodeID = DefaultNodeID
	}
	return &Generator{
		Mutex:  &sync.Mutex{},
		epoch:  opts.Epoch,
		dc:     opts.DataCenterID,
		dcBits: opts.DataCenterBits,

		node:     opts.NodeID,
		nodeBits: opts.NodeBits,

		seq:     0,
		seqBits: opts.SequenceBits,
	}, nil

}
func NewDataCenterNode(dcID, nodeID uint64) (*Generator, error) {
	options := &Options{
		DataCenterID: dcID,
		NodeID:       nodeID,
	}
	return options.build()
}

// NewGenerator ,use default Generator
func NewGenerator() (*Generator, error) {
	opts := new(Options)
	return opts.build()
}

func (g *Generator) maxSeq() uint64 {
	return mask(g.seqBits)
}

func (g *Generator) nodeShift() uint64 {
	return g.seqBits
}
func (g *Generator) datacenterShift() uint64 {
	return g.nodeBits + g.seqBits
}
func (g *Generator) timestampShift() uint64 {
	return g.dcBits + g.nodeBits + g.seqBits
}

//compute is not safe for concurrency.
//you should call the Next() or NextN go guarantee safety
func (g *Generator) compute() uint64 {
	ts := uint64(time.Now().UnixNano()/nanosInMilli) - g.epoch
	if ts <= g.lastTS {
		g.seq += 1
		if g.seq > g.maxSeq() {
			g.lastTS++
			g.seq = 0
		}
	} else {
		g.seq = 0
		g.lastTS = ts

	}
	id := (g.lastTS << g.timestampShift()) | g.dc<<g.datacenterShift() | g.node<<g.nodeShift() | g.seq
	return id
}

// Next , generate one id
func (g *Generator) Next() uint64 {
	g.Lock()
	id := g.compute()
	g.Unlock()
	return id
}

//NextN , generate specific amount of ids.
func (g *Generator) NextN(n int64) []uint64 {
	g.Lock()
	ids := make([]uint64, 0, n)
	var i int64
	for i = 0; i < n; i++ {
		ids = append(ids, g.compute())
	}
	g.Unlock()
	return ids
}
func (g *Generator) Detail(num uint64) (ts uint64, dcID uint64, nodeID uint64, seq uint64, genAt time.Time) {
	ts = num >> g.timestampShift()
	seq = mask(g.seqBits) & num
	nodeID = mask(g.seqBits+g.nodeBits) & num >> g.nodeShift()
	dcID = mask(g.seqBits+g.nodeBits+g.dcBits) & num >> g.datacenterShift()
	timeSec := int64(ts+g.epoch) / 1000
	timeNanoSec := (int64(ts+g.epoch) % 1000) * 1e6
	genAt = time.Unix(timeSec, timeNanoSec)
	return
}
