// Package model 定义持久化实体。
//
// 与前端、Handler 列表/详情 JSON 统一约定：
//
//	createTime — 创建时间
//	updateTime — 更改时间
//	createdBy  — 创建人用户 ID（无则省略）
//	status     — 状态（整型等依业务；system_logs.status 为请求结果等字符串语义）
//
// 软删除使用 GORM 的 DeletedAt，默认不参与 JSON。需要对外展示删除时间时用字段名 deleteTime，在 Handler 中格式化。
//
// 可嵌入复用块见 base_model.go：AuditModel、TenantScoped、CreatorOptional（并非每张表都具备全部字段）。
package model
