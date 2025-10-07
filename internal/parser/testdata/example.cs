using System;
using System.Collections.Generic;

namespace UserManagement
{
    /// <summary>
    /// 用户管理模块
    /// 提供用户创建、查询和管理功能
    /// </summary>
    public class UserManager
    {
        private List<User> users;
        
        /// <summary>
        /// 初始化用户管理器
        /// </summary>
        public UserManager()
        {
            this.users = new List<User>();
        }
        
        /// <summary>
        /// 添加用户到管理器
        /// </summary>
        /// <param name="user">User 实例</param>
        /// <returns>用户总数</returns>
        public int AddUser(User user)
        {
            this.users.Add(user);
            return this.users.Count;
        }
        
        /// <summary>
        /// 获取所有用户
        /// </summary>
        /// <returns>用户列表</returns>
        public List<User> GetUsers()
        {
            return new List<User>(this.users);
        }
        
        /// <summary>
        /// 根据ID查找用户
        /// </summary>
        /// <param name="userId">用户ID</param>
        /// <returns>找到的用户或null</returns>
        public User FindUserById(int userId)
        {
            foreach (User user in this.users)
            {
                if (user.Id == userId)
                {
                    return user;
                }
            }
            return null;
        }
    }
    
    /// <summary>
    /// 用户类
    /// 表示系统中的用户实体
    /// </summary>
    public class User
    {
        public int Id { get; set; }
        public string Name { get; set; }
        public string Email { get; set; }
        
        /// <summary>
        /// 初始化用户
        /// </summary>
        /// <param name="id">用户ID</param>
        /// <param name="name">用户名</param>
        /// <param name="email">邮箱地址</param>
        public User(int id, string name, string email)
        {
            this.Id = id;
            this.Name = name;
            this.Email = email;
        }
        
        /// <summary>
        /// 获取用户信息
        /// </summary>
        /// <returns>用户信息字符串</returns>
        public string GetInfo()
        {
            return $"{this.Name} ({this.Email})";
        }
    }
    
    /// <summary>
    /// 工具类
    /// 提供用户相关的工具方法
    /// </summary>
    public static class UserUtils
    {
        /// <summary>
        /// 创建用户实例
        /// </summary>
        /// <param name="id">用户ID</param>
        /// <param name="name">用户名</param>
        /// <param name="email">邮箱地址</param>
        /// <returns>用户实例</returns>
        public static User CreateUser(int id, string name, string email)
        {
            return new User(id, name, email);
        }
        
        /// <summary>
        /// 验证用户数据
        /// </summary>
        /// <param name="user">用户对象</param>
        /// <returns>验证结果</returns>
        public static bool ValidateUser(User user)
        {
            return user != null && 
                   !string.IsNullOrEmpty(user.Name) &&
                   !string.IsNullOrEmpty(user.Email) && 
                   user.Email.Contains("@");
        }
    }
}
