/*
 * Go doesn't like the function pointers, like in
 * MediaInfo's headers, so we need to provide a
 * tiny wrapper.
 * */

#include <stddef.h>
#include <MediaInfoDLL/MediaInfoDLL.h>

void mediainfo_c_init()
{
    MediaInfoDLL_Load();
}

void *mediainfo_c_open(char *filename)
{
    void *handle;
    int mret;

    handle = MediaInfo_New();
    if (!handle)
        return NULL;

    mret = MediaInfo_Open(handle, filename);
    if (!mret)
        return NULL;

    return handle;
}

char *mediainfo_c_get(void *opaque, char *key,
                      size_t stream, enum MediaInfo_stream_t type)
{
    return (char *) MediaInfo_Get(opaque, type, stream, key,
                                  MediaInfo_Info_Text, MediaInfo_Info_Name);
}

char *mediainfo_c_option(void *opaque, char *key, char *value)
{
    return (char *) MediaInfo_Option(opaque, key, value);
}

char *mediainfo_c_inform(void *opaque, size_t stream)
{
    return (char *) MediaInfo_Inform(opaque, stream);
}

void mediainfo_c_close(void *opaque)
{
    MediaInfo_Close(opaque);
}
