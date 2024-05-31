# awesome-proxy

发现Kubernetes一个很有意思的功能，就是通过代理来访问集群内部的服务。

比如说有个nginx服务，部署清单如下：

```yaml
apiVersion: v1
kind: Service
metadata:
  name: nginx
  namespaces: default
spec:
  ports:
    - name: http
      port: 80
      protocol: TCP
      targetPort: 80
  selector:
    app: nginx
  type: ClusterIP
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
  namespace: default
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - image: nginx:alpine
          name: nginx
          ports:
            - containerPort: 80
              name: http
              protocol: TCP
```

在本地通过`kubectl proxy`启动代理服务。

然后就可以通过Kubernetes Service Proxy访问这个服务了。

`http:localhost:8001/api/v1/namespaces/default/services/http:nginx:80/proxy/`

非常有意思的是，这个代理会修改HTTP的响应数据，会将响应数据中和静态资源相关的链接进行重写.

举个例子：

假设HTML源码有这么一行引入CSS资源的代码：`<link rel="stylesheet" href="/normalize.css">`。

经过代理服务处理后，会变成：`<link rel="stylesheet" href="/api/v1/namespaces/default/services/http:nginx:80/proxy/normalize.css">`

也就是说对静态资源（准确来说是绝对路径的静态资源）添加api代理的前缀。

至于哪些标签会被处理，源码是这么定义的：

```go
var atomsToAttrs = map[atom.Atom]sets.String{
	atom.A:          sets.NewString("href"),
	atom.Applet:     sets.NewString("codebase"),
	atom.Area:       sets.NewString("href"),
	atom.Audio:      sets.NewString("src"),
	atom.Base:       sets.NewString("href"),
	atom.Blockquote: sets.NewString("cite"),
	atom.Body:       sets.NewString("background"),
	atom.Button:     sets.NewString("formaction"),
	atom.Command:    sets.NewString("icon"),
	atom.Del:        sets.NewString("cite"),
	atom.Embed:      sets.NewString("src"),
	atom.Form:       sets.NewString("action"),
	atom.Frame:      sets.NewString("longdesc", "src"),
	atom.Head:       sets.NewString("profile"),
	atom.Html:       sets.NewString("manifest"),
	atom.Iframe:     sets.NewString("longdesc", "src"),
	atom.Img:        sets.NewString("longdesc", "src", "usemap"),
	atom.Input:      sets.NewString("src", "usemap", "formaction"),
	atom.Ins:        sets.NewString("cite"),
	atom.Link:       sets.NewString("href"),
	atom.Object:     sets.NewString("classid", "codebase", "data", "usemap"),
	atom.Q:          sets.NewString("cite"),
	atom.Script:     sets.NewString("src"),
	atom.Source:     sets.NewString("src"),
	atom.Video:      sets.NewString("poster", "src"),

	// TODO: css URLs hidden in style elements.
}
```

## 项目介绍

在我日常的运维工作中，经常会遇到修改web软件的访问路径的需求。

- 很多web软件，其静态资源是写绝对路径的，而且不支持配置上下文路径（Context Path） 或者叫 基路径（Base Path），而我们又想给这个web配置一个访问路径，如`/xxx-admin`。
- 这个项目就可以轻易实现上诉需求。

因此实现一个类似的功能。其实就是搬了Kubernetes Service Proxy的代码过来，然后做了一些修改。

核心代码在Kubernetes源码中：

- staging/src/k8s.io/apimachinery/pkg/util/proxy/transport.go

### 配置

```yaml
server:
  listen: :8080 # 监听端口
proxy:
  - location: /test/ # 代理路径
    proxyPass: http://127.0.0.1:8000/ # 代理后端服务

```

### 运行

```bash
awesome-proxy --config proxy.yaml
```

### 局限

如果是js拼接的地址，这个代理服务就无能为力了。