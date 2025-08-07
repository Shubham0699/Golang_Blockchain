package main

import "fmt"
import "bytes"

type block struct{
	timestamp int64
	data  []byte
	PrevBlockHash []byte
	Hash  []byte
}

func (b *Block) SetHash() {
headers:=[][]byte{

	b.PrevBlockhash,
	b.Data,
	[]byte(strconv.FormatInt(b.Timestamp,10))

}
glue:=[]byte()
result := bytes.Join(headers, glue)
}
