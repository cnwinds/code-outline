//! 用户管理模块
//! 提供用户创建、查询和管理功能

use std::collections::HashMap;

/// 用户结构体
/// 表示系统中的用户实体
#[derive(Debug, Clone)]
pub struct User {
    pub id: u32,
    pub name: String,
    pub email: String,
}

impl User {
    /// 创建新的用户实例
    /// 
    /// # Arguments
    /// 
    /// * `id` - 用户ID
    /// * `name` - 用户名
    /// * `email` - 邮箱地址
    /// 
    /// # Returns
    /// 
    /// 返回新的用户实例
    pub fn new(id: u32, name: String, email: String) -> Self {
        User { id, name, email }
    }
    
    /// 获取用户信息
    /// 
    /// # Returns
    /// 
    /// 返回格式化的用户信息字符串
    pub fn get_info(&self) -> String {
        format!("{} ({})", self.name, self.email)
    }
    
    /// 验证用户数据
    /// 
    /// # Returns
    /// 
    /// 返回验证结果
    pub fn is_valid(&self) -> bool {
        !self.name.is_empty() && 
        !self.email.is_empty() && 
        self.email.contains('@')
    }
}

/// 用户管理器
/// 负责管理用户集合
pub struct UserManager {
    users: Vec<User>,
}

impl UserManager {
    /// 创建新的用户管理器
    /// 
    /// # Returns
    /// 
    /// 返回新的用户管理器实例
    pub fn new() -> Self {
        UserManager {
            users: Vec::new(),
        }
    }
    
    /// 添加用户到管理器
    /// 
    /// # Arguments
    /// 
    /// * `user` - 要添加的用户
    /// 
    /// # Returns
    /// 
    /// 返回用户总数
    pub fn add_user(&mut self, user: User) -> usize {
        self.users.push(user);
        self.users.len()
    }
    
    /// 获取所有用户
    /// 
    /// # Returns
    /// 
    /// 返回用户列表的克隆
    pub fn get_users(&self) -> Vec<User> {
        self.users.clone()
    }
    
    /// 根据ID查找用户
    /// 
    /// # Arguments
    /// 
    /// * `user_id` - 用户ID
    /// 
    /// # Returns
    /// 
    /// 返回找到的用户或None
    pub fn find_user_by_id(&self, user_id: u32) -> Option<&User> {
        self.users.iter().find(|user| user.id == user_id)
    }
    
    /// 获取用户统计信息
    /// 
    /// # Returns
    /// 
    /// 返回包含统计信息的HashMap
    pub fn get_stats(&self) -> HashMap<String, usize> {
        let mut stats = HashMap::new();
        stats.insert("total_users".to_string(), self.users.len());
        stats.insert("valid_users".to_string(), 
                    self.users.iter().filter(|u| u.is_valid()).count());
        stats
    }
}

/// 用户工具函数
pub mod user_utils {
    use super::User;
    
    /// 创建用户实例
    /// 
    /// # Arguments
    /// 
    /// * `id` - 用户ID
    /// * `name` - 用户名
    /// * `email` - 邮箱地址
    /// 
    /// # Returns
    /// 
    /// 返回新的用户实例
    pub fn create_user(id: u32, name: &str, email: &str) -> User {
        User::new(id, name.to_string(), email.to_string())
    }
    
    /// 验证用户数据
    /// 
    /// # Arguments
    /// 
    /// * `user` - 用户对象
    /// 
    /// # Returns
    /// 
    /// 返回验证结果
    pub fn validate_user(user: &User) -> bool {
        user.is_valid()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_user_creation() {
        let user = User::new(1, "张三".to_string(), "zhangsan@example.com".to_string());
        assert_eq!(user.id, 1);
        assert_eq!(user.name, "张三");
        assert_eq!(user.email, "zhangsan@example.com");
    }
    
    #[test]
    fn test_user_validation() {
        let valid_user = User::new(1, "张三".to_string(), "zhangsan@example.com".to_string());
        let invalid_user = User::new(2, "".to_string(), "invalid-email".to_string());
        
        assert!(valid_user.is_valid());
        assert!(!invalid_user.is_valid());
    }
}
