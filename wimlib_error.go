package cwimkit

import (
	"fmt"
)

type WimlibError int

func (err WimlibError) Error() string {
	return fmt.Sprintf("wimlib: error %d: %s", err, wimlibErrorString(err))
}

var (
	ErrAlreadyLocked                  = WimlibError(1)
	ErrDecompression                  = WimlibError(2)
	ErrFuse                           = WimlibError(6)
	ErrGlobHadNoMatches               = WimlibError(8)
	ErrImageCount                     = WimlibError(10)
	ErrImageNameCollision             = WimlibError(11)
	ErrInsufficientPrivileges         = WimlibError(12)
	ErrIntegrity                      = WimlibError(13)
	ErrInvalidCaptureConfig           = WimlibError(14)
	ErrInvalidChunkSize               = WimlibError(15)
	ErrInvalidCompressionType         = WimlibError(16)
	ErrInvalidHeader                  = WimlibError(17)
	ErrInvalidImage                   = WimlibError(18)
	ErrInvalidIntegrityTable          = WimlibError(19)
	ErrInvalidLookupTableEntry        = WimlibError(20)
	ErrInvalidMetadataResource        = WimlibError(21)
	ErrInvalidOverlay                 = WimlibError(23)
	ErrInvalidParam                   = WimlibError(24)
	ErrInvalidPartNumber              = WimlibError(25)
	ErrInvalidPipableWim              = WimlibError(26)
	ErrInvalidReparseData             = WimlibError(27)
	ErrInvalidResourceHash            = WimlibError(28)
	ErrInvalidUtf16String             = WimlibError(30)
	ErrInvalidUtf8String              = WimlibError(31)
	ErrIsDirectory                    = WimlibError(32)
	ErrIsSplitWim                     = WimlibError(33)
	ErrLink                           = WimlibError(35)
	ErrMetadataNotFound               = WimlibError(36)
	ErrMkdir                          = WimlibError(37)
	ErrMqueue                         = WimlibError(38)
	ErrNoMem                          = WimlibError(39)
	ErrNotDir                         = WimlibError(40)
	ErrNotEmpty                       = WimlibError(41)
	ErrNotARegularFile                = WimlibError(42)
	ErrNotAWimFile                    = WimlibError(43)
	ErrNotPipable                     = WimlibError(44)
	ErrNoFilename                     = WimlibError(45)
	ErrNtfs3G                         = WimlibError(46)
	ErrOpen                           = WimlibError(47)
	ErrOpenDir                        = WimlibError(48)
	ErrPathDoesNotExist               = WimlibError(49)
	ErrRead                           = WimlibError(50)
	ErrReadLink                       = WimlibError(51)
	ErrRename                         = WimlibError(52)
	ErrReparsePointFixupFailed        = WimlibError(54)
	ErrResourceNotFound               = WimlibError(55)
	ErrResourceOrder                  = WimlibError(56)
	ErrSetAttributes                  = WimlibError(57)
	ErrSetReparseData                 = WimlibError(58)
	ErrSetSecurity                    = WimlibError(59)
	ErrSetShortName                   = WimlibError(60)
	ErrSetTimestamps                  = WimlibError(61)
	ErrSplitInvalid                   = WimlibError(62)
	ErrStat                           = WimlibError(63)
	ErrUnexpectedEndOfFile            = WimlibError(65)
	ErrUnicodeStringNotRepresentable  = WimlibError(66)
	ErrUnknownVersion                 = WimlibError(67)
	ErrUnsupported                    = WimlibError(68)
	ErrUnsupportedFile                = WimlibError(69)
	ErrWimIsReadonly                  = WimlibError(71)
	ErrWrite                          = WimlibError(72)
	ErrXml                            = WimlibError(73)
	ErrWimIsEncrypted                 = WimlibError(74)
	ErrWIMBoot                        = WimlibError(75)
	ErrAbortedByProgress              = WimlibError(76)
	ErrUnknownProgressStatus          = WimlibError(77)
	ErrMknod                          = WimlibError(78)
	ErrMountedImageIsBusy             = WimlibError(79)
	ErrNotAMountPoint                 = WimlibError(80)
	ErrNotPermittedToUnmount          = WimlibError(81)
	ErrFveLockedVolume                = WimlibError(82)
	ErrUnableToReadCaptureConfig      = WimlibError(83)
	ErrWimIsIncomplete                = WimlibError(84)
	ErrCompactionNotPossible          = WimlibError(85)
	ErrImageHasMultipleReferences     = WimlibError(86)
	ErrDuplicateExportedImage         = WimlibError(87)
	ErrConcurrentModificationDetected = WimlibError(88)
	ErrSnapshotFailure                = WimlibError(89)
	ErrInvalidXAttr                   = WimlibError(90)
	ErrSetXAttr                       = WimlibError(91)
)

