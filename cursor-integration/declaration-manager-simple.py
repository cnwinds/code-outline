#!/usr/bin/env python3
"""
声明管理工具 - 简化版本（无 Unicode 字符）
"""

import json
import os
import sys
import subprocess
import argparse
from pathlib import Path
from typing import Dict, List, Any, Optional
import time
from datetime import datetime

class DeclarationManager:
    def __init__(self, project_path: str = "."):
        self.project_path = Path(project_path).resolve()
        self.contextgen_path = self._find_contextgen()
        
    def _find_contextgen(self) -> Optional[Path]:
        """查找 contextgen 可执行文件"""
        current_dir = Path(__file__).parent
        project_root = current_dir.parent.parent  # 回到项目根目录
        
        possible_paths = [
            project_root / "contextgen.exe",
            project_root / "contextgen",
            project_root / "build" / "contextgen.exe",
            project_root / "build" / "contextgen",
            self.project_path / "contextgen.exe",
            self.project_path / "contextgen",
            self.project_path / "build" / "contextgen.exe",
            self.project_path / "build" / "contextgen",
        ]
        
        for path in possible_paths:
            if path.exists() and path.is_file():
                return path
                
        # 在 PATH 中查找
        try:
            result = subprocess.run(["where" if os.name == "nt" else "which", "contextgen"], 
                                 capture_output=True, text=True, timeout=5)
            if result.returncode == 0:
                return Path(result.stdout.strip().split('\n')[0])
        except:
            pass
            
        return None
    
    def get_all_declarations(self, use_cache: bool = True) -> Dict[str, Any]:
        """获取所有文件的声明内容"""
        if not self.contextgen_path:
            raise FileNotFoundError("未找到 contextgen 可执行文件")
        
        print("获取所有文件声明...")
        
        # 使用 contextgen generate 命令获取所有数据
        cmd = [
            str(self.contextgen_path),
            "generate",
            "--path", str(self.project_path),
            "--output", "temp_context.json"
        ]
        
        try:
            result = subprocess.run(cmd, capture_output=True, text=True, timeout=120)
            
            if result.returncode != 0:
                raise RuntimeError(f"获取声明失败: {result.stderr}")
            
            # 从生成的文件读取数据
            output_file = self.project_path / "temp_context.json"
            if output_file.exists():
                with open(output_file, 'r', encoding='utf-8') as f:
                    declarations = json.load(f)
                # 清理临时文件
                output_file.unlink()
            else:
                raise RuntimeError("无法生成项目上下文文件")
            
            # 增强数据结构
            enhanced_declarations = {
                "timestamp": time.strftime("%Y-%m-%d %H:%M:%S"),
                "project_path": str(self.project_path),
                "total_files": len(declarations.get("files", {})),
                "declarations": declarations,
                "summary": self._generate_summary(declarations)
            }
            
            print(f"成功获取 {enhanced_declarations['total_files']} 个文件的声明")
            return enhanced_declarations
            
        except subprocess.TimeoutExpired:
            raise RuntimeError("获取声明超时，请检查项目大小")
        except Exception as e:
            raise RuntimeError(f"获取声明失败: {e}")
    
    def get_file_declarations(self, file_path: str, use_cache: bool = True) -> Dict[str, Any]:
        """获取指定文件的声明内容"""
        file_path = Path(file_path)
        if not file_path.is_absolute():
            file_path = self.project_path / file_path
        
        if not self.contextgen_path:
            raise FileNotFoundError("未找到 contextgen 可执行文件")
        
        print(f"获取文件声明: {file_path.name}")
        
        # 使用 contextgen generate 命令获取指定文件数据
        # 注意：contextgen 不支持单文件模式，这里生成整个项目然后过滤
        cmd = [
            str(self.contextgen_path),
            "generate",
            "--path", str(self.project_path),
            "--output", "temp_file_context.json"
        ]
        
        try:
            result = subprocess.run(cmd, capture_output=True, text=True, timeout=60)
            
            if result.returncode != 0:
                raise RuntimeError(f"获取文件声明失败: {result.stderr}")
            
            # 从生成的文件读取数据并过滤指定文件
            output_file = self.project_path / "temp_file_context.json"
            if output_file.exists():
                with open(output_file, 'r', encoding='utf-8') as f:
                    all_data = json.load(f)
                
                # 过滤指定文件
                target_file_path = str(file_path.relative_to(self.project_path))
                declarations = {
                    "files": {
                        target_file_path: all_data.get("files", {}).get(target_file_path, {})
                    }
                }
                
                # 清理临时文件
                output_file.unlink()
            else:
                raise RuntimeError("无法生成项目上下文文件")
            
            # 增强数据结构
            enhanced_declarations = {
                "timestamp": time.strftime("%Y-%m-%d %H:%M:%S"),
                "file_path": str(file_path),
                "file_name": file_path.name,
                "declarations": declarations,
                "summary": self._generate_file_summary(declarations, str(file_path))
            }
            
            print(f"成功获取文件声明: {file_path.name}")
            return enhanced_declarations
            
        except subprocess.TimeoutExpired:
            raise RuntimeError("获取文件声明超时")
        except Exception as e:
            raise RuntimeError(f"获取文件声明失败: {e}")
    
    def create_project_declarations(self, output_file: str = "project_declarations.json") -> Dict[str, Any]:
        """创建整个项目的声明内容文件"""
        print("创建项目声明文件...")
        
        # 获取所有声明
        all_declarations = self.get_all_declarations()
        
        # 创建项目声明文件
        project_declarations = {
            "project_info": {
                "name": all_declarations.get("declarations", {}).get("projectName", "未知项目"),
                "path": str(self.project_path),
                "created_at": time.strftime("%Y-%m-%d %H:%M:%S"),
                "total_files": all_declarations.get("total_files", 0)
            },
            "declarations": all_declarations.get("declarations", {}),
            "summary": all_declarations.get("summary", {}),
            "file_index": self._create_file_index(all_declarations.get("declarations", {}))
        }
        
        # 保存到文件
        output_path = self.project_path / output_file
        with open(output_path, 'w', encoding='utf-8') as f:
            json.dump(project_declarations, f, ensure_ascii=False, indent=2)
        
        print(f"项目声明文件已创建: {output_path}")
        
        # 显示统计信息
        self._display_project_summary(project_declarations)
        
        return project_declarations
    
    def update_file_declarations(self, file_path: str, output_file: str = "updated_declarations.json") -> Dict[str, Any]:
        """更新指定文件的声明内容"""
        file_path = Path(file_path)
        if not file_path.is_absolute():
            file_path = self.project_path / file_path
        
        print(f"更新文件声明: {file_path.name}")
        
        # 获取文件声明
        file_declarations = self.get_file_declarations(str(file_path), use_cache=False)
        
        # 创建更新记录
        update_record = {
            "timestamp": time.strftime("%Y-%m-%d %H:%M:%S"),
            "file_path": str(file_path),
            "file_name": file_path.name,
            "action": "update",
            "declarations": file_declarations.get("declarations", {}),
            "summary": file_declarations.get("summary", {}),
            "changes": self._detect_changes(file_declarations)
        }
        
        # 保存更新记录
        output_path = self.project_path / output_file
        with open(output_path, 'w', encoding='utf-8') as f:
            json.dump(update_record, f, ensure_ascii=False, indent=2)
        
        print(f"文件声明已更新: {output_path}")
        
        # 显示更新摘要
        self._display_update_summary(update_record)
        
        return update_record
    
    def _generate_summary(self, declarations: Dict[str, Any]) -> Dict[str, Any]:
        """生成声明摘要"""
        files = declarations.get("files", {})
        
        summary = {
            "total_files": len(files),
            "total_symbols": 0,
            "languages": set(),
            "file_types": {},
            "complexity_analysis": {
                "high_complexity_files": [],
                "large_files": []
            }
        }
        
        for file_path, file_info in files.items():
            # 统计符号数量
            symbols = file_info.get("symbols", [])
            summary["total_symbols"] += len(symbols)
            
            # 统计语言
            file_ext = Path(file_path).suffix.lower()
            language = self._detect_language(file_ext)
            if language:
                summary["languages"].add(language)
            
            # 统计文件类型
            summary["file_types"][file_ext] = summary["file_types"].get(file_ext, 0) + 1
            
            # 复杂度分析
            if len(symbols) > 20:
                summary["complexity_analysis"]["high_complexity_files"].append({
                    "file": file_path,
                    "symbol_count": len(symbols)
                })
            
            # 大文件分析
            file_size = file_info.get("fileSize", 0)
            if file_size > 10000:  # 10KB
                summary["complexity_analysis"]["large_files"].append({
                    "file": file_path,
                    "size": file_size
                })
        
        summary["languages"] = list(summary["languages"])
        return summary
    
    def _generate_file_summary(self, declarations: Dict[str, Any], file_path: str) -> Dict[str, Any]:
        """生成文件摘要"""
        files = declarations.get("files", {})
        file_info = files.get(file_path, {})
        
        symbols = file_info.get("symbols", [])
        
        return {
            "file_name": Path(file_path).name,
            "total_symbols": len(symbols),
            "symbol_types": self._categorize_symbols(symbols),
            "file_size": file_info.get("fileSize", 0),
            "last_modified": file_info.get("lastModified", ""),
            "purpose": file_info.get("purpose", "")
        }
    
    def _categorize_symbols(self, symbols: List[Dict[str, Any]]) -> Dict[str, int]:
        """分类符号类型"""
        categories = {
            "functions": 0,
            "classes": 0,
            "variables": 0,
            "constants": 0,
            "other": 0
        }
        
        for symbol in symbols:
            prototype = symbol.get("prototype", "").lower()
            if "func" in prototype or "function" in prototype:
                categories["functions"] += 1
            elif "class" in prototype or "struct" in prototype:
                categories["classes"] += 1
            elif "var" in prototype or "let" in prototype or "const" in prototype:
                if "const" in prototype:
                    categories["constants"] += 1
                else:
                    categories["variables"] += 1
            else:
                categories["other"] += 1
        
        return categories
    
    def _detect_language(self, file_ext: str) -> Optional[str]:
        """检测文件语言"""
        language_map = {
            '.go': 'Go',
            '.js': 'JavaScript',
            '.jsx': 'JavaScript',
            '.ts': 'TypeScript',
            '.tsx': 'TypeScript',
            '.py': 'Python',
            '.java': 'Java',
            '.cs': 'C#',
            '.rs': 'Rust',
            '.cpp': 'C++',
            '.cc': 'C++',
            '.cxx': 'C++',
            '.hpp': 'C++',
            '.c': 'C',
            '.h': 'C'
        }
        return language_map.get(file_ext)
    
    def _create_file_index(self, declarations: Dict[str, Any]) -> Dict[str, Any]:
        """创建文件索引"""
        files = declarations.get("files", {})
        
        index = {
            "by_language": {},
            "by_complexity": {
                "simple": [],
                "medium": [],
                "complex": []
            },
            "by_purpose": {}
        }
        
        for file_path, file_info in files.items():
            # 按语言分类
            file_ext = Path(file_path).suffix.lower()
            language = self._detect_language(file_ext)
            if language:
                if language not in index["by_language"]:
                    index["by_language"][language] = []
                index["by_language"][language].append(file_path)
            
            # 按复杂度分类
            symbols = file_info.get("symbols", [])
            if len(symbols) <= 5:
                index["by_complexity"]["simple"].append(file_path)
            elif len(symbols) <= 15:
                index["by_complexity"]["medium"].append(file_path)
            else:
                index["by_complexity"]["complex"].append(file_path)
            
            # 按用途分类
            purpose = file_info.get("purpose", "")
            if "main" in purpose.lower() or "entry" in purpose.lower():
                if "main" not in index["by_purpose"]:
                    index["by_purpose"]["main"] = []
                index["by_purpose"]["main"].append(file_path)
            elif "test" in purpose.lower():
                if "test" not in index["by_purpose"]:
                    index["by_purpose"]["test"] = []
                index["by_purpose"]["test"].append(file_path)
        
        return index
    
    def _detect_changes(self, file_declarations: Dict[str, Any]) -> List[Dict[str, Any]]:
        """检测文件变化"""
        return [
            {
                "type": "declaration_update",
                "description": "文件声明已更新",
                "timestamp": time.strftime("%Y-%m-%d %H:%M:%S")
            }
        ]
    
    def _display_project_summary(self, project_declarations: Dict[str, Any]):
        """显示项目摘要"""
        print("\n项目声明摘要:")
        project_info = project_declarations["project_info"]
        summary = project_declarations["summary"]
        
        print(f"  项目名称: {project_info['name']}")
        print(f"  文件数量: {project_info['total_files']}")
        print(f"  符号总数: {summary.get('total_symbols', 0)}")
        print(f"  使用语言: {', '.join(summary.get('languages', []))}")
        print(f"  复杂文件: {len(summary.get('complexity_analysis', {}).get('high_complexity_files', []))}")
    
    def _display_update_summary(self, update_record: Dict[str, Any]):
        """显示更新摘要"""
        print(f"\n文件更新摘要:")
        print(f"  文件: {update_record['file_name']}")
        print(f"  时间: {update_record['timestamp']}")
        print(f"  操作: {update_record['action']}")
        print(f"  变化: {len(update_record['changes'])} 项")

