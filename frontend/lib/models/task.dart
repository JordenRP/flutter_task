import 'package:flutter/material.dart';

class Task {
  final int id;
  final String title;
  final String description;
  final bool completed;
  final int userId;
  final DateTime dueDate;
  final int priority;
  final DateTime createdAt;
  final DateTime updatedAt;

  Task({
    required this.id,
    required this.title,
    required this.description,
    required this.completed,
    required this.userId,
    required this.dueDate,
    required this.priority,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Task.fromJson(Map<String, dynamic> json) {
    return Task(
      id: json['id'],
      title: json['title'],
      description: json['description'],
      completed: json['completed'],
      userId: json['user_id'],
      dueDate: DateTime.parse(json['due_date']),
      priority: json['priority'],
      createdAt: DateTime.parse(json['created_at']),
      updatedAt: DateTime.parse(json['updated_at']),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'description': description,
      'completed': completed,
      'user_id': userId,
      'due_date': dueDate.toIso8601String(),
      'priority': priority,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }

  String get priorityText {
    switch (priority) {
      case 0:
        return 'Низкий';
      case 1:
        return 'Средний';
      case 2:
        return 'Высокий';
      default:
        return 'Неизвестный';
    }
  }

  Color get priorityColor {
    switch (priority) {
      case 0:
        return Colors.green;
      case 1:
        return Colors.orange;
      case 2:
        return Colors.red;
      default:
        return Colors.grey;
    }
  }
} 