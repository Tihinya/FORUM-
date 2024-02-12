export class SendMessageEvent {
    constructor(message, sender_id, receiver_id) {
        this.message = message
        this.sender_id = sender_id
        this.receiver_id = receiver_id
    }
}
export class NewMessageEvent {
    constructor(message, sender_id, sent_date, receiver_id) {
        this.message = message
        this.sender_id = sender_id
        this.receiver_id = receiver_id
        this.sent_date = sent_date
    }
}

export class RequestMessageHistoryEvent {
    constructor(receiver_id) {
        this.receiver_id = receiver_id
    }
}

export class IsTypingEvent {
    constructor(sender_id, receiver_id, typing_status) {
        this.sender_id = sender_id
        this.receiver_id = receiver_id
        this.typing_status = typing_status
    }
}

export class Event {
    constructor(type, payload) {
        this.type = type
        this.payload = payload
    }
}