def main():
    parser = argparse.ArgumentParser(description="声明管理工具")
    parser.add_argument("action", choices=["get-all", "get-file", "create-project", "update-file"], 
                       help="操作类型")
    parser.add_argument("--path", default=".", help="项目路径")
    parser.add_argument("--file", help="指定文件路径")
    parser.add_argument("--output", help="输出文件")
    parser.add_argument("--no-cache", action="store_true", help="不使用缓存")
    parser.add_argument("--verbose", action="store_true", help="详细输出")
    
    args = parser.parse_args()
    
    try:
        manager = DeclarationManager(args.path)
        
        if args.action == "get-all":
            result = manager.get_all_declarations(not args.no_cache)
            output_file = args.output or "all_declarations.json"
            
        elif args.action == "get-file":
            if not args.file:
                raise ValueError("必须指定文件路径")
            result = manager.get_file_declarations(args.file, not args.no_cache)
            output_file = args.output or "file_declarations.json"
            
        elif args.action == "create-project":
            result = manager.create_project_declarations(args.output or "project_declarations.json")
            output_file = args.output or "project_declarations.json"
            
        elif args.action == "update-file":
            if not args.file:
                raise ValueError("必须指定文件路径")
            result = manager.update_file_declarations(args.file, args.output or "updated_declarations.json")
            output_file = args.output or "updated_declarations.json"
        
        # 保存结果
        if args.action != "create-project" and args.action != "update-file":
            output_path = manager.project_path / output_file
            with open(output_path, 'w', encoding='utf-8') as f:
                json.dump(result, f, ensure_ascii=False, indent=2)
            print(f"结果已保存到: {output_path}")
        
        if args.verbose:
            print("\n详细结果:")
            print(json.dumps(result, ensure_ascii=False, indent=2))
        
        print(f"\n{args.action} 操作完成！")
        
    except Exception as e:
        print(f"错误: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()
