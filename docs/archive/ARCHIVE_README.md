# 文档归档说明

此目录用于记录已归档的旧文档信息，这些文档已被新文档替代。

## 归档历史

### 2025年6月17日 - 文档结构重组

在 v3.2.1 版本中，项目进行了全面的文档结构重组，以下旧文档已被新文档替代：

| 已归档文档 | 状态 | 替代文档 |
|-----------|------|---------|
| DOCUMENT_UPDATES.md | 已删除 | docs/USER_GUIDE.md, docs/API.md |
| requirements-document.md | 已删除 | docs/REQUIREMENTS.md |
| HARDCODING_GUIDELINES.md | 已删除 | docs/DEVELOPMENT.md |
| TECHNICAL_UPDATES.md | 已删除 | docs/CHANGELOG.md |
| configuration-test-plan.md | 已删除 | docs/CONFIG.md |
| placeholders.md | 已删除 | docs/USER_GUIDE.md |
| README.md | 已删除 | 根目录 README.md |

### 2025年6月17日 - 清理过期文档

清理了 docs 目录中的过期文档文件：

| 已删除文件 | 状态 | 说明 |
|-----------|------|------|
| docs/DOCUMENT_UPDATES.md | 已删除 | 内容过时，引用已不存在的 backup 目录 |
| docs/requirements-document.md | 已删除 | 内容与 DOCUMENT_UPDATES.md 完全相同，已过期 |

## 当前文档结构

```
dir-monitor-go/
├── README.md                 # 项目概述和快速开始
├── CHANGELOG.md             # 版本更新历史
├── LICENSE                  # 开源许可证
└── docs/                   # 文档目录
    ├── USER_GUIDE.md       # 用户使用指南
    ├── CONFIG.md           # 配置参考
    ├── CONFIG_EXAMPLE.md   # 配置示例
    ├── API.md              # API文档
    ├── DEVELOPMENT.md      # 开发指南
    ├── DEPLOYMENT.md       # 部署指南
    ├── FAQ.md              # 常见问题
    ├── DESIGN.md           # 设计文档
    ├── CODE_REVIEW.md      # 代码审查报告
    ├── INDEX.md            # 文档索引
    └── archive/            # 文档归档
        └── ARCHIVE_README.md # 本文件
```

## 注意事项

- 归档文档中的信息可能已过时，不建议参考
- 请使用根目录和docs目录下的新文档获取最新信息
- 如需查看历史变更，请参考CHANGELOG.md

## 归档策略

- 保留归档说明文档，记录文档变更历史
- 删除过期的文档文件，减少项目体积
- 在归档说明中记录文档映射关系，便于追溯