package cwimkit

import (
	"errors"
	"fmt"

	wasm "github.com/tetratelabs/wazero/api"
)

type vtable struct {
	callCtors wasm.Function
	callDtors wasm.Function

	free   wasm.Function
	malloc wasm.Function
	calloc wasm.Function

	fdOpen wasm.Function

	eventPortNew  wasm.Function
	eventPortFree wasm.Function

	wimlibAddEmptyImage                    wasm.Function
	wimlibAddImage                         wasm.Function
	wimlibAddImageMultisource              wasm.Function
	wimlibAddTree                          wasm.Function
	wimlibCreateNewWim                     wasm.Function
	wimlibDeleteImage                      wasm.Function
	wimlibDeletePath                       wasm.Function
	wimlibExportImage                      wasm.Function
	wimlibExtractImage                     wasm.Function
	wimlibExtractImageFromPipe             wasm.Function
	wimlibExtractImageFromPipeWithProgress wasm.Function
	wimlibExtractPathlist                  wasm.Function
	wimlibExtractPaths                     wasm.Function
	wimlibExtractXmlData                   wasm.Function
	wimlibFree                             wasm.Function
	wimlibGetCompressionTypeString         wasm.Function
	wimlibGetErrorString                   wasm.Function
	wimlibGetImageDescription              wasm.Function
	wimlibGetImageName                     wasm.Function
	wimlibGetImageProperty                 wasm.Function
	wimlibGetVersion                       wasm.Function
	wimlibGetVersionString                 wasm.Function
	wimlibGetWimInfo                       wasm.Function
	wimlibGetXmlData                       wasm.Function
	wimlibGlobalInit                       wasm.Function
	wimlibGlobalCleanup                    wasm.Function
	wimlibImageNameInUse                   wasm.Function
	wimlibIterateDirTree                   wasm.Function
	wimlibIterateLookupTable               wasm.Function
	wimlibJoin                             wasm.Function
	wimlibJoinWithProgress                 wasm.Function
	wimlibMountImage                       wasm.Function
	wimlibOpenWim                          wasm.Function
	wimlibOpenWimWithProgress              wasm.Function
	wimlibOverwrite                        wasm.Function
	wimlibPrintAvailableImages             wasm.Function
	wimlibPrintHeader                      wasm.Function
	wimlibReferenceResourceFiles           wasm.Function
	wimlibReferenceResources               wasm.Function
	wimlibReferenceTemplateImage           wasm.Function
	wimlibRegisterProgressFunction         wasm.Function
	wimlibRenamePath                       wasm.Function
	wimlibResolveImage                     wasm.Function
	wimlibSetErrorFile                     wasm.Function
	wimlibSetErrorFileByName               wasm.Function
	wimlibSetImageDescripton               wasm.Function
	wimlibSetImageFlags                    wasm.Function
	wimlibSetImageName                     wasm.Function
	wimlibSetImageProperty                 wasm.Function
	wimlibSetMemoryAllocator               wasm.Function
	wimlibSetOutputChunkSize               wasm.Function
	wimlibSetOutputPackChunkSize           wasm.Function
	wimlibSetOutputCompressionType         wasm.Function
	wimlibSetOutputPackCompressionType     wasm.Function
	wimlibSetPrintErrors                   wasm.Function
	wimlibSetWimInfo                       wasm.Function
	wimlibSplit                            wasm.Function
	wimlibVerifyWim                        wasm.Function
	wimlibUnmountImage                     wasm.Function
	wimlibUnmountImageWithProgress         wasm.Function
	wimlibUpdateImage                      wasm.Function
	wimlibWrite                            wasm.Function
	wimlibWriteToFd                        wasm.Function
	wimlibSetDefaultCompressionLevel       wasm.Function
	wimlibGetCompressorNeededMemory        wasm.Function
	wimlibCreateCompressor                 wasm.Function
	wimlibCompress                         wasm.Function
	wimlibFreeCompressor                   wasm.Function
	wimlibCreateDecompressor               wasm.Function
	wimlibDecompress                       wasm.Function
	wimlibFreeDecompressor                 wasm.Function
}

