package cwimkit

import (
	"io"

	"github.com/go-wimkit/cwimkit/internal/csize"
)

type resourceEntry struct {
	/** If this blob is not missing, then this is the uncompressed size of
	 * this blob in bytes.  */
	UncompressedSize uint64

	/** If this blob is located in a non-solid WIM resource, then this is
	 * the compressed size of that resource.  */
	CompressedSize uint64

	/** If this blob is located in a non-solid WIM resource, then this is
	 * the offset of that resource within the WIM file containing it.  If
	 * this blob is located in a solid WIM resource, then this is the offset
	 * of this blob within that solid resource when uncompressed.  */
	Offset uint64

	/** If this blob is located in a WIM resource, then this is the SHA-1
	 * message digest of the blob's uncompressed contents.  */
	Hash [20]byte

	/** If this blob is located in a WIM resource, then this is the part
	 * number of the WIM file containing it.  */
	PartNumber uint32

	/** If this blob is not missing, then this is the number of times this
	 * blob is referenced over all images in the WIM.  This number is not
	 * guaranteed to be correct.  */
	ReferenceCount uint32

	/** Flags */
	Flags uint32

	/** If this blob is located in a solid WIM resource, then this is the
	 * offset of that solid resource within the WIM file containing it.  */
	RawResourceOffsetInWim uint64

	/** If this blob is located in a solid WIM resource, then this is the
	 * compressed size of that solid resource.  */
	RawResourceCompressedSize uint64

	/** If this blob is located in a solid WIM resource, then this is the
	 * uncompressed size of that solid resource.  */
	RawResourceUncompressedSize uint64

	Reserved [8]byte
}

type streamEntryPacked struct {
	Name Pointer

	Reserved0 [4]byte

	Resource resourceEntry

	Reserved [32]byte
}

var streamEntrySize = csize.OfMust(streamEntryPacked{})

type StreamEntry struct {
	Packed streamEntryPacked

	api *API
}

func newStreamEntry(api *API, ptr Pointer) (*StreamEntry, error) {
	streamEntry := &StreamEntry{api: api}

	if err := api.readStruct(ptr, &streamEntry.Packed); err != nil {
		return nil, err
	}

	return streamEntry, nil
}

type wimlibTimeSpec struct {
	TvSec  uint64
	TvNSec uint64
}

const WimlibGuidLen = 16

type wimlibObjectID struct {
	ObjectID      [WimlibGuidLen]byte
	BirthVolumeID [WimlibGuidLen]byte
	BirthObjectID [WimlibGuidLen]byte
	DomainID      [WimlibGuidLen]byte
}

type dirEntryPacked struct {
	Filename  Pointer
	ShortName Pointer
	FullPath  Pointer

	Depth int32

	SecurityDescriptor     Pointer
	SecurityDescriptorSize int32

	Attributes uint32

	ReparseTag uint32

	LinksNum uint32

	NamedStreamsNum uint32

	HardlinkGroupID uint64

	CreationTime   wimlibTimeSpec
	LastWriteTime  wimlibTimeSpec
	LastAccessTime wimlibTimeSpec

	UnixUID  uint32
	UnixGID  uint32
	UnixMode uint32
	UnixRDev uint32

	ObjectID wimlibObjectID

	Reserved [48]byte
}

type DirEntry struct {
	Packed dirEntryPacked

	streamsPtr Pointer

	api *API
}

func (e *DirEntry) Filename() (string, error) {
	return e.api.readZString(e.Packed.Filename)
}

func (e *DirEntry) ShortName() (string, error) {
	return e.api.readZString(e.Packed.ShortName)
}

func (e *DirEntry) FullPath() (string, error) {
	return e.api.readZString(e.Packed.FullPath)
}

func (e *DirEntry) StreamsNum() int {
	return int(e.Packed.NamedStreamsNum + 1)
}

func (e *DirEntry) Stream(i int) (*StreamEntry, error) {
	if i > e.StreamsNum() {
		return nil, io.EOF
	}

	ptr := e.streamsPtr + Pointer(uint32(i)*streamEntrySize)

	return newStreamEntry(e.api, ptr)
}

func newDirEntryFromPtr(api *API, ptr Pointer) (*DirEntry, error) {
	dent := &DirEntry{api: api}

	if err := api.readStruct(ptr, &dent.Packed); err != nil {
		return nil, err
	}

	dent.streamsPtr = ptr + Pointer(csize.OfMust(&dent.Packed))

	return dent, nil
}
