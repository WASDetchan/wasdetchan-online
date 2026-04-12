package auth

type Capabilities []Capability

type Capability uint32

const (
	CapabilityAdmin Capability = 1 + iota
	CapabilityTokenManip
	CapabilityReceiptWrite
	CapabilityReceiptRead
	capabilityCount
)

func AdminCapabilities() Capabilities {
	return append(AuthenticatedCapabilities(), CapabilityAdmin)
}

func AuthenticatedCapabilities() Capabilities {
	return Capabilities{
		CapabilityTokenManip,
		CapabilityReceiptWrite,
		CapabilityReceiptRead,
	}
}

func (cs Capabilities) Has(c Capability) bool {
	for _, cp := range cs {
		if c == cp {
			return true
		}
	}
	return false
}

func (cs *Capabilities) ReadBytes(bytes []byte) {
	num := uint32(0)
	for i, b := range bytes {
		num += uint32(b) << (i & 3 * 8)
		if i%4 == 3 {
			*cs = append(*cs, Capability(num))
			num = 0
		}
	}
}

func (cs *Capabilities) IntoBytes() []byte {
	bytes := make([]byte, len(*cs)*4)
	i := 0
	for _, c := range *cs {
		for _ = range 4 {
			bytes[i] = byte(c & 255)
			c >>= 8
			i++
		}
	}
	return bytes
}
