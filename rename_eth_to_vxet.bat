@echo off
echo 开始将eth目录重命名为vxet目录...

REM 确保vxet目录存在
if not exist "vxet" mkdir vxet

REM 复制主要目录结构
echo 正在复制目录结构...
for /d %%d in (eth\*) do (
    if not exist "vxet\%%~nd" mkdir "vxet\%%~nd"
    echo 创建目录: vxet\%%~nd
)

REM 复制文件
echo 正在复制文件...
for /r eth %%f in (*.go) do (
    set "relpath=%%f"
    setlocal enabledelayedexpansion
    set "relpath=!relpath:eth\=vxet\!"
    echo 复制: %%f 到 !relpath!
    copy "%%f" "!relpath!" > nul
    endlocal
)

echo 复制完成。

echo.
echo 重要提示：
echo 1. 所有文件已复制，但需要手动修改导入路径和包声明。
echo 2. 请参考RENAME_ETH_TO_VXET.md文件获取更多信息。
echo.

pause 