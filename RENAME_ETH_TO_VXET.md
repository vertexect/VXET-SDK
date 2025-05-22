# 将eth包重命名为vxet包

本文档描述了将`eth`包重命名为`vxet`包的过程。这是VXET-SDK项目从以太坊代码库转变为VXET专用代码库的一部分。

## 已完成的工作

1. 创建了基本的目录结构：
   - `vxet/` - 主包目录
   - `vxet/ethconfig/` - 配置目录
   - `vxet/downloader/` - 区块下载器目录
   - `vxet/gasprice/` - Gas价格估算器目录
   - `vxet/protocols/` - 网络协议目录

2. 修改了主要文件：
   - `interfaces.go` - 将`ethereum`包改为`vxetchain`包
   - `vxet/backend.go` - 创建并修改了主要的后端服务文件
   - `vxet/ethconfig/config.go` - 创建并修改了配置文件

## 待完成工作

1. 复制并修改`eth/`目录下的所有文件到`vxet/`目录：
   ```
   for file in $(find eth -type f -name "*.go"); do
     mkdir -p "vxet/$(dirname $file | cut -d/ -f2-)"
     cp "$file" "vxet/$(dirname $file | cut -d/ -f2-)/$(basename $file)"
   done
   ```

2. 修改所有Go文件中的导入路径：
   - 将`github.com/ethereum/go-ethereum/eth`改为`github.com/VXETChain/VXET-SDK/vxet`
   - 将`github.com/ethereum/go-ethereum`改为`github.com/VXETChain/VXET-SDK`
   - 将包声明从`package eth`改为`package vxet`
   - 更新版权声明，添加`// Copyright 2023-2024 The VXET-SDK Authors`

3. 更新所有引用`eth`包的其他包：
   - `cmd/`目录下的命令行工具
   - `node/`目录下的节点服务
   - `internal/`目录下的内部包

4. 修改所有文档和注释，将"Ethereum"替换为"VXET"

## 注意事项

1. 这是一个大规模的重命名操作，可能会影响许多文件和包。
2. 需要仔细测试所有功能，确保重命名后的代码能够正常工作。
3. 可能需要更新构建脚本和测试用例。

## 完成后的验证

1. 确保所有导入路径正确
2. 编译整个项目，确保没有编译错误
3. 运行测试套件，确保所有测试通过
4. 启动节点，确保基本功能正常 