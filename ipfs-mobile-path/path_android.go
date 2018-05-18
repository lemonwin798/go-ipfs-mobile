package mobilePath

/*
#cgo LDFLAGS: -landroid -llog -lEGL -lGLESv2
#include <android/log.h>
#include <jni.h>
#include <stdlib.h>
#include <errno.h>

#define CLASS_NAME "android/app/NativeActivity"

#define LOG_INFO(...) __android_log_print(ANDROID_LOG_INFO, "Go", __VA_ARGS__)

void getExternalFilesDir(uintptr_t java_vm, uintptr_t jni_env, jobject ctx)
{
	JavaVM* vm = (JavaVM*)java_vm;
	JNIEnv* env = (JNIEnv*)jni_env;

	jclass cls_Env = (*env)->FindClass(env, CLASS_NAME );
    jmethodID mid = (*env)->GetMethodID(env, cls_Env, "getExternalFilesDir", "(Ljava/lang/String;)Ljava/io/File;" );
    jobject obj_File = (*env)->CallObjectMethod(env, ctx, mid, NULL );
    jclass cls_File = (*env)->FindClass(env, "java/io/File" );
    jmethodID mid_getPath = (*env)->GetMethodID(env, cls_File, "getPath", "()Ljava/lang/String;" );
    jstring strPath = (jstring) (*env)->CallObjectMethod(env, obj_File, mid_getPath );


   // jstring strPath = getExternalFilesDirJString( env );
    const char* path = (*env)->GetStringUTFChars(env, strPath, NULL );
    //std::string s( path );

	if (setenv("ExFileDIR", path, 1) != 0) {
		LOG_INFO("setenv(\"ExFileDIR\", \"%s\", 1) failed: %d", path, errno);
	}

    (*env)->ReleaseStringUTFChars(env, strPath, path );
    (*env)->DeleteLocalRef(env, strPath );
    //vm->DetachCurrentThread();
}
*/
import "C"
import (
	"log"
	"os"
	"unsafe"

	"golang.org/x/mobile/app"
)

func initExternStorageFilePath() error {
	err := app.App_RunOnJVM(func(vm, env, ctx uintptr) error {

		C.getExternalFilesDir(C.uintptr_t(vm), C.uintptr_t(env), C.jobject(ctx))

		n := C.CString("ExFileDIR")
		os.Setenv("ExFileDIR", C.GoString(C.getenv(n)))
		C.free(unsafe.Pointer(n))

		return nil
	})
	if err != nil {
		log.Fatalf("asset: %v", err)
	}
	return err
}

func getExternStorageFilePath() string {
	return os.Getenv("ExFileDIR") + "/"
}

func getStorageCachePath() string {
	return os.Getenv("TMPDIR") + "/"
}
