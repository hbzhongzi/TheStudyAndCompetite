# 迁移指南（Migration Guide）

本文件为 `go-backend/docs` 中迁移相关文档的合并与索引页，聚合了迁移流程、手动步骤、MySQL8 兼容性与修复记录。

已合并/来源文件（详情已归档）：

- `FINAL_MIGRATION_GUIDE.md`（详见归档）
- `MANUAL_MIGRATION_GUIDE.md`（详见归档）
- `MIGRATION_ISSUE_SUMMARY.md`（详见归档）
- `MYSQL8_MIGRATION_GUIDE.md`（详见归档）
- `TEACHER_ID_MIGRATION_FIX.md`（详见归档）

快速指引：

1. 备份现有数据库并导出数据快照。
2. 在测试环境执行 `sql/update_competition_database_complete.sql`（或 `migrations/` 下相关脚本）。
3. 运行迁移脚本并验证数据完整性；如遇问题参见归档的 `MIGRATION_ISSUE_SUMMARY.md`。
4. 测试应用行为（用户登录、项目查询、教师绑定等核心功能）。

归档位置与恢复：详细步骤已复制到 `go-backend/backups/docs_archive/`，如需恢复原始内容可在归档目录中找到完整文档。

---

合并的归档内容（来源汇总）

来源：`go-backend/backups/docs_archive/FINAL_MIGRATION_GUIDE.md`

```
`FINAL_MIGRATION_GUIDE.md` 的原始详细内容已归档于 backups，具体迁移步骤与脚本请参考归档副本。
```

来源：`go-backend/backups/docs_archive/MANUAL_MIGRATION_GUIDE.md`

```
`MANUAL_MIGRATION_GUIDE.md` 的原始详细内容已归档于 backups，包含在自动迁移失败时的手动步骤。
```

来源：`go-backend/backups/docs_archive/MIGRATION_ISSUE_SUMMARY.md`

```
`MIGRATION_ISSUE_SUMMARY.md` 已归档，列出迁移时的已知问题、可用脚本与排查方法。
```

来源：`go-backend/backups/docs_archive/MYSQL8_MIGRATION_GUIDE.md`

```
`MYSQL8_MIGRATION_GUIDE.md` 已归档，包含 MySQL 8.0 专用的迁移脚本说明与注意事项。
```

来源：`go-backend/backups/docs_archive/TEACHER_ID_MIGRATION_FIX.md`

```
`TEACHER_ID_MIGRATION_FIX.md` 已归档，包含为 `projects` 表添加 `teacher_id` 的 SQL 示例与验证步骤。
```

注：归档副本中保存了原文摘录或完整内容（视文件而定）。如需将其中任何档案的完整正文并入本文件，我可逐个提取并插入完整段落。
