/**
 * 用户管理模块
 * 提供用户创建、查询和管理功能
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/**
 * 用户结构体
 * 表示系统中的用户实体
 */
typedef struct {
    int id;
    char* name;
    char* email;
} User;

/**
 * 用户管理器结构体
 * 负责管理用户集合
 */
typedef struct {
    User* users;
    int count;
    int capacity;
} UserManager;

/**
 * 创建新的用户实例
 * 
 * @param id 用户ID
 * @param name 用户名
 * @param email 邮箱地址
 * @return User* 新的用户实例或NULL
 */
User* create_user(int id, const char* name, const char* email) {
    User* user = (User*)malloc(sizeof(User));
    if (user == NULL) {
        return NULL;
    }
    
    user->id = id;
    user->name = strdup(name);
    user->email = strdup(email);
    
    if (user->name == NULL || user->email == NULL) {
        free_user(user);
        return NULL;
    }
    
    return user;
}

/**
 * 释放用户内存
 * 
 * @param user 用户对象
 */
void free_user(User* user) {
    if (user != NULL) {
        free(user->name);
        free(user->email);
        free(user);
    }
}

/**
 * 获取用户信息
 * 
 * @param user 用户对象
 * @param buffer 输出缓冲区
 * @param size 缓冲区大小
 * @return int 成功返回0，失败返回-1
 */
int get_user_info(const User* user, char* buffer, size_t size) {
    if (user == NULL || buffer == NULL) {
        return -1;
    }
    
    return snprintf(buffer, size, "%s (%s)", user->name, user->email);
}

/**
 * 验证用户数据
 * 
 * @param user 用户对象
 * @return int 有效返回1，无效返回0
 */
int is_user_valid(const User* user) {
    if (user == NULL || user->name == NULL || user->email == NULL) {
        return 0;
    }
    
    return (strlen(user->name) > 0 && 
            strlen(user->email) > 0 && 
            strchr(user->email, '@') != NULL);
}

/**
 * 创建新的用户管理器
 * 
 * @param initial_capacity 初始容量
 * @return UserManager* 新的用户管理器或NULL
 */
UserManager* create_user_manager(int initial_capacity) {
    UserManager* manager = (UserManager*)malloc(sizeof(UserManager));
    if (manager == NULL) {
        return NULL;
    }
    
    manager->users = (User*)malloc(sizeof(User) * initial_capacity);
    if (manager->users == NULL) {
        free(manager);
        return NULL;
    }
    
    manager->count = 0;
    manager->capacity = initial_capacity;
    
    return manager;
}

/**
 * 添加用户到管理器
 * 
 * @param manager 用户管理器
 * @param user 要添加的用户
 * @return int 成功返回用户总数，失败返回-1
 */
int add_user(UserManager* manager, User* user) {
    if (manager == NULL || user == NULL) {
        return -1;
    }
    
    if (manager->count >= manager->capacity) {
        // 扩展容量
        int new_capacity = manager->capacity * 2;
        User* new_users = (User*)realloc(manager->users, sizeof(User) * new_capacity);
        if (new_users == NULL) {
            return -1;
        }
        manager->users = new_users;
        manager->capacity = new_capacity;
    }
    
    manager->users[manager->count] = *user;
    manager->count++;
    
    return manager->count;
}

/**
 * 根据ID查找用户
 * 
 * @param manager 用户管理器
 * @param user_id 用户ID
 * @return User* 找到的用户或NULL
 */
User* find_user_by_id(const UserManager* manager, int user_id) {
    if (manager == NULL) {
        return NULL;
    }
    
    for (int i = 0; i < manager->count; i++) {
        if (manager->users[i].id == user_id) {
            return &manager->users[i];
        }
    }
    
    return NULL;
}

/**
 * 释放用户管理器内存
 * 
 * @param manager 用户管理器
 */
void free_user_manager(UserManager* manager) {
    if (manager != NULL) {
        for (int i = 0; i < manager->count; i++) {
            free(manager->users[i].name);
            free(manager->users[i].email);
        }
        free(manager->users);
        free(manager);
    }
}
