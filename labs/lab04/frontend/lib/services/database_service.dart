import 'package:sqflite/sqflite.dart';
import 'package:path/path.dart';
import '../models/user.dart';

class DatabaseService {
  static Database? _database;
  static const String _dbName = 'lab04_app.db';
  static const int _version = 1;

  // database getter
  static Future<Database> get database async {
    if (_database != null) return _database!;
    _database = await _initDatabase();
    return _database!;
  }

  // _initDatabase method
  static Future<Database> _initDatabase() async {
    // Get the databases path
    String databasesPath = await getDatabasesPath();
    String path = join(databasesPath, _dbName);

    // Open database with version and callbacks
    return await openDatabase(
      path,
      version: _version,
      onCreate: _onCreate,
      onUpgrade: _onUpgrade,
    );
  }

  // _onCreate method
  static Future<void> _onCreate(Database db, int version) async {
    // Create users table
    await db.execute('''
      CREATE TABLE users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        email TEXT NOT NULL UNIQUE,
        created_at TEXT NOT NULL,
        updated_at TEXT NOT NULL
      )
    ''');

    // Create posts table
    await db.execute('''
      CREATE TABLE posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        published INTEGER NOT NULL DEFAULT 0,
        created_at TEXT NOT NULL,
        updated_at TEXT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
      )
    ''');
  }

  // _onUpgrade method
  static Future<void> _onUpgrade(
      Database db, int oldVersion, int newVersion) async {
    // For now, you can leave this empty or add migration logic later
    if (oldVersion < 2) {
      // Example: Add new column in future version
      // await db.execute('ALTER TABLE users ADD COLUMN avatar TEXT');
    }
  }

  // User CRUD operations

  // createUser method
  static Future<User> createUser(CreateUserRequest request) async {
    final db = await database;
    final now = DateTime.now().toIso8601String();

    final id = await db.insert('users', {
      'name': request.name,
      'email': request.email,
      'created_at': now,
      'updated_at': now,
    });

    return User(
      id: id,
      name: request.name,
      email: request.email,
      createdAt: DateTime.parse(now),
      updatedAt: DateTime.parse(now),
    );
  }

  // getUser method
  static Future<User?> getUser(int id) async {
    final db = await database;
    final List<Map<String, dynamic>> maps = await db.query(
      'users',
      where: 'id = ?',
      whereArgs: [id],
    );

    if (maps.isEmpty) return null;

    final map = maps.first;
    return User(
      id: map['id'],
      name: map['name'],
      email: map['email'],
      createdAt: DateTime.parse(map['created_at']),
      updatedAt: DateTime.parse(map['updated_at']),
    );
  }

  // getAllUsers method
  static Future<List<User>> getAllUsers() async {
    final db = await database;
    final List<Map<String, dynamic>> maps = await db.query(
      'users',
      orderBy: 'created_at DESC',
    );

    return List.generate(maps.length, (i) {
      final map = maps[i];
      return User(
        id: map['id'],
        name: map['name'],
        email: map['email'],
        createdAt: DateTime.parse(map['created_at']),
        updatedAt: DateTime.parse(map['updated_at']),
      );
    });
  }

  // updateUser method
  static Future<User> updateUser(int id, Map<String, dynamic> updates) async {
    final db = await database;
    final now = DateTime.now().toIso8601String();
    
    // Add updated_at to updates
    updates['updated_at'] = now;

    await db.update(
      'users',
      updates,
      where: 'id = ?',
      whereArgs: [id],
    );

    // Get updated user
    final updatedUser = await getUser(id);
    if (updatedUser == null) {
      throw Exception('User not found after update');
    }

    return updatedUser;
  }

  // deleteUser method
  static Future<void> deleteUser(int id) async {
    final db = await database;
    await db.delete(
      'users',
      where: 'id = ?',
      whereArgs: [id],
    );
  }

  // getUserCount method
  static Future<int> getUserCount() async {
    final db = await database;
    final result = await db.rawQuery('SELECT COUNT(*) as count FROM users');
    return Sqflite.firstIntValue(result) ?? 0;
  }

  // searchUsers method
  static Future<List<User>> searchUsers(String query) async {
    final db = await database;
    final searchTerm = '%$query%';
    
    final List<Map<String, dynamic>> maps = await db.query(
      'users',
      where: 'name LIKE ? OR email LIKE ?',
      whereArgs: [searchTerm, searchTerm],
      orderBy: 'name ASC',
    );

    return List.generate(maps.length, (i) {
      final map = maps[i];
      return User(
        id: map['id'],
        name: map['name'],
        email: map['email'],
        createdAt: DateTime.parse(map['created_at']),
        updatedAt: DateTime.parse(map['updated_at']),
      );
    });
  }

  // Database utility methods

  // closeDatabase method
  static Future<void> closeDatabase() async {
    if (_database != null) {
      await _database!.close();
      _database = null;
    }
  }

  // clearAllData method
  static Future<void> clearAllData() async {
    final db = await database;
    await db.delete('posts');
    await db.delete('users');
  }

  // getDatabasePath method
  static Future<String> getDatabasePath() async {
    String databasesPath = await getDatabasesPath();
    return join(databasesPath, _dbName);
  }
}
