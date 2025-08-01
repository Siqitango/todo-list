// API 服务配置
const API_BASE_URL = 'http://localhost:8000';

// API 服务类
class TodoApiService {
    static async request(method, url, data = null) {
        const options = {
            method,
            headers: {
                'Content-Type': 'application/json',
            },
        };

        if (data) {
            options.body = JSON.stringify(data);
        }

        try {
            const response = await fetch(`${API_BASE_URL}${url}`, options);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            // 对于删除操作，不期望返回数据
            if (method === 'DELETE') {
                return { success: true };
            }

            return await response.json();
        } catch (error) {
            console.error('API request failed:', error);
            throw error;
        }
    }

    // 创建新任务
    static async createTodo(todo) {
        return this.request('POST', '/api/todo', {
            title: todo.text,
            description: '', // 前端目前没有提供描述字段
            priority: this.mapPriority(todo.priority)
        });
    }

    // 获取所有任务
    static async getTodos(page = 1, pageSize = 10, priority = null, status = null) {
        let url = `/api/todos?page=${page}&pageSize=${pageSize}`;
        if (priority !== null) {
            url += `&priority=${this.mapPriority(priority)}`;
        }
        if (status !== null) {
            url += `&status=${status}`;
        }
        return this.request('GET', url);
    }

    // 获取单个任务
    static async getTodo(id) {
        return this.request('GET', `/api/todo/${id}`);
    }

    // 更新任务
    static async updateTodo(id, updates) {
        const data = {};
        if (updates.text !== undefined) data.title = updates.text;
        if (updates.completed !== undefined) data.status = updates.completed ? 2 : 1;
        if (updates.priority !== undefined) data.priority = this.mapPriority(updates.priority);

        return this.request('PUT', `/api/todo/${id}`, data);
    }

    // 删除任务
    static async deleteTodo(id) {
        return this.request('DELETE', `/api/todo/${id}`);
    }

    // 映射前端优先级到后端枚举
    static mapPriority(priority) {
        switch (priority) {
            case 'high': return 3;
            case 'medium': return 2;
            case 'low': return 1;
            default: return 2;
        }
    }

    // 映射后端优先级到前端
    static mapBackPriority(priority) {
        switch (priority) {
            case 3: return 'high';
            case 2: return 'medium';
            case 1: return 'low';
            default: return 'medium';
        }
    }
}

// 任务数据管理
class TodoManager {
    constructor() {
        this.todos = [];
        this.currentFilter = 'all';
        this.currentPage = 1;
        this.pageSize = 10;
        this.total = 0;
    }

    // 加载任务
    async loadTodos() {
        try {
            // 始终获取所有任务，然后在前端进行筛选
            const response = await TodoApiService.getTodos(
                this.currentPage,
                this.pageSize,
                null,
                null
            );

            // 转换后端数据格式为前端使用的格式
            this.todos = response.todos.map(todo => ({
                id: todo.id,
                text: todo.title,
                description: todo.description,
                completed: todo.status === 2,
                priority: TodoApiService.mapBackPriority(todo.priority),
                createdAt: todo.createdAt,
                updatedAt: todo.updatedAt
            }));

            this.total = response.total;
            this.renderTodos();
            this.updateStats();
        } catch (error) {
            console.error('Failed to load todos:', error);
            // 显示错误消息给用户
            const todoList = document.getElementById('todo-list');
            todoList.innerHTML = '<div class="p-8 text-center text-red-500"><i class="fa fa-exclamation-circle text-3xl mb-3"></i><p>加载任务失败，请重试</p></div>';
        }
    }

    // 添加新任务
    async addTodo(text, priority = 'medium') {
        try {
            const todoData = {
                text,
                priority
            };

            const response = await TodoApiService.createTodo(todoData);

            // 转换后端响应为前端使用的格式
            const newTodo = {
                id: response.id,
                text: response.title,
                description: response.description,
                completed: response.status === 2,
                priority: TodoApiService.mapBackPriority(response.priority),
                createdAt: response.createdAt,
                updatedAt: response.updatedAt
            };

            this.todos.unshift(newTodo); // 添加到列表开头
            this.renderTodos();
            this.updateStats();
            return newTodo;
        } catch (error) {
            console.error('Failed to add todo:', error);
            alert('添加任务失败，请重试');
        }
    }

    // 切换任务完成状态
    async toggleTodo(id) {
        try {
            const todo = this.todos.find(t => t.id === id);
            if (!todo) return;

            const updates = {
                completed: !todo.completed
            };

            // 先更新本地UI，提升用户体验
            this.todos = this.todos.map(t =>
                t.id === id ? { ...t, completed: !t.completed } : t
            );
            this.renderTodos();
            this.updateStats();

            // 再更新后端数据
            await TodoApiService.updateTodo(id, updates);
        } catch (error) {
            console.error('Failed to toggle todo:', error);
            // 回滚UI状态
            this.todos = this.todos.map(t =>
                t.id === id ? { ...t, completed: !t.completed } : t
            );
            this.renderTodos();
            this.updateStats();
            alert('更新任务状态失败，请重试');
        }
    }

