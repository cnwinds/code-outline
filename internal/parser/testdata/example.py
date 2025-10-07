# -*- coding: utf-8 -*-
"""
用户管理模块
提供用户创建、查询和管理功能
"""

class UserManager:
    """用户管理器类"""
    
    def __init__(self):
        """初始化用户管理器"""
        self.users = []
    
    def add_user(self, user):
        """添加用户到管理器
        
        Args:
            user: User 实例
            
        Returns:
            int: 用户总数
        """
        self.users.append(user)
        return len(self.users)
    
    def get_users(self):
        """获取所有用户
        
        Returns:
            list: 用户列表
        """
        return self.users.copy()
    
    def find_user_by_id(self, user_id):
        """根据ID查找用户
        
        Args:
            user_id: 用户ID
            
        Returns:
            User or None: 找到的用户或None
        """
        for user in self.users:
            if user.id == user_id:
                return user
        return None

class User:
    """用户类"""
    
    def __init__(self, user_id, name, email):
        """初始化用户
        
        Args:
            user_id: 用户ID
            name: 用户名
            email: 邮箱地址
        """
        self.id = user_id
        self.name = name
        self.email = email
    
    def get_info(self):
        """获取用户信息
        
        Returns:
            str: 用户信息字符串
        """
        return f"{self.name} ({self.email})"

def create_user(user_id, name, email):
    """创建用户实例
    
    Args:
        user_id: 用户ID
        name: 用户名
        email: 邮箱地址
        
    Returns:
        User: 用户实例
    """
    return User(user_id, name, email)
