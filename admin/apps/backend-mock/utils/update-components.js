// 临时脚本：将所有菜单组件指向现有的 profile 页面
const fs = require('fs');
const path = require('path');

// 读取 menu-data.ts
const filePath = path.join(__dirname, 'menu-data.ts');
let content = fs.readFileSync(filePath, 'utf8');

// 组件路径映射
const componentMap = {
  '/system/user/list': '/_core/profile/index',
  '/system/permission/role': '/_core/profile/index',
  '/system/permission/menu': '/_core/profile/index',
  '/system/permission/assignment': '/_core/profile/index',
  '/system/config/basic': '/_core/profile/base-setting',
  '/system/config/dict': '/_core/profile/index',
  '/system/config/log': '/_core/profile/index',
  '/info/category/list': '/_core/profile/index',
  '/info/category/config': '/_core/profile/index',
  '/info/list/publish': '/_core/profile/index',
  '/info/list/review': '/_core/profile/index',
  '/info/list/archive': '/_core/profile/index',
  '/info/content/article': '/_core/profile/index',
  '/info/content/resource': '/_core/profile/index',
  '/info/content/comment': '/_core/profile/index',
  '/profile/info': '/_core/profile/index',
  '/profile/password': '/_core/profile/password-setting',
  '/profile/message': '/_core/profile/notification-setting',
  '/profile/history': '/_core/profile/index',
};

// 替换组件路径
Object.entries(componentMap).forEach(([oldPath, newPath]) => {
  const regex = new RegExp(oldPath.replace(/\//g, '\\/'), 'g');
  content = content.replace(regex, newPath);
});

// 写回文件
fs.writeFileSync(filePath, content);
console.log('组件路径更新完成！');
