package plugins

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/pkg/errors"
)

func init() {
	caddy.RegisterModule(ZipFS{})
}

type ZipFS struct {
	Path      string `json:"path"`
	zipReader *zip.ReadCloser
}

func (ZipFS) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "caddy.fs.zip",
		New: func() caddy.Module { return new(ZipFS) },
	}
}

func (z *ZipFS) Cleanup() error {
	if z.zipReader == nil {
		return nil
	}
	return z.zipReader.Close()
}

func (z *ZipFS) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	if !d.Next() { // skip block beginning
		return d.ArgErr()
	}

	for nesting := d.Nesting(); d.NextBlock(nesting); {
		switch d.Val() {
		case "path":
			if !d.AllArgs(&z.Path) {
				return d.ArgErr()
			}
		default:
			return d.Errf("%s not a valid caddy.fs.zip option", d.Val())
		}
	}

	return nil
}

func (fs *ZipFS) Provision(ctx caddy.Context) error {
	var err error
	fs.zipReader, err = zip.OpenReader(fs.Path)
	if err != nil {
		return errors.Wrapf(err, "opening file %s", fs.Path)
	}

	return nil
}

func (z *ZipFS) Open(name string) (fs.File, error) {
	ret, err := z.zipReader.Open(name)
	if err != nil {
		return nil, err
	}
	stat, err := ret.Stat()
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return ret.(fs.ReadDirFile), nil
	}
	return &File{
		File: ret,
	}, nil
}

func (z *ZipFS) Stat(name string) (fs.FileInfo, error) {
	f, err := z.zipReader.Open(name)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return stat, nil
}

// Interface guards
var (
	_ fs.StatFS             = (*ZipFS)(nil)
	_ caddyfile.Unmarshaler = (*ZipFS)(nil)
	_ caddy.Provisioner     = (*ZipFS)(nil)
	_ caddy.CleanerUpper    = (*ZipFS)(nil)
)

type File struct {
	fs.File
	offset int64
}

func (f *File) Read(b []byte) (int, error) {
	data := make([]byte, f.offset+int64(len(b)))
	_, err := f.File.Read(data)
	if err != nil && !errors.Is(err, io.EOF) {
		return 0, err
	}

	bytesRead := copy(b, data[f.offset:])
	f.offset += int64(len(b))

	return bytesRead, nil
}

func (f File) Seek(offset int64, whence int) (int64, error) {
	stat, err := f.File.Stat()
	if err != nil {
		return 0, err
	}

	var newOffset int64
	switch whence {
	case io.SeekStart:
		newOffset = offset
	case io.SeekCurrent:
		newOffset = f.offset + offset
	case io.SeekEnd:
		newOffset = stat.Size() + offset
	default:
		return 0, errors.New("Unknown Seek Method")
	}
	if newOffset > stat.Size() || newOffset < 0 {
		return 0, fmt.Errorf("invalid offset %d", offset)
	}
	f.offset = newOffset
	return newOffset, nil
}

// notes(gfyrag): Caddy internally cast to io.ReadSeeker, so we need to provide an implementation
var _ io.ReadSeeker = (*File)(nil)
