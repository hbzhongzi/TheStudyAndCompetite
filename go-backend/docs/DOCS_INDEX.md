# 文档索引（DOCS_INDEX）

此文件为 `go-backend/docs` 的主索引，列出建议保留的主文档、拟合并的文档以及建议归档的文件清单和后续操作步骤。

## 建议保留为主文档（合并后为规范位置）

- `README.md` — 总览与快速入口（合并 `BACKEND_DOCUMENTATION.md` 要点）
- `API_DOCUMENTATION.md` — API 主文档（合并 `API_SUMMARY.md`、`PROJECT_API_DOCUMENTATION.md`、`API_PATH_FIX_SUMMARY.md`、以及相关增强说明）
- `STARTUP_GUIDE.md` — 启动与快速运行（合并 `QUICK_START_README.md`、`START_EXAMPLES.md`）
- `MIGRATION_GUIDE.md`（新）— 所有迁移相关内容合并（`FINAL_MIGRATION_GUIDE.md`、`MANUAL_MIGRATION_GUIDE.md`、`MYSQL8_MIGRATION_GUIDE.md`、`MIGRATION_ISSUE_SUMMARY.md` 等）
- `ADMIN_GUIDE.md`（新）— 管理后台相关（合并 `ADMIN_BACKEND_SUMMARY.md`、`ADMIN_PERMISSION_FIX.md`、`ADMIN_QUICK_FIX.md`）
- `PROJECT_GUIDE.md`（新）— 项目管理模块文档（合并 `PROJECT_MODULE_SUMMARY.md`、`PROJECT_STATUS_SUMMARY.md`、`PROJECT_STRUCTURE.md` 等）
- `TOKEN_FIXES.md`（新）— Token 与鉴权相关修复合集（合并 `TOKEN_EXPIRATION_FIX.md`、`TOKEN_FIX_COMPLETE.md`、`TOKEN_REFRESH_FIX.md`）
- `TROUBLESHOOTING_GUIDE.md` — 故障排查（保留或合并重复项）

## 建议归档（移动到 `backups/docs_archive/`）

- API 概要类重复说明：`API_SUMMARY.md`、`API_PATH_FIX_SUMMARY.md`、`UPDATED_API_WITH_NEW_FEATURES.md`、`UPDATED_PROJECT_API.md`
- 快速启动/示例重复：`QUICK_START_README.md`、`START_EXAMPLES.md`
- 历史/临时修复记录：`REFACTOR_README.md`、`ADMIN_QUICK_FIX.md`（如已合并）等
- 根目录下的 `API_ENHANCEMENT_SUMMARY.md`（如其内容已并入 `API_DOCUMENTATION.md`）

## 操作流程（变更预览阶段）

1. 生成并审阅本 `DOCS_INDEX.md` 与归档清单 `backups/docs_archive/ARCHIVE_LIST.md`（本次已生成）
2. 确认合并映射与归档名单后，由我执行：
   - 在 `docs/` 下创建合并后的新文件（如 `MIGRATION_GUIDE.md`、`ADMIN_GUIDE.md` 等），并把要点合并进去；
   - 将建议归档的文件移动到 `go-backend/backups/docs_archive/`；
   - 在归档目录生成 `ARCHIVE_LIST.md` 记录所有被移动文件及原始路径；
   - 更新 `docs/README.md`（或 `DOCS_INDEX.md`）的链接索引。

3. 操作完成后会提供变更预览（每个移动/合并步骤的具体文件列表），等待你最终确认后提交。

## 说明

本文件为预览索引，暂未对原始文件做任何修改或移动。请审阅上述保留/归档建议，并回复“开始执行”或列出你想调整的合并策略。

----
小提示：合并时我会保留每个原始文件的头部作者/日期信息（如存在），并在合并文档中注明来源文件路径。
