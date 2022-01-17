package crypto

type Block struct {
	Sum  [16]byte
	IV   [8]byte
	Data []byte
}

func StreamToBlock(stream []byte) Block {
	return Block{}
}

func (b *Block) ToStream() []byte {
	return nil
}
