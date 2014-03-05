Github/Bitbucket webhook tools
==

[![Qiniu Logo](http://qiniutek.com/images/logo-2.png)](http://qiniu.com/)

基于 Github/Bitbucket webhook 做了个小工具。

它能够做什么？简单来说，它就是一个让 Github/Bitbucket repo 在某个分支发生 push 行为的时候，自动触发一段脚本。

# 用法

```
go run webhook.go xxx.conf
```

这样就启动 webhook 服务器了。其中 conf 文件格式如下：

```
{
    "bind": ":9876",
    "items": [
    {
        "repo": "https://github.com/qiniu/docs.qiniu.com",
        "branch": "master",
        "script": "update-qiniu-docs.sh"
    },
    {
        "repo": "https://bitbucket.org/Wuvist/angelbot/",
        "branch": "master",
        "script": "restart-angelbot.sh"
    }
]}
```

这个样例是真实的。它设置了 1 个 hook 脚本，在 https://github.com/qiniu/docs.qiniu.com 这个 repo 的 master 有变化时，自动执行 update-qiniu-docs.sh 脚本。


# 与 Github 关联

在你的 repo 首页（例如 https://github.com/qiniu/docs.qiniu.com ），点 Settings，再进入 Service Hooks，再进入 WebHook URLs，这里你就可以设置你的 WebHook URL 了，比如 http://example.com:9876/ 。

配置好后，再确定 webhook 已经启动，你就可以尝试向 repo push 一些修改，看看能不能执行相应的脚本了。
