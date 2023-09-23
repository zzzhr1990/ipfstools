package getter

import (
	"errors"
	"net/http"
	"time"

	"github.com/guonaihong/gout"
	"golang.org/x/net/context"

	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"

	//cid "github.com/ipfs/go-cid"
	legacy "github.com/ipfs/go-ipld-legacy"
	"github.com/ipfs/go-merkledag"
	dagpb "github.com/ipld/go-codec-dagpb"

	//ipldformat "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

type HttpGetter struct {
	decoder *legacy.Decoder
}

func NewHttpGetter() *HttpGetter {
	d := legacy.NewDecoder()
	d.RegisterCodec(cid.DagProtobuf, dagpb.Type.PBNode, merkledag.ProtoNodeConverter)
	d.RegisterCodec(cid.Raw, basicnode.Prototype.Bytes, merkledag.RawNodeConverter)
	return &HttpGetter{
		decoder: d,
	}
}

func (h *HttpGetter) Getblock(cidString string) (legacy.UniversalNode, error) {
	c, err := cid.Decode(cidString)
	if err != nil {
		return nil, err
	}
	url := "https://http-file-proxy-v6.2dland.cn/raw/" + cidString
	body := []byte{}
	code := 0
	err = gout.NewWithOpt(gout.WithClient(&http.Client{}), gout.WithTimeout(time.Second*30)).GET(url).BindBody(&body).Code(&code).Do()
	if err != nil {
		return nil, err
	}
	if code != http.StatusOK && code != http.StatusPartialContent {
		message := string(body)
		if message == "" {
			message = http.StatusText(code)
		}
		return nil, errors.New("http get error: " + message)
	}
	block, err := blocks.NewBlockWithCid(body, c)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	node, err := h.decoder.DecodeNode(ctx, block)
	if err != nil {
		return nil, err
	}
	return node, nil
}
