import 'dart:convert';
import 'package:http/http.dart' as http;
import 'auth_service.dart';

class TaskNotification {
  final int id;
  final int userId;
  final int taskId;
  final String message;
  final DateTime createdAt;
  final bool read;

  TaskNotification({
    required this.id,
    required this.userId,
    required this.taskId,
    required this.message,
    required this.createdAt,
    required this.read,
  });

  factory TaskNotification.fromJson(Map<String, dynamic> json) {
    return TaskNotification(
      id: json['id'],
      userId: json['user_id'],
      taskId: json['task_id'],
      message: json['message'],
      createdAt: DateTime.parse(json['created_at']),
      read: json['read'],
    );
  }
}

class NotificationService {
  static const baseUrl = 'http://localhost:8080/api/notifications';

  Future<Map<String, String>> _getHeaders() async {
    final token = await AuthService.getToken();
    return {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer $token',
    };
  }

  Future<List<TaskNotification>> getNotifications() async {
    final response = await http.get(
      Uri.parse(baseUrl),
      headers: await _getHeaders(),
    );

    if (response.statusCode == 200) {
      final List<dynamic> data = jsonDecode(response.body);
      return data.map((json) => TaskNotification.fromJson(json)).toList();
    } else {
      throw Exception('Failed to load notifications');
    }
  }

  Future<void> markAsRead(int id) async {
    final response = await http.post(
      Uri.parse('$baseUrl/$id/read'),
      headers: await _getHeaders(),
    );

    if (response.statusCode != 200) {
      throw Exception('Failed to mark notification as read');
    }
  }
} 