func newVTable(mod wasm.Module) (vtable, error) {
	var errs []error

	getFun := func(name string) wasm.Function {
		f := mod.ExportedFunction(name)
		if f == nil {
			errs = append(errs, fmt.Errorf("%s: %w", name, ErrNoFunction))
			return nil
		}
		return f
	}

	vt := vtable{
		callCtors: getFun("__wasm_call_ctors"),
		callDtors: getFun("__wasm_call_dtors"),

		free:   getFun("free"),
		malloc: getFun("malloc"),
		calloc: getFun("calloc"),

		fdOpen: getFun("fdopen"),

		eventPortNew:  getFun("event_port_new"),
		eventPortFree: getFun("event_port_free"),

		wimlibAddEmptyImage:                    getFun("wimlib_add_empty_image"),
		wimlibAddImage:                         getFun("wimlib_add_image"),
		wimlibAddImageMultisource:              getFun("wimlib_add_image_multisource"),
		wimlibAddTree:                          getFun("wimlib_add_tree"),
		wimlibCreateNewWim:                     getFun("wimlib_create_new_wim"),
		wimlibDeleteImage:                      getFun("wimlib_delete_image"),
		wimlibDeletePath:                       getFun("wimlib_delete_path"),
		wimlibExportImage:                      getFun("wimlib_export_image"),
		wimlibExtractImage:                     getFun("wimlib_extract_image"),
		wimlibExtractImageFromPipe:             getFun("wimlib_extract_image_from_pipe"),
		wimlibExtractImageFromPipeWithProgress: getFun("wimlib_extract_image_from_pipe_with_progress"),
		wimlibExtractPathlist:                  getFun("wimlib_extract_pathlist"),
		wimlibExtractPaths:                     getFun("wimlib_extract_paths"),
		wimlibExtractXmlData:                   getFun("wimlib_extract_xml_data"),
		wimlibFree:                             getFun("wimlib_free"),
		wimlibGetCompressionTypeString:         getFun("wimlib_get_compression_type_string"),
		wimlibGetErrorString:                   getFun("wimlib_get_error_string"),
		wimlibGetImageDescription:              getFun("wimlib_get_image_description"),
		wimlibGetImageName:                     getFun("wimlib_get_image_name"),
		wimlibGetImageProperty:                 getFun("wimlib_get_image_property"),
		wimlibGetVersion:                       getFun("wimlib_get_version"),
		wimlibGetVersionString:                 getFun("wimlib_get_version_string"),
		wimlibGetWimInfo:                       getFun("wimlib_get_wim_info"),
		wimlibGetXmlData:                       getFun("wimlib_get_xml_data"),
		wimlibGlobalInit:                       getFun("wimlib_global_init"),
		wimlibGlobalCleanup:                    getFun("wimlib_global_cleanup"),
		wimlibImageNameInUse:                   getFun("wimlib_image_name_in_use"),
		wimlibIterateDirTree:                   getFun("wimlib_iterate_dir_tree2"),
		wimlibIterateLookupTable:               getFun("wimlib_iterate_lookup_table"),
		wimlibJoin:                             getFun("wimlib_join"),
		wimlibJoinWithProgress:                 getFun("wimlib_join_with_progress"),
		wimlibMountImage:                       getFun("wimlib_mount_image"),
		wimlibOpenWim:                          getFun("wimlib_open_wim"),
		wimlibOpenWimWithProgress:              getFun("wimlib_open_wim_with_progress"),
		wimlibOverwrite:                        getFun("wimlib_overwrite"),
		wimlibPrintAvailableImages:             getFun("wimlib_print_available_images"),
		wimlibPrintHeader:                      getFun("wimlib_print_header"),
		wimlibReferenceResourceFiles:           getFun("wimlib_reference_resource_files"),
		wimlibReferenceResources:               getFun("wimlib_reference_resources"),
		wimlibReferenceTemplateImage:           getFun("wimlib_reference_template_image"),
		wimlibRegisterProgressFunction:         getFun("wimlib_register_progress_function"),
		wimlibRenamePath:                       getFun("wimlib_rename_path"),
		wimlibResolveImage:                     getFun("wimlib_resolve_image"),
		wimlibSetErrorFile:                     getFun("wimlib_set_error_file"),
		wimlibSetErrorFileByName:               getFun("wimlib_set_error_file_by_name"),
		wimlibSetImageDescripton:               getFun("wimlib_set_image_descripton"),
		wimlibSetImageFlags:                    getFun("wimlib_set_image_flags"),
		wimlibSetImageName:                     getFun("wimlib_set_image_name"),
		wimlibSetImageProperty:                 getFun("wimlib_set_image_property"),
		wimlibSetMemoryAllocator:               getFun("wimlib_set_memory_allocator"),
		wimlibSetOutputChunkSize:               getFun("wimlib_set_output_chunk_size"),
		wimlibSetOutputPackChunkSize:           getFun("wimlib_set_output_pack_chunk_size"),
		wimlibSetOutputCompressionType:         getFun("wimlib_set_output_compression_type"),
		wimlibSetOutputPackCompressionType:     getFun("wimlib_set_output_pack_compression_type"),
		wimlibSetPrintErrors:                   getFun("wimlib_set_print_errors"),
		wimlibSetWimInfo:                       getFun("wimlib_set_wim_info"),
		wimlibSplit:                            getFun("wimlib_split"),
		wimlibVerifyWim:                        getFun("wimlib_verify_wim"),
		wimlibUnmountImage:                     getFun("wimlib_unmount_image"),
		wimlibUnmountImageWithProgress:         getFun("wimlib_unmount_image_with_progress"),
		wimlibUpdateImage:                      getFun("wimlib_update_image"),
		wimlibWrite:                            getFun("wimlib_write"),
		wimlibWriteToFd:                        getFun("wimlib_write_to_fd"),
		wimlibSetDefaultCompressionLevel:       getFun("wimlib_set_default_compression_level"),
		wimlibGetCompressorNeededMemory:        getFun("wimlib_get_compressor_needed_memory"),
		wimlibCreateCompressor:                 getFun("wimlib_create_compressor"),
		wimlibCompress:                         getFun("wimlib_compress"),
		wimlibFreeCompressor:                   getFun("wimlib_free_compressor"),
		wimlibCreateDecompressor:               getFun("wimlib_create_decompressor"),
		wimlibDecompress:                       getFun("wimlib_decompress"),
		wimlibFreeDecompressor:                 getFun("wimlib_free_decompressor"),
	}
	if len(errs) > 0 {
		return vtable{}, errors.Join(errs...)
	}
	return vt, nil
}
