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

func newNodeBuilder(output uint64) nodeBuilder {
	final := output == 0
	// transitions := make([]transition, 0, 10)
	return nodeBuilder{final, nil, output}
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
	lastKey []byte
	stack   []nodeBuilder
}

// New creates a new builder.
func New() (*Builder, error) {
	return &Builder{
		lastKey: []byte(""),
		stack:   make([]nodeBuilder, 0, 100),
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
		b.stack = append(b.stack, newNodeBuilder(value))
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
