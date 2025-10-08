#!/usr/bin/env python3
"""
code-outline 声明管理工具安装脚本
"""

import json
import os
import sys
import shutil
from pathlib import Path
import platform

class SpecKitInstaller:
    def __init__(self):
        self.system = platform.system().lower()
        self.cursor_config_dir = self._get_cursor_config_dir()
        self.tools_dir = Path(__file__).parent
        
    def _get_cursor_config_dir(self) -> Path:
        """获取 Cursor 配置目录"""
        if self.system == "windows":
            config_dir = Path.home() / "AppData" / "Roaming" / "Cursor" / "User"
        elif self.system == "darwin":  # macOS
            config_dir = Path.home() / "Library" / "Application Support" / "Cursor" / "User"
        else:  # Linux
            config_dir = Path.home() / ".config" / "Cursor" / "User"
        
        return config_dir
    
    def check_cursor_installation(self) -> bool:
        """检查 Cursor 是否已安装"""
        return self.cursor_config_dir.exists()
    
    def install_tools(self) -> bool:
        """安装工具到 Cursor"""
        try:
            print("🚀 开始安装 code-outline 声明管理工具...")
            
            # 检查 Cursor 安装
            if not self.check_cursor_installation():
                print("❌ 未找到 Cursor 安装，请先安装 Cursor 编辑器")
                return False
            
            # 创建配置目录
            external_tools_dir = self.cursor_config_dir / "externalTools"
            external_tools_dir.mkdir(exist_ok=True)
            
            # 复制工具文件
            tools_to_copy = [
                "cursor-code-outline.json"
            ]
            
            for tool in tools_to_copy:
                src = self.tools_dir / tool
                dst = external_tools_dir / tool
                
                if src.exists():
                    shutil.copy2(src, dst)
                    print(f"✅ 已复制: {tool}")
                else:
                    print(f"⚠️  文件不存在: {tool}")
            
            # 创建启动脚本
            self._create_launcher_scripts(external_tools_dir)
            
            # 更新 Cursor 配置
            self._update_cursor_config()
            
            print("✅ 工具安装完成！")
            print("💡 请重启 Cursor 编辑器以加载新工具")
            
            return True
            
        except Exception as e:
            print(f"❌ 安装失败: {e}")
            return False
    
    def _create_launcher_scripts(self, tools_dir: Path):
        """创建启动脚本"""
        # 不再需要创建启动脚本，因为工具直接通过 Cursor 的 External Tools 功能调用
        print("✅ 工具配置已准备就绪")
    
    def _update_cursor_config(self):
        """更新 Cursor 配置"""
        try:
            # 读取现有配置
            config_file = self.cursor_config_dir / "settings.json"
            if config_file.exists():
                with open(config_file, 'r', encoding='utf-8') as f:
                    config = json.load(f)
            else:
                config = {}
            
            # 添加外部工具配置
            if "externalTools" not in config:
                config["externalTools"] = {}
            
            # 读取工具配置
            tools_config_file = self.tools_dir / "cursor-code-outline.json"
            if tools_config_file.exists():
                with open(tools_config_file, 'r', encoding='utf-8') as f:
                    tools_config = json.load(f)
                
                # 合并工具配置
                for tool in tools_config.get("tools", []):
                    tool_name = tool["name"]
                    config["externalTools"][tool_name] = {
                        "command": tool["command"],
                        "args": tool["args"],
                        "cwd": tool["cwd"],
                        "category": tool["category"]
                    }
                
                # 保存配置
                with open(config_file, 'w', encoding='utf-8') as f:
                    json.dump(config, f, ensure_ascii=False, indent=2)
                
                print("✅ 已更新 Cursor 配置")
            
        except Exception as e:
            print(f"⚠️  配置更新失败: {e}")
    
    def uninstall_tools(self) -> bool:
        """卸载工具"""
        try:
            print("🗑️  开始卸载工具...")
            
            external_tools_dir = self.cursor_config_dir / "externalTools"
            if external_tools_dir.exists():
                # 删除工具文件
                tools_to_remove = [
                    "cursor-code-outline.json"
                ]
                
                for tool in tools_to_remove:
                    tool_file = external_tools_dir / tool
                    if tool_file.exists():
                        tool_file.unlink()
                        print(f"✅ 已删除: {tool}")
                
                # 清理空目录
                try:
                    external_tools_dir.rmdir()
                    print("✅ 已清理配置目录")
                except OSError:
                    pass  # 目录不为空，忽略
            
            print("✅ 工具卸载完成！")
            return True
            
        except Exception as e:
            print(f"❌ 卸载失败: {e}")
            return False
    
    def check_installation(self) -> bool:
        """检查安装状态"""
        external_tools_dir = self.cursor_config_dir / "externalTools"
        config_file = external_tools_dir / "cursor-code-outline.json"
        
        if config_file.exists():
            print("✅ 工具已安装")
            return True
        else:
            print("❌ 工具未安装")
            return False

def main():
    import argparse
    
    parser = argparse.ArgumentParser(description="code-outline 声明管理工具安装器")
    parser.add_argument("action", choices=["install", "uninstall", "check"], 
                       help="操作类型")
    parser.add_argument("--verbose", action="store_true", help="详细输出")
    
    args = parser.parse_args()
    
    installer = SpecKitInstaller()
    
    if args.action == "install":
        success = installer.install_tools()
        sys.exit(0 if success else 1)
        
    elif args.action == "uninstall":
        success = installer.uninstall_tools()
        sys.exit(0 if success else 1)
        
    elif args.action == "check":
        success = installer.check_installation()
        sys.exit(0 if success else 1)

if __name__ == "__main__":
    main()
