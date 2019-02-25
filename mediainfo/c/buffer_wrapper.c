#include <unistd.h>
#include <stddef.h>
#include <MediaInfoDLL/MediaInfoDLL.h>

void *buffer_c_init() {
	void *mi;

	MediaInfoDLL_Load();
	if (!MediaInfoDLL_IsLoaded()) {
		return NULL;
	}

	mi = MediaInfo_New();
	if (!mi) {
		return NULL;
	}

	return mi;
}
//
//int buffer_c_read(void *mi, int buffer, size_t size) {
//	if (!buffer_c_open(mi, size, 0)) {
//		return 0;
//	}
//
//	unsigned char data[1024];
//	ssize_t len;
//
//	while ((len = read(buffer, data, 1024)) > 0) {
//		size_t cont = MediaInfo_Open_Buffer_Continue(mi, data, len);
//		if (cont != 0) {
//			break;
//		}
//	}
//
//	return 1;
//}

size_t buffer_c_open(void *mi, size_t size, size_t offset) {
	MediaInfo_Option(mi, "File_IsSeekable", "0");
	return MediaInfo_Open_Buffer_Init(mi, size, offset);
}

size_t buffer_c_continue(void *mi, unsigned char *data, size_t len) {
	return MediaInfo_Open_Buffer_Continue(mi, data, len);
}

char *buffer_c_get(void *mi, char *key, size_t stream, enum MediaInfo_stream_t type) {
	return (char *) MediaInfo_Get(mi, type, stream, key,
								  MediaInfo_Info_Text, MediaInfo_Info_Name);
}

char *buffer_c_option(void *mi, char *key, char *value) {
	return (char *) MediaInfo_Option(mi, key, value);
}

char *buffer_c_inform(void *mi, size_t stream) {
	return (char *) MediaInfo_Inform(mi, stream);
}

void buffer_c_close(void *mi) {
	MediaInfo_Open_Buffer_Finalize(mi);
}
