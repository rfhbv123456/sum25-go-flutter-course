import 'package:shared_preferences/shared_preferences.dart';
import 'dart:convert';

class PreferencesService {
  static SharedPreferences? _prefs;

  // init method
  static Future<void> init() async {
    _prefs = await SharedPreferences.getInstance();
  }

  // setString method
  static Future<void> setString(String key, String value) async {
    if (_prefs == null) {
      throw Exception('PreferencesService not initialized. Call init() first.');
    }
    await _prefs!.setString(key, value);
  }

  // getString method
  static String? getString(String key) {
    if (_prefs == null) {
      throw Exception('PreferencesService not initialized. Call init() first.');
    }
    return _prefs!.getString(key);
  }

  // setInt method
  static Future<void> setInt(String key, int value) async {
    if (_prefs == null) {
      throw Exception('PreferencesService not initialized. Call init() first.');
    }
    await _prefs!.setInt(key, value);
  }

  // getInt method
  static int? getInt(String key) {
    if (_prefs == null) {
      throw Exception('PreferencesService not initialized. Call init() first.');
    }
    return _prefs!.getInt(key);
  }

  // setBool method
  static Future<void> setBool(String key, bool value) async {
    if (_prefs == null) {
      throw Exception('PreferencesService not initialized. Call init() first.');
    }
    await _prefs!.setBool(key, value);
  }

  // getBool method
  static bool? getBool(String key) {
    if (_prefs == null) {
      throw Exception('PreferencesService not initialized. Call init() first.');
    }
    return _prefs!.getBool(key);
  }

  // setStringList method
  static Future<void> setStringList(String key, List<String> value) async {
    if (_prefs == null) {
      throw Exception('PreferencesService not initialized. Call init() first.');
    }
    await _prefs!.setStringList(key, value);
  }

  // getStringList method
  static List<String>? getStringList(String key) {
    if (_prefs == null) {
      throw Exception('PreferencesService not initialized. Call init() first.');
    }
    return _prefs!.getStringList(key);
  }

  // setObject method
  static Future<void> setObject(String key, Map<String, dynamic> value) async {
    if (_prefs == null) {
      throw Exception('PreferencesService not initialized. Call init() first.');
    }
    final jsonString = jsonEncode(value);
    await _prefs!.setString(key, jsonString);
  }

  // getObject method
  static Map<String, dynamic>? getObject(String key) {
    if (_prefs == null) {
      throw Exception('PreferencesService not initialized. Call init() first.');
    }
    final jsonString = _prefs!.getString(key);
    if (jsonString == null) return null;
    
    try {
      return jsonDecode(jsonString) as Map<String, dynamic>;
    } catch (e) {
      return null;
    }
  }

  // remove method
  static Future<void> remove(String key) async {
    if (_prefs == null) {
      throw Exception('PreferencesService not initialized. Call init() first.');
    }
    await _prefs!.remove(key);
  }

  // clear method
  static Future<void> clear() async {
    if (_prefs == null) {
      throw Exception('PreferencesService not initialized. Call init() first.');
    }
    await _prefs!.clear();
  }

  // containsKey method
  static bool containsKey(String key) {
    if (_prefs == null) {
      throw Exception('PreferencesService not initialized. Call init() first.');
    }
    return _prefs!.containsKey(key);
  }

  // getAllKeys method
  static Set<String> getAllKeys() {
    if (_prefs == null) {
      throw Exception('PreferencesService not initialized. Call init() first.');
    }
    return _prefs!.getKeys();
  }
}
