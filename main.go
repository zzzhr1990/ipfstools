package main

import (
	"fmt"
	"ipfstools/getter"

	"github.com/multiformats/go-multicodec"
	//cid "github.com/ipfs/go-cid"
	// ipldformat "github.com/ipfs/go-ipld-format"
	// "github.com/ipfs/go-ipld-legacy"
)

func main() {
	fileCid := "bafybeifoetl2q4drlxvoschmwe22oc2dgvmq7vtzeqax4ihcxkdt3bwtvy"
	httpGetter := getter.NewHttpGetter()
	node, err := httpGetter.Getblock(fileCid)
	if err != nil {
		panic(err)
	}

	typeCode := node.Cid().Type()
	// convert code to hex string
	typeCodeStr := fmt.Sprintf("0x%x", typeCode)
	println(typeCodeStr)
	if multicodec.Code(typeCode) == multicodec.DagPb {
		println("dagpb")
		// extract data from node
		size, err := node.Size()
		if err != nil {
			panic(err)
		}
		fmt.Printf("size: %d, len: %v\n", size, node.Length())

		for _, link := range node.Links() {
			fmt.Printf("link: %+v\n", link.Cid.String())
		}
	} else {
		println("not dagpb")
	}
}
