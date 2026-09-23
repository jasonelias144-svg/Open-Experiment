import React, { useState, useEffect, useRef } from 'react';
import { StyleSheet, Text, View, FlatList, TextInput, TouchableOpacity, SafeAreaView } from 'react-native';
export default function ConfluenceRoom({ roomId, clientId, senderType }) {
  const [messages, setMessages] = useState([]);
  const [inputText, setInputText] = useState('');
  const ws = useRef(null);

  useEffect(() => {
    // Mount connection interface string dynamically
    ws.current = new WebSocket(`ws://localhost:8080/ws?room=${roomId}&type=${senderType}&id=${clientId}`);

    ws.current.onmessage = (e) => {
      const incomingMessage = JSON.parse(e.data);
      setMessages((prev) => [...prev, incomingMessage]);
    };

    return () => {
      if (ws.current) ws.current.close();
    };
  }, [roomId]);

  const transmitMessage = () => {
    if (!inputText.trim()) return;
    
    const packet = {
      content: inputText,
      contextType: 'philosophical', // Statically mapping a target segment space for initial seed test
    };

    if (ws.current && ws.current.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify(packet));
      setInputText('');
    }
  };

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.headerTitle}>MATRIX ROOM // {roomId.toUpperCase()}</Text>
      </View>
      <FlatList
        data={messages}
        keyExtractor={(item, index) => index.toString()}
        renderItem={({ item }) => (
          <View style={[styles.bubble, item.senderType === 'ai' ? styles.aiBubble : styles.humanBubble]}>
            <Text style={styles.metaText}>[{item.senderType.toUpperCase()}] ID: {item.senderID}</Text>
            <Text style={styles.bodyText}>{item.content}</Text>
          </View>
        )}
      />
      <View style={styles.inputContainer}>
        <TextInput
          style={styles.input}
          value={inputText}
          onChangeText={setInputText}
          placeholder="Transmit into matrix..."
          placeholderTextColor="#666"
        />
        <TouchableOpacity style={styles.sendButton} onPress={transmitMessage}>
          <Text style={styles.sendText}>➡️</Text>
        </TouchableOpacity>
      </View>
    </SafeAreaView>
  );
}
const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#0A0A0F' },
  header: { padding: 15, borderBottomWidth: 1, borderBottomColor: '#00FF66' },
  headerTitle: { color: '#00FF66', fontFamily: 'monospace', fontWeight: 'bold' },
  bubble: { padding: 12, marginVertical: 6, marginHorizontal: 12, borderRadius: 6 },
  humanBubble: { backgroundColor: '#141420', alignSelf: 'flex-start' },
  aiBubble: { backgroundColor: '#002B11', alignSelf: 'flex-end', borderColor: '#00FF66', borderWidth: 0.5 },
  metaText: { fontSize: 10, color: '#666', fontFamily: 'monospace', marginBottom: 4 },
  bodyText: { fontSize: 14, color: '#E0E0E0', fontFamily: 'monospace' },
  inputContainer: { flexDirection: 'row', padding: 12, backgroundColor: '#101018' },
  input: { flex: 1, backgroundColor: '#050508', color: '#00FF66', fontFamily: 'monospace', paddingHorizontal: 12, borderRadius: 4, height: 40 },
  sendButton: { justifyContent: 'center', alignItems: 'center', paddingHorizontal: 15 },
  sendText: { fontSize: 18 }
});
