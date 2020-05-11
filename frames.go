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

// DataQuality are the names of observatories.
// Though these are names, the format describes it as data quality.
type DataQuality uint32

const (
	// TAMA300 is located at the Mitaka campus of the National Astronomical Observatory of Japan
	TAMA300 DataQuality = 1 << (2 * iota)
	// Virgo experiment at the European Gravitational Observatory.
	Virgo
	// GEO600 is located near Sarstedt in the South of Hanover, Germany.
	GEO600
	// LIGOHanford2km is on the DOE Hanford Site located near Richland, Washington.
	// During the Initial and Enhanced LIGO phases, a half-length interferometer
	// operated in parallel with the main interferometer. For this 2 km interferometer,
	// the Fabry–Pérot arm cavities had the same optical finesse, and, thus,
	// half the storage time as the 4 km interferometers
	LIGOHanford2km
	// LIGOHanford4km is on the DOE Hanford Site located near Richland, Washington.
	// This is the advanced LIGO experiment. The primary interferometer consists of
	// two beam lines of 4 km length which form a power-recycled Michelson
	// interferometer with Gires–Tournois etalon arms.
	LIGOHanford4km
	// LIGOLivingston4km located in Livingston, Louisiana.
	// LIGO Livingston Observatory houses one laser interferometer in the
	// primary configuration.
	LIGOLivingston4km
	// LIGOCaltech was a 40-meter prototype.
	LIGOCaltech
	// ALLEGRO was a ground-based, cryogenic resonant Weber bar, gravitational-wave
	// detector run by Warren Johnson, et al. at Louisiana State University in
	// Baton Rouge, Louisiana.
	ALLEGRO
	// AURIGA is an ultracryogenic resonant bar gravitational wave detector near Padua, Italy.
	AURIGA
	// EXPLORER is based in Geneva, Switzerland.  CERN, INFN.
	// EXPLORER is a cylinder of Al5056 weighing 2300 kg; it is 3 m long and it has a diameter of 60 cm
	// It is cooled at the temperature of liquid helium (4.2 K) and it operates at
	// the temperature of 2 K, which is reached by lowering the pressure on the liquid helium reservoir.
	// Its resonance frequencies are around 906 and 923 Hz.
	EXPLORER
	// NIOBE was a 1500 kg Nb antenna, located in Perth, Australia, cooled at 5K,
	// and equipped with a parametric transducer and FET amplifier
	NIOBE
	// NAUTILUS was a 2260 kg Al antenna, located in Frascati, Italy, cooled at
	// 130mK with liquid helium dilution refrigerator, and equipped with a capacitive
	// transducer and SQUID amplifier.
	NAUTILUS
)

// FileHeader is the header of a GWF file containing metadata, notably,
// version, endianness, and checksums.
type FileHeader struct {
	Magic       [5]byte  // IGWD\0 file magic
	Version     int8     // e.g. 8: https://dcc.ligo.org/public/0000/T970130/002/T970130-v2.pdf
	Minor       uint8    // Minor version of software that wrote this file; 255 means beta
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
	Length   uint64   // Length of this structure _including_ the byte count of Length.
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

// FrameHeader ...
type FrameHeader struct {
	// CommonHeader immediately proceeds
	NameLen      uint16      // Length of Name including \0
	Name         []byte      // Name of project or other experiment description (e.g., GEO; LIGO; VIRGO; TAMA;...)
	Run          int32       // Run number (number < 0 reserved for simulated data); monotonic for experimental runs.
	Frame        uint32      // Frame number, monotonically increasing until end of run, re-starting from 0 with each new run.
	DataQuality  DataQuality // A logical 32-bit word to denote top level quality of data. Lowest order bits are reserved in pairs for the various GW detectors
	StartGPS     uint32      // Frame start time in GPS Seconds.
	Residual     uint32      // Frame start time residual, integer nanoseconds.
	LeapSeconds  uint16      // Number of leap seconds between GPS and UTC.
	FrameSeconds float64     // Frame length in seconds.
	// TODO(goller): learn about dem PTR_STRUCT
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

// TOCHeader contains the amount of data and frames within the TableOfContents.
type TOCHeader struct {
	CommonHeader
	Seconds int16  // From the first FrameH in the file ; TODO(goller): unknown what this is
	Frames  uint32 // Number of frames in this file.
}

// TableOfContents enables indexing of key structures.
type TableOfContents struct {
	TOCHeader
	DataQuality []uint32
	GTimeS      []uint32
	GTimeN      []uint32
	DT          []float64
}
