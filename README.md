# game-frame-probe
a game framework architecture

## 简介
- Actor

使用Actor作为服务间网络通信框架, 优点是弱化多服务器节点之间的远程调用概念, 任何消息发送只需要知道PID就行了;
任何消息都是异步的,且在同一个Actor中消息非并发(队列), actor真是解决高并发的好模型;

- 服务发现

支持架构多种多个服务器的集群, 使用etcd做服务发现, 支持Gate后多个类型的Node, 每个类型还可以有多个Node, 目前在选取节点的时候使用随机的策略; 

- Gate处理所有请求

有gate转发所有来自用户和服务器的消息, 需要配置路由根据CMD分发服务器; gate服务器需要从服务发现获取所有的ServerNode; 

- Node

每个ServerNode都有一个ID和类型, 在客户端发送消息来的时候, Gate会选举一个Node与客户端通信; 
Node不需要拥有其他Node的服务, 任何外部通信都需要经过Gate, 由Gate分发, 这样做的好处是: 只需要Gate管理所有Node, 更稳定的服务发现; 
服务器之间弱化长连接概念, 消息只管发送, 而不关心是否已连接, 类似于UDP; 这样方便任何服务器随时上线下不会影响Gate对其的通信, 只要服务器上线就一定能收到消息;

- Hubs底层网络框架

Gate与客户端连接时候Hubs框架, 支持多个网络协议 如Tcp,Ws ; 
Hubs框架支持优雅重启,为Gate未来实现优雅重启带来方便;

- Game服务器

每一个游戏房间一个Actor, Actor让并发更简单;
Game服务器需要管理所有的Actor, 负责分配Actor去服务客户端; 当actor中没有玩家在线的时候, 此actor会自动销毁;
