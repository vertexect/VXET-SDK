# PowerShell脚本，用于更新vxet目录中的Go文件

Write-Host "开始更新vxet目录中的Go文件..." -ForegroundColor Green

# 检查vxet目录是否存在
if (-not (Test-Path -Path "vxet")) {
    Write-Host "错误: vxet目录不存在，请先运行rename_eth_to_vxet.bat" -ForegroundColor Red
    exit 1
}

# 更新所有Go文件
$files = Get-ChildItem -Path "vxet" -Filter "*.go" -Recurse

foreach ($file in $files) {
    Write-Host "处理文件: $($file.FullName)" -ForegroundColor Yellow
    
    # 读取文件内容
    $content = Get-Content -Path $file.FullName -Raw
    
    # 更新导入路径
    $content = $content -replace 'github.com/ethereum/go-ethereum/eth', 'github.com/VXETChain/VXET-SDK/vxet'
    $content = $content -replace 'github.com/ethereum/go-ethereum', 'github.com/VXETChain/VXET-SDK'
    
    # 更新包声明
    $content = $content -replace '^package eth', 'package vxet'
    
    # 添加VXET版权声明
    if (-not ($content -match 'Copyright 2023-2024 The VXET-SDK Authors')) {
        $content = "// Copyright 2023-2024 The VXET-SDK Authors`n$content"
    }
    
    # 更新库声明
    $content = $content -replace '// This file is part of the go-ethereum library.', '// This file is part of the VXET-SDK library.'
    
    # 写回文件
    $content | Set-Content -Path $file.FullName
}

Write-Host "更新完成。" -ForegroundColor Green
Write-Host "注意: 可能还需要手动检查和修复一些特定的引用。" -ForegroundColor Yellow 