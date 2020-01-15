// import fullRoutes from '../router/fullpath';
/**
 * 根据权限获取动态路由
 * @param permission 后端返回的权限
 * @returns {Array}
 */
const getActiveRouters = (fullRoutes, permission) => {
  // const permission = {home: 1, parent: 1, child1: 1, child2: 1};
  /**
   * input like this:
   * {
   *     home: 1,
   *     child1: 0
   * }
   */
  const filterFn = (item) => {
    return permission[item.key] === 1;
  };
  if (Array.isArray(fullRoutes)) {
    let routersRes = [];
    fullRoutes.map((item) => {
      if (!item.hasChildren && !item.key) {
        routersRes.push(item);
      } else if (!item.hasChildren && item.key) {
        if (permission[item.key]) {
          routersRes.push(item);
        }
      } else if (item.hasChildren) {
        let childrenArr = item.children.filter(filterFn);
        if (childrenArr.length !== 0) {
          delete item.children;
          item.children = childrenArr;
          routersRes.push(item);
        }
      }
      return null;
    });
    return routersRes;
  }
};
export {
  getActiveRouters,
};
