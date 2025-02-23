import 'package:flutter/material.dart';
import '../services/task_service.dart';
import '../services/notification_service.dart';
import '../services/category_service.dart';
import '../models/task.dart';
import '../models/category.dart';
import '../widgets/notification_widget.dart';
import '../widgets/category_widget.dart';
import 'package:intl/intl.dart';

class TaskScreen extends StatefulWidget {
  const TaskScreen({super.key});

  @override
  TaskScreenState createState() => TaskScreenState();
}

class TaskScreenState extends State<TaskScreen> {
  final TaskService _taskService = TaskService();
  final NotificationService _notificationService = NotificationService();
  final CategoryService _categoryService = CategoryService();
  final _titleController = TextEditingController();
  final _descriptionController = TextEditingController();
  List<Task> _tasks = [];
  List<TaskNotification> _notifications = [];
  Task? _editingTask;
  DateTime _selectedDueDate = DateTime.now().add(const Duration(days: 1));
  int _selectedPriority = 0;
  bool _isLoadingNotifications = false;
  Category? _selectedCategory;
  bool _isCategoryPanelVisible = false;

  @override
  void initState() {
    super.initState();
    _loadTasks();
    _loadNotifications();
    _startPeriodicUpdate();
  }

  void _startPeriodicUpdate() {
    Future.delayed(const Duration(minutes: 1), () {
      if (mounted) {
        _loadNotifications();
        _startPeriodicUpdate();
      }
    });
  }