    // 删除任务
    async deleteTodo(id) {
        try {
            await TodoApiService.deleteTodo(id);
            // 更新本地数据
            this.todos = this.todos.filter(t => t.id !== id);
            this.renderTodos();
            this.updateStats();
        } catch (error) {
            console.error('Failed to delete todo:', error);
            alert('删除任务失败，请重试');
        }
    }

    // 清除所有已完成任务
    async clearCompleted() {
        if (!confirm('确定要清除所有已完成的任务吗？')) return;

        try {
            // 获取所有已完成任务的ID
            const completedIds = this.todos.filter(t => t.completed).map(t => t.id);

            // 批量删除
            for (const id of completedIds) {
                await TodoApiService.deleteTodo(id);
            }

            // 更新本地数据
            this.todos = this.todos.filter(t => !t.completed);
            this.renderTodos();
            this.updateStats();
        } catch (error) {
            console.error('Failed to clear completed todos:', error);
            alert('清除已完成任务失败，请重试');
        }
    }

    // 设置当前过滤器
    async setFilter(filter) {
        this.currentFilter = filter;
        this.currentPage = 1; // 重置为第一页
        await this.loadTodos();
        // 更新筛选按钮状态
        document.querySelectorAll('.filter-btn').forEach(btn => {
            btn.classList.toggle('active', btn.dataset.filter === filter);
            btn.classList.toggle('bg-primary', btn.dataset.filter === filter);
            btn.classList.toggle('text-white', btn.dataset.filter === filter);
            btn.classList.toggle('bg-gray-200', btn.dataset.filter !== filter);
        });
    }

    // 获取筛选后的任务
    getFilteredTodos() {
        switch (this.currentFilter) {
            case 'active':
                return this.todos.filter(todo => !todo.completed);
            case 'completed':
                return this.todos.filter(todo => todo.completed);
            case 'high':
                return this.todos.filter(todo => todo.priority === 'high');
            case 'medium':
                return this.todos.filter(todo => todo.priority === 'medium');
            case 'low':
                return this.todos.filter(todo => todo.priority === 'low');
            default:
                return this.todos;
        }
    }

    // 获取优先级样式
    getPriorityClass(priority) {
        switch (priority) {
            case 'high':
                return 'bg-red-50 border-red-100';
            case 'medium':
                return 'bg-yellow-50 border-yellow-100';
            case 'low':
                return 'bg-green-50 border-green-100';
            default:
                return '';
        }
    }

    // 获取优先级标签
    getPriorityBadge(priority) {
        switch (priority) {
            case 'high':
                return '<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">高优先级</span>';
            case 'medium':
                return '<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">中优先级</span>';
            case 'low':
                return '<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">低优先级</span>';
            default:
                return '';
        }
    }

    // 渲染任务列表
    renderTodos() {
        const todoList = document.getElementById('todo-list');
        const filteredTodos = this.getFilteredTodos();

        if (filteredTodos.length === 0) {
            todoList.innerHTML = '<div class="p-8 text-center text-gray-400"><i class="fa fa-tasks text-5xl mb-3 opacity-30"></i><p>暂无任务，添加一个新任务开始吧！</p></div>';
            return;
        }

        todoList.innerHTML = filteredTodos.map(todo => {
            const priorityClass = this.getPriorityClass(todo.priority);
            const priorityBadge = this.getPriorityBadge(todo.priority);
            const completedClass = todo.completed ? 'line-through text-gray-400' : '';
            const checkboxClass = todo.completed ? 'bg-green-500 text-white' : 'border-gray-300 hover:border-primary';

            return `
                <div class="todo-item p-4 flex items-center gap-3 ${priorityClass} border-l-4 border-l-${todo.priority === 'high' ? 'red' : todo.priority === 'medium' ? 'yellow' : 'green'}-400 hover:bg-gray-50 transition-colors"
                     data-id="${todo.id}">
                    <button class="toggle-btn flex-shrink-0 w-6 h-6 rounded-full ${checkboxClass} border-2 flex items-center justify-center transition-all duration-300"
                            aria-label="${todo.completed ? '标记为未完成' : '标记为已完成'}">
                        ${todo.completed ? '<i class="fa fa-check text-sm"></i>' : ''}
                    </button>
                    <div class="flex-grow min-w-0">
                        <p class="text-sm ${completedClass} break-words">${todo.text}</p>
                        <div class="flex items-center mt-1 gap-2 text-xs text-gray-500">
                            ${priorityBadge}
                            <span>${new Date(todo.createdAt).toLocaleDateString()}</span>
                            ${todo.completed ? '<span class="bg-green-100 text-green-800 px-2 py-0.5 rounded-full text-xs font-medium">已完成</span>' : ''}
                        </div>
                    </div>
                    <button class="delete-btn text-gray-400 hover:text-red-500 transition-colors p-1"
                            aria-label="删除任务">
                        <i class="fa fa-trash-o"></i>
                    </button>
                </div>
            `;
        }).join('');

        // 更新任务计数
        document.getElementById('task-count').textContent = this.todos.length;

        // 添加事件监听器
        this.addEventListeners();
    }

