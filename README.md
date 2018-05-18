# go-ipfs-mobile
Version of ipfs mobile

go-ipfs-mobile是go-ipfs的移动版，当前版本并没有提供在手机上的文件存储服务，仅提供用于在手机上获取ipfs数据。

当前仅提供了android版功能，期待ios...

工程：

ipfs-mobile-lib
sdk库与接口，提供接口Api_InitNode，Api_CloseNode， Api_Get, Api_Catching

Api_InitNode用于初始化ipfs节点。

Api_CloseNode在应用关闭时使用，清除环境

Api_Get用于从ipfs获取文件数据并保存

Api_Catching用于从ipfs获取文件数据到内存

note:

1.
手机上运行ipfs节点，需要可读写目录权限，在gomobile代码里，仅提供了getCacheDir()的内部存储路径，我不喜欢，所以增加了path模块用于获取外部存储路径，
但是gomobile提供的RunOnJVM只在internal内部使用，所以只能给gomobile的app.go中增加了App_RunOnJVM，方便在外部调用，所以如果你要编译go-ipfs-mobile，需要在gomobile的app.go代码中把这个函数加上去

func App_RunOnJVM(fn func(vm, env, ctx uintptr) error) error {
	return mobileinit.RunOnJVM(fn)
}

2.
example-android是一个go的安卓工程，可以编译成apk运行。这里我借鉴了gomobile的flappy示例，只是将asset读取sprite.png换成了从ipfs读取，所以如果要测试example工程，需要提前在pc运行一个ipfs节点，将flappy下的sprite.png添加到ipfs里，如下命令：
运行ipfs节点
$ ipfs daemon
添加sprite.png
$ ipfs add sprite.png


安装：

$ go get -d github.com/lemonwin798/go-ipfs-mobile

编译：

工程生成依赖于gomobile，

生成sdk库
$ gomobile bind -target=android github.com/lemonwin798/go-ipfs-mobile/ipfs-mobile-lib

生成android示例
$ gomobile build -target=android github.com/lemonwin798/go-ipfs-mobile/example-android


