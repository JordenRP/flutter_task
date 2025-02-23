import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/category.dart';
import '../models/task.dart';
import 'auth_service.dart';

class CategoryService {
  static const baseUrl = 'http://localhost:8080/api/categories';

  Future<Map<String, String>> _getHeaders() async {
    final token = await AuthService.getToken();
    return {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer $token',
    };
  }

  Future<List<Category>> getCategories() async {
    final response = await http.get(
      Uri.parse(baseUrl),
      headers: await _getHeaders(),
    );

    if (response.statusCode == 200) {
      final List<dynamic> data = jsonDecode(response.body);
      return data.map((json) => Category.fromJson(json)).toList();
    } else {
      throw Exception('Failed to load categories');
    }
  }

  Future<Category> createCategory(String name) async {
    final response = await http.post(
      Uri.parse(baseUrl),
      headers: await _getHeaders(),
      body: jsonEncode({
        'name': name,
      }),
    );

    if (response.statusCode == 201) {
      return Category.fromJson(jsonDecode(response.body));
    } else {
      throw Exception('Failed to create category');
    }
  }

  Future<void> deleteCategory(int id) async {
    final response = await http.delete(
      Uri.parse('$baseUrl/$id'),
      headers: await _getHeaders(),
    );

    if (response.statusCode != 204) {
      throw Exception('Failed to delete category');
    }
  }

  Future<List<Task>> getTasksByCategory(int categoryId) async {
    final response = await http.get(
      Uri.parse('$baseUrl/$categoryId/tasks'),
      headers: await _getHeaders(),
    );

    if (response.statusCode == 200) {
      final List<dynamic> data = jsonDecode(response.body);
      return data.map((json) => Task.fromJson(json)).toList();
    } else {
      throw Exception('Failed to load tasks for category');
    }
  }

  Future<void> updateTaskCategory(int taskId, int categoryId) async {
    final response = await http.put(
      Uri.parse('$baseUrl/tasks/$taskId'),
      headers: await _getHeaders(),
      body: jsonEncode({
        'category_id': categoryId,
      }),
    );

    if (response.statusCode != 200) {
      throw Exception('Failed to update task category');
    }
  }
} 