package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time" // 引入 time 包

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Comment 结构体定义了评论的数据模型
type Comment struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserName  string    `json:"userName"`
	Content   string    `json:"content"`
	// === 核心修复：将 CreatedAt 类型从 gorm.DeletedAt 修改为 time.Time ===
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"` // 自动创建时间，更符合前端预期
}

// 统一响应结构体
type APIResponse struct {
	Code int         `json:"code"` // 0 表示成功，其他表示失败
	Msg  string      `json:"msg"`  // 返回信息
	Data interface{} `json:"data"` // 返回数据
}

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("comments.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}

	db.AutoMigrate(&Comment{}) // GORM 会根据 Comment 结构体自动创建/更新表
	log.Println("数据库连接成功并已执行迁移。")

	// === 修改路由：匹配新的 URL 和方法 ===
	http.HandleFunc("/comment/get", getCommentsHandler)   // 获取评论的 URL
	http.HandleFunc("/comment/add", addCommentHandler)    // 假设添加评论的 URL
	http.HandleFunc("/comment/delete", deleteCommentHandler) // 假设删除评论的 URL

	port := ":8080"
	log.Printf("服务器正在端口 %s 上监听...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// getCommentsHandler 处理 GET /comment/get 请求
func getCommentsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    if r.Method == http.MethodOptions { // 处理 CORS 预检请求
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        w.WriteHeader(http.StatusOK)
        return
    }
    w.Header().Set("Access-Control-Allow-Origin", "*") // 对于非 OPTIONS 请求也设置 CORS

    if r.Method != http.MethodGet {
        writeJSONResponse(w, http.StatusMethodNotAllowed, 1000, "不支持的方法", nil)
        return
    }

    pageStr := r.URL.Query().Get("page")
    sizeStr := r.URL.Query().Get("size")

    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        page = 1 // 默认页码为 1
    }

    size, err := strconv.Atoi(sizeStr)
    if err != nil || size == 0 { // size=-1 表示所有评论，其他非正数或解析失败也给个默认值
        size = 10 // 默认每页容量为 10
    }

    var allComments []Comment
	// 确保 result 变量在 if 语句外部声明
    var result *gorm.DB
    result = db.Order("id asc").Find(&allComments) // 假设按 ID 升序排序

    if result.Error != nil {
        log.Printf("获取评论失败: %v\n", result.Error)
        writeJSONResponse(w, http.StatusInternalServerError, 1001, "无法获取评论", nil)
        return
    }

    var commentsToShow []Comment

    if size == -1 {
        commentsToShow = allComments // 返回所有评论
    } else {
        // 实现分页逻辑
        start := (page - 1) * size
        end := start + size

        if start >= len(allComments) { // 如果起始索引超出总评论数
            commentsToShow = []Comment{} // 返回空数组
        } else {
            if end > len(allComments) { // 如果结束索引超出总评论数
                end = len(allComments)
            }
            commentsToShow = allComments[start:end]
        }
    }

    log.Printf("GET /comment/get: 从数据库返回 %d 条评论 (page=%d, size=%d)\n", len(commentsToShow), page, size)
    writeJSONResponse(w, http.StatusOK, 0, "成功", commentsToShow)
}


// addCommentHandler 封装统一响应格式
func addCommentHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    if r.Method == http.MethodOptions { // 处理 CORS 预检请求
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        w.WriteHeader(http.StatusOK)
        return
    }
    w.Header().Set("Access-Control-Allow-Origin", "*") // 对于非 OPTIONS 请求也设置 CORS

    if r.Method != http.MethodPost {
        writeJSONResponse(w, http.StatusMethodNotAllowed, 1000, "不支持的方法", nil)
        return
    }

    var comment Comment
    if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
        writeJSONResponse(w, http.StatusBadRequest, 1002, "无效的请求体: "+err.Error(), nil)
        return
    }

    // 确保 result 变量在 if 语句外部声明
    var result *gorm.DB 

    result = db.Create(&comment) // 移除 := 操作符

    if result.Error != nil {
        log.Printf("添加评论失败: %v\n", result.Error)
        writeJSONResponse(w, http.StatusInternalServerError, 1003, "无法添加评论", nil)
        return
    }

    if result.RowsAffected == 0 {
        log.Println("添加评论成功但RowsAffected为0，可能发生异常")
        writeJSONResponse(w, http.StatusInternalServerError, 1007, "添加评论失败或未影响任何行", nil)
        return
    }

    writeJSONResponse(w, http.StatusCreated, 0, "成功添加", comment)
    log.Printf("POST /comment/add: 添加评论 ID %d: %s\n", comment.ID, comment.UserName)
}

// deleteCommentHandler 封装统一响应格式
func deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    if r.Method == http.MethodOptions { // 处理 CORS 预检请求
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        w.WriteHeader(http.StatusOK)
        return
    }
    w.Header().Set("Access-Control-Allow-Origin", "*") // 对于非 OPTIONS 请求也设置 CORS

    if r.Method != http.MethodDelete {
        writeJSONResponse(w, http.StatusMethodNotAllowed, 1000, "不支持的方法", nil)
        return
    }

    // 假设删除的ID是通过查询参数传递的，例如 /comment/delete?id=123
    idStr := r.URL.Query().Get("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        writeJSONResponse(w, http.StatusBadRequest, 1004, "无效的评论ID", nil)
        return
    }

    // 确保 result 变量在 if 语句外部声明
    var result *gorm.DB 

    result = db.Delete(&Comment{}, id) 
    
    if result.Error != nil {
        log.Printf("删除评论失败: %v\n", result.Error)
        writeJSONResponse(w, http.StatusInternalServerError, 1005, "无法删除评论", nil)
        return
    }

    if result.RowsAffected == 0 {
        writeJSONResponse(w, http.StatusNotFound, 1006, "未找到评论", nil)
        return
    }

    writeJSONResponse(w, http.StatusOK, 0, "成功删除", nil) // 删除成功，data 可以为 nil
    log.Printf("DELETE /comment/delete: 删除评论 ID %d\n", id)
}

// 统一 JSON 响应函数
func writeJSONResponse(w http.ResponseWriter, statusCode int, code int, msg string, data interface{}) {
    w.WriteHeader(statusCode)
    response := APIResponse{
        Code: code,
        Msg:  msg,
        Data: data,
    }
    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("编码 JSON 响应失败: %v", err)
    }
}
