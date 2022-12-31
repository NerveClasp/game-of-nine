import { writable, type Writable } from 'svelte/store';
import type { ChatMessage, IncomingMessage, Message, OutgoingMessage, Player } from './types';

const messageStore: Writable<IncomingMessage | null> = writable(null);

export const chatMessages: Writable<ChatMessage[]> = writable([]);

export const clientUid = writable('');

const socket =
	typeof window !== 'undefined'
		? new WebSocket(
				`ws://${typeof window !== 'undefined' ? document.location.host : 'localhost:7331'}/ws`
		  )
		: { onopen: () => {}, onmessage: () => {}, readyState: -1, send: console.log };

// Connection opened
socket.onopen = (e: Event) => {
	console.log('Websocket is open', e);
};

// Listen for messages
socket.onmessage = (e: MessageEvent) => {
	try {
		let message: IncomingMessage = JSON.parse(e.data);
		chatMessages.update((messages) => {
			if (message?.type !== 'chat') return messages;
			return [...messages, message];
		});
		messageStore.set(message);
	} catch (error) {
		console.log('Whoa, got something strange from socket', error);
	}
};

const sendMessage = (message: string | Message, player: Player) => {
	if (!message) return;

	if (socket.readyState <= 1) {
		const outgoingMessage: OutgoingMessage =
			typeof message === 'object' ? { ...message, player } : { type: 'chat', message, player };

		socket.send(JSON.stringify(outgoingMessage));
	}
};

export default {
	subscribe: messageStore.subscribe,
	sendMessage
};
