import 'package:flutter/material.dart';

class Category {
  final int id;
  final String name;
  final int userId;
  final DateTime createdAt;

  Category({
    required this.id,
    required this.name,
    required this.userId,
    required this.createdAt,
  });

  factory Category.fromJson(Map<String, dynamic> json) {
    return Category(
      id: json['id'],
      name: json['name'],
      userId: json['user_id'],
      createdAt: DateTime.parse(json['created_at']),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'user_id': userId,
      'created_at': createdAt.toIso8601String(),
    };
  }
} 