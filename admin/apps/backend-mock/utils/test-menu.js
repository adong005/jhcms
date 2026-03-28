// 测试菜单数据
const { MOCK_MENUS, MOCK_MENU_LIST } = require('./menu-data.ts');

console.log('=== 用户菜单数据测试 ===\n');

// 测试不同用户的菜单
MOCK_MENUS.forEach(userMenu => {
  console.log(`用户: ${userMenu.username} (角色: ${userMenu.roles.join(', ')})`);
  console.log(`菜单数量: ${userMenu.menus.length}`);
  
  // 打印一级菜单
  userMenu.menus.forEach(menu => {
    console.log(`  - ${menu.meta?.title || menu.name} (${menu.path})`);
  });
  console.log('');
});

console.log('\n=== 菜单列表数据测试 ===\n');
console.log(`总菜单数量: ${MOCK_MENU_LIST.length}`);

// 统计各类型菜单数量
let menuCount = 0;
let catalogCount = 0;

MOCK_MENU_LIST.forEach(item => {
  if (item.type === 'menu') menuCount++;
  if (item.type === 'catalog') catalogCount++;
});

console.log(`菜单项: ${menuCount}`);
console.log(`目录项: ${catalogCount}`);
