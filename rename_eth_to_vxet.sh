#!/bin/bash

echo "开始将eth目录重命名为vxet目录..."

# 确保vxet目录存在
mkdir -p vxet

# 复制文件并保持目录结构
echo "正在复制文件和目录结构..."
for file in $(find eth -type f -name "*.go"); do
  target_dir="vxet/$(dirname $file | cut -d/ -f2-)"
  mkdir -p "$target_dir"
  cp "$file" "$target_dir/$(basename $file)"
  echo "复制: $file 到 $target_dir/$(basename $file)"
done

echo "复制完成。"

echo ""
echo "重要提示："
echo "1. 所有文件已复制，但需要手动修改导入路径和包声明。"
echo "2. 请参考RENAME_ETH_TO_VXET.md文件获取更多信息。"
echo ""

# 创建一个简单的查找和替换脚本
cat > update_imports.sh << 'EOF'
#!/bin/bash

# 更新所有Go文件中的导入路径
find vxet -type f -name "*.go" | xargs sed -i 's|github.com/ethereum/go-ethereum/eth|github.com/VXETChain/VXET-SDK/vxet|g'
find vxet -type f -name "*.go" | xargs sed -i 's|github.com/ethereum/go-ethereum|github.com/VXETChain/VXET-SDK|g'
find vxet -type f -name "*.go" | xargs sed -i 's|^package eth|package vxet|g'

# 添加VXET版权声明
find vxet -type f -name "*.go" | xargs sed -i '1s|^|// Copyright 2023-2024 The VXET-SDK Authors\n|'
find vxet -type f -name "*.go" | xargs sed -i 's|// This file is part of the go-ethereum library.|// This file is part of the VXET-SDK library.|g'

echo "导入路径和包声明已更新。"
EOF

chmod +x update_imports.sh
echo "已创建update_imports.sh脚本，可用于自动更新导入路径和包声明。" 