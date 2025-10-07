/**
 * 用户管理模块
 * 提供用户创建、查询和管理功能
 */
public class UserManager {
    private java.util.List<User> users;
    
    /**
     * 初始化用户管理器
     */
    public UserManager() {
        this.users = new java.util.ArrayList<>();
    }
    
    /**
     * 添加用户到管理器
     * 
     * @param user User 实例
     * @return int 用户总数
     */
    public int addUser(User user) {
        this.users.add(user);
        return this.users.size();
    }
    
    /**
     * 获取所有用户
     * 
     * @return List<User> 用户列表
     */
    public java.util.List<User> getUsers() {
        return new java.util.ArrayList<>(this.users);
    }
    
    /**
     * 根据ID查找用户
     * 
     * @param userId 用户ID
     * @return User 找到的用户或null
     */
    public User findUserById(int userId) {
        for (User user : this.users) {
            if (user.getId() == userId) {
                return user;
            }
        }
        return null;
    }
}

/**
 * 用户类
 * 表示系统中的用户实体
 */
class User {
    private int id;
    private String name;
    private String email;
    
    /**
     * 初始化用户
     * 
     * @param id 用户ID
     * @param name 用户名
     * @param email 邮箱地址
     */
    public User(int id, String name, String email) {
        this.id = id;
        this.name = name;
        this.email = email;
    }
    
    /**
     * 获取用户ID
     * 
     * @return int 用户ID
     */
    public int getId() {
        return this.id;
    }
    
    /**
     * 获取用户名
     * 
     * @return String 用户名
     */
    public String getName() {
        return this.name;
    }
    
    /**
     * 获取邮箱地址
     * 
     * @return String 邮箱地址
     */
    public String getEmail() {
        return this.email;
    }
    
    /**
     * 获取用户信息
     * 
     * @return String 用户信息字符串
     */
    public String getInfo() {
        return this.name + " (" + this.email + ")";
    }
}

/**
 * 工具类
 * 提供用户相关的工具方法
 */
class UserUtils {
    /**
     * 创建用户实例
     * 
     * @param id 用户ID
     * @param name 用户名
     * @param email 邮箱地址
     * @return User 用户实例
     */
    public static User createUser(int id, String name, String email) {
        return new User(id, name, email);
    }
    
    /**
     * 验证用户数据
     * 
     * @param user 用户对象
     * @return boolean 验证结果
     */
    public static boolean validateUser(User user) {
        return user != null && 
               user.getName() != null && 
               !user.getName().trim().isEmpty() &&
               user.getEmail() != null && 
               user.getEmail().contains("@");
    }
}
