package commands

import (
	"io"
	"strings"

	cmds "github.com/ipfs/go-ipfs/commands"
	core "github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreunix"
	path "github.com/ipfs/go-ipfs/path"
	tar "github.com/ipfs/go-ipfs/tar"
)

var TarCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Utility functions for tar files in IPFS.",
	},

	Subcommands: map[string]*cmds.Command{
		"add": tarAddCmd,
		"cat": tarCatCmd,
	},
}

var tarAddCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Import a tar file into ipfs.",
		ShortDescription: `
'ipfs tar add' will parse a tar file and create a merkledag structure to
represent it.
`,
	},

	Arguments: []cmds.Argument{
		cmds.FileArg("file", true, false, "Tar file to add.").EnableStdin(),
	},
	Run: func(req cmds.Request, res cmds.Response) {
		nd, err := req.InvocContext().GetNode()
		if err != nil {
			res.SetError(err, cmds.ErrNormal)
			return
		}

		fi, err := req.Files().NextFile()
		if err != nil {
			res.SetError(err, cmds.ErrNormal)
			return
		}

		node, err := tar.ImportTar(fi, nd.DAG)
		if err != nil {
			res.SetError(err, cmds.ErrNormal)
			return
		}

		k, err := node.Key()
		if err != nil {
			res.SetError(err, cmds.ErrNormal)
			return
		}

		fi.FileName()
		res.SetOutput(&coreunix.AddedObject{
			Name: fi.FileName(),
			Hash: k.B58String(),
		})
	},
	Type: coreunix.AddedObject{},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: func(res cmds.Response) (io.Reader, error) {
			o := res.Output().(*coreunix.AddedObject)
			return strings.NewReader(o.Hash + "\n"), nil
		},
	},
}

var tarCatCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Export a tar file from IPFS.",
		ShortDescription: `
'ipfs tar cat' will export a tar file from a previously imported one in IPFS.
`,
	},

	Arguments: []cmds.Argument{
		cmds.StringArg("path", true, false, "IPFS path of archive to export."),
	},
	Run: func(req cmds.Request, res cmds.Response) {
		nd, err := req.InvocContext().GetNode()
		if err != nil {
			res.SetError(err, cmds.ErrNormal)
			return
		}

		p, err := path.ParsePath(req.Arguments()[0])
		if err != nil {
			res.SetError(err, cmds.ErrNormal)
			return
		}

		root, err := core.Resolve(req.Context(), nd, p)
		if err != nil {
			res.SetError(err, cmds.ErrNormal)
			return
		}

		r, err := tar.ExportTar(req.Context(), root, nd.DAG)
		if err != nil {
			res.SetError(err, cmds.ErrNormal)
			return
		}

		res.SetOutput(r)
	},
}
