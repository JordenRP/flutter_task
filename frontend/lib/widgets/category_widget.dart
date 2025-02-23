import 'package:flutter/material.dart';
import '../models/category.dart';
import '../services/category_service.dart';

class CategoryWidget extends StatefulWidget {
  final Function(Category?) onCategorySelected;
  final Category? selectedCategory;

  const CategoryWidget({
    Key? key,
    required this.onCategorySelected,
    this.selectedCategory,
  }) : super(key: key);

  @override
  CategoryWidgetState createState() => CategoryWidgetState();
}

class CategoryWidgetState extends State<CategoryWidget> {
  final CategoryService _categoryService = CategoryService();
  final TextEditingController _categoryController = TextEditingController();
  List<Category> _categories = [];
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    _loadCategories();
  }

  Future<void> _loadCategories() async {
    setState(() {
      _isLoading = true;
    });
    try {
      final categories = await _categoryService.getCategories();
      setState(() {
        _categories = categories;
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _isLoading = false;
      });
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(e.toString())),
        );
      }
    }
  }

  Future<void> _createCategory() async {
    if (_categoryController.text.isEmpty) return;

    try {
      final category = await _categoryService.createCategory(_categoryController.text);
      setState(() {
        _categories.add(category);
        _categoryController.clear();
      });
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(e.toString())),
        );
      }
    }
  }

  Future<void> _deleteCategory(Category category) async {
    try {
      await _categoryService.deleteCategory(category.id);
      setState(() {
        _categories.removeWhere((c) => c.id == category.id);
        if (widget.selectedCategory?.id == category.id) {
          widget.onCategorySelected(null);
        }
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
    return Container(
      width: 250,
      color: Theme.of(context).cardColor,
      child: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(8.0),
            child: Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: _categoryController,
                    decoration: const InputDecoration(
                      labelText: 'Новая категория',
                      border: OutlineInputBorder(),
                    ),
                  ),
                ),
                IconButton(
                  icon: const Icon(Icons.add),
                  onPressed: _createCategory,
                ),
              ],
            ),
          ),
          const Divider(),
          ListTile(
            title: const Text('Все задачи'),
            selected: widget.selectedCategory == null,
            onTap: () => widget.onCategorySelected(null),
          ),
          if (_isLoading)
            const Center(child: CircularProgressIndicator())
          else
            Expanded(
              child: ListView.builder(
                itemCount: _categories.length,
                itemBuilder: (context, index) {
                  final category = _categories[index];
                  return ListTile(
                    title: Text(category.name),
                    selected: widget.selectedCategory?.id == category.id,
                    trailing: IconButton(
                      icon: const Icon(Icons.delete),
                      onPressed: () => _deleteCategory(category),
                    ),
                    onTap: () => widget.onCategorySelected(category),
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
    _categoryController.dispose();
    super.dispose();
  }
} 