  Future<void> _loadTasks() async {
    try {
      final tasks = await _taskService.getTasks(category: _selectedCategory);
      setState(() {
        _tasks = tasks;
      });
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(e.toString())),
        );
      }
    }
  }

  Future<void> _loadNotifications() async {
    setState(() {
      _isLoadingNotifications = true;
    });
    try {
      final notifications = await _notificationService.getNotifications();
      if (mounted) {
        setState(() {
          _notifications = notifications;
          _isLoadingNotifications = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _isLoadingNotifications = false;
        });
      }
    }
  }

  Future<void> _selectDueDate() async {
    final DateTime? picked = await showDatePicker(
      context: context,
      initialDate: _selectedDueDate,
      firstDate: DateTime.now(),
      lastDate: DateTime.now().add(const Duration(days: 365)),
    );
    if (picked != null) {
      setState(() {
        _selectedDueDate = picked;
      });
    }
  }

  Future<void> _createOrUpdateTask() async {
    try {
      if (_editingTask != null) {
        final updatedTask = await _taskService.updateTask(
          _editingTask!.id,
          _titleController.text,
          _descriptionController.text,
          _editingTask!.completed,
          _selectedDueDate,
          _selectedPriority,
          category: _selectedCategory,
        );
        setState(() {
          final index = _tasks.indexWhere((t) => t.id == _editingTask!.id);
          _tasks[index] = updatedTask;
          _editingTask = null;
          _selectedCategory = null;
        });
      } else {
        final task = await _taskService.createTask(
          _titleController.text,
          _descriptionController.text,
          _selectedDueDate,
          _selectedPriority,
          category: _selectedCategory,
        );
        setState(() {
          _tasks.insert(0, task);
        });
      }
      _titleController.clear();
      _descriptionController.clear();
      setState(() {
        _selectedPriority = 0;
        _selectedDueDate = DateTime.now().add(const Duration(days: 1));
        _selectedCategory = null;
      });
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(e.toString())),
        );
      }
    }
  }

  void _editTask(Task task) {
    setState(() {
      _editingTask = task;
      _titleController.text = task.title;
      _descriptionController.text = task.description;
      _selectedDueDate = task.dueDate;
      _selectedPriority = task.priority;
      _selectedCategory = task.category;
    });
  }

  void _cancelEdit() {
    setState(() {
      _editingTask = null;
      _titleController.clear();
      _descriptionController.clear();
      _selectedDueDate = DateTime.now().add(const Duration(days: 1));
      _selectedPriority = 0;
      _selectedCategory = null;
    });
  }

  Future<void> _toggleTaskCompletion(Task task) async {
    try {
      final updatedTask = await _taskService.updateTask(
        task.id,
        task.title,
        task.description,
        !task.completed,
        task.dueDate,
        task.priority,
      );
      setState(() {
        final index = _tasks.indexWhere((t) => t.id == task.id);
        _tasks[index] = updatedTask;
      });
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(e.toString())),
        );
      }
    }
  }

  Future<void> _deleteTask(Task task) async {
    try {
      await _taskService.deleteTask(task.id);
      setState(() {
        _tasks.removeWhere((t) => t.id == task.id);
      });
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(e.toString())),
        );
      }
    }
  }

  void _showNotifications() {
    showModalBottomSheet(
      context: context,
      builder: (context) => SizedBox(
        height: MediaQuery.of(context).size.height * 0.7,
        child: Column(
          children: [
            AppBar(
              title: const Text('Оповещения'),
              automaticallyImplyLeading: false,
              actions: [
                IconButton(
                  icon: const Icon(Icons.close),
                  onPressed: () => Navigator.pop(context),
                ),
              ],
            ),
            Expanded(
              child: NotificationWidget(),
            ),
          ],
        ),
      ),
    );
  }

  void _toggleCategoryPanel() {
    setState(() {
      _isCategoryPanelVisible = !_isCategoryPanelVisible;
    });
  }

  void _onCategorySelected(Category? category) {
    setState(() {
      _selectedCategory = category;
    });
    _loadTasks();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Задачи'),
        actions: [
          if (_editingTask != null)
            IconButton(
              icon: const Icon(Icons.cancel),
              onPressed: _cancelEdit,
            ),
          IconButton(
            icon: const Icon(Icons.category),
            onPressed: _toggleCategoryPanel,
          ),
          Stack(
            children: [
              IconButton(
                icon: const Icon(Icons.notifications),
                onPressed: _showNotifications,
              ),
              if (_notifications.where((n) => !n.read).isNotEmpty)
                Positioned(
                  right: 0,
                  top: 0,
                  child: Container(
                    padding: const EdgeInsets.all(2),
                    decoration: BoxDecoration(
                      color: Colors.red,
                      borderRadius: BorderRadius.circular(10),
                    ),
                    constraints: const BoxConstraints(
                      minWidth: 16,
                      minHeight: 16,
                    ),
                    child: Text(
                      _notifications.where((n) => !n.read).length.toString(),
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 10,
                      ),
                      textAlign: TextAlign.center,
                    ),
                  ),
                ),
            ],
          ),
        ],
      ),
      body: Row(
        children: [
          if (_isCategoryPanelVisible)
            CategoryWidget(
              onCategorySelected: _onCategorySelected,
              selectedCategory: _selectedCategory,
            ),
          Expanded(
            child: Column(
              children: [
                Padding(
                  padding: const EdgeInsets.all(8.0),
                  child: Column(
                    children: [
                      TextField(
                        controller: _titleController,
                        decoration: const InputDecoration(
                          labelText: 'Название',
                        ),
                      ),
                      TextField(
                        controller: _descriptionController,
                        decoration: const InputDecoration(
                          labelText: 'Описание',
                        ),
                      ),
                      ListTile(
                        title: const Text('Срок выполнения'),
                        subtitle: Text(DateFormat('dd.MM.yyyy').format(_selectedDueDate)),
                        trailing: const Icon(Icons.calendar_today),
                        onTap: _selectDueDate,
                      ),
                      DropdownButtonFormField<int>(
                        value: _selectedPriority,
                        decoration: const InputDecoration(
                          labelText: 'Приоритет',
                        ),
                        items: const [
                          DropdownMenuItem(value: 0, child: Text('Низкий')),
                          DropdownMenuItem(value: 1, child: Text('Средний')),
                          DropdownMenuItem(value: 2, child: Text('Высокий')),
                        ],
                        onChanged: (value) {
                          setState(() {
                            _selectedPriority = value!;
                          });
                        },
                      ),
                      const SizedBox(height: 8),
                      FutureBuilder<List<Category>>(
                        future: _categoryService.getCategories(),
                        builder: (context, snapshot) {
                          if (snapshot.connectionState == ConnectionState.waiting) {
                            return const CircularProgressIndicator();
                          }
                          if (snapshot.hasError) {
                            return Text('Ошибка: ${snapshot.error}');
                          }
                          final categories = snapshot.data ?? [];
                          
                          // Проверяем, существует ли выбранная категория в списке
                          final selectedExists = _selectedCategory == null || 
                              categories.any((c) => c.id == _selectedCategory!.id);
                          
                          // Если выбранная категория не существует, сбрасываем её
                          if (!selectedExists && _selectedCategory != null) {
                            Future.microtask(() => setState(() {
                              _selectedCategory = null;
                            }));
                          }
                          
                          return DropdownButtonFormField<Category?>(
                            value: selectedExists ? _selectedCategory : null,
                            decoration: const InputDecoration(
                              labelText: 'Категория',
                              border: OutlineInputBorder(),
                            ),
                            items: [
                              const DropdownMenuItem<Category?>(
                                value: null,
                                child: Text('Без категории'),
                              ),
                              ...categories.map((category) {
                                return DropdownMenuItem<Category?>(
                                  value: category,
                                  child: Text(category.name),
                                );
                              }).toList(),
                            ],
                            onChanged: (Category? newValue) {
                              setState(() {
                                _selectedCategory = newValue;
                              });
                            },
                            selectedItemBuilder: (BuildContext context) {
                              return [null, ...categories].map<Widget>((Category? category) {
                                if (category == null) {
                                  return const Text('Без категории');
                                }
                                return Text(category.name);
                              }).toList();
                            },
                          );
                        },
                      ),
                      const SizedBox(height: 16),
                      ElevatedButton(
                        onPressed: _createOrUpdateTask,
                        child: Text(_editingTask == null ? 'Добавить задачу' : 'Обновить задачу'),
                      ),
                    ],
                  ),
                ),
                Expanded(
                  child: ListView.builder(
                    itemCount: _tasks.length,
                    itemBuilder: (context, index) {
                      final task = _tasks[index];
                      return Card(
                        margin: const EdgeInsets.symmetric(horizontal: 8.0, vertical: 4.0),
                        child: ListTile(
                          title: Text(task.title),
                          subtitle: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(task.description),
                              const SizedBox(height: 4),
                              Row(
                                children: [
                                  Icon(Icons.calendar_today, size: 16, color: task.dueDate.isBefore(DateTime.now()) ? Colors.red : Colors.grey),
                                  const SizedBox(width: 4),
                                  Text(DateFormat('dd.MM.yyyy').format(task.dueDate)),
                                  const SizedBox(width: 8),
                                  Container(
                                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                                    decoration: BoxDecoration(
                                      color: task.priorityColor.withOpacity(0.2),
                                      borderRadius: BorderRadius.circular(12),
                                    ),
                                    child: Text(
                                      task.priorityText,
                                      style: TextStyle(color: task.priorityColor),
                                    ),
                                  ),
                                  if (task.category != null) ...[
                                    const SizedBox(width: 8),
                                    Container(
                                      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                                      decoration: BoxDecoration(
                                        color: Colors.blue.withOpacity(0.2),
                                        borderRadius: BorderRadius.circular(12),
                                      ),
                                      child: Text(
                                        task.category!.name,
                                        style: const TextStyle(color: Colors.blue),
                                      ),
                                    ),
                                  ],
                                ],
                              ),
                            ],
                          ),
                          leading: Checkbox(
                            value: task.completed,
                            onChanged: (_) => _toggleTaskCompletion(task),
                          ),
                          trailing: Container(
                            width: 100,
                            child: Row(
                              mainAxisSize: MainAxisSize.min,
                              children: [
                                IconButton(
                                  icon: const Icon(Icons.edit, color: Colors.blue),
                                  onPressed: () => _editTask(task),
                                ),
                                IconButton(
                                  icon: const Icon(Icons.delete, color: Colors.red),
                                  onPressed: () => _deleteTask(task),
                                ),
                              ],
                            ),
                          ),
                        ),
                      );
                    },
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  @override
  void dispose() {
    _titleController.dispose();
    _descriptionController.dispose();
    super.dispose();
  }
} 