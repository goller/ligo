package ligo

// Library indicates the software that produced the GWF file.
type Library int8

const (
	// LibUnknown indicates the GWF file was not producted by the C or CPP frame libraries.
	LibUnknown Library = iota
	// LibFrameL is the C frame library. http://lappweb.in2p3.fr/virgo/FrameL/
	LibFrameL
	// LibFrameCPP is the CPP frame library.  http://software.ligo.org/lscsoft/source/
	LibFrameCPP
)

// Checksum indicates the type of checksum used in the GWF file.
type Checksum int8

const (
	// SumNone indicates there is no checksum.
	SumNone Checksum = iota
	// SumCRC indicates there is a POSIX.2 checksum.
	SumCRC
)

// FileHeader is the header of a GWF file containing metadata, notably,
// version, endianness, and checksums.
type FileHeader struct {
	Magic       [5]byte  // IGWD\0 file magic
	Version     int8     // e.g. 8: https://dcc.ligo.org/public/0000/T970130/002/T970130-v2.pdf
	Minor       int8     // Minor version of software that wrote this file; 255 means beta
	SizeInt16   int8     // Size of 16-bit int originating hardware
	SizeInt32   int8     // Size of 32-bit int on originating hardware
	SizeInt64   int8     // Size of 64-bit int on originating hardware
	SizeFloat32 int8     // Size of 32-bit float on originating hardware
	SizeFloat64 int8     // Size of 64-bit float on originating hardware
	Endian2     [2]byte  // 2 bytes containing 0x1234. Used to determine byte order
	Endian4     [4]byte  // 4 bytes containing 0x12345678. Used to determine byte order
	Endian8     [8]byte  // 8 bytes containing 0x123456789abcdef. Used to determine byte order
	Pi32        float32  // IEEE float representation of pi
	Pi64        float64  // IEEE double representation of pi
	Library     Library  // Library is the producers of this GWF file.
	Checksum    Checksum // Checksum describes the checksum recorded at the EOF struct.
}

// CommonHeader contains the common elements of all frame structures.
type CommonHeader struct {
	Length   uint8    // Length of this structure _including_ the byte count of Length.
	Checksum Checksum // Checksum describes the checksum recorded at the FrameFooter.
	Class    uint8    // Class of the frame
	Instance uint32   // Count of this class of structure within current frame or current file, starting from 0
}

// FileFooter describes the frames and checksums of the file.
type FileFooter struct {
	CommonHeader
	NumFrames      uint32 // NumFrames is the total number of frames in the file.
	Bytes          uint64 // Bytes is the total number of bytes in a file; 0 if not computed.
	SeekTOC        uint64 // Bytes to seek from EOF to reach the address of the table of contents; 0 if no TOC.
	HeaderChecksum uint32 // HeaderChecksum of the FileHeader; 0 if no checksum.
	FooterChecksum uint32 // FootChecksum of NumFrames, Bytes, SeekTOC, and HeaderChecksum; 0 if no checksum.
	FileChecksum   uint32 // FileChecksum of the entire file _except_ FileCheck (meaning, all but the last 8 bytes).
}

// DictHeader is a dictionary-type structure.
type DictHeader struct {
	// CommonHeader immediately proceeds with CommonHeader.Class == 1
	NameLen    uint16 // Length of Name including \0
	Name       []byte // Name of structure being described by this dictionary structure.
	Class      uint16 // Class number of structure being described
	CommentLen uint16 // Length of Comment including \0
	Comment    []byte // Comment describing the frame.
	Checksum   uint32 // Structure checksum starting with the "length" variable including Comment
}

// DictElement is a dictionary-type structure.
type DictElement struct {
	// CommonHeader immediately proceeds with CommonHeader.Class == 2
	NameLen    uint16 // Length of Name including \0
	Name       []byte // Name of structure being described by this dictionary structure.
	ClassLen   uint16 // Length of Class including \0
	Class      []byte // Literally contains “CHAR”, “INT_2U”,...
	CommentLen uint16 // Length of Comment including \0
	Comment    []byte // Comment describing the frame.
	Checksum   uint32 // Structure checksum starting with the "length" variable including Comment
}

// FrameFooter ...
type FrameFooter struct {
	CommonHeader
	Run      int32  // Run number; same as in Frame Header run number datum.
	Frame    uint32 // Frame number, monotonically increasing until end of run.
	StartGPS uint32 // Frame start time in GPS Seconds.
	Residual uint32 // Frame start time residual, integer nanoseconds.
	Checksum uint32 // Structure checksum starting with the "length" variable including Residual.
}
