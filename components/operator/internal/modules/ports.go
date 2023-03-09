package modules

type PortAllocator interface {
	NextPort() int32
}
type PortAllocatorFn func() int32

func (fn PortAllocatorFn) NextPort() int32 {
	return fn()
}
func StaticPortAllocator(v int32) PortAllocatorFn {
	return func() int32 {
		return v
	}
}

type portRangeAllocator struct {
	startRange int32
	offset     int32
}

func (p *portRangeAllocator) NextPort() int32 {
	p.offset++
	return p.startRange + p.offset
}

func NewPortRangeAllocator(startRange int32) *portRangeAllocator {
	return &portRangeAllocator{
		startRange: startRange,
	}
}
