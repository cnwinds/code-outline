---
description: 增量更新项目的代码结构概览，检测文件变更并只重新解析已修改的文件，提高更新效率。
scripts:
  sh: ./build/code-outline update --compact --path . --files $FILES
  ps: .\build\code-outline.exe update --compact --path . --files $FILES
---
  
输入参数：
$FILES - 指定要更新的文件，用逗号分隔（如：to/path/file1.go,path/file2.js）相对于项目根目录的路径