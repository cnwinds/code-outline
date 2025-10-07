/**
 * 用户管理模块
 * 提供用户创建、查询和管理功能
 */

#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
#include <memory>

/**
 * 用户类
 * 表示系统中的用户实体
 */
class User {
private:
    int id;
    std::string name;
    std::string email;

public:
    /**
     * 构造函数
     * 
     * @param id 用户ID
     * @param name 用户名
     * @param email 邮箱地址
     */
    User(int id, const std::string& name, const std::string& email)
        : id(id), name(name), email(email) {}
    
    /**
     * 获取用户ID
     * 
     * @return int 用户ID
     */
    int getId() const {
        return id;
    }
    
    /**
     * 获取用户名
     * 
     * @return const std::string& 用户名
     */
    const std::string& getName() const {
        return name;
    }
    
    /**
     * 获取邮箱地址
     * 
     * @return const std::string& 邮箱地址
     */
    const std::string& getEmail() const {
        return email;
    }
    
    /**
     * 获取用户信息
     * 
     * @return std::string 用户信息字符串
     */
    std::string getInfo() const {
        return name + " (" + email + ")";
    }
    
    /**
     * 验证用户数据
     * 
     * @return bool 验证结果
     */
    bool isValid() const {
        return !name.empty() && 
               !email.empty() && 
               email.find('@') != std::string::npos;
    }
};

/**
 * 用户管理器类
 * 负责管理用户集合
 */
class UserManager {
private:
    std::vector<std::unique_ptr<User>> users;

public:
    /**
     * 添加用户到管理器
     * 
     * @param user 要添加的用户
     * @return size_t 用户总数
     */
    size_t addUser(std::unique_ptr<User> user) {
        users.push_back(std::move(user));
        return users.size();
    }
    
    /**
     * 获取所有用户
     * 
     * @return const std::vector<std::unique_ptr<User>>& 用户列表
     */
    const std::vector<std::unique_ptr<User>>& getUsers() const {
        return users;
    }
    
    /**
     * 根据ID查找用户
     * 
     * @param userId 用户ID
     * @return User* 找到的用户或nullptr
     */
    User* findUserById(int userId) const {
        auto it = std::find_if(users.begin(), users.end(),
            [userId](const std::unique_ptr<User>& user) {
                return user->getId() == userId;
            });
        
        return (it != users.end()) ? it->get() : nullptr;
    }
    
    /**
     * 获取用户统计信息
     * 
     * @return std::pair<size_t, size_t> 总用户数和有效用户数
     */
    std::pair<size_t, size_t> getStats() const {
        size_t totalUsers = users.size();
        size_t validUsers = std::count_if(users.begin(), users.end(),
            [](const std::unique_ptr<User>& user) {
                return user->isValid();
            });
        
        return std::make_pair(totalUsers, validUsers);
    }
};

/**
 * 用户工具类
 * 提供用户相关的工具方法
 */
class UserUtils {
public:
    /**
     * 创建用户实例
     * 
     * @param id 用户ID
     * @param name 用户名
     * @param email 邮箱地址
     * @return std::unique_ptr<User> 用户实例
     */
    static std::unique_ptr<User> createUser(int id, 
                                           const std::string& name, 
                                           const std::string& email) {
        return std::make_unique<User>(id, name, email);
    }
    
    /**
     * 验证用户数据
     * 
     * @param user 用户对象
     * @return bool 验证结果
     */
    static bool validateUser(const User& user) {
        return user.isValid();
    }
    
    /**
     * 打印用户信息
     * 
     * @param user 用户对象
     */
    static void printUserInfo(const User& user) {
        std::cout << "用户信息: " << user.getInfo() << std::endl;
    }
};

/**
 * 用户管理命名空间
 * 包含用户管理相关的功能
 */
namespace UserManagement {
    /**
     * 用户验证器类
     * 提供用户数据验证功能
     */
    class UserValidator {
    public:
        /**
         * 验证用户名
         * 
         * @param name 用户名
         * @return bool 验证结果
         */
        static bool validateName(const std::string& name) {
            return !name.empty() && name.length() >= 2;
        }
        
        /**
         * 验证邮箱地址
         * 
         * @param email 邮箱地址
         * @return bool 验证结果
         */
        static bool validateEmail(const std::string& email) {
            return !email.empty() && 
                   email.find('@') != std::string::npos &&
                   email.find('.') != std::string::npos;
        }
    };
}
