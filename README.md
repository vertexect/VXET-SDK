# VXET-SDK

VXET-SDK 是 VXET 区块链的官方开发工具包，提供了与 VXET 区块链交互的完整工具集。

## 功能特点

- 完整的 VXET 区块链客户端实现
- 多种命令行工具，便于开发和调试
- 丰富的 API 接口，方便应用开发
- 灵活的共识机制支持
- 完善的 RPC 服务

## 系统要求

- Go 1.20 或更高版本
- 适用于 Linux、MacOS 和 Windows 系统

## 安装指南

### 从源码构建

1. 克隆仓库

```bash
git clone https://github.com/VXETChain/VXET-SDK.git
cd VXET-SDK
```

2. 编译 VXET 客户端

```bash
make vxet
```

编译完成后，可执行文件将位于 `build/bin/` 目录中。

### 使用预编译版本

请访问 [发布页面](https://github.com/VXETChain/VXET-SDK/releases) 下载适合您系统的预编译版本。

## 使用方法

### 启动 VXET 客户端

```bash
./build/bin/vxet
```

### 开发工具

安装开发所需工具：

```bash
make devtools
```

### 运行测试

```bash
make test
```

## 项目结构

- `cmd/`: 包含各种命令行工具
  - `vxet/`: VXET 客户端
  - `abigen/`: 智能合约 ABI 生成工具
  - 其他实用工具
- `common/`: 通用工具和数据结构
  - `hexutil/`: 十六进制工具
  - `math/`: 数学运算工具
  - `prque/`: 优先队列实现
- `consensus/`: 共识算法实现
  - `clique/`: PoA 共识算法
- `vxetclient/`: 客户端实现
  - `gethclient/`: 以太坊客户端兼容接口
- `vxet/`: 核心区块链实现
  - `protocols/`: 网络协议实现
  - `gasprice/`: Gas 价格管理
  - `downloader/`: 区块下载器
- `params/`: 网络参数配置
- `rpc/`: RPC 服务实现
- `p2p/`: P2P 网络层实现
- `vxetdb/`: 区块链数据库
- `node/`: 节点服务
- `console/`: 控制台接口
- `graphql/`: GraphQL API 实现
- `signer/`: 交易签名工具

## 许可证

本项目采用 GPL-3.0 许可证，详情请参阅 [COPYING](COPYING) 文件。

## 贡献指南

欢迎提交 Pull Requests 和 Issues。在提交代码前，请确保通过所有测试并遵循项目的代码规范。

## 联系方式

- 官方网站：[VXET Chain](https://www.vertexect.com/)
- 问题反馈：请在 GitHub 上提交 Issue 