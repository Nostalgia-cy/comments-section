import React, { useState, useEffect } from 'react'; 
import './App.css'; // 确保你的样式被正确引入


function App() {

  const [comments, setComments] = useState([
    { id: 1, userName: 'tmy', content: '这是第一条评论' },
    { id: 2, userName: 'tmy', content: '这是第二条评论' },
    { id: 3, userName: 'tmy', content: '这是第三条评论' },
    { id: 4, userName: 'tmy', content: '这是第四条评论' },
  ]);

  const [username, setUsername] = useState(''); // 用于存储用户名输入
  const [commentContent, setCommentContent] = useState(''); // 用于存储评论内容输入

  const [currentPage, setCurrentPage] = useState(1); // 当前页码
  const [commentsPerPage] = useState(5); // 每页显示的评论数量


  // 计算总页数
  const totalPages = Math.ceil(comments.length / commentsPerPage);

  // 计算当前页评论的起始和结束索引
  const indexOfLastComment = currentPage * commentsPerPage;
  const indexOfFirstComment = indexOfLastComment - commentsPerPage;

  // 获取当前页需要显示的评论
  const currentComments = comments.slice(indexOfFirstComment, indexOfLastComment);

  const handleAddComment = () => { // 不再是 async 函数
    if (username.trim() === '' || commentContent.trim() === '') {
      alert('用户名和评论内容不能为空！');
      return;
    }

    // 客户端生成一个简单的 ID
    const newId = comments.length > 0 ? Math.max(...comments.map(c => c.id)) + 1 : 1;
    
    const newComment = { // 用于本地状态的新评论数据
      id: newId,
      userName: username,
      content: commentContent,
    };

    setComments((prevComments) => [...prevComments, newComment]); // 更新评论列表
    
    // 添加新评论后，跳转到包含新评论的最后一页
    const newTotalComments = comments.length + 1; // 加上新添加的评论数量
    const newTotalPagesAfterAdd = Math.ceil(newTotalComments / commentsPerPage);
    setCurrentPage(newTotalPagesAfterAdd); // 跳转到新评论所在的最后一页

    setUsername('');
    setCommentContent('');
    console.log('新评论已添加:', newComment);
  };

  const handleDeleteComment = (idToDelete) => { // 不再是 async 函数
    setComments((prevComments) => {
      const updatedComments = prevComments.filter((comment) => comment.id !== idToDelete);
      
      // 删除评论后，调整当前页码
      const newTotalPages = Math.ceil(updatedComments.length / commentsPerPage);
      if (newTotalPages === 0) { // 如果删除后没有评论了，页码回到1
          setCurrentPage(1);
      } else if (currentPage > newTotalPages) { // 如果当前页码超出了新总页数，调整到最后一页
          setCurrentPage(newTotalPages);
      }
      return updatedComments;
    });
    console.log('评论已删除，ID:', idToDelete);
  };

  return (
    <div className="container">
      {/* 评论输入框部分 */}
      <div className="comment-input-section">
        <h3>评论区</h3>
        <div className="form-group">
          <label htmlFor="username">user name</label>
          <input
            type="text"
            id="username"
            name="username"
            placeholder="请输入您的用户名"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
        </div>
        <div className="form-group">
          <label htmlFor="commentContent">comment content</label>
          <textarea
            id="commentContent"
            name="commentContent"
            placeholder="请输入您的评论"
            value={commentContent}
            onChange={(e) => setCommentContent(e.target.value)}
          ></textarea>
        </div>
        <div className="submit-button-container">
          <button
            type="submit"
            className="submit-button"
            onClick={handleAddComment}
          >
            提交
          </button>
        </div>
      </div>

      {/* 下方评论展示部分 */}
      <div className="comments-list-section">
        <h3>所有评论</h3>
        {/* 使用 currentComments 动态生成评论项 */}
        {currentComments.map((comment) => (
          <div key={comment.id} className="comment-item">
            <div className="user-name">{comment.userName}</div>
            <div className="comment-content">{comment.content}</div>
            <button
              className="delete-button"
              onClick={() => handleDeleteComment(comment.id)}
            >
              删除
            </button>
          </div>
        ))}
        {/* 当 currentComments 为空但 comments 不为空时显示提示 */}
        {currentComments.length === 0 && comments.length > 0 && (
          <p style={{ textAlign: 'center', color: '#666' }}>当前页没有评论，请尝试切换页面。</p>
        )}
        {/* 当 comments 总数为 0 时显示提示 */}
        {comments.length === 0 && (
          <p style={{ textAlign: 'center', color: '#666' }}>暂无评论，快来添加第一条评论吧！</p>
        )}

        {/* 分页控制按钮 */}
        {comments.length > 0 && ( // 只有当有评论时才显示分页控制
          <div className="pagination-controls" style={{ textAlign: 'center', marginTop: '20px' }}>
            <button
              onClick={() => setCurrentPage(prevPage => prevPage - 1)}
              disabled={currentPage === 1} // 当在第一页时禁用“上一页”按钮
              className="pagination-button"
              style={{ padding: '8px 15px', margin: '0 10px', backgroundColor: '#007bff', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer' }}
            >
              上一页
            </button>
            <span>第 {currentPage} 页 / 共 {totalPages} 页</span>
            <button
              onClick={() => setCurrentPage(prevPage => prevPage + 1)}
              disabled={currentPage === totalPages || totalPages === 0} // 当在最后一页或没有评论时禁用“下一页”按钮
              className="pagination-button"
              style={{ padding: '8px 15px', margin: '0 10px', backgroundColor: '#007bff', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer' }}
            >
              下一页
            </button>
          </div>
        )}
      </div>
    </div>
  );
}

export default App;