    // 添加事件监听器
    addEventListeners() {
        // 使用事件委托方式绑定事件
        const todoList = document.getElementById('todo-list');
        
        // 切换任务状态
        todoList.addEventListener('click', (e) => {
            if (e.target.closest('.toggle-btn')) {
                const btn = e.target.closest('.toggle-btn');
                const id = btn.closest('.todo-item').dataset.id;
                this.toggleTodo(id);
            }
        });
        
        // 删除任务
        todoList.addEventListener('click', (e) => {
            if (e.target.closest('.delete-btn')) {
                const btn = e.target.closest('.delete-btn');
                const id = btn.closest('.todo-item').dataset.id;
                this.deleteTodo(id);
            }
        });
    }

    // 更新统计信息
    updateStats() {
        const total = this.todos.length;
        const completed = this.todos.filter(todo => todo.completed).length;
        const active = total - completed;
        const completionRate = total > 0 ? Math.round((completed / total) * 100) : 0;
        const highPriority = this.todos.filter(todo => todo.priority === 'high').length;
        const mediumPriority = this.todos.filter(todo => todo.priority === 'medium').length;
        const lowPriority = this.todos.filter(todo => todo.priority === 'low').length;

        // 更新统计面板
        document.getElementById('total-tasks').textContent = total;
        document.getElementById('completed-tasks').textContent = completed;
        document.getElementById('active-tasks').textContent = active;
        document.getElementById('completion-progress').style.width = `${completionRate}%`;
        document.getElementById('completion-rate').textContent = `${completionRate}%`;
        document.getElementById('high-priority-count').textContent = highPriority;
        document.getElementById('medium-priority-count').textContent = mediumPriority;
        document.getElementById('low-priority-count').textContent = lowPriority;
    }
}

// 初始化应用
async function initApp() {
    const todoManager = new TodoManager();
    const addTodoForm = document.getElementById('add-todo-form');
    const todoInput = document.getElementById('todo-input');
    const prioritySelect = document.getElementById('priority-select');
    const filterButtons = document.querySelectorAll('.filter-btn');
    const clearCompletedBtn = document.getElementById('clear-completed');
    const statsBtn = document.getElementById('stats-btn');
    const closeStatsBtn = document.getElementById('close-stats');
    const closeStatsModalBtn = document.getElementById('close-stats-btn');
    const statsModal = document.getElementById('stats-modal');
    const themeToggle = document.getElementById('theme-toggle');

    // 加载初始任务列表
    await todoManager.loadTodos();

    // 添加新任务
    addTodoForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const text = todoInput.value.trim();
        const priority = prioritySelect.value;

        if (text) {
            await todoManager.addTodo(text, priority);
            todoInput.value = '';
            todoInput.focus();
        }
    });

    // 筛选任务
    filterButtons.forEach(btn => {
        btn.addEventListener('click', async () => {
            await todoManager.setFilter(btn.dataset.filter);
        });
    });

    // 清除已完成任务
    clearCompletedBtn.addEventListener('click', async () => {
        await todoManager.clearCompleted();
    });

    // 打开统计面板
    statsBtn.addEventListener('click', () => {
        statsModal.classList.remove('opacity-0', 'pointer-events-none');
        statsModal.querySelector('div').classList.remove('translate-y-8');
    });

    // 关闭统计面板
    closeStatsBtn.addEventListener('click', closeStats);
    closeStatsModalBtn.addEventListener('click', closeStats);
    statsModal.addEventListener('click', (e) => {
        if (e.target === statsModal) {
            closeStats();
        }
    });

    function closeStats() {
        statsModal.classList.add('opacity-0', 'pointer-events-none');
        statsModal.querySelector('div').classList.add('translate-y-8');
    }

    // 主题切换功能 (简单实现)
    let isDarkMode = false;
    themeToggle.addEventListener('click', () => {
        isDarkMode = !isDarkMode;
        const icon = themeToggle.querySelector('i');

        if (isDarkMode) {
            document.body.classList.add('bg-gray-900', 'text-white');
            document.body.classList.remove('bg-gradient-to-br', 'from-light', 'to-gray-100', 'text-dark');
            icon.classList.remove('fa-moon-o');
            icon.classList.add('fa-sun-o');
        } else {
            document.body.classList.remove('bg-gray-900', 'text-white');
            document.body.classList.add('bg-gradient-to-br', 'from-light', 'to-gray-100', 'text-dark');
            icon.classList.remove('fa-sun-o');
            icon.classList.add('fa-moon-o');
        }
    });
}

// 当DOM加载完成后初始化应用
document.addEventListener('DOMContentLoaded', initApp);