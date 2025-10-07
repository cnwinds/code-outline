// 用户管理模块
class UserManager {
    constructor() {
        this.users = [];
    }

    // 添加用户
    addUser(user) {
        this.users.push(user);
        return this.users.length;
    }

    // 获取所有用户
    getUsers() {
        return [...this.users];
    }

    // 根据ID查找用户
    findUserById(id) {
        return this.users.find(user => user.id === id);
    }
}

// 用户类
class User {
    constructor(id, name, email) {
        this.id = id;
        this.name = name;
        this.email = email;
    }

    // 获取用户信息
    getInfo() {
        return `${this.name} (${this.email})`;
    }
}

// 工具函数
const createUser = (id, name, email) => {
    return new User(id, name, email);
};

// 导出模块
module.exports = { UserManager, User, createUser };
