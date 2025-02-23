import 'package:flutter/material.dart';
import '../services/task_service.dart';
import '../models/task.dart';

class TaskScreen extends StatefulWidget {
  const TaskScreen({super.key});

  @override
  TaskScreenState createState() => TaskScreenState();
}

class TaskScreenState extends State<TaskScreen> {
  final TaskService _taskService = TaskService();
  final _titleController = TextEditingController();
  final _descriptionController = TextEditingController();
  List<Task> _tasks = [];
  Task? _editingTask;

  @override
  void initState() {
    super.initState();
    _loadTasks();
  }

  Future<void> _loadTasks() async {
    try {
      final tasks = await _taskService.getTasks();
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

  Future<void> _createOrUpdateTask() async {
    try {
      if (_editingTask != null) {
        final updatedTask = await _taskService.updateTask(
          _editingTask!.id,
          _titleController.text,
          _descriptionController.text,
          _editingTask!.completed,
        );
        setState(() {
          final index = _tasks.indexWhere((t) => t.id == _editingTask!.id);
          _tasks[index] = updatedTask;
          _editingTask = null;
        });
      } else {
        final task = await _taskService.createTask(
          _titleController.text,
          _descriptionController.text,
        );
        setState(() {
          _tasks.insert(0, task);
        });
      }
      _titleController.clear();
      _descriptionController.clear();
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
    });
  }

  void _cancelEdit() {
    setState(() {
      _editingTask = null;
      _titleController.clear();
      _descriptionController.clear();
    });
  }

  Future<void> _toggleTaskCompletion(Task task) async {
    try {
      final updatedTask = await _taskService.updateTask(
        task.id,
        task.title,
        task.description,
        !task.completed,
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
        ],
      ),
      body: Column(
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
                    subtitle: Text(task.description),
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
    );
  }

  @override
  void dispose() {
    _titleController.dispose();
    _descriptionController.dispose();
    super.dispose();
  }
} 