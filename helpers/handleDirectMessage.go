package helpers

import (
    "fmt"
    "time"

    "reddit_engine/messages"
)

func (e *EngineActor) HandleSendDirectMessage(msg *SendDirectMessage) (*SendDirectMessageResponse, error) {
    if !e.Users[msg.FromUser] || !e.Users[msg.ToUser] {
        return nil, fmt.Errorf("invalid users: sender or recipient not found")
    }

    messageID := fmt.Sprintf("dm_%d", len(e.DirectMessages)+1)
    dm := &DirectMessage{
        ID:        messageID,
        FromUser:  msg.FromUser,
        ToUser:    msg.ToUser,
        Content:   msg.Content,
        Timestamp: time.Now().Unix(),
        IsRead:    false,
    }

    e.DirectMessages[messageID] = dm
    
    if e.UserInbox[msg.ToUser] == nil {
        e.UserInbox[msg.ToUser] = make([]string, 0)
    }
    e.UserInbox[msg.ToUser] = append(e.UserInbox[msg.ToUser], messageID)

    return &SendDirectMessageResponse{
        Success: true,
        Message: dm,
    }, nil
}

func (e *EngineActor) HandleGetDirectMessages(msg *GetDirectMessages) (*GetDirectMessagesResponse, error) {
    if !e.Users[msg.Username] {
        return nil, fmt.Errorf("user not found: %s", msg.Username)
    }

    messageIDs := e.UserInbox[msg.Username]
    messages := make([]*DirectMessage, 0, len(messageIDs))

    for _, msgID := range messageIDs {
        if dm, exists := e.DirectMessages[msgID]; exists {
            messages = append(messages, dm)
        }
    }

    return &GetDirectMessagesResponse{Messages: messages}, nil
}

func (e *EngineActor) HandleReplyDirectMessage(msg *messages.ReplyDirectMessage) (*messages.ReplyDirectMessageResponse, error) {
    // Check if original message exists
    originalMsg, exists := e.DirectMessages[msg.ReplyToMessageId]
    if !exists {
        return &messages.ReplyDirectMessageResponse{Success: false}, fmt.Errorf("original message not found")
    }

    // Verify sender exists and is part of the conversation
    if !e.Users[msg.FromUser] {
        return &messages.ReplyDirectMessageResponse{Success: false}, fmt.Errorf("user not found")
    }

    if msg.FromUser != originalMsg.ToUser && msg.FromUser != originalMsg.FromUser {
        return &messages.ReplyDirectMessageResponse{Success: false}, fmt.Errorf("unauthorized to reply to this message")
    }

    // Create new message
    messageID := fmt.Sprintf("dm_%d", len(e.DirectMessages)+1)
    toUser := originalMsg.FromUser
    if msg.FromUser == originalMsg.FromUser {
        toUser = originalMsg.ToUser
    }

    dm := &DirectMessage{
        ID:        messageID,
        FromUser:  msg.FromUser,
        ToUser:    toUser,
        Content:   msg.Content,
        Timestamp: time.Now().Unix(),
        IsRead:    false,
    }

    e.DirectMessages[messageID] = dm
    
    if e.UserInbox[toUser] == nil {
        e.UserInbox[toUser] = make([]string, 0)
    }
    e.UserInbox[toUser] = append(e.UserInbox[toUser], messageID)

    return &messages.ReplyDirectMessageResponse{
        Success: true,
        Message: &messages.DirectMessage{
            Id:        dm.ID,
            FromUser:  dm.FromUser,
            ToUser:    dm.ToUser,
            Content:   dm.Content,
            Timestamp: dm.Timestamp,
            IsRead:    dm.IsRead,
        },
    }, nil
}
