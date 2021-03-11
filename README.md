# CloudChef Gateway

Gateway是[上海骞云信息科技有限公司](http://www.cloudchef.io/)研发，集成了Proxy Agent、[Guacd](https://github.com/glyptodon/guacamole-client)、[Prometheus](https://github.com/prometheus/prometheus)、[Consul](https://github.com/hashicorp/consul)开源组件，以实现跨两个隔离网络的云资源管理。

## 一、架构图

![gateway README](C:\Users\lenovo\Desktop\gateway README.png)

## 二、工作原理

在完成Gateway的安装并启动Gateway后，Proxy Agent会向SmartCMP主动发起注册（这就需要部署了Gateway的机器能够访问内网），注册请求被通过后，ProxyAgent会与SmartCMP的4993端口建立起TCP的长连接，并基于TCP通过云网关协议进行通信（所有通信都是被加密的）。通过云网关，SmartCMP就可以纳管您私有环境下的云资源。





## **三、部署建议**

我们为CentOS7.2-CentOS7.9提供了一键部署脚本，您可以在登录到SmartCMP，从基础设施-云网关管理-添加获取安装脚本完成一键部署安装。

使用其他操作系统的用户，我们向您提供了手动编译部署的文档，您可以通过该文档自行完成Gateway的部署。