func wimlibErrorString(err WimlibError) string {
	switch err {
	case WimlibError(1):
		return "The WIM is already locked for writing"
	case WimlibError(2):
		return "The WIM contains invalid compressed data"
	case WimlibError(6):
		return "An error was returned by fuse_main()"
	case WimlibError(8):
		return "The provided file glob did not match any files"
	case WimlibError(10):
		return "Inconsistent image count among the metadata resources, the WIM header, and/or the XML data"
	case WimlibError(11):
		return "Tried to add an image with a name that is already in use"
	case WimlibError(12):
		return "The user does not have sufficient privileges"
	case WimlibError(13):
		return "The WIM file is corrupted (failed integrity check)"
	case WimlibError(14):
		return "The contents of the capture configuration file were invalid"
	case WimlibError(15):
		return "The compression chunk size was unrecognized"
	case WimlibError(16):
		return "The compression type was unrecognized"
	case WimlibError(17):
		return "The WIM header was invalid"
	case WimlibError(18):
		return "Tried to select an image that does not exist in the WIM"
	case WimlibError(19):
		return "The WIM's integrity table is invalid"
	case WimlibError(20):
		return "An entry in the WIM's lookup table is invalid"
	case WimlibError(21):
		return "The metadata resource is invalid"
	case WimlibError(23):
		return "Conflicting files in overlay when creating a WIM image"
	case WimlibError(24):
		return "An invalid parameter was given"
	case WimlibError(25):
		return "The part number or total parts of the WIM is invalid"
	case WimlibError(26):
		return "The pipable WIM is invalid"
	case WimlibError(27):
		return "The reparse data of a reparse point was invalid"
	case WimlibError(28):
		return "The SHA-1 message digest of a WIM resource did not match the expected value"
	case WimlibError(30):
		return "A string was not a valid UTF-8 string"
	case WimlibError(31):
		return "A string was not a valid UTF-16 string"
	case WimlibError(32):
		return "One of the specified paths to delete was a directory"
	case WimlibError(33):
		return "The WIM is part of a split WIM, which is not supported for this operation"
	case WimlibError(35):
		return "Failed to create a hard or symbolic link when extracting a file from the WIM"
	case WimlibError(36):
		return "The WIM does not contain image metadata; it only contains file data"
	case WimlibError(37):
		return "Failed to create a directory"
	case WimlibError(38):
		return "Failed to create or use a POSIX message queue"
	case WimlibError(39):
		return "Ran out of memory"
	case WimlibError(40):
		return "Expected a directory"
	case WimlibError(41):
		return "Directory was not empty"
	case WimlibError(42):
		return "One of the specified paths to extract did not correspond to a regular file"
	case WimlibError(43):
		return "The file did not begin with the magic characters that identify a WIM file"
	case WimlibError(44):
		return "The WIM is not identified with a filename"
	case WimlibError(45):
		return "The WIM was not captured such that it can be applied from a pipe"
	case WimlibError(46):
		return "NTFS-3G encountered an error (check errno)"
	case WimlibError(47):
		return "Failed to open a file"
	case WimlibError(48):
		return "Failed to open a directory"
	case WimlibError(49):
		return "The path does not exist in the WIM image"
	case WimlibError(50):
		return "Could not read data from a file"
	case WimlibError(51):
		return "Could not read the target of a symbolic link"
	case WimlibError(52):
		return "Could not rename a file"
	case WimlibError(54):
		return "Unable to complete reparse point fixup"
	case WimlibError(55):
		return "A file resource needed to complete the operation was missing from the WIM"
	case WimlibError(56):
		return "The components of the WIM were arranged in an unexpected order"
	case WimlibError(57):
		return "Failed to set attributes on extracted file"
	case WimlibError(58):
		return "Failed to set reparse data on extracted file"
	case WimlibError(59):
		return "Failed to set file owner, group, or other permissions on extracted file"
	case WimlibError(60):
		return "Failed to set short name on extracted file"
	case WimlibError(61):
		return "Failed to set timestamps on extracted file"
	case WimlibError(62):
		return "The WIM is part of an invalid split WIM"
	case WimlibError(63):
		return "Could not read the metadata for a file or directory"
	case WimlibError(65):
		return "Unexpectedly reached the end of the file"
	case WimlibError(66):
		return "A Unicode string could not be represented in the current locale's encoding"
	case WimlibError(67):
		return "The WIM file is marked with an unknown version number"
	case WimlibError(68):
		return "The requested operation is unsupported"
	case WimlibError(69):
		return "A file in the directory tree to archive was not of a supported type"
	case WimlibError(71):
		return "The WIM is read-only (file permissions, header flag, or split WIM)"
	case WimlibError(72):
		return "Failed to write data to a file"
	case WimlibError(73):
		return "The XML data of the WIM is invalid"
	case WimlibError(74):
		return "The WIM file (or parts of it) is encrypted"
	case WimlibError(75):
		return "Failed to set WIMBoot pointer data"
	case WimlibError(76):
		return "The operation was aborted by the library user"
	case WimlibError(77):
		return "The user-provided progress function returned an unrecognized value"
	case WimlibError(78):
		return "Unable to create a special file (e.g. device node or socket)"
	case WimlibError(79):
		return "There are still files open on the mounted WIM image"
	case WimlibError(80):
		return "There is not a WIM image mounted on the directory"
	case WimlibError(81):
		return "The current user does not have permission to unmount the WIM image"
	case WimlibError(82):
		return "The volume must be unlocked before it can be used"
	case WimlibError(83):
		return "The capture configuration file could not be read"
	case WimlibError(84):
		return "The WIM file is incomplete"
	case WimlibError(85):
		return "The WIM file cannot be compacted because of its format, its layout, or the write parameters specified by the user"
	case WimlibError(86):
		return "The WIM image cannot be modified because it is currently referenced from multiple places"
	case WimlibError(87):
		return "The destination WIM already contains one of the source images"
	case WimlibError(88):
		return "A file being added to a WIM image was concurrently modified"
	case WimlibError(89):
		return "Unable to create a filesystem snapshot"
	case WimlibError(90):
		return "An extended attribute entry in the WIM image is invalid"
	case WimlibError(91):
		return "Failed to set an extended attribute on an extracted file"
	default:
		return "unknown"
	}
}
