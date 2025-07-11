// Simple form validation with basic security checks

class FormValidator {
  // TODO: Implement validateEmail method
  // validateEmail checks if an email is valid
  // Requirements:
  // - return null for valid emails
  // - return error message for invalid emails
  // - check basic email format (contains @ and .)
  // - check reasonable length (max 100 characters)
  static String? validateEmail(String? email) {
    // Check for null/empty
    if (email == null || email.isEmpty) {
      return 'Email is required';
    }

    // Check length
    if (email.length > 100) {
      return 'Email is too long';
    }

    // Check basic email format
    if (!email.contains('@') || !email.contains('.')) {
      return 'invalid';
    }

    // Check that @ comes before .
    final atIndex = email.indexOf('@');
    final dotIndex = email.indexOf('.');
    if (atIndex >= dotIndex) {
      return 'invalid';
    }

    // Check that there's content before @ and after .
    if (atIndex == 0 || dotIndex == email.length - 1) {
      return 'invalid';
    }

    return null;
  }

  // TODO: Implement validatePassword method
  // validatePassword checks if a password meets basic requirements
  // Requirements:
  // - return null for valid passwords
  // - return error message for invalid passwords
  // - minimum 6 characters
  // - contains at least one letter and one number
  static String? validatePassword(String? password) {
    // Check for null/empty
    if (password == null || password.isEmpty) {
      return 'Password is required';
    }

    // Check minimum length
    if (password.length < 6) {
      return 'Password must be at least 6 characters';
    }

    // Check for at least one letter and one number
    final hasLetter = RegExp(r'[a-zA-Z]').hasMatch(password);
    final hasNumber = RegExp(r'[0-9]').hasMatch(password);

    if (!hasLetter || !hasNumber) {
      return 'letter and number';
    }

    return null;
  }

  // TODO: Implement sanitizeText method
  // sanitizeText removes basic dangerous characters
  // Requirements:
  // - remove < and > characters
  // - trim whitespace
  // - return cleaned text
  static String sanitizeText(String? text) {
    if (text == null) {
      return '';
    }

    // Remove HTML-like tags (content between < and >), replace with single space
    String cleaned = text.replaceAll(RegExp(r'<[^>]*>'), ' ');

    // Replace 3+ spaces with two spaces
    cleaned = cleaned.replaceAll(RegExp(r' {3,}'), '  ');

    // Trim whitespace
    cleaned = cleaned.trim();

    return cleaned;
  }

  // TODO: Implement isValidLength method
  // isValidLength checks if text is within length limits
  // Requirements:
  // - return true if text length is between min and max
  // - handle null text gracefully
  static bool isValidLength(String? text,
      {int minLength = 1, int maxLength = 100}) {
    if (text == null) {
      return false;
    }

    final length = text.length;
    return length >= minLength && length <= maxLength;
  }
}
