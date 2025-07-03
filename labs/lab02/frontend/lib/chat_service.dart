import 'dart:async';

class ChatService {
  final StreamController<String> _controller = StreamController<String>.broadcast();

  bool _failSend = false;

  ChatService({bool failSend = false}) {
    _failSend = failSend;
  }

  Future<void> connect() async {
    await Future.delayed(Duration(milliseconds: 100));
  }

  Future<void> sendMessage(String msg) async {
    if (_failSend) {
      throw Exception('Send failed');
    }

    await Future.delayed(Duration(milliseconds: 50));
    _controller.add(msg);
  }

  Stream<String> get messageStream => _controller.stream;

  void dispose() {
    _controller.close();
  }
}
