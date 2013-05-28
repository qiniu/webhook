Github.com webhook tools
==

[![Build Status](https://travis-ci.org/qiniu/webhook.png?branch=master)](https://travis-ci.org/qiniu/webhook)

[![Qiniu Logo](http://qiniutek.com/images/logo-2.png)](http://qiniu.com/)

基于 Github webhook 做了个小工具。

它能够做什么？简单来说，它就是一个让 Github repo 在某个分支发生 push 行为的时候，自动触发一段脚本。

# 用法

```
go run webhook.go xxx.conf
```

其中 conf 文件格式如下：

```
{
    "bind": ":9876",
    "items": [
    {
        "repo": "https://github.com/qiniu/docs.qiniu.com",
        "branch": "master",
        "script": "update-qiniu-docs.sh"
    }
]}
```

这个样例是真实的。它设置了 1 个 hook 脚本，在 https://github.com/qiniu/docs.qiniu.com 这个 repo 的 master 有变化时，自动执行 update-qiniu-docs.sh 脚本。

