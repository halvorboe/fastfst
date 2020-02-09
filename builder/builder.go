package builder

type transition struct {
	in   byte
	addr uint64
}
type nodeBuilder struct {
	final       bool
	transitions []transition
	output      uint64
}

func newNodeBuilder() nodeBuilder {
	builder := nodeBuilder{transitions: make([]transition, 0, 32)}
	builder.Reset()
	return builder
}

func (nb *nodeBuilder) Compile() {

}

func (nb *nodeBuilder) Reset() {
	nb.final = false
	nb.transitions = nb.transitions[:0]
	nb.output = 0
}

// Builder builds a FST.
type Builder struct {
	lastKey  []byte
	stack    []*nodeBuilder
	builders []nodeBuilder
	existing map[int]int
}

// New creates a new builder.
func New() (*Builder, error) {
	builders := make([]nodeBuilder, 0, 100)
	for i := 0; i < 100; i++ {
		builders = append(builders, newNodeBuilder())
	}
	return &Builder{
		lastKey:  []byte(""),
		stack:    make([]*nodeBuilder, 0, 64),
		builders: builders,
		existing: make(map[int]int, 1000),
	}, nil
}

// Insert inserts a key-value pair into the fst.
func (b *Builder) Insert(key []byte, value uint64) error {
	// 1. Verify that the key is valid.
	// 2. Find the prefix.
	prefixLength := b.findPrefixLength(key)
	// 3. Freeze nodes until stack contains only the prefix.
	b.compileTo(prefixLength)
	b.stack = b.stack[:prefixLength]
	// 4. Append the suffix to the stack.
	for i := prefixLength; i < len(key); i++ {
		if i > 0 {
			b.stack[i-1].transitions = append(b.stack[i-1].transitions, transition{in: key[i]})
		}
		b.stack = append(b.stack, &b.builders[len(b.stack)])
		b.stack[len(b.stack)-1].Reset()
	}
	// Update last key
	b.lastKey = key
	return nil
}

func (b *Builder) compileTo(size int) {
	for i := len(b.stack) - 1; i >= size; i-- {
		b.stack[i].Compile()
	}
}

func (b *Builder) findPrefixLength(key []byte) int {
	maxLength := min(len(key), len(b.lastKey))
	for i := 0; i < maxLength; i++ {
		if key[i] != b.lastKey[i] {
			return i
		}
	}
	return maxLength
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Close cleans up after the builder.
func (b *Builder) Close() error {
	b.compileTo(0)
	return